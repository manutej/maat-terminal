package datasource

import (
	"context"
	"fmt"
	"os"

	"github.com/manutej/maat-terminal/internal/graph"
)

// DataSource is the interface for loading graph data from various sources.
// Following Commandment #7 (Composition): Thin API clients, unified interface.
type DataSource interface {
	// Name returns the data source identifier
	Name() string

	// Load fetches nodes and edges from the source
	Load(ctx context.Context) ([]graph.Node, []graph.Edge, error)

	// SupportsRefresh returns true if the source can be refreshed
	SupportsRefresh() bool
}

// Config holds configuration for data sources
type Config struct {
	// ProjectPath is the local path to scan (for git/files)
	ProjectPath string

	// GitHubRepo is the GitHub repository (owner/repo format)
	GitHubRepo string

	// GitHubToken is the personal access token for GitHub API
	GitHubToken string

	// UseMock if true, uses mock data instead of real sources
	UseMock bool
}

// Loader orchestrates loading from multiple data sources
type Loader struct {
	sources []DataSource
}

// NewLoader creates a new data source loader
func NewLoader(sources ...DataSource) *Loader {
	return &Loader{sources: sources}
}

// LoadAll loads data from all configured sources and merges results
func (l *Loader) LoadAll(ctx context.Context) ([]graph.Node, []graph.Edge, error) {
	var allNodes []graph.Node
	var allEdges []graph.Edge

	for _, source := range l.sources {
		nodes, edges, err := source.Load(ctx)
		if err != nil {
			// Log error but continue with other sources
			fmt.Fprintf(os.Stderr, "Error loading from %s: %v\n", source.Name(), err)
			continue
		}
		fmt.Fprintf(os.Stderr, "Loaded %d nodes from %s\n", len(nodes), source.Name())
		allNodes = append(allNodes, nodes...)
		allEdges = append(allEdges, edges...)
	}

	return allNodes, allEdges, nil
}

// AddSource adds a new data source
func (l *Loader) AddSource(source DataSource) {
	l.sources = append(l.sources, source)
}
