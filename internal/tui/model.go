package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/manutej/maat-terminal/internal/graph"
)

// NOTE: Pane concept removed in favor of single-pane design with ViewMode cycling.
// Tab key cycles between Graph/Details/Relations views (full screen each).

// Model holds ALL state (Commandment #1: Immutable Truth)
// No global variables, no pointer mutations in Update
type Model struct {
	// Display state (simplified view of graph data)
	focusedNode string
	nodes       []DisplayNode
	edges       []DisplayEdge

	// UI State
	currentView     ViewMode        // Graph, Details, or Relations (full-screen views)
	filterMode      FilterMode      // Controls which node types are shown (default: FilterProjects)
	statusFilter    StatusFilter    // Controls which statuses are shown (default: StatusAll)
	collapsed       map[string]bool // Tracks which projects/nodes are collapsed
	navStack        NavigationStack
	ready           bool
	width           int
	height          int
	selectedRelIdx  int    // Index of selected relation in Relations view (for drill-down)
	relationsScroll int    // Scroll offset for relations list
	graphScroll     int    // Scroll offset for graph view (line-based)
	searchMode      bool   // True when in search/filter mode (/ key)
	searchQuery     string // Current search query for filtering

	// Components
	viewport viewport.Model
	help     help.Model
	keys     KeyMap

	// Application State
	data         interface{}
	err          error
	loading      bool
	confirmation *ConfirmationRequest
}

// ConfirmationRequest represents a pending external write (Commandment #10: Sovereignty)
type ConfirmationRequest struct {
	Action  string
	Execute func() error
}

// NewModel creates the initial model state
func NewModel() Model {
	return Model{
		// Display state
		focusedNode: "",
		nodes:       make([]DisplayNode, 0),
		edges:       make([]DisplayEdge, 0),

		// UI State
		currentView: ViewGraph,              // Start in Graph view (full screen)
		filterMode:  FilterProjects,         // Start with filtered view (much more usable!)
		collapsed:   make(map[string]bool),  // All projects start expanded
		navStack:    NewNavigationStack(),
		ready:       false,
		width:       80,
		height:      24,

		// Components
		viewport: viewport.New(80, 24),
		help:     help.New(),
		keys:     DefaultKeyMap(),

		// Application State
		data:         nil,
		err:          nil,
		loading:      true,
		confirmation: nil,
	}
}

// NewModelWithData creates a model with pre-loaded data from data sources
func NewModelWithData(nodes []graph.Node, edges []graph.Edge, projectPath string) Model {
	m := NewModel()

	// Convert graph nodes to display nodes
	displayNodes := make([]DisplayNode, len(nodes))
	for i, node := range nodes {
		displayNodes[i] = DisplayNode{
			ID:          node.ID,
			Type:        node.Type,
			Title:       node.Title(),
			Status:      node.Status(),
			Description: node.Description(),
			Priority:    node.Priority(),
			Labels:      node.Labels(),
		}
	}

	// Convert graph edges to display edges
	displayEdges := make([]DisplayEdge, len(edges))
	for i, edge := range edges {
		displayEdges[i] = DisplayEdge{
			FromID:   edge.FromID,
			ToID:     edge.ToID,
			Relation: edge.Relation,
		}
	}

	m.nodes = displayNodes
	m.edges = displayEdges
	m.loading = false

	// Set focus to first node if available
	if len(displayNodes) > 0 {
		m.focusedNode = displayNodes[0].ID
	}

	return m
}

// WithSize returns a new Model with updated dimensions
func (m Model) WithSize(width, height int) Model {
	m.width = width
	m.height = height
	m.viewport.Width = width
	m.viewport.Height = height - 3 // Reserve space for status bar
	return m
}

// WithData returns a new Model with updated data
func (m Model) WithData(data interface{}) Model {
	m.data = data
	m.loading = false
	m.err = nil
	return m
}

// WithError returns a new Model with an error
func (m Model) WithError(err error) Model {
	m.err = err
	m.loading = false
	return m
}

// WithLoading returns a new Model in loading state
func (m Model) WithLoading(loading bool) Model {
	m.loading = loading
	return m
}

// WithConfirmation returns a new Model with a pending confirmation
func (m Model) WithConfirmation(req *ConfirmationRequest) Model {
	m.confirmation = req
	if req != nil {
		m.currentView = ViewConfirm
	}
	return m
}

// WithView returns a new Model with a different view mode
func (m Model) WithView(view ViewMode) Model {
	m.currentView = view
	return m
}

// PushView navigates down (Enter key)
func (m Model) PushView(newView ViewMode) Model {
	m.navStack = m.navStack.Push(m.currentView)
	m.currentView = newView
	return m
}

// PopView navigates up (Esc key)
func (m Model) PopView() Model {
	newStack, previousView, ok := m.navStack.Pop()
	if !ok {
		// Stack empty, stay in current view
		return m
	}
	m.navStack = newStack
	m.currentView = previousView
	return m
}

// WithReady returns a new Model with the ready state set.
func (m Model) WithReady(ready bool) Model {
	m.ready = ready
	return m
}

// WithNodes returns a new Model with display nodes set.
func (m Model) WithNodes(nodes []DisplayNode) Model {
	m.nodes = nodes
	if len(nodes) > 0 && m.focusedNode == "" {
		m.focusedNode = nodes[0].ID
	}
	return m
}

// WithEdges returns a new Model with display edges set.
func (m Model) WithEdges(edges []DisplayEdge) Model {
	m.edges = edges
	return m
}

// WithFocusedNode returns a new Model with the focused node set.
func (m Model) WithFocusedNode(nodeID string) Model {
	m.focusedNode = nodeID
	m.selectedRelIdx = 0 // Reset relation selection when focus changes
	return m
}

// GetFocusedNode returns the currently focused display node, if any.
func (m Model) GetFocusedNode() (DisplayNode, bool) {
	if m.focusedNode == "" || len(m.nodes) == 0 {
		return DisplayNode{}, false
	}
	for _, node := range m.nodes {
		if node.ID == m.focusedNode {
			return node, true
		}
	}
	return DisplayNode{}, false
}

// GetEdgesFrom returns edges originating from a node.
func (m Model) GetEdgesFrom(nodeID string) []DisplayEdge {
	var result []DisplayEdge
	for _, edge := range m.edges {
		if edge.FromID == nodeID {
			result = append(result, edge)
		}
	}
	return result
}

// GetNodeByID returns a node by its ID.
func (m Model) GetNodeByID(nodeID string) (DisplayNode, bool) {
	for _, node := range m.nodes {
		if node.ID == nodeID {
			return node, true
		}
	}
	return DisplayNode{}, false
}

// IsReady returns whether the model is ready for display.
func (m Model) IsReady() bool {
	return m.ready
}

// WithFilterMode returns a new Model with updated filter mode.
func (m Model) WithFilterMode(mode FilterMode) Model {
	m.filterMode = mode
	return m
}

// WithStatusFilter returns a new Model with updated status filter.
func (m Model) WithStatusFilter(filter StatusFilter) Model {
	m.statusFilter = filter
	return m
}

// GetStatusFilter returns the current status filter.
func (m Model) GetStatusFilter() StatusFilter {
	return m.statusFilter
}

// GetFilteredNodes returns nodes filtered by the current filter mode, status filter, and search query.
func (m Model) GetFilteredNodes() []DisplayNode {
	allowedTypes := m.filterMode.Types()

	// Build type filter set
	var typeSet map[string]bool
	if allowedTypes != nil {
		typeSet = make(map[string]bool)
		for _, t := range allowedTypes {
			typeSet[string(t)] = true
		}
	}

	// Normalize search query for case-insensitive matching
	searchLower := strings.ToLower(m.searchQuery)

	filtered := make([]DisplayNode, 0)
	for _, node := range m.nodes {
		// Apply type filter
		if typeSet != nil && !typeSet[string(node.Type)] {
			continue
		}

		// Apply status filter (for nodes that have status - issues, PRs)
		// Projects are always shown as parents, even if their children are filtered
		if node.Type == graph.NodeTypeIssue || node.Type == graph.NodeTypePR {
			if !m.statusFilter.MatchesStatus(node.Status) {
				continue
			}
		}

		// Apply search query filter (if active)
		if searchLower != "" {
			titleLower := strings.ToLower(node.Title)
			if !strings.Contains(titleLower, searchLower) {
				continue
			}
		}

		filtered = append(filtered, node)
	}
	return filtered
}

// GetFilteredEdges returns edges that connect filtered nodes.
func (m Model) GetFilteredEdges() []DisplayEdge {
	filteredNodes := m.GetFilteredNodes()
	nodeSet := make(map[string]bool)
	for _, node := range filteredNodes {
		nodeSet[node.ID] = true
	}

	filtered := make([]DisplayEdge, 0)
	for _, edge := range m.edges {
		if nodeSet[edge.FromID] && nodeSet[edge.ToID] {
			filtered = append(filtered, edge)
		}
	}
	return filtered
}

// GetFilterMode returns the current filter mode.
func (m Model) GetFilterMode() FilterMode {
	return m.filterMode
}

// WithSelectedRelIdx returns a new Model with updated relation selection index.
func (m Model) WithSelectedRelIdx(idx int) Model {
	m.selectedRelIdx = idx
	return m
}

// GetRelationsList returns the list of relations for the focused node.
// Returns a slice of (targetNodeID, relationName, isOutgoing) tuples.
func (m Model) GetRelationsList() []RelationItem {
	node, ok := m.GetFocusedNode()
	if !ok {
		return nil
	}

	var relations []RelationItem

	// Outgoing edges first
	for _, edge := range m.edges {
		if edge.FromID == node.ID {
			if targetNode, ok := m.GetNodeByID(edge.ToID); ok {
				relations = append(relations, RelationItem{
					NodeID:     edge.ToID,
					NodeTitle:  targetNode.Title,
					NodeType:   targetNode.Type,
					Relation:   string(edge.Relation),
					IsOutgoing: true,
				})
			}
		}
	}

	// Incoming edges
	for _, edge := range m.edges {
		if edge.ToID == node.ID {
			if sourceNode, ok := m.GetNodeByID(edge.FromID); ok {
				relations = append(relations, RelationItem{
					NodeID:     edge.FromID,
					NodeTitle:  sourceNode.Title,
					NodeType:   sourceNode.Type,
					Relation:   string(edge.Relation),
					IsOutgoing: false,
				})
			}
		}
	}

	return relations
}

// RelationItem represents a single relation in the relations list.
type RelationItem struct {
	NodeID     string
	NodeTitle  string
	NodeType   graph.NodeType
	Relation   string
	IsOutgoing bool
}

// moveRelationUp moves the selection up in the Relations view.
func (m Model) moveRelationUp() Model {
	relations := m.GetRelationsList()
	if len(relations) == 0 {
		return m
	}

	newIdx := m.selectedRelIdx - 1
	if newIdx < 0 {
		newIdx = len(relations) - 1 // Wrap to bottom
	}
	return m.WithSelectedRelIdx(newIdx)
}

// moveRelationDown moves the selection down in the Relations view.
func (m Model) moveRelationDown() Model {
	relations := m.GetRelationsList()
	if len(relations) == 0 {
		return m
	}

	newIdx := m.selectedRelIdx + 1
	if newIdx >= len(relations) {
		newIdx = 0 // Wrap to top
	}
	return m.WithSelectedRelIdx(newIdx)
}

// jumpToSelectedRelation jumps to the selected relation's node and switches to Graph view.
func (m Model) jumpToSelectedRelation() Model {
	relations := m.GetRelationsList()
	if len(relations) == 0 || m.selectedRelIdx >= len(relations) {
		return m
	}

	// Get selected relation
	rel := relations[m.selectedRelIdx]

	// Jump to the related node
	m = m.WithFocusedNode(rel.NodeID)

	// Switch to Graph view to see the node in context
	m = m.WithView(ViewGraph)

	// Reset relation selection for next time
	m = m.WithSelectedRelIdx(0)

	return m
}

// IsCollapsed returns true if the node is collapsed (children hidden)
func (m Model) IsCollapsed(nodeID string) bool {
	return m.collapsed[nodeID]
}

// WithGraphScroll returns a new Model with updated graph scroll position.
func (m Model) WithGraphScroll(offset int) Model {
	if offset < 0 {
		offset = 0
	}
	m.graphScroll = offset
	return m
}

// GetGraphScroll returns the current graph scroll offset.
func (m Model) GetGraphScroll() int {
	return m.graphScroll
}

// WithSearchMode returns a new Model with search mode enabled/disabled.
func (m Model) WithSearchMode(enabled bool) Model {
	m.searchMode = enabled
	if !enabled {
		m.searchQuery = ""
	}
	return m
}

// WithSearchQuery returns a new Model with updated search query.
func (m Model) WithSearchQuery(query string) Model {
	m.searchQuery = query
	return m
}

// IsSearchMode returns true if search/filter mode is active.
func (m Model) IsSearchMode() bool {
	return m.searchMode
}

// GetSearchQuery returns the current search query.
func (m Model) GetSearchQuery() string {
	return m.searchQuery
}

// ToggleCollapse toggles the collapsed state of a node
func (m Model) ToggleCollapse(nodeID string) Model {
	// Create a new map to maintain immutability
	newCollapsed := make(map[string]bool)
	for k, v := range m.collapsed {
		newCollapsed[k] = v
	}
	newCollapsed[nodeID] = !newCollapsed[nodeID]
	m.collapsed = newCollapsed
	return m
}

// HasChildren returns true if the node has children in the graph
func (m Model) HasChildren(nodeID string) bool {
	for _, edge := range m.edges {
		if edge.FromID == nodeID && isHierarchicalEdgeType(edge.Relation) {
			return true
		}
	}
	return false
}

// isHierarchicalEdgeType checks if edge represents parent-child relationship
func isHierarchicalEdgeType(relation graph.EdgeType) bool {
	switch relation {
	case graph.EdgeOwns, graph.EdgeImplements, graph.EdgeModifies:
		return true
	default:
		return false
	}
}
