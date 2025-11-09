# gmv

[![CI](https://github.com/ishrq/gmv/actions/workflows/ci.yml/badge.svg)](https://github.com/ishrq/gmv/actions/workflows/ci.yml)
[![Release](https://github.com/ishrq/gmv/actions/workflows/release.yml/badge.svg)](https://github.com/ishrq/gmv/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/ishrq/gmv)](https://goreportcard.com/report/github.com/ishrq/gmv)

A powerful CLI tool for batch renaming files and directories using your `$EDITOR`.

## Features

- **Batch rename** files and directories in your text editor
- **Overwrite protection** - warns and prompts before overwriting existing files
- **Automatic swap handling** - swap file names without conflicts using temporary files
- **Dry-run mode** - preview changes before applying them
- **Operation logging** - keeps a temporary log of all rename operations
- **Comprehensive validation** - prevents directory moves, detects duplicates
- **Cross-platform** - works on Linux, macOS, BSD systems, and Android

## Installation

### Using Go

```bash
go install github.com/ishrq/gmv@latest
```

### Pre-built binaries

Download pre-built binaries from the [releases page](https://github.com/ishrq/gmv/releases).

### From source

```bash
git clone https://github.com/ishrq/gmv.git
cd gmv
make build
sudo make install
```

## Usage

### Basic Examples

```bash
# Rename all files in current directory
gmv *

# Rename specific files
gmv file1.go file2.go

# Rename specific file types
gmv *.{pdf,epub}

# Rename files in a subdirectory
gmv src/*.jsx

# Rename all directories
gmv */

# Rename files in all subdirectories
gmv */*
```

### Options

```bash
# Preview changes without applying (dry-run)
gmv --dry-run *

# Skip overwrite confirmation prompts
gmv --force *
gmv -f *

# Display help
gmv --help
gmv -h
```

## How It Works

1. **gmv** opens a temporary file in your `$EDITOR` with a list of files to rename
2. Edit the filenames as needed (one per line)
3. Save and exit the editor
4. **gmv** validates the changes and applies the renames

### File Swaps

**gmv** automatically handles file swaps and cycles using temporary files:

```bash
# Original files
file-a.txt
file-b.txt

# After editing: swap the names
file-b.txt
file-a.txt

# gmv handles this automatically!
```

Even complex cycles work seamlessly:
```
a → b → c → a  # No problem!
```

### Overwrite Protection

If renaming would overwrite files not in the original list, **gmv** will:
- Display a warning with affected files
- Prompt for confirmation (unless `--force` is used)
- Allow you to cancel the operation

File swaps within your rename list are always safe and won't trigger warnings.

## Environment Variables

- `$EDITOR` - Your preferred text editor (defaults to `vi` or `nano`)

## Validation

**gmv** validates all edits before applying changes:

- ✅ Line count must match the original file list
- ✅ Files cannot be moved to different directories
- ✅ Duplicate target filenames are not allowed (except in swaps)
- ✅ Empty or deleted lines will cause an error

## Operation Logs

All rename operations are logged to `/tmp/gmv-log-YYYYMMDD-HHMMSS` for your records and potential undo operations.

## Building from Source

### Prerequisites

- Go 1.25 or later
- Make

### Build Commands

```bash
make              # Build binary and man page
make build        # Build binary only
make man          # Generate man page
make test         # Run tests
make install      # Install to system
make clean        # Clean build artifacts
```

## Similar Tools

- [mmv](https://github.com/itchyny/mmv) - Move/copy/link multiple files
- [vidir](https://joeyh.name/code/moreutils/) - Edit directories in a text editor
- [rename](https://www.nongnu.org/renameutils/) - Rename files using patterns

## License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.
