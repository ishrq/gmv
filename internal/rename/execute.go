package rename

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ExecuteRenames performs the rename operations or prints them in dry-run mode
func ExecuteRenames(plan []RenameOp, dryRun bool) error {
	for _, op := range plan {
		if dryRun {
			fmt.Printf("%s -> %s\n", op.From, op.To)
		} else {
			if err := os.Rename(op.From, op.To); err != nil {
				return fmt.Errorf("failed to rename %s to %s: %w", op.From, op.To, err)
			}
		}
	}
	return nil
}

// WriteLog creates a log file with all rename operations
func WriteLog(plan []RenameOp) (string, error) {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "unknown"
	}

	// Create log file
	timestamp := time.Now().Format("20060102-150405")
	logPath := filepath.Join(os.TempDir(), fmt.Sprintf("gmv-log-%s", timestamp))

	logFile, err := os.Create(logPath)
	if err != nil {
		return "", fmt.Errorf("failed to create log file: %w", err)
	}
	defer logFile.Close()

	// Write header
	header := fmt.Sprintf("# gmv operation log - %s\n", time.Now().Format("2006-01-02 15:04:05"))
	header += fmt.Sprintf("# Working directory: %s\n\n", cwd)

	if _, err := logFile.WriteString(header); err != nil {
		return "", fmt.Errorf("failed to write log header: %w", err)
	}

	// Write operations
	for _, op := range plan {
		line := fmt.Sprintf("%s -> %s\n", op.From, op.To)
		if _, err := logFile.WriteString(line); err != nil {
			return "", fmt.Errorf("failed to write log entry: %w", err)
		}
	}

	return logPath, nil
}
