package main

import (
	"context"
	"fmt"
	"os"

	"github.com/manutej/maat-terminal/internal/datasource"
)

func main() {
	teamID := os.Getenv("LINEAR_TEAM_ID")
	if teamID == "" {
		teamID = "bee0badb-31e3-4d7a-b18d-7c7d16c4eb9f" // Ceti-luxor default
	}

	if os.Getenv("LINEAR_API_KEY") == "" {
		fmt.Println("ERROR: LINEAR_API_KEY environment variable not set")
		fmt.Println("Get one from: Linear → Settings → API → Personal API Keys")
		os.Exit(1)
	}

	fmt.Println("Testing Linear datasource...")
	fmt.Printf("Team ID: %s\n\n", teamID)

	source := datasource.NewLinearSource(teamID)
	ctx := context.Background()

	nodes, edges, err := source.Load(ctx)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Loaded %d nodes and %d edges\n\n", len(nodes), len(edges))

	// Print issues
	fmt.Println("=== LINEAR ISSUES ===")
	issueCount := 0
	for _, node := range nodes {
		if node.Type == "Issue" {
			issueCount++
			title := node.Title()
			status := node.Status()
			fmt.Printf("  [%s] %s - %s\n", status, node.ID, title)
		}
	}
	fmt.Printf("\nTotal issues: %d\n", issueCount)

	// Print projects
	fmt.Println("\n=== LINEAR PROJECTS ===")
	projectCount := 0
	for _, node := range nodes {
		if node.Type == "Project" {
			projectCount++
			title := node.Title()
			fmt.Printf("  %s - %s\n", node.ID, title)
		}
	}
	fmt.Printf("\nTotal projects: %d\n", projectCount)

	// Print some edges
	fmt.Println("\n=== EDGES (first 10) ===")
	for i, edge := range edges {
		if i >= 10 {
			fmt.Printf("  ... and %d more edges\n", len(edges)-10)
			break
		}
		fmt.Printf("  %s -[%s]-> %s\n", edge.FromID, edge.Relation, edge.ToID)
	}

	fmt.Println("\n✅ Linear datasource working!")
}
