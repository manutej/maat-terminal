package tui

import (
	"encoding/json"

	"github.com/manutej/maat-terminal/internal/graph"
)

// DisplayNode is a simplified node representation for TUI display.
// It extracts common display fields from the graph.Node JSON data.
type DisplayNode struct {
	ID          string
	Type        graph.NodeType
	Title       string
	Description string
	Status      string
	Priority    int
	Labels      []string
	URL         string // Link to source (Linear, GitHub, etc.)
	Identifier  string // Short identifier (e.g., CET-352 for Linear issues)
	Project     string // Parent project name
}

// IssueData represents the JSON data structure for Issue nodes.
type IssueData struct {
	Title       string   `json:"title"`
	Identifier  string   `json:"identifier"`
	Description string   `json:"description"`
	Status      string   `json:"status"`
	Priority    int      `json:"priority"`
	Labels      []string `json:"labels"`
	Assignee    string   `json:"assignee"`
	URL         string   `json:"url"`
	Project     string   `json:"project"`
}

// PRData represents the JSON data structure for PR nodes.
type PRData struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Number      int    `json:"number"`
	Author      string `json:"author"`
	URL         string `json:"url"`
}

// CommitData represents the JSON data structure for Commit nodes.
type CommitData struct {
	Message string `json:"message"`
	Author  string `json:"author"`
	Hash    string `json:"hash"`
	Date    string `json:"date"`
}

// FileData represents the JSON data structure for File nodes.
type FileData struct {
	Path     string `json:"path"`
	Language string `json:"language"`
	Lines    int    `json:"lines"`
}

// NodeToDisplayNode converts a graph.Node to a DisplayNode for TUI display.
func NodeToDisplayNode(node graph.Node) DisplayNode {
	display := DisplayNode{
		ID:   node.ID,
		Type: node.Type,
	}

	switch node.Type {
	case graph.NodeTypeIssue:
		var data IssueData
		if err := json.Unmarshal(node.Data, &data); err == nil {
			display.Title = data.Title
			display.Identifier = data.Identifier
			display.Description = data.Description
			display.Status = data.Status
			display.Priority = data.Priority
			display.Labels = data.Labels
			display.URL = data.URL
			display.Project = data.Project
		}

	case graph.NodeTypePR:
		var data PRData
		if err := json.Unmarshal(node.Data, &data); err == nil {
			display.Title = data.Title
			display.Description = data.Description
			display.Status = data.Status
		}

	case graph.NodeTypeCommit:
		var data CommitData
		if err := json.Unmarshal(node.Data, &data); err == nil {
			display.Title = data.Message
			display.Description = data.Author
		}

	case graph.NodeTypeFile:
		var data FileData
		if err := json.Unmarshal(node.Data, &data); err == nil {
			display.Title = data.Path
			display.Description = data.Language
		}

	default:
		// Try to extract a title from generic JSON
		var generic map[string]interface{}
		if err := json.Unmarshal(node.Data, &generic); err == nil {
			if title, ok := generic["title"].(string); ok {
				display.Title = title
			} else if name, ok := generic["name"].(string); ok {
				display.Title = name
			} else {
				display.Title = node.ID
			}
		}
	}

	// Fallback if title is still empty
	if display.Title == "" {
		display.Title = node.ID
	}

	return display
}

// NodesToDisplayNodes converts a slice of graph.Node to DisplayNodes.
func NodesToDisplayNodes(nodes []graph.Node) []DisplayNode {
	displayNodes := make([]DisplayNode, 0, len(nodes))
	for _, node := range nodes {
		displayNodes = append(displayNodes, NodeToDisplayNode(node))
	}
	return displayNodes
}

// DisplayEdge is a simplified edge representation for TUI display.
type DisplayEdge struct {
	FromID   string
	ToID     string
	Relation graph.EdgeType
}

// EdgeToDisplayEdge converts a graph.Edge to a DisplayEdge.
func EdgeToDisplayEdge(edge graph.Edge) DisplayEdge {
	return DisplayEdge{
		FromID:   edge.FromID,
		ToID:     edge.ToID,
		Relation: edge.Relation,
	}
}

// EdgesToDisplayEdges converts a slice of graph.Edge to DisplayEdges.
func EdgesToDisplayEdges(edges []graph.Edge) []DisplayEdge {
	displayEdges := make([]DisplayEdge, 0, len(edges))
	for _, edge := range edges {
		displayEdges = append(displayEdges, EdgeToDisplayEdge(edge))
	}
	return displayEdges
}
