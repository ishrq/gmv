package test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/ishrq/gmv/internal/rename"
)

// setupTestFiles creates temporary files and directories for testing
func setupTestFiles(t *testing.T, files []string) (string, func()) {
	tmpDir, err := os.MkdirTemp("", "gmv-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	for _, file := range files {
		fullPath := filepath.Join(tmpDir, file)

		// Check if it's a directory (ends with /)
		if file[len(file)-1] == '/' {
			if err := os.MkdirAll(fullPath, 0755); err != nil {
				t.Fatalf("Failed to create directory %s: %v", file, err)
			}
		} else {
			// Create parent directories if needed
			dir := filepath.Dir(fullPath)
			if err := os.MkdirAll(dir, 0755); err != nil {
				t.Fatalf("Failed to create parent dir for %s: %v", file, err)
			}
			if err := os.WriteFile(fullPath, []byte("test"), 0644); err != nil {
				t.Fatalf("Failed to create file %s: %v", file, err)
			}
		}
	}

	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	return tmpDir, cleanup
}

// fileExists checks if a file or directory exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func TestSimpleRenames(t *testing.T) {
	files := []string{"file1.txt", "file2.txt", "file3.txt"}
	tmpDir, cleanup := setupTestFiles(t, files)
	defer cleanup()

	// Simulate name changes
	original := []string{
		filepath.Join(tmpDir, "file1.txt"),
		filepath.Join(tmpDir, "file2.txt"),
		filepath.Join(tmpDir, "file3.txt"),
	}
	edited := []string{
		filepath.Join(tmpDir, "renamed1.txt"),
		filepath.Join(tmpDir, "renamed2.txt"),
		filepath.Join(tmpDir, "renamed3.txt"),
	}

	// Validate, plan, and execute
	if err := rename.ValidateEdits(original, edited); err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	plan, err := rename.BuildRenamePlan(original, edited)
	if err != nil {
		t.Fatalf("Build plan failed: %v", err)
	}

	if len(plan) != 3 {
		t.Errorf("Expected 3 rename operations, got %d", len(plan))
	}

	if err := rename.ExecuteRenames(plan, false); err != nil {
		t.Fatalf("Execute renames failed: %v", err)
	}

	// Verify renames
	for _, newName := range edited {
		if !fileExists(newName) {
			t.Errorf("File %s does not exist after rename", newName)
		}
	}
	for _, oldName := range original {
		if fileExists(oldName) {
			t.Errorf("File %s still exists after rename", oldName)
		}
	}
}

func TestDirectoryRenames(t *testing.T) {
	dirs := []string{"dir1/", "dir2/", "dir3/"}
	tmpDir, cleanup := setupTestFiles(t, dirs)
	defer cleanup()

	original := []string{
		filepath.Join(tmpDir, "dir1"),
		filepath.Join(tmpDir, "dir2"),
		filepath.Join(tmpDir, "dir3"),
	}
	edited := []string{
		filepath.Join(tmpDir, "folder1"),
		filepath.Join(tmpDir, "folder2"),
		filepath.Join(tmpDir, "folder3"),
	}

	if err := rename.ValidateEdits(original, edited); err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	plan, err := rename.BuildRenamePlan(original, edited)
	if err != nil {
		t.Fatalf("Build plan failed: %v", err)
	}

	if err := rename.ExecuteRenames(plan, false); err != nil {
		t.Fatalf("Execute renames failed: %v", err)
	}

	// Verify directory renames
	for _, newName := range edited {
		if !fileExists(newName) {
			t.Errorf("Directory %s does not exist after rename", newName)
		}
	}
}

func TestMixedFilesAndDirectories(t *testing.T) {
	items := []string{"file1.txt", "dir1/", "file2.go", "dir2/"}
	tmpDir, cleanup := setupTestFiles(t, items)
	defer cleanup()

	original := []string{
		filepath.Join(tmpDir, "file1.txt"),
		filepath.Join(tmpDir, "dir1"),
		filepath.Join(tmpDir, "file2.go"),
		filepath.Join(tmpDir, "dir2"),
	}
	edited := []string{
		filepath.Join(tmpDir, "renamed.txt"),
		filepath.Join(tmpDir, "folder"),
		filepath.Join(tmpDir, "main.go"),
		filepath.Join(tmpDir, "src"),
	}

	if err := rename.ValidateEdits(original, edited); err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	plan, err := rename.BuildRenamePlan(original, edited)
	if err != nil {
		t.Fatalf("Build plan failed: %v", err)
	}

	if err := rename.ExecuteRenames(plan, false); err != nil {
		t.Fatalf("Execute renames failed: %v", err)
	}

	// Verify all renames
	for _, newName := range edited {
		if !fileExists(newName) {
			t.Errorf("%s does not exist after rename", newName)
		}
	}
}

func TestSimpleSwap(t *testing.T) {
	files := []string{"fileA.txt", "fileB.txt"}
	tmpDir, cleanup := setupTestFiles(t, files)
	defer cleanup()

	original := []string{
		filepath.Join(tmpDir, "fileA.txt"),
		filepath.Join(tmpDir, "fileB.txt"),
	}
	edited := []string{
		filepath.Join(tmpDir, "fileB.txt"),
		filepath.Join(tmpDir, "fileA.txt"),
	}

	if err := rename.ValidateEdits(original, edited); err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	plan, err := rename.BuildRenamePlan(original, edited)
	if err != nil {
		t.Fatalf("Build plan failed: %v", err)
	}

	// Should have extra operations for the swap (using temp file)
	if len(plan) < 2 {
		t.Errorf("Expected at least 2 operations for swap, got %d", len(plan))
	}

	if err := rename.ExecuteRenames(plan, false); err != nil {
		t.Fatalf("Execute renames failed: %v", err)
	}

	// Both files should still exist (just swapped)
	if !fileExists(filepath.Join(tmpDir, "fileA.txt")) {
		t.Error("fileA.txt does not exist after swap")
	}
	if !fileExists(filepath.Join(tmpDir, "fileB.txt")) {
		t.Error("fileB.txt does not exist after swap")
	}
}

func TestMultipleSwaps(t *testing.T) {
	files := []string{"file1.txt", "file2.txt", "file3.txt", "file4.txt"}
	tmpDir, cleanup := setupTestFiles(t, files)
	defer cleanup()

	original := []string{
		filepath.Join(tmpDir, "file1.txt"),
		filepath.Join(tmpDir, "file2.txt"),
		filepath.Join(tmpDir, "file3.txt"),
		filepath.Join(tmpDir, "file4.txt"),
	}
	// Create two swaps: 1<->2 and 3<->4
	edited := []string{
		filepath.Join(tmpDir, "file2.txt"),
		filepath.Join(tmpDir, "file1.txt"),
		filepath.Join(tmpDir, "file4.txt"),
		filepath.Join(tmpDir, "file3.txt"),
	}

	if err := rename.ValidateEdits(original, edited); err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	plan, err := rename.BuildRenamePlan(original, edited)
	if err != nil {
		t.Fatalf("Build plan failed: %v", err)
	}

	if err := rename.ExecuteRenames(plan, false); err != nil {
		t.Fatalf("Execute renames failed: %v", err)
	}

	// All files should still exist
	for _, file := range original {
		if !fileExists(file) {
			t.Errorf("File %s does not exist after swaps", file)
		}
	}
}

func TestCyclicSwap(t *testing.T) {
	files := []string{"a.txt", "b.txt", "c.txt"}
	tmpDir, cleanup := setupTestFiles(t, files)
	defer cleanup()

	original := []string{
		filepath.Join(tmpDir, "a.txt"),
		filepath.Join(tmpDir, "b.txt"),
		filepath.Join(tmpDir, "c.txt"),
	}
	// Create cycle: a->b, b->c, c->a
	edited := []string{
		filepath.Join(tmpDir, "b.txt"),
		filepath.Join(tmpDir, "c.txt"),
		filepath.Join(tmpDir, "a.txt"),
	}

	if err := rename.ValidateEdits(original, edited); err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	plan, err := rename.BuildRenamePlan(original, edited)
	if err != nil {
		t.Fatalf("Build plan failed: %v", err)
	}

	if err := rename.ExecuteRenames(plan, false); err != nil {
		t.Fatalf("Execute renames failed: %v", err)
	}

	// All files should exist after cycle
	for _, file := range original {
		if !fileExists(file) {
			t.Errorf("File %s does not exist after cyclic swap", file)
		}
	}
}

func TestLargeNumberOfFiles(t *testing.T) {
	// Create 1000 files
	numFiles := 1000
	files := make([]string, numFiles)
	for i := 0; i < numFiles; i++ {
		files[i] = fmt.Sprintf("file%d.txt", i)
	}

	tmpDir, cleanup := setupTestFiles(t, files)
	defer cleanup()

	original := make([]string, numFiles)
	edited := make([]string, numFiles)
	for i := 0; i < numFiles; i++ {
		original[i] = filepath.Join(tmpDir, fmt.Sprintf("file%d.txt", i))
		edited[i] = filepath.Join(tmpDir, fmt.Sprintf("renamed%d.txt", i))
	}

	if err := rename.ValidateEdits(original, edited); err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	plan, err := rename.BuildRenamePlan(original, edited)
	if err != nil {
		t.Fatalf("Build plan failed: %v", err)
	}

	if len(plan) != numFiles {
		t.Errorf("Expected %d rename operations, got %d", numFiles, len(plan))
	}

	if err := rename.ExecuteRenames(plan, false); err != nil {
		t.Fatalf("Execute renames failed: %v", err)
	}

	// Spot check some files
	for i := 0; i < 10; i++ {
		if !fileExists(edited[i]) {
			t.Errorf("Renamed file %s does not exist", edited[i])
		}
	}
}

func TestLineCountMismatch(t *testing.T) {
	files := []string{"file1.txt", "file2.txt"}
	tmpDir, cleanup := setupTestFiles(t, files)
	defer cleanup()

	original := []string{
		filepath.Join(tmpDir, "file1.txt"),
		filepath.Join(tmpDir, "file2.txt"),
	}
	// Only one edited line (missing one)
	edited := []string{
		filepath.Join(tmpDir, "renamed.txt"),
	}

	err := rename.ValidateEdits(original, edited)
	if err == nil {
		t.Fatal("Expected line count mismatch error, got nil")
	}

	if err.Error() != "line count mismatch: expected 2 lines, got 1 lines" {
		t.Errorf("Wrong error message: %v", err)
	}
}

func TestDuplicateTargetNames(t *testing.T) {
	files := []string{"file1.txt", "file2.txt"}
	tmpDir, cleanup := setupTestFiles(t, files)
	defer cleanup()

	original := []string{
		filepath.Join(tmpDir, "file1.txt"),
		filepath.Join(tmpDir, "file2.txt"),
	}
	// Both renamed to same name
	edited := []string{
		filepath.Join(tmpDir, "duplicate.txt"),
		filepath.Join(tmpDir, "duplicate.txt"),
	}

	err := rename.ValidateEdits(original, edited)
	if err == nil {
		t.Fatal("Expected duplicate filename error, got nil")
	}

	if err.Error() != fmt.Sprintf("duplicate target filename: %s", edited[0]) {
		t.Errorf("Wrong error message: %v", err)
	}
}

func TestDirectoryChangeNotAllowed(t *testing.T) {
	files := []string{"file.txt"}
	tmpDir, cleanup := setupTestFiles(t, files)
	defer cleanup()

	original := []string{
		filepath.Join(tmpDir, "file.txt"),
	}
	// Try to move to a subdirectory
	edited := []string{
		filepath.Join(tmpDir, "subdir", "file.txt"),
	}

	err := rename.ValidateEdits(original, edited)
	if err == nil {
		t.Fatal("Expected directory change error, got nil")
	}
}

func TestNoChanges(t *testing.T) {
	files := []string{"file1.txt", "file2.txt"}
	tmpDir, cleanup := setupTestFiles(t, files)
	defer cleanup()

	original := []string{
		filepath.Join(tmpDir, "file1.txt"),
		filepath.Join(tmpDir, "file2.txt"),
	}
	// No changes - same names
	edited := []string{
		filepath.Join(tmpDir, "file1.txt"),
		filepath.Join(tmpDir, "file2.txt"),
	}

	if err := rename.ValidateEdits(original, edited); err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	plan, err := rename.BuildRenamePlan(original, edited)
	if err != nil {
		t.Fatalf("Build plan failed: %v", err)
	}

	// Should have no operations
	if len(plan) != 0 {
		t.Errorf("Expected 0 operations for no changes, got %d", len(plan))
	}
}

func TestOverwriteDetection(t *testing.T) {
	// Create test files including one that will be overwritten
	files := []string{"file1.txt", "file2.txt", "existing.txt"}
	tmpDir, cleanup := setupTestFiles(t, files)
	defer cleanup()

	// Only file1 and file2 are in the rename list
	original := []string{
		filepath.Join(tmpDir, "file1.txt"),
		filepath.Join(tmpDir, "file2.txt"),
	}
	// Try to rename file1 to existing.txt (which exists but is NOT in our list)
	edited := []string{
		filepath.Join(tmpDir, "existing.txt"),
		filepath.Join(tmpDir, "file2.txt"),
	}

	if err := rename.ValidateEdits(original, edited); err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	plan, err := rename.BuildRenamePlan(original, edited)
	if err != nil {
		t.Fatalf("Build plan failed: %v", err)
	}

	// Check for overwrites
	overwrites := rename.CheckOverwrites(plan, original)

	// Should detect that existing.txt would be overwritten
	if len(overwrites) != 1 {
		t.Errorf("Expected 1 overwrite, got %d", len(overwrites))
	}

	if len(overwrites) > 0 && overwrites[0] != filepath.Join(tmpDir, "existing.txt") {
		t.Errorf("Expected overwrite of existing.txt, got %s", overwrites[0])
	}
}

func TestSwapNotFlaggedAsOverwrite(t *testing.T) {
	files := []string{"fileA.txt", "fileB.txt"}
	tmpDir, cleanup := setupTestFiles(t, files)
	defer cleanup()

	original := []string{
		filepath.Join(tmpDir, "fileA.txt"),
		filepath.Join(tmpDir, "fileB.txt"),
	}
	// Swap the names
	edited := []string{
		filepath.Join(tmpDir, "fileB.txt"),
		filepath.Join(tmpDir, "fileA.txt"),
	}

	if err := rename.ValidateEdits(original, edited); err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	plan, err := rename.BuildRenamePlan(original, edited)
	if err != nil {
		t.Fatalf("Build plan failed: %v", err)
	}

	// Check for overwrites - swaps should NOT be flagged
	overwrites := rename.CheckOverwrites(plan, original)

	if len(overwrites) != 0 {
		t.Errorf("Swaps should not be flagged as overwrites, got %d overwrites", len(overwrites))
	}
}

func TestMultipleOverwrites(t *testing.T) {
	// Create files including multiple that will be overwritten
	files := []string{"file1.txt", "file2.txt", "existing1.txt", "existing2.txt"}
	tmpDir, cleanup := setupTestFiles(t, files)
	defer cleanup()

	original := []string{
		filepath.Join(tmpDir, "file1.txt"),
		filepath.Join(tmpDir, "file2.txt"),
	}
	// Try to overwrite both existing files
	edited := []string{
		filepath.Join(tmpDir, "existing1.txt"),
		filepath.Join(tmpDir, "existing2.txt"),
	}

	if err := rename.ValidateEdits(original, edited); err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	plan, err := rename.BuildRenamePlan(original, edited)
	if err != nil {
		t.Fatalf("Build plan failed: %v", err)
	}

	overwrites := rename.CheckOverwrites(plan, original)

	// Should detect both overwrites
	if len(overwrites) != 2 {
		t.Errorf("Expected 2 overwrites, got %d", len(overwrites))
	}
}

func TestNoOverwriteWhenTargetDoesNotExist(t *testing.T) {
	files := []string{"file1.txt"}
	tmpDir, cleanup := setupTestFiles(t, files)
	defer cleanup()

	original := []string{
		filepath.Join(tmpDir, "file1.txt"),
	}
	// Rename to a file that doesn't exist
	edited := []string{
		filepath.Join(tmpDir, "newfile.txt"),
	}

	if err := rename.ValidateEdits(original, edited); err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	plan, err := rename.BuildRenamePlan(original, edited)
	if err != nil {
		t.Fatalf("Build plan failed: %v", err)
	}

	overwrites := rename.CheckOverwrites(plan, original)

	// Should not detect any overwrites
	if len(overwrites) != 0 {
		t.Errorf("Expected 0 overwrites, got %d", len(overwrites))
	}
}

// These would need to be copied or the package structure changed to allow testing
func validateEdits(original, edited []string) error {
	// This is a placeholder - in real implementation, you'd import from main
	// or refactor to make testable
	return nil
}

func buildRenamePlan(original, edited []string) ([]struct{ From, To string }, error) {
	return nil, nil
}

func executeRenames(plan []struct{ From, To string }, dryRun bool) error {
	return nil
}
