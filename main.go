package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type RenameOp struct {
	From string
	To   string
}

func parseArgs() (files []string, dryRun bool, err error) {
	args := os.Args[1:]

	for _, arg := range args {
		if arg == "--dry-run" {
			dryRun = true
		} else {
			files = append(files, arg)
		}
	}

	if len(files) == 0 {
		return nil, false, fmt.Errorf("no files specified")
	}

	return files, dryRun, nil
}

func validateFiles(files []string) error {
	seen := make(map[string]bool)

	for _, file := range files {
		// Check for duplicates
		if seen[file] {
			return fmt.Errorf("duplicate file specified: %s", file)
		}
		seen[file] = true

		// Check if file exists
		if _, err := os.Stat(file); os.IsNotExist(err) {
			return fmt.Errorf("file does not exist: %s", file)
		}
	}

	return nil
}

func createTempFile(files []string) (string, error) {
	tmpFile, err := os.CreateTemp("", "gmv-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tmpFile.Close()

	// Write each file path on its own line
	for _, file := range files {
		if _, err := tmpFile.WriteString(file + "\n"); err != nil {
			return "", fmt.Errorf("failed to write file path: %w", err)
		}
	}

	return tmpFile.Name(), nil
}

func launchEditor(filepath string) error {
	// Get editor from environment or use fallback
	editor := os.Getenv("EDITOR")
	if editor == "" {
		// Try vi first, then nano
		if _, err := exec.LookPath("vi"); err == nil {
			editor = "vi"
		} else if _, err := exec.LookPath("nano"); err == nil {
			editor = "nano"
		} else {
			return fmt.Errorf("no editor found: $EDITOR not set and neither vi nor nano are available")
		}
	}

	// Launch the editor
	cmd := exec.Command(editor, filepath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("editor exited with error: %w", err)
	}

	return nil
}

func parseEdited(filepath string) ([]string, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read edited file: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	var edited []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			edited = append(edited, line)
		}
	}

	return edited, nil
}

func validateEdits(original, edited []string) error {
	// Check line count matches
	if len(original) != len(edited) {
		return fmt.Errorf("line count mismatch: expected %d lines, got %d lines", len(original), len(edited))
	}

	// Track target filenames to detect duplicates
	targets := make(map[string]bool)

	for i := 0; i < len(original); i++ {
		origPath := original[i]
		editPath := edited[i]

		// Check that directory hasn't changed
		origDir := filepath.Dir(origPath)
		editDir := filepath.Dir(editPath)

		if origDir != editDir {
			return fmt.Errorf("cannot move files to different directories: %s -> %s", origPath, editPath)
		}

		// Check for duplicate target filenames
		if targets[editPath] {
			return fmt.Errorf("duplicate target filename: %s", editPath)
		}
		targets[editPath] = true
	}

	return nil
}

func buildRenamePlan(original, edited []string) ([]RenameOp, error) {
	var plan []RenameOp

	for i := 0; i < len(original); i++ {
		// Skip if no change
		if original[i] == edited[i] {
			continue
		}

		plan = append(plan, RenameOp{
			From: original[i],
			To:   edited[i],
		})
	}

	return plan, nil
}

func detectCycles(plan []RenameOp) [][]string {
	// Build adjacency map: from -> to
	graph := make(map[string]string)
	for _, op := range plan {
		graph[op.From] = op.To
	}

	visited := make(map[string]bool)
	recStack := make(map[string]bool)
	var cycles [][]string

	var dfs func(node string, path []string) bool
	dfs = func(node string, path []string) bool {
		visited[node] = true
		recStack[node] = true
		path = append(path, node)

		if next, exists := graph[node]; exists {
			if recStack[next] {
				// Found a cycle, extract it from path
				cycleStart := -1
				for i, n := range path {
					if n == next {
						cycleStart = i
						break
					}
				}
				if cycleStart != -1 {
					cycle := make([]string, len(path)-cycleStart)
					copy(cycle, path[cycleStart:])
					cycles = append(cycles, cycle)
				}
				return true
			} else if !visited[next] {
				if dfs(next, path) {
					return true
				}
			}
		}

		recStack[node] = false
		return false
	}

	// Run DFS from each unvisited node
	for _, op := range plan {
		if !visited[op.From] {
			dfs(op.From, []string{})
		}
	}

	return cycles
}

func main() {
	files, dryRun, err := parseArgs()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if err := validateFiles(files); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	tempFilePath, err := createTempFile(files)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if err := launchEditor(tempFilePath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	editedFiles, err := parseEdited(tempFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if err := validateEdits(files, editedFiles); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	plan, err := buildRenamePlan(files, editedFiles)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	cycles := detectCycles(plan)
	if len(cycles) > 0 {
		fmt.Printf("Detected cycles:\n")
		for _, cycle := range cycles {
			fmt.Printf("  %v\n", cycle)
		}
	}

	fmt.Printf("Rename plan:\n")
	for _, op := range plan {
		fmt.Printf("  %s -> %s\n", op.From, op.To)
	}
	fmt.Printf("Dry run: %v\n", dryRun)
	fmt.Printf("Temp file: %s\n", tempFilePath)
}
