package rename

import (
	"fmt"
	"path/filepath"
	"time"
)

// BuildRenamePlan creates a plan for renaming files, handling cycles with temp files
func BuildRenamePlan(original, edited []string) ([]RenameOp, error) {
	initialPlan := []RenameOp{}
	renameMap := make(map[string]string) // from -> to mapping

	for i := 0; i < len(original); i++ {
		// Skip if no change
		if original[i] == edited[i] {
			continue
		}

		initialPlan = append(initialPlan, RenameOp{
			From: original[i],
			To:   edited[i],
		})
		renameMap[original[i]] = edited[i]
	}

	// Detect cycles
	cycles := DetectCycles(initialPlan)

	// If no cycles, return the initial plan
	if len(cycles) == 0 {
		return initialPlan, nil
	}

	// Handle cycles by using temp files
	finalPlan := []RenameOp{}
	handledInCycle := make(map[string]bool)

	for _, cycle := range cycles {
		if len(cycle) == 0 {
			continue
		}

		// Generate temp filename
		firstFile := cycle[0]
		dir := filepath.Dir(firstFile)
		tempName := filepath.Join(dir, fmt.Sprintf(".gmv_temp_%d", time.Now().UnixNano()))

		// Mark all files in cycle as handled
		for _, file := range cycle {
			handledInCycle[file] = true
		}

		// Step 1: Move first file to temp
		finalPlan = append(finalPlan, RenameOp{
			From: firstFile,
			To:   tempName,
		})

		// Step 2: Move rest of the cycle
		for i := 1; i < len(cycle); i++ {
			from := cycle[i]
			to := renameMap[from]
			finalPlan = append(finalPlan, RenameOp{
				From: from,
				To:   to,
			})
		}

		// Step 3: Move temp to final destination
		finalPlan = append(finalPlan, RenameOp{
			From: tempName,
			To:   renameMap[firstFile],
		})
	}

	// Add non-cycle operations
	for _, op := range initialPlan {
		if !handledInCycle[op.From] {
			finalPlan = append(finalPlan, op)
		}
	}

	return finalPlan, nil
}

// DetectCycles finds cycles in rename operations using DFS
func DetectCycles(plan []RenameOp) [][]string {
	// Build adjacency map: from -> to
	graph := make(map[string]string)
	for _, op := range plan {
		graph[op.From] = op.To
	}

	visited := make(map[string]bool)
	recStack := make(map[string]bool)
	var cycles [][]string

	var dfs func(node string, path []string) bool
	dfs = func(node string, path []string) bool {
		visited[node] = true
		recStack[node] = true
		path = append(path, node)

		if next, exists := graph[node]; exists {
			if recStack[next] {
				// Found a cycle, extract it from path
				cycleStart := -1
				for i, n := range path {
					if n == next {
						cycleStart = i
						break
					}
				}
				if cycleStart != -1 {
					cycle := make([]string, len(path)-cycleStart)
					copy(cycle, path[cycleStart:])
					cycles = append(cycles, cycle)
				}
				return true
			} else if !visited[next] {
				if dfs(next, path) {
					return true
				}
			}
		}

		recStack[node] = false
		return false
	}

	// Run DFS from each unvisited node
	for _, op := range plan {
		if !visited[op.From] {
			dfs(op.From, []string{})
		}
	}

	return cycles
}
