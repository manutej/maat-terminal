package datasource

import (
	"context"

	"github.com/manutej/maat-terminal/internal/graph"
	"github.com/manutej/maat-terminal/internal/tui"
)

// MockSource provides the existing mock data for testing/demo purposes.
type MockSource struct{}

// NewMockSource creates a mock data source
func NewMockSource() *MockSource {
	return &MockSource{}
}

// Name returns the data source identifier
func (m *MockSource) Name() string {
	return "mock"
}

// SupportsRefresh returns false - mock data is static
func (m *MockSource) SupportsRefresh() bool {
	return false
}

// Load returns the existing mock graph data
func (m *MockSource) Load(ctx context.Context) ([]graph.Node, []graph.Edge, error) {
	nodes, edges := tui.GetMockGraph()
	return nodes, edges, nil
}
