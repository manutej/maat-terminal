package tui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/manutej/maat-terminal/internal/graph"
	"github.com/manutej/maat-terminal/internal/tui/styles"
)

// RenderGraph renders the knowledge graph as a clean, navigable tree list.
// This replaces the broken canvas-based approach with a much more usable design.
// Pure function following Commandment #1 (Immutable Truth).
func RenderGraph(m Model, maxWidth int) string {
	// Get filtered nodes and edges
	nodes := m.GetFilteredNodes()
	edges := m.GetFilteredEdges()

	if len(nodes) == 0 {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Render("No nodes match current filter. Press 'f' to change filter.")
	}

	// Build the tree structure
	tree := buildTree(nodes, edges)

	// Render the tree
	var result strings.Builder

	// Header with filter info
	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Accent).
		Bold(true)
	countStyle := lipgloss.NewStyle().
		Foreground(styles.Muted)

	result.WriteString(headerStyle.Render(fmt.Sprintf("Filter: %s", m.filterMode.String())))
	result.WriteString(countStyle.Render(fmt.Sprintf(" (%d nodes)", len(nodes))))
	result.WriteString("\n\n")

	// Render tree nodes
	for i, root := range tree.Roots {
		isLast := i == len(tree.Roots)-1
		result.WriteString(renderTreeNode(root, tree, m, "", isLast, maxWidth))
	}

	return result.String()
}

// TreeStructure holds the hierarchical representation of nodes
type TreeStructure struct {
	Roots    []string            // Root node IDs (no parents)
	Children map[string][]string // Parent -> Children mapping
	Nodes    map[string]DisplayNode
}

// buildTree creates a hierarchical tree from nodes and edges
func buildTree(nodes []DisplayNode, edges []DisplayEdge) TreeStructure {
	tree := TreeStructure{
		Roots:    make([]string, 0),
		Children: make(map[string][]string),
		Nodes:    make(map[string]DisplayNode),
	}

	// Index all nodes
	for _, node := range nodes {
		tree.Nodes[node.ID] = node
	}

	// Build parent-child relationships
	hasParent := make(map[string]bool)
	for _, edge := range edges {
		// Only consider "owns", "implements", "modifies" as parent-child
		if isHierarchicalEdge(edge.Relation) {
			if _, fromExists := tree.Nodes[edge.FromID]; fromExists {
				if _, toExists := tree.Nodes[edge.ToID]; toExists {
					tree.Children[edge.FromID] = append(tree.Children[edge.FromID], edge.ToID)
					hasParent[edge.ToID] = true
				}
			}
		}
	}

	// Find roots (nodes without parents)
	for _, node := range nodes {
		if !hasParent[node.ID] {
			tree.Roots = append(tree.Roots, node.ID)
		}
	}

	// Sort roots by type priority, then by title
	sort.Slice(tree.Roots, func(i, j int) bool {
		ni := tree.Nodes[tree.Roots[i]]
		nj := tree.Nodes[tree.Roots[j]]
		if typePriority(ni.Type) != typePriority(nj.Type) {
			return typePriority(ni.Type) < typePriority(nj.Type)
		}
		return ni.Title < nj.Title
	})

	// Sort children of each node by type, then by status, then by title
	for parent := range tree.Children {
		sort.Slice(tree.Children[parent], func(i, j int) bool {
			ni := tree.Nodes[tree.Children[parent][i]]
			nj := tree.Nodes[tree.Children[parent][j]]
			// First sort by type (projects before issues, etc.)
			if typePriority(ni.Type) != typePriority(nj.Type) {
				return typePriority(ni.Type) < typePriority(nj.Type)
			}
			// Then sort by status (In Progress → Backlog → Done)
			if statusPriority(ni.Status) != statusPriority(nj.Status) {
				return statusPriority(ni.Status) < statusPriority(nj.Status)
			}
			// Finally sort by title alphabetically
			return ni.Title < nj.Title
		})
	}

	return tree
}

// isHierarchicalEdge returns true if the edge represents a parent-child relationship
func isHierarchicalEdge(relation graph.EdgeType) bool {
	switch relation {
	case graph.EdgeOwns, graph.EdgeImplements, graph.EdgeModifies:
		return true
	default:
		return false
	}
}

// typePriority returns sort priority for node types (lower = higher priority)
func typePriority(t graph.NodeType) int {
	switch t {
	case graph.NodeTypeService:
		return 0
	case graph.NodeTypeProject:
		return 1
	case graph.NodeTypeIssue:
		return 2
	case graph.NodeTypePR:
		return 3
	case graph.NodeTypeCommit:
		return 4
	case graph.NodeTypeFile:
		return 5
	default:
		return 99
	}
}

// statusPriority returns sort priority for issue statuses (lower = higher priority)
// Active work shows first, upcoming work second, completed work last
func statusPriority(status string) int {
	s := strings.ToLower(status)
	switch s {
	case "in progress", "in_progress", "started", "in review":
		return 0 // Show first - active work
	case "backlog", "todo", "pending", "triage":
		return 1 // Show second - upcoming work
	case "done", "completed", "merged", "closed":
		return 2 // Show last - completed work
	case "blocked", "canceled", "cancelled":
		return 3 // Show at bottom - blocked/cancelled
	default:
		return 1 // Default to backlog priority
	}
}

// renderTreeNode renders a single node and its children recursively
// Now supports collapsed state from model
func renderTreeNode(nodeID string, tree TreeStructure, m Model, prefix string, isLast bool, maxWidth int) string {
	node, exists := tree.Nodes[nodeID]
	if !exists {
		return ""
	}

	var result strings.Builder

	// Determine tree connector
	connector := "├── "
	if isLast {
		connector = "└── "
	}

	// Build the node line
	isFocused := nodeID == m.focusedNode
	isCollapsed := m.IsCollapsed(nodeID)
	hasChildren := len(tree.Children[nodeID]) > 0

	// Collapse/expand indicator for nodes with children
	var collapseIcon string
	if hasChildren {
		if isCollapsed {
			collapseIcon = "▸ " // Collapsed - right arrow
		} else {
			collapseIcon = "▾ " // Expanded - down arrow
		}
	} else {
		collapseIcon = "  " // No children - spacing
	}

	// Type icon
	icon := getTypeIcon(node.Type)

	// Status indicator with color
	status := getStatusIndicator(node.Status)
	statusColor := getStatusColor(node.Status)

	// Title (truncate if needed)
	title := node.Title
	maxTitleLen := maxWidth - len(prefix) - len(connector) - 15 // Reserve space for icons, status, etc.
	if maxTitleLen < 10 {
		maxTitleLen = 10
	}
	if len(title) > maxTitleLen {
		title = title[:maxTitleLen-3] + "..."
	}

	// Status text for display
	statusText := ""
	if node.Status != "" {
		statusText = fmt.Sprintf(" [%s]", node.Status)
	}

	// Build the line content
	lineContent := fmt.Sprintf("%s%s%s %s%s", collapseIcon, icon, status, title, statusText)

	// Apply styling
	var lineStyle lipgloss.Style
	if isFocused {
		lineStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(styles.Accent).
			Background(lipgloss.Color("236"))
	} else {
		// Color status text differently
		lineStyle = lipgloss.NewStyle().
			Foreground(getTypeColor(node.Type))
	}

	// Status styling (applied separately for non-focused items)
	statusStyle := lipgloss.NewStyle().Foreground(statusColor).Faint(true)

	// Tree prefix styling
	prefixStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	result.WriteString(prefixStyle.Render(prefix + connector))
	if isFocused {
		result.WriteString(lineStyle.Render(lineContent))
	} else {
		// Render with colored status
		baseContent := fmt.Sprintf("%s%s%s %s", collapseIcon, icon, status, title)
		result.WriteString(lineStyle.Render(baseContent))
		if statusText != "" {
			result.WriteString(statusStyle.Render(statusText))
		}
	}
	result.WriteString("\n")

	// Render children only if not collapsed
	if !isCollapsed {
		children := tree.Children[nodeID]
		childPrefix := prefix
		if isLast {
			childPrefix += "    "
		} else {
			childPrefix += "│   "
		}

		for i, childID := range children {
			childIsLast := i == len(children)-1
			result.WriteString(renderTreeNode(childID, tree, m, childPrefix, childIsLast, maxWidth))
		}
	}

	return result.String()
}

// getTypeIcon returns an emoji icon for the node type
func getTypeIcon(t graph.NodeType) string {
	switch t {
	case graph.NodeTypeProject:
		return "📦"
	case graph.NodeTypeIssue:
		return "🔹"
	case graph.NodeTypePR:
		return "🔀"
	case graph.NodeTypeCommit:
		return "💾"
	case graph.NodeTypeFile:
		return "📄"
	case graph.NodeTypeService:
		return "⚙️"
	default:
		return "❓"
	}
}

// getStatusIndicator returns a compact status indicator with Linear status support
func getStatusIndicator(status string) string {
	// Normalize to lowercase for comparison
	s := strings.ToLower(status)
	switch s {
	case "done", "merged", "completed", "closed":
		return "[✓]"
	case "in progress", "in_progress", "open", "started", "in review":
		return "[◐]"
	case "backlog", "todo", "pending", "triage":
		return "[○]"
	case "draft":
		return "[◌]"
	case "blocked", "canceled", "cancelled":
		return "[✗]"
	default:
		return "[-]"
	}
}

// getStatusColor returns the color for a status
func getStatusColor(status string) lipgloss.Color {
	s := strings.ToLower(status)
	switch s {
	case "done", "merged", "completed", "closed":
		return lipgloss.Color("42") // Green
	case "in progress", "in_progress", "open", "started", "in review":
		return lipgloss.Color("214") // Orange
	case "backlog", "todo", "pending", "triage":
		return lipgloss.Color("240") // Gray
	case "blocked", "canceled", "cancelled":
		return lipgloss.Color("196") // Red
	default:
		return lipgloss.Color("252")
	}
}

// getTypeTag returns a short type tag
func getTypeTag(t graph.NodeType) string {
	tagStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Faint(true)

	switch t {
	case graph.NodeTypeProject:
		return tagStyle.Render("")
	case graph.NodeTypeIssue:
		return tagStyle.Render("")
	case graph.NodeTypePR:
		return tagStyle.Render("")
	case graph.NodeTypeCommit:
		return tagStyle.Render("")
	case graph.NodeTypeFile:
		return tagStyle.Render("")
	case graph.NodeTypeService:
		return tagStyle.Render("")
	default:
		return ""
	}
}

// getTypeColor returns the color for a node type
func getTypeColor(t graph.NodeType) lipgloss.Color {
	switch t {
	case graph.NodeTypeProject:
		return lipgloss.Color("33") // Blue
	case graph.NodeTypeIssue:
		return lipgloss.Color("214") // Orange
	case graph.NodeTypePR:
		return lipgloss.Color("135") // Purple
	case graph.NodeTypeCommit:
		return lipgloss.Color("250") // Gray
	case graph.NodeTypeFile:
		return lipgloss.Color("70") // Green
	case graph.NodeTypeService:
		return lipgloss.Color("45") // Cyan
	default:
		return lipgloss.Color("252")
	}
}

// RenderGraphList renders nodes as a simple flat list (legacy fallback)
func RenderGraphList(m Model) string {
	nodes := m.GetFilteredNodes()
	if len(nodes) == 0 {
		return "No nodes to display"
	}

	var result strings.Builder
	for _, node := range nodes {
		icon := getTypeIcon(node.Type)
		status := getStatusIndicator(node.Status)
		focused := ""
		if node.ID == m.focusedNode {
			focused = " ← FOCUSED"
		}
		line := fmt.Sprintf("%s %s %s%s\n", icon, status, node.Title, focused)
		result.WriteString(line)
	}
	return result.String()
}
