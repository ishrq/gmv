package rename

import (
	"fmt"
	"os"
	"strings"
)

func CreateTempFile(files []string) (string, error) {
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

func ParseEdited(filepath string) ([]string, error) {
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
