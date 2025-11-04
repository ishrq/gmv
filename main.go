package main

import (
	"fmt"
	"os"
)

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

func createTempFile() (string, error) {
	tmpFile, err := os.CreateTemp("", "gmv-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tmpFile.Close()

	header := `# Edit filenames below. Save and exit to apply changes.
# Lines will be mapped in order to the original files.
# Do not delete lines - this will cause an error.
`
	if _, err := tmpFile.WriteString(header); err != nil {
		return "", fmt.Errorf("failed to write header: %w", err)
	}

	return tmpFile.Name(), nil
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

	tempFilePath, err := createTempFile()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Files to rename: %v\n", files)
	fmt.Printf("Dry run: %v\n", dryRun)
	fmt.Printf("Temp file created: %s\n", tempFilePath)
}
