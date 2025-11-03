# gmv

A CLI tool for batch renaming files using your $EDITOR.

## Usage

Example commands:

```bash
gmv test-file.go        # rename test-file.go in the editor
gmv *                   # rename all files in the editor
gmv *.{pdf,epub}        # rename all pdf and epub files
gmv */                  # rename all directories
gmv test-dir/*.txt      # rename all text files in test-dir
gmv */*                 # rename all files in all directories
gmv --dry-run *         # print changes without applying
gmv --help              # print help
```

## Installation

```bash
go install github.com/ishrq/gmv@latest
```
