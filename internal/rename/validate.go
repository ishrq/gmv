package rename

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ValidateFiles(files []string) error {
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

func ValidateEdits(original, edited []string) error {
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

func CheckOverwrites(plan []RenameOp, originalFiles []string) []string {
	// Create a set of original files for quick lookup
	originals := make(map[string]bool)
	for _, file := range originalFiles {
		originals[file] = true
	}

	var overwrites []string

	for _, op := range plan {
		// Skip temp files (used for swaps)
		baseName := filepath.Base(op.To)
		if strings.HasPrefix(baseName, ".gmv_temp_") {
			continue
		}

		// Check if target exists and is NOT in the original list
		if _, err := os.Stat(op.To); err == nil {
			if !originals[op.To] {
				// File exists and is not in our rename list - would be overwritten!
				overwrites = append(overwrites, op.To)
			}
		}
	}

	return overwrites
}
