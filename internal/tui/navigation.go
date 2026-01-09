package tui

// HandleNavigation is the main navigation handler that routes keys to specific handlers.
// Pure function following Commandment #1 (Immutable Truth).
func (m Model) HandleNavigation(key string) Model {
	switch key {
	case "h":
		return m.moveLeft()
	case "j":
		return m.moveDown()
	case "k":
		return m.moveUp()
	case "l":
		return m.moveRight()
	default:
		return m
	}
}

// moveLeft implements h key - navigate to parent node.
func (m Model) moveLeft() Model {
	if len(m.nodes) == 0 || m.focusedNode == "" {
		return m
	}

	// Find parent nodes (edges pointing TO current node)
	parents := getParentNodes(m.focusedNode, m.GetFilteredEdges())
	if len(parents) > 0 {
		// Check if parent is in filtered set
		for _, parentID := range parents {
			if m.isNodeInFilter(parentID) {
				return m.WithFocusedNode(parentID)
			}
		}
	}

	// No navigable parent - stay at current
	return m
}

// moveRight implements l key - navigate to first child node.
func (m Model) moveRight() Model {
	if len(m.nodes) == 0 || m.focusedNode == "" {
		return m
	}

	// Find child nodes (edges pointing FROM current node)
	children := getChildNodes(m.focusedNode, m.GetFilteredEdges())
	if len(children) > 0 {
		// Check if child is in filtered set
		for _, childID := range children {
			if m.isNodeInFilter(childID) {
				return m.WithFocusedNode(childID)
			}
		}
	}

	// No navigable child - stay at current
	return m
}

// moveUp implements k key - navigate to previous node in tree order.
func (m Model) moveUp() Model {
	filteredNodes := m.GetFilteredNodes()
	if len(filteredNodes) == 0 || m.focusedNode == "" {
		return m
	}

	// Build tree and get flattened list
	tree := buildTree(filteredNodes, m.GetFilteredEdges())
	flatList := flattenTreeWithCollapse(tree, m)

	// Find current index and move up
	currentIdx := -1
	for i, id := range flatList {
		if id == m.focusedNode {
			currentIdx = i
			break
		}
	}

	var newIdx int
	if currentIdx > 0 {
		newIdx = currentIdx - 1
	} else if len(flatList) > 0 {
		// Already at top - wrap to bottom
		newIdx = len(flatList) - 1
	} else {
		return m
	}

	// Update focused node and adjust scroll to ensure visibility
	m = m.WithFocusedNode(flatList[newIdx])
	m = m.ensureFocusVisible(newIdx, len(flatList))
	return m
}

// moveDown implements j key - navigate to next node in tree order.
func (m Model) moveDown() Model {
	filteredNodes := m.GetFilteredNodes()
	if len(filteredNodes) == 0 || m.focusedNode == "" {
		return m
	}

	// Build tree and get flattened list
	tree := buildTree(filteredNodes, m.GetFilteredEdges())
	flatList := flattenTreeWithCollapse(tree, m)

	// Find current index and move down
	currentIdx := -1
	for i, id := range flatList {
		if id == m.focusedNode {
			currentIdx = i
			break
		}
	}

	var newIdx int
	if currentIdx >= 0 && currentIdx < len(flatList)-1 {
		newIdx = currentIdx + 1
	} else if len(flatList) > 0 {
		// Already at bottom - wrap to top
		newIdx = 0
	} else {
		return m
	}

	// Update focused node and adjust scroll to ensure visibility
	m = m.WithFocusedNode(flatList[newIdx])
	m = m.ensureFocusVisible(newIdx, len(flatList))
	return m
}

// flattenTree returns node IDs in tree traversal order (depth-first)
func flattenTree(tree TreeStructure) []string {
	result := make([]string, 0, len(tree.Nodes))
	visited := make(map[string]bool)

	var visit func(nodeID string)
	visit = func(nodeID string) {
		if visited[nodeID] {
			return
		}
		visited[nodeID] = true
		result = append(result, nodeID)

		for _, childID := range tree.Children[nodeID] {
			visit(childID)
		}
	}

	// Visit all roots
	for _, rootID := range tree.Roots {
		visit(rootID)
	}

	// Add any unvisited nodes (orphans not in tree)
	for id := range tree.Nodes {
		if !visited[id] {
			result = append(result, id)
		}
	}

	return result
}

// isNodeInFilter checks if a node ID is in the current filtered set
func (m Model) isNodeInFilter(nodeID string) bool {
	for _, node := range m.GetFilteredNodes() {
		if node.ID == nodeID {
			return true
		}
	}
	return false
}

// getParentNodes returns all parent nodes (nodes with edges pointing TO this node).
// Used for h key (move left) to follow parent relationships.
func getParentNodes(nodeID string, edges []DisplayEdge) []string {
	parents := make([]string, 0)
	seen := make(map[string]bool)

	for _, edge := range edges {
		if edge.ToID == nodeID && !seen[edge.FromID] {
			parents = append(parents, edge.FromID)
			seen[edge.FromID] = true
		}
	}

	return parents
}

// getChildNodes returns all child nodes (nodes with edges pointing FROM this node).
// Used for l key (move right) to follow child relationships.
func getChildNodes(nodeID string, edges []DisplayEdge) []string {
	children := make([]string, 0)
	seen := make(map[string]bool)

	for _, edge := range edges {
		if edge.FromID == nodeID && !seen[edge.ToID] {
			children = append(children, edge.ToID)
			seen[edge.ToID] = true
		}
	}

	return children
}

// flattenTreeWithCollapse returns visible node IDs respecting collapsed state
func flattenTreeWithCollapse(tree TreeStructure, m Model) []string {
	result := make([]string, 0, len(tree.Nodes))
	visited := make(map[string]bool)

	var visit func(nodeID string)
	visit = func(nodeID string) {
		if visited[nodeID] {
			return
		}
		visited[nodeID] = true
		result = append(result, nodeID)

		// Only visit children if not collapsed
		if !m.IsCollapsed(nodeID) {
			for _, childID := range tree.Children[nodeID] {
				visit(childID)
			}
		}
	}

	// Visit all roots
	for _, rootID := range tree.Roots {
		visit(rootID)
	}

	// Add any unvisited nodes (orphans not in tree)
	for id := range tree.Nodes {
		if !visited[id] {
			result = append(result, id)
		}
	}

	return result
}

// ensureFocusVisible adjusts scroll to keep focused item visible
func (m Model) ensureFocusVisible(focusedIdx int, totalItems int) Model {
	// Calculate visible area (reserve 4 lines for header/footer)
	visibleLines := m.height - 6
	if visibleLines < 5 {
		visibleLines = 5
	}

	// Ensure scroll keeps focused item visible with some context
	padding := 2 // Keep 2 lines of context above/below

	// If focused item is above visible area, scroll up
	if focusedIdx < m.graphScroll+padding {
		newScroll := focusedIdx - padding
		if newScroll < 0 {
			newScroll = 0
		}
		return m.WithGraphScroll(newScroll)
	}

	// If focused item is below visible area, scroll down
	if focusedIdx >= m.graphScroll+visibleLines-padding {
		newScroll := focusedIdx - visibleLines + padding + 1
		if newScroll < 0 {
			newScroll = 0
		}
		// Don't scroll past the end
		maxScroll := totalItems - visibleLines
		if maxScroll < 0 {
			maxScroll = 0
		}
		if newScroll > maxScroll {
			newScroll = maxScroll
		}
		return m.WithGraphScroll(newScroll)
	}

	return m
}
