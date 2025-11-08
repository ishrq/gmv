package main

import (
	"fmt"
	"os"

	"github.com/ishrq/gmv/internal/rename"
)

func printHelp() {
	help := `gmv - Batch rename files using $EDITOR

	USAGE:
	gmv [OPTIONS] <files>...

	OPTIONS:
	--dry-run    Print changes without applying them
	--help, -h   Show this help message

	EXAMPLES:
	gmv test-file.go        # Rename test-file.go in the editor
	gmv *                   # Rename all files in the editor
	gmv *.{pdf,epub}        # Rename all pdf and epub files
	gmv */                  # Rename all directories
	gmv test-dir/*.txt      # Rename all text files in test-dir
	gmv */*                 # Rename all files in all directories
	gmv --dry-run *         # Preview changes without applying
	gmv --help              # Print help

	DESCRIPTION:
	gmv opens your $EDITOR with a list of files to rename. Edit the filenames,
	save and exit. The files will be renamed accordingly. File swaps are
	automatically handled using temporary files.

	A log of all rename operations is saved in your system's temp directory.
	`
	fmt.Print(help)
}

func parseArgs() (files []string, dryRun bool, err error) {
	args := os.Args[1:]

	// Check for help flag first
	for _, arg := range args {
		if arg == "--help" || arg == "-h" {
			printHelp()
			os.Exit(0)
		}
	}

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

func main() {
	files, dryRun, err := parseArgs()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if err := rename.ValidateFiles(files); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	tempFilePath, err := rename.CreateTempFile(files)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if err := rename.LaunchEditor(tempFilePath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	editedFiles, err := rename.ParseEdited(tempFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if err := rename.ValidateEdits(files, editedFiles); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	plan, err := rename.BuildRenamePlan(files, editedFiles)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Check if there are any changes
	if len(plan) == 0 {
		fmt.Println("No files were renamed.")
		os.Exit(0)
	}

	if err := rename.ExecuteRenames(plan, dryRun); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if !dryRun {
		logPath, err := rename.WriteLog(plan)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to write log: %v\n", err)
		} else {
			fmt.Printf("Successfully renamed files.\n")
			fmt.Printf("A log file is saved at %s\n", logPath)
		}
	}
}
