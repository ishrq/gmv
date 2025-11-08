package main

import "time"

func generateManPage() string {
	year := time.Now().Format("2006")
	month := time.Now().Format("January 2006")

	return `.TH GMV 1 "` + month + `" "gmv 1.0" "User Commands"
.SH NAME
gmv \- batch rename files using $EDITOR
.SH SYNOPSIS
.B gmv
[\fIOPTIONS\fR]
.I files...
.SH DESCRIPTION
.B gmv
is a command-line tool for batch renaming files using your preferred text editor.
It opens a list of files in your $EDITOR, allowing you to edit the filenames.
Upon saving and exiting, the files are renamed accordingly.
.PP
File swaps are automatically handled using temporary files to avoid conflicts.
A log of all rename operations is saved in the system's temporary directory.
.PP
The program includes safety checks for file overwrites. If renaming would overwrite
files that are not part of the rename operation, a warning is displayed and user
confirmation is required (unless the
.B \-\-force
flag is used).
.SH OPTIONS
.TP
.B \-\-dry\-run
Preview the changes without actually renaming any files.
Shows what would be renamed without applying the operations.
Also displays warnings for potential overwrites.
.TP
.B \-\-force, \-f
Skip confirmation prompt when files would be overwritten.
Use with caution as this can lead to data loss.
.TP
.B \-\-help, \-h
Display help information and exit.
.SH EXAMPLES
.TP
.B gmv test-file.go
Rename a single file in the editor.
.TP
.B gmv *
Rename all files in the current directory.
.TP
.B gmv *.{pdf,epub}
Rename all PDF and EPUB files.
.TP
.B gmv */
Rename all directories in the current directory.
.TP
.B gmv test-dir/*.txt
Rename all text files in the test-dir directory.
.TP
.B gmv */*
Rename all files in all subdirectories.
.TP
.B gmv \-\-dry\-run *
Preview rename operations without applying them.
.TP
.B gmv \-\-force *
Skip overwrite confirmation prompts.
.SH ENVIRONMENT
.TP
.B EDITOR
The text editor to use for editing filenames.
If not set, defaults to
.B vi
or
.B nano
(whichever is available).
.SH FILES
.TP
.I /tmp/gmv-log-YYYYMMDD-HHMMSS
Log files containing records of rename operations.
Each log file includes a timestamp and the working directory
where the operations were performed.
.SH EXIT STATUS
.TP
.B 0
Success
.TP
.B 1
Error occurred (invalid arguments, file not found, validation failed, etc.)
.SH NOTES
.PP
.B gmv
validates all edits before applying any renames:
.IP \(bu 2
Line count must match the original file list
.IP \(bu 2
Files cannot be moved to different directories
.IP \(bu 2
Duplicate target filenames are not allowed (except in swap operations)
.IP \(bu 2
Empty lines or deleted lines will cause an error
.PP
When files are swapped (e.g., file1 \(-> file2 and file2 \(-> file1),
.B gmv
automatically uses temporary files to handle the operation safely.
.PP
.B Overwrite Protection
.PP
If renaming would overwrite existing files that are not part of the rename list,
.B gmv
will display a warning with the list of files that would be affected and prompt
for confirmation. File swaps and cycles within the rename list are allowed and
will not trigger the overwrite warning. Use
.B \-\-force
to bypass the confirmation prompt.
.SH BUGS
Report bugs at: https://github.com/ishrq/gmv/issues
.SH AUTHOR
Written by Ishraque Alam.
.SH COPYRIGHT
Copyright \(co ` + year + `. License: MIT
.SH SEE ALSO
.BR mv (1),
.BR rename (1),
.BR vidir (1)
`
}
