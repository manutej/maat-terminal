package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/manutej/maat-terminal/internal/graph"
	"github.com/manutej/maat-terminal/internal/tui/styles"
)

// View renders the entire UI - PURE FUNCTION from state.
// SINGLE-PANE DESIGN: Tab cycles between Graph/Details/Relations views.
// No split panes - full terminal width/height for focused content.
func (m Model) View() string {
	// Handle not-ready state with loading message
	if !m.ready {
		return m.renderLoadingScreen()
	}

	// Handle confirmation dialog overlay
	if m.confirmation != nil {
		return m.renderConfirmDialog()
	}

	// Render current view mode (full screen)
	return m.renderCurrentView()
}

// renderLoadingScreen shows a loading message while waiting for window size.
func (m Model) renderLoadingScreen() string {
	loadingMsg := styles.LoadingStyle.Render("Initializing MAAT...")

	// Center the loading message
	return styles.LoadingContainerStyle.
		Width(m.width).
		Height(m.height).
		Render(loadingMsg)
}

// renderCurrentView renders the full-screen view based on currentView mode.
func (m Model) renderCurrentView() string {
	// Reserve space for status bar (2 lines)
	contentHeight := m.height - 2

	// Render content based on current view mode
	var content string
	switch m.currentView {
	case ViewGraph:
		content = m.renderGraphView(m.width, contentHeight)
	case ViewDetails:
		content = m.renderDetailsView(m.width, contentHeight)
	case ViewRelations:
		content = m.renderRelationsView(m.width, contentHeight)
	default:
		content = m.renderGraphView(m.width, contentHeight)
	}

	// Render status bar
	statusBar := m.renderStatusBar()

	// Stack content and status bar vertically
	return lipgloss.JoinVertical(
		lipgloss.Left,
		content,
		statusBar,
	)
}

// renderGraphView renders the full-screen hierarchical graph view.
func (m Model) renderGraphView(width, height int) string {
	var builder strings.Builder

	// View title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.Accent).
		Width(width).
		Align(lipgloss.Center).
		MarginBottom(1)

	builder.WriteString(titleStyle.Render("📊 Knowledge Graph"))
	builder.WriteString("\n")

	// Render graph with full terminal width
	if len(m.nodes) == 0 {
		noDataMsg := styles.LoadingStyle.Render("No nodes loaded. Press 'r' to refresh.")
		builder.WriteString(lipgloss.NewStyle().
			Width(width).
			Height(height - 3).
			Align(lipgloss.Center, lipgloss.Center).
			Render(noDataMsg))
	} else {
		// Use hierarchical tree rendering with FULL WIDTH (no pane constraint)
		graphViz := RenderGraph(m, width-4) // -4 for padding

		// Apply scrolling - split into lines and show only visible portion
		lines := strings.Split(graphViz, "\n")
		visibleHeight := height - 4 // Reserve for title and margins

		// Calculate scroll bounds
		scrollStart := m.graphScroll
		if scrollStart < 0 {
			scrollStart = 0
		}
		if scrollStart >= len(lines) {
			scrollStart = 0
		}

		scrollEnd := scrollStart + visibleHeight
		if scrollEnd > len(lines) {
			scrollEnd = len(lines)
		}

		// Show only visible lines
		if scrollStart < len(lines) {
			visibleLines := lines[scrollStart:scrollEnd]
			builder.WriteString(strings.Join(visibleLines, "\n"))
		}

		// Show scroll indicator if content is scrolled
		if len(lines) > visibleHeight {
			scrollInfo := lipgloss.NewStyle().
				Foreground(styles.Muted).
				Faint(true).
				Render(fmt.Sprintf("\n[%d-%d of %d lines]", scrollStart+1, scrollEnd, len(lines)))
			builder.WriteString(scrollInfo)
		}
	}

	return builder.String()
}

// renderDetailsView renders the full-screen details view for focused node.
func (m Model) renderDetailsView(width, height int) string {
	var builder strings.Builder

	// View title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.Accent).
		Width(width).
		Align(lipgloss.Center).
		MarginBottom(1)

	builder.WriteString(titleStyle.Render("📝 Node Details"))
	builder.WriteString("\n")

	// Get focused node
	node, ok := m.GetFocusedNode()
	if !ok {
		noSelectionMsg := styles.PaneContentStyle.Render("No node selected. Press Tab to view Graph and select a node.")
		builder.WriteString(lipgloss.NewStyle().
			Width(width).
			Height(height - 3).
			Align(lipgloss.Center, lipgloss.Center).
			Render(noSelectionMsg))
		return builder.String()
	}

	// Render detailed node information (centered, max 80 chars wide)
	contentWidth := 80
	if width < 80 {
		contentWidth = width - 4
	}

	detailsBox := m.renderNodeDetailsExpanded(node, contentWidth)
	centeredDetails := lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		Render(detailsBox)

	builder.WriteString(centeredDetails)

	return builder.String()
}

// renderRelationsView renders the full-screen relationship view with interactive selection.
func (m Model) renderRelationsView(width, height int) string {
	var builder strings.Builder

	// View title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.Accent).
		Width(width).
		Align(lipgloss.Center).
		MarginBottom(1)

	builder.WriteString(titleStyle.Render("🔗 Relationships (j/k to select, Enter to jump)"))
	builder.WriteString("\n")

	// Get focused node
	node, ok := m.GetFocusedNode()
	if !ok {
		noSelectionMsg := styles.PaneContentStyle.Render("No node selected. Press Tab to view Graph and select a node.")
		builder.WriteString(lipgloss.NewStyle().
			Width(width).
			Height(height - 3).
			Align(lipgloss.Center, lipgloss.Center).
			Render(noSelectionMsg))
		return builder.String()
	}

	// Render interactive relationship list
	contentWidth := 100
	if width < 100 {
		contentWidth = width - 4
	}

	relationsBox := m.renderInteractiveRelationsList(node, contentWidth)
	centeredRelations := lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		Render(relationsBox)

	builder.WriteString(centeredRelations)

	return builder.String()
}

// renderInteractiveRelationsList renders relations with selection highlighting.
func (m Model) renderInteractiveRelationsList(node DisplayNode, maxWidth int) string {
	var lines []string

	// Header with node context
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.Accent)
	lines = append(lines, headerStyle.Render(fmt.Sprintf("Relationships for: %s", node.Title)))
	lines = append(lines, "")

	relations := m.GetRelationsList()

	if len(relations) == 0 {
		noRelStyle := lipgloss.NewStyle().
			Foreground(styles.Muted).
			Italic(true)
		lines = append(lines, noRelStyle.Render("No relationships found for this node."))
		return strings.Join(lines, "\n")
	}

	// Group by direction
	var outgoing, incoming []RelationItem
	for _, rel := range relations {
		if rel.IsOutgoing {
			outgoing = append(outgoing, rel)
		} else {
			incoming = append(incoming, rel)
		}
	}

	idx := 0

	// Outgoing relationships
	if len(outgoing) > 0 {
		outgoingStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(styles.Primary)
		lines = append(lines, outgoingStyle.Render("→ Outgoing Relations:"))
		lines = append(lines, "")

		for _, rel := range outgoing {
			line := m.renderRelationLine(rel, idx, maxWidth)
			lines = append(lines, line)
			idx++
		}
		lines = append(lines, "")
	}

	// Incoming relationships
	if len(incoming) > 0 {
		incomingStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(styles.Secondary)
		lines = append(lines, incomingStyle.Render("← Incoming Relations:"))
		lines = append(lines, "")

		for _, rel := range incoming {
			line := m.renderRelationLine(rel, idx, maxWidth)
			lines = append(lines, line)
			idx++
		}
	}

	// Summary and instructions
	lines = append(lines, "")
	summaryStyle := lipgloss.NewStyle().Foreground(styles.Muted)
	lines = append(lines, summaryStyle.Render(fmt.Sprintf(
		"Total: %d outgoing, %d incoming | j/k: navigate | Enter: jump to selected",
		len(outgoing),
		len(incoming),
	)))

	return strings.Join(lines, "\n")
}

// renderRelationLine renders a single relation with selection highlighting.
func (m Model) renderRelationLine(rel RelationItem, idx int, maxWidth int) string {
	isSelected := idx == m.selectedRelIdx

	// Style based on selection
	var lineStyle lipgloss.Style
	if isSelected {
		lineStyle = lipgloss.NewStyle().
			Background(styles.Primary).
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true).
			Width(maxWidth - 4)
	} else {
		lineStyle = lipgloss.NewStyle().
			Foreground(styles.Foreground)
	}

	// Build relation display
	icon := getNodeIcon(rel.NodeType)
	arrow := "→"
	if !rel.IsOutgoing {
		arrow = "←"
	}

	relTypeStyle := lipgloss.NewStyle().Foreground(styles.Accent)
	if isSelected {
		relTypeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
	}

	// Format: [idx] icon Title ← relation
	content := fmt.Sprintf("  %s %s %s %s",
		icon,
		truncate(rel.NodeTitle, 40),
		arrow,
		relTypeStyle.Render(rel.Relation),
	)

	if isSelected {
		content = "▶ " + content[2:] // Replace leading spaces with indicator
	}

	return lineStyle.Render(content)
}

// renderSearchBar renders the search input bar when in search mode.
func (m Model) renderSearchBar() string {
	// Search prompt style
	promptStyle := lipgloss.NewStyle().
		Foreground(styles.Accent).
		Bold(true)

	// Input style with cursor
	inputStyle := lipgloss.NewStyle().
		Foreground(styles.Foreground)

	// Hint style
	hintStyle := lipgloss.NewStyle().
		Foreground(styles.Muted).
		Faint(true)

	// Count matching nodes
	filteredNodes := m.GetFilteredNodes()
	countText := fmt.Sprintf("(%d matches)", len(filteredNodes))

	// Build search bar content
	content := fmt.Sprintf("%s %s%s  %s  %s",
		promptStyle.Render("/"),
		inputStyle.Render(m.searchQuery),
		inputStyle.Render("█"), // Cursor
		hintStyle.Render(countText),
		hintStyle.Render("Enter:select | Esc:cancel"),
	)

	return styles.RenderStatusBar(content, m.width)
}

// renderStatusBar renders the bottom status bar with view indicator.
func (m Model) renderStatusBar() string {
	// If in search mode, show search input prominently
	if m.searchMode {
		return m.renderSearchBar()
	}

	var parts []string

	// Show current view mode with clear indicator
	viewText := styles.StatusBarKeyStyle.Render(fmt.Sprintf("[%s]", m.currentView.String()))
	parts = append(parts, viewText)

	// Show filter mode in Graph view
	if m.currentView == ViewGraph {
		filterText := styles.StatusBarTextStyle.Render(fmt.Sprintf("Type: %s", m.filterMode.String()))
		parts = append(parts, filterText)

		// Show status filter if not "All"
		if m.statusFilter != StatusAll {
			statusFilterText := styles.StatusBarKeyStyle.Render(fmt.Sprintf("Status: %s", m.statusFilter.String()))
			parts = append(parts, statusFilterText)
		}

		// Show active search query if any
		if m.searchQuery != "" {
			searchText := styles.StatusBarKeyStyle.Render(fmt.Sprintf("Search: \"%s\"", m.searchQuery))
			parts = append(parts, searchText)
		}
	}

	// Show focused node if any
	if node, ok := m.GetFocusedNode(); ok {
		nodeText := styles.StatusBarTextStyle.Render(fmt.Sprintf("→ %s", truncate(node.Title, 25)))
		parts = append(parts, nodeText)
	}

	// Show loading indicator
	if m.loading {
		loadingText := styles.StatusBarLoadingStyle.Render("Loading...")
		parts = append(parts, loadingText)
	}

	// Show error if any
	if m.err != nil {
		errText := styles.StatusBarErrorStyle.Render(fmt.Sprintf("Error: %v", m.err))
		parts = append(parts, errText)
	}

	// Add key hints on the right (updated for filter and search)
	var keyHints string
	switch m.currentView {
	case ViewGraph:
		keyHints = styles.StatusBarTextStyle.Render("/:search | f:type | s:status | jk:nav | Enter:toggle | q:quit")
	case ViewDetails:
		keyHints = styles.StatusBarTextStyle.Render("Tab:Relations | Esc:back | q:quit")
	case ViewRelations:
		relations := m.GetRelationsList()
		if len(relations) > 0 {
			keyHints = styles.StatusBarTextStyle.Render(fmt.Sprintf("jk:select (%d/%d) | Enter:jump | Tab:Graph | q:quit", m.selectedRelIdx+1, len(relations)))
		} else {
			keyHints = styles.StatusBarTextStyle.Render("Tab:Graph | q:quit")
		}
	default:
		keyHints = styles.StatusBarTextStyle.Render("Tab:view | Esc:back | q:quit")
	}

	// Join left and right parts
	leftContent := strings.Join(parts, " | ")

	// Calculate spacing to right-align key hints
	leftLen := lipgloss.Width(leftContent)
	rightLen := lipgloss.Width(keyHints)
	spacing := m.width - leftLen - rightLen - 4 // -4 for padding
	if spacing < 2 {
		spacing = 2
	}

	fullContent := leftContent + strings.Repeat(" ", spacing) + keyHints

	return styles.RenderStatusBar(fullContent, m.width)
}

// renderConfirmDialog renders the confirmation dialog overlay.
func (m Model) renderConfirmDialog() string {
	if m.confirmation == nil {
		return m.renderCurrentView()
	}

	// Render dialog box
	dialogStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Accent).
		Padding(1, 2).
		Width(50).
		Align(lipgloss.Center)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.Foreground).
		MarginBottom(1)

	contentStyle := lipgloss.NewStyle().
		Foreground(styles.Foreground)

	buttonStyle := lipgloss.NewStyle().
		MarginTop(1)

	yesButton := lipgloss.NewStyle().
		Background(styles.Accent).
		Foreground(lipgloss.Color("#000000")).
		Padding(0, 2).
		Bold(true).
		Render("[y] Yes")

	noButton := lipgloss.NewStyle().
		Background(styles.Muted).
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(0, 2).
		Render("[n] No")

	dialog := dialogStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Center,
			titleStyle.Render("Confirm Action"),
			contentStyle.Render(m.confirmation.Action),
			buttonStyle.Render(
				lipgloss.JoinHorizontal(lipgloss.Top, yesButton, "  ", noButton),
			),
		),
	)

	// Center dialog on screen
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		dialog,
	)
}

// renderNodeDetailsExpanded renders comprehensive node details (for Details view).
func (m Model) renderNodeDetailsExpanded(node DisplayNode, maxWidth int) string {
	var lines []string

	// Node icon and title (large and prominent)
	icon := getNodeIcon(node.Type)
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.Accent).
		Underline(true)

	// Show identifier if available (e.g., CET-352)
	titleText := node.Title
	if node.Identifier != "" {
		titleText = fmt.Sprintf("[%s] %s", node.Identifier, node.Title)
	}
	lines = append(lines, titleStyle.Render(fmt.Sprintf("%s %s", icon, titleText)))
	lines = append(lines, "")

	// Type and Project badges on same line
	typeStyle := lipgloss.NewStyle().
		Background(styles.Primary).
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(0, 2).
		Bold(true)

	badgeLine := typeStyle.Render(fmt.Sprintf("Type: %s", node.Type))
	if node.Project != "" {
		projectStyle := lipgloss.NewStyle().
			Background(styles.Secondary).
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 2)
		badgeLine += "  " + projectStyle.Render(fmt.Sprintf("📦 %s", node.Project))
	}
	lines = append(lines, badgeLine)
	lines = append(lines, "")

	// Status with color and icon
	if node.Status != "" {
		statusColor := styles.StatusColor(node.Status)
		statusStyle := lipgloss.NewStyle().
			Foreground(statusColor).
			Bold(true)
		statusIcon := getStatusIconLarge(node.Status)
		lines = append(lines, statusStyle.Render(fmt.Sprintf("%s Status: %s", statusIcon, node.Status)))
	}

	// Priority with color and badge
	if node.Priority > 0 {
		priorityColor := styles.PriorityColor(node.Priority)
		priorityStyle := lipgloss.NewStyle().
			Foreground(priorityColor).
			Bold(true)
		priorityLabel := getPriorityLabel(node.Priority)
		lines = append(lines, priorityStyle.Render(fmt.Sprintf("🔥 Priority: %s", priorityLabel)))
	}

	lines = append(lines, "")

	// Description (wrapped to maxWidth)
	if node.Description != "" {
		descStyle := lipgloss.NewStyle().
			Foreground(styles.Foreground).
			Width(maxWidth)
		lines = append(lines, descStyle.Render("Description:"))
		lines = append(lines, descStyle.Render(wrapText(node.Description, maxWidth-4)))
	}

	// Labels as badges
	if len(node.Labels) > 0 {
		lines = append(lines, "")
		var labelParts []string
		labelParts = append(labelParts, "🏷  Labels: ")
		for _, label := range node.Labels {
			labelStyle := lipgloss.NewStyle().
				Background(styles.Secondary).
				Foreground(lipgloss.Color("#FFFFFF")).
				Padding(0, 1)
			labelParts = append(labelParts, labelStyle.Render(label)+" ")
		}
		lines = append(lines, strings.Join(labelParts, ""))
	}

	// Related nodes preview (quick glance at connections)
	relations := m.GetRelationsList()
	if len(relations) > 0 {
		lines = append(lines, "")
		lines = append(lines, "")
		relHeader := lipgloss.NewStyle().
			Bold(true).
			Foreground(styles.Secondary)
		lines = append(lines, relHeader.Render(fmt.Sprintf("🔗 Related (%d connections):", len(relations))))

		// Show first 5 relations as preview
		maxPreview := 5
		if len(relations) < maxPreview {
			maxPreview = len(relations)
		}

		for i := 0; i < maxPreview; i++ {
			rel := relations[i]
			relIcon := getNodeIcon(rel.NodeType)
			arrow := "→"
			if !rel.IsOutgoing {
				arrow = "←"
			}
			relLine := fmt.Sprintf("  %s %s %s (%s)", relIcon, truncate(rel.NodeTitle, 30), arrow, rel.Relation)
			lines = append(lines, lipgloss.NewStyle().Foreground(styles.Muted).Render(relLine))
		}

		if len(relations) > maxPreview {
			moreStyle := lipgloss.NewStyle().Foreground(styles.Muted).Italic(true)
			lines = append(lines, moreStyle.Render(fmt.Sprintf("  ... and %d more (Tab to Relations view)", len(relations)-maxPreview)))
		}
	}

	// URL link (if available)
	if node.URL != "" {
		lines = append(lines, "")
		urlStyle := lipgloss.NewStyle().
			Foreground(styles.Accent).
			Underline(true)
		linkLabel := lipgloss.NewStyle().
			Foreground(styles.Muted).
			Bold(true)
		lines = append(lines, linkLabel.Render("🔗 Link: ")+urlStyle.Render(node.URL))
	}

	// ID (faint, at bottom)
	lines = append(lines, "")
	idStyle := lipgloss.NewStyle().Foreground(styles.Muted).Faint(true)
	lines = append(lines, idStyle.Render(fmt.Sprintf("ID: %s", node.ID)))

	// Helpful hint if no description
	if node.Description == "" && node.Type == graph.NodeTypeIssue {
		lines = append(lines, "")
		hintStyle := lipgloss.NewStyle().
			Foreground(styles.Muted).
			Italic(true)
		lines = append(lines, hintStyle.Render("💡 Description not loaded. Use the link above to view full details in Linear."))
	}

	return strings.Join(lines, "\n")
}


// Helper functions

// getNodeIcon returns an icon character for a node type.
func getNodeIcon(nodeType graph.NodeType) string {
	switch nodeType {
	case graph.NodeTypeIssue:
		return "🐛"
	case graph.NodeTypePR:
		return "🔀"
	case graph.NodeTypeCommit:
		return "💾"
	case graph.NodeTypeFile:
		return "📄"
	case graph.NodeTypeProject:
		return "📦"
	case graph.NodeTypeService:
		return "⚙️"
	default:
		return "❓"
	}
}

// getStatusIconLarge returns a larger icon for status (Details view).
func getStatusIconLarge(status string) string {
	switch status {
	case "done", "merged", "completed":
		return "✅"
	case "in_progress", "open", "in progress":
		return "🔄"
	case "todo", "pending":
		return "📋"
	case "blocked":
		return "🚫"
	case "canceled":
		return "❌"
	default:
		return "⚪"
	}
}

// getPriorityLabel returns a human-readable priority label.
func getPriorityLabel(priority int) string {
	switch priority {
	case 1:
		return "Urgent"
	case 2:
		return "High"
	case 3:
		return "Medium"
	default:
		return "Low"
	}
}

// truncate shortens a string to max length with ellipsis.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// wrapText wraps text to fit within maxWidth.
func wrapText(text string, maxWidth int) string {
	if len(text) <= maxWidth {
		return text
	}

	var lines []string
	words := strings.Fields(text)
	currentLine := ""

	for _, word := range words {
		if len(currentLine)+len(word)+1 > maxWidth {
			if currentLine != "" {
				lines = append(lines, strings.TrimSpace(currentLine))
				currentLine = word
			} else {
				// Word itself is longer than maxWidth
				lines = append(lines, word[:maxWidth])
				currentLine = word[maxWidth:]
			}
		} else {
			if currentLine == "" {
				currentLine = word
			} else {
				currentLine += " " + word
			}
		}
	}

	if currentLine != "" {
		lines = append(lines, strings.TrimSpace(currentLine))
	}

	return strings.Join(lines, "\n")
}
