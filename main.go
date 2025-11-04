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
		if seen[file] {
			return fmt.Errorf("duplicate file specified: %s", file)
		}
		seen[file] = true

		if _, err := os.Stat(file); os.IsNotExist(err) {
			return fmt.Errorf("file does not exist: %s", file)
		}
	}

	return nil
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

	fmt.Printf("Files to rename: %v\n", files)
	fmt.Printf("Dry run: %v\n", dryRun)
}
