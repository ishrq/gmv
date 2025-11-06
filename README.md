# gmv

A CLI tool for batch renaming files using your $EDITOR.

## Usage

Example commands:

```bash
gmv *                   # Rename all files in current directory
gmv *.{pdf,epub}        # Rename all PDF and EPUB files
gmv test-dir/*.txt      # Rename all text files in test-dir
gmv file1.go file2.go   # Rename specific files
gmv */                  # Rename all directories
gmv --dry-run *         # Preview changes without applying
gmv --help
```

## Installation

```bash
go install github.com/ishrq/gmv@latest
```
