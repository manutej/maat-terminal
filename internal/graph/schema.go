package graph

import (
	"encoding/json"
	"time"
)

// NodeType represents the type of entity in the knowledge graph
type NodeType string

const (
	NodeTypeIssue   NodeType = "Issue"
	NodeTypePR      NodeType = "PR"
	NodeTypeCommit  NodeType = "Commit"
	NodeTypeFile    NodeType = "File"
	NodeTypeProject NodeType = "Project"
	NodeTypeService NodeType = "Service"
)

// EdgeType represents the relationship between nodes
type EdgeType string

const (
	EdgeBlocks     EdgeType = "blocks"
	EdgeRelated    EdgeType = "related"
	EdgeImplements EdgeType = "implements"
	EdgeCalls      EdgeType = "calls"
	EdgeOwns       EdgeType = "owns"
	EdgeModifies   EdgeType = "modifies"
	EdgeMentions   EdgeType = "mentions"
	EdgeParentOf   EdgeType = "parent_of"
)

// Role represents access level (from ADR-006 IDP spec)
type Role string

const (
	RoleExec Role = "exec"
	RoleLead Role = "lead"
	RoleIC   Role = "ic"
)

// Node represents a graph node with arbitrary JSON data
type Node struct {
	ID       string          `json:"id"`
	Type     NodeType        `json:"type"`
	Source   string          `json:"source"`
	Data     json.RawMessage `json:"data"`
	Metadata NodeMetadata    `json:"metadata"`
}

// NodeMetadata contains tracking and access control information
type NodeMetadata struct {
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedBy   string    `json:"created_by"`    // user | ai:<session_id>
	AccessLevel Role      `json:"access_level"`  // exec | lead | ic
	SyncedAt    time.Time `json:"synced_at"`     // Last API sync
}

// Edge represents a directed relationship between two nodes
type Edge struct {
	ID       string       `json:"id"`
	FromID   string       `json:"from_id"`
	ToID     string       `json:"to_id"`
	Relation EdgeType     `json:"relation"`
	Metadata EdgeMetadata `json:"metadata,omitempty"`
}

// EdgeMetadata contains optional relationship metadata
type EdgeMetadata struct {
	CreatedAt time.Time              `json:"created_at,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

// NodeFilter provides filtering for node queries
type NodeFilter struct {
	Types        []NodeType
	Sources      []string
	UpdatedAfter time.Time
}

// ValidateNodeType checks if a string is a valid NodeType
func ValidateNodeType(t string) bool {
	switch NodeType(t) {
	case NodeTypeIssue, NodeTypePR, NodeTypeCommit, NodeTypeFile, NodeTypeProject, NodeTypeService:
		return true
	default:
		return false
	}
}

// ValidateEdgeType checks if a string is a valid EdgeType
func ValidateEdgeType(t string) bool {
	switch EdgeType(t) {
	case EdgeBlocks, EdgeRelated, EdgeImplements, EdgeCalls, EdgeOwns, EdgeModifies, EdgeMentions, EdgeParentOf:
		return true
	default:
		return false
	}
}

// Helper methods to extract common fields from Data JSON

// Title extracts the title field from node data
// Falls back to "name" field for Projects and Services
func (n *Node) Title() string {
	var data map[string]interface{}
	if err := json.Unmarshal(n.Data, &data); err != nil {
		return n.ID // Fallback to ID if JSON parsing fails
	}
	// Try "title" first (Issues, PRs, Commits)
	if title, ok := data["title"].(string); ok {
		return title
	}
	// Fallback to "name" (Projects, Services)
	if name, ok := data["name"].(string); ok {
		return name
	}
	// Last resort: try "path" (Files)
	if path, ok := data["path"].(string); ok {
		return path
	}
	return n.ID // Ultimate fallback
}

// Description extracts the description field from node data
func (n *Node) Description() string {
	var data map[string]interface{}
	if err := json.Unmarshal(n.Data, &data); err != nil {
		return ""
	}
	if desc, ok := data["description"].(string); ok {
		return desc
	}
	return ""
}

// Status extracts the status field from node data
func (n *Node) Status() string {
	var data map[string]interface{}
	if err := json.Unmarshal(n.Data, &data); err != nil {
		return ""
	}
	if status, ok := data["status"].(string); ok {
		return status
	}
	return ""
}

// Priority extracts the priority field from node data
func (n *Node) Priority() int {
	var data map[string]interface{}
	if err := json.Unmarshal(n.Data, &data); err != nil {
		return 0
	}
	if priority, ok := data["priority"].(float64); ok {
		return int(priority)
	}
	return 0
}

// Labels extracts the labels field from node data
func (n *Node) Labels() []string {
	var data map[string]interface{}
	if err := json.Unmarshal(n.Data, &data); err != nil {
		return nil
	}
	if labelsRaw, ok := data["labels"].([]interface{}); ok {
		labels := make([]string, 0, len(labelsRaw))
		for _, l := range labelsRaw {
			if label, ok := l.(string); ok {
				labels = append(labels, label)
			}
		}
		return labels
	}
	return nil
}
