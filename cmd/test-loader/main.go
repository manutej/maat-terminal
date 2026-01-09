// +build ignore

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/manutej/maat-terminal/internal/datasource"
)

func main() {
	// Test with different paths
	paths := []string{
		"/Users/manu/Documents/LUXOR/MAAT",
		"/Users/manu/Documents/LUXOR",
	}

	// Allow path from command line
	if len(os.Args) > 1 {
		paths = []string{os.Args[1]}
	}

	ctx := context.Background()

	for _, path := range paths {
		fmt.Printf("\n========================================\n")
		fmt.Printf("Testing path: %s\n", path)
		fmt.Printf("========================================\n")

		// Test git scanner
		fmt.Println("\n=== GitScanner ===")
		gitScanner := datasource.NewGitScanner(path)
		gitScanner.SetMaxCommits(5)

		nodes, edges, err := gitScanner.Load(ctx)
		if err != nil {
			fmt.Printf("Git error: %v\n", err)
		} else {
			fmt.Printf("Git nodes: %d, edges: %d\n", len(nodes), len(edges))
			for i, n := range nodes {
				if i < 3 {
					fmt.Printf("  - %s: %s\n", n.Type, n.Title())
				}
			}
		}

		// Test file scanner
		fmt.Println("\n=== FileScanner ===")
		fileScanner := datasource.NewFileScanner(path, fmt.Sprintf("project:%s", path))
		fileScanner.SetMaxFiles(10)

		nodes, edges, err = fileScanner.Load(ctx)
		if err != nil {
			fmt.Printf("File error: %v\n", err)
		} else {
			fmt.Printf("File nodes: %d, edges: %d\n", len(nodes), len(edges))
			for i, n := range nodes {
				if i < 3 {
					fmt.Printf("  - %s: %s\n", n.Type, n.Title())
				}
			}
		}
	}
}
