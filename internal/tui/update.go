package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// Init initializes the model (Bubble Tea lifecycle)
func (m Model) Init() tea.Cmd {
	// If model already has data (loaded from main.go), don't fetch mock data
	if len(m.nodes) > 0 {
		return nil
	}
	return fetchData()
}

// Update handles all messages (Commandment #1: VALUE receiver, no pointer mutation)
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Window resize - also handles initial ready state
	case tea.WindowSizeMsg:
		m = m.WithSize(msg.Width, msg.Height)
		if !m.ready {
			m = m.WithReady(true)
			// Only fetch mock data if no data was pre-loaded
			if len(m.nodes) == 0 {
				return m, fetchData()
			}
		}
		return m, nil

	// Keyboard input
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	// Custom messages
	case DataLoadedMsg:
		return m.WithData(msg.Data), nil

	case GraphDataLoadedMsg:
		// Load graph nodes and edges into model
		m = m.WithNodes(msg.Nodes).WithEdges(msg.Edges).WithLoading(false)
		return m, nil

	case ErrorOccurred:
		return m.WithError(msg.Err), nil

	case RefreshRequested:
		return m.WithLoading(true), refreshData()

	case AIInvoked:
		// Commandment #6: Human Contact - AI requires explicit Ctrl+A
		// Placeholder for Phase 4+ AI integration
		return m.WithData("AI invoked - feature coming soon"), nil

	case ConfirmationRequested:
		// Commandment #10: Sovereignty - external writes require confirmation
		return m.WithConfirmation(&ConfirmationRequest{
			Action:  msg.Action,
			Execute: msg.Execute,
		}), nil

	case ConfirmationAccepted:
		if m.confirmation != nil {
			req := m.confirmation
			return m.WithConfirmation(nil), executeConfirmedAction(req.Execute)
		}
		return m, nil

	case ConfirmationRejected:
		return m.WithConfirmation(nil).PopView(), nil

	case NavigateDown:
		// Commandment #4: Navigation Monopoly - Enter drills down
		return m.PushView(ViewDetails), nil

	case NavigateUp:
		// Commandment #4: Navigation Monopoly - Esc backs out
		return m.PopView(), nil
	}

	return m, nil
}

// handleKeyPress processes keyboard input
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Handle confirmation view separately
	if m.currentView == ViewConfirm {
		return m.handleConfirmationKeys(msg)
	}

	// Handle search mode input
	if m.searchMode {
		return m.handleSearchInput(msg)
	}

	// Global keybindings
	switch {
	case key.Matches(msg, m.keys.Quit):
		// Commandment #9: Terminal Citizenship - Ctrl+C exits
		return m, tea.Quit

	case key.Matches(msg, m.keys.Enter):
		// Drill down - behavior depends on view
		if m.currentView == ViewRelations {
			// Jump to selected relation's node
			return m.jumpToSelectedRelation(), nil
		}
		// In Graph view, toggle collapse for projects/nodes with children
		if m.currentView == ViewGraph {
			if m.HasChildren(m.focusedNode) {
				return m.ToggleCollapse(m.focusedNode), nil
			}
			// For leaf nodes (issues), show details
			return m.WithView(ViewDetails), nil
		}
		return m.Update(NavigateDown{})

	case key.Matches(msg, m.keys.Back):
		// Back up
		if m.navStack.IsEmpty() {
			// At top level, Esc does nothing
			return m, nil
		}
		return m.Update(NavigateUp{})

	case key.Matches(msg, m.keys.Refresh):
		return m.Update(RefreshRequested{})

	case key.Matches(msg, m.keys.AI):
		return m.Update(AIInvoked{})

	case key.Matches(msg, m.keys.Up):
		// k key - behavior depends on view
		if m.currentView == ViewRelations {
			return m.moveRelationUp(), nil
		}
		return m.HandleNavigation("k"), nil

	case key.Matches(msg, m.keys.Down):
		// j key - behavior depends on view
		if m.currentView == ViewRelations {
			return m.moveRelationDown(), nil
		}
		return m.HandleNavigation("j"), nil

	case key.Matches(msg, m.keys.Left):
		// h key - move focus left (spatial)
		return m.HandleNavigation("h"), nil

	case key.Matches(msg, m.keys.Right):
		// l key - move focus right (spatial)
		return m.HandleNavigation("l"), nil
	}

	// Handle Tab for view cycling (single-pane design: Graph → Details → Relations)
	switch msg.String() {
	case "tab":
		// Cycle forward through views
		m = m.WithView(m.currentView.CycleView())
		return m, nil
	case "shift+tab":
		// Cycle backward through views
		var newView ViewMode
		switch m.currentView {
		case ViewGraph:
			newView = ViewRelations
		case ViewDetails:
			newView = ViewGraph
		case ViewRelations:
			newView = ViewDetails
		default:
			newView = ViewGraph
		}
		m = m.WithView(newView)
		return m, nil
	case "f":
		// Cycle filter mode (only in Graph view)
		if m.currentView == ViewGraph {
			m = m.WithFilterMode(m.filterMode.CycleFilter())
			// Reset focus to first filtered node if current focus is filtered out
			filteredNodes := m.GetFilteredNodes()
			if len(filteredNodes) > 0 {
				found := false
				for _, node := range filteredNodes {
					if node.ID == m.focusedNode {
						found = true
						break
					}
				}
				if !found {
					m = m.WithFocusedNode(filteredNodes[0].ID)
				}
			}
		}
		return m, nil
	case "s":
		// Cycle status filter (only in Graph view)
		if m.currentView == ViewGraph {
			m = m.WithStatusFilter(m.statusFilter.CycleStatusFilter())
			// Reset scroll and focus if current focus is filtered out
			m = m.WithGraphScroll(0)
			filteredNodes := m.GetFilteredNodes()
			if len(filteredNodes) > 0 {
				found := false
				for _, node := range filteredNodes {
					if node.ID == m.focusedNode {
						found = true
						break
					}
				}
				if !found {
					m = m.WithFocusedNode(filteredNodes[0].ID)
				}
			}
		}
		return m, nil
	case "/":
		// Enter search mode (only in Graph view)
		if m.currentView == ViewGraph {
			m = m.WithSearchMode(true)
		}
		return m, nil
	}

	return m, nil
}

// handleSearchInput processes input while in search/filter mode
func (m Model) handleSearchInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		// Exit search mode and clear query
		return m.WithSearchMode(false), nil

	case tea.KeyEnter:
		// Exit search mode but keep filter active
		m.searchMode = false
		// Focus on first matching node if any
		filteredNodes := m.GetFilteredNodes()
		if len(filteredNodes) > 0 {
			m = m.WithFocusedNode(filteredNodes[0].ID)
			m = m.WithGraphScroll(0)
		}
		return m, nil

	case tea.KeyBackspace:
		// Remove last character from query
		if len(m.searchQuery) > 0 {
			m = m.WithSearchQuery(m.searchQuery[:len(m.searchQuery)-1])
		}
		return m, nil

	case tea.KeyRunes:
		// Add typed characters to query
		m = m.WithSearchQuery(m.searchQuery + string(msg.Runes))
		// Auto-focus on first matching node
		filteredNodes := m.GetFilteredNodes()
		if len(filteredNodes) > 0 {
			m = m.WithFocusedNode(filteredNodes[0].ID)
			m = m.WithGraphScroll(0)
		}
		return m, nil

	case tea.KeyCtrlC:
		return m, tea.Quit
	}

	return m, nil
}

// handleConfirmationKeys processes keys in confirmation view
func (m Model) handleConfirmationKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y", "enter":
		return m.Update(ConfirmationAccepted{})
	case "n", "N", "esc":
		return m.Update(ConfirmationRejected{})
	case "ctrl+c", "q":
		return m, tea.Quit
	}
	return m, nil
}
