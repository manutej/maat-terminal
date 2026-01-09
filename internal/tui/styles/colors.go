// Package styles provides Lipgloss styling for the MAAT TUI.
// All colors and styles are defined here for consistent visual hierarchy.
package styles

import "github.com/charmbracelet/lipgloss"

// Adaptive colors for light/dark terminals
var (
	// Primary palette
	Primary   = lipgloss.AdaptiveColor{Light: "#5A56E0", Dark: "#7C78FF"}
	Secondary = lipgloss.AdaptiveColor{Light: "#6E6AE0", Dark: "#8E8AFF"}
	Accent    = lipgloss.AdaptiveColor{Light: "#00D084", Dark: "#00E898"}

	// Status colors (Linear-inspired)
	StatusTodo       = lipgloss.Color("#6B7280") // Gray
	StatusInProgress = lipgloss.Color("#F59E0B") // Amber
	StatusDone       = lipgloss.Color("#10B981") // Green
	StatusCanceled   = lipgloss.Color("#EF4444") // Red
	StatusBlocked    = lipgloss.Color("#8B5CF6") // Purple

	// Priority colors
	PriorityUrgent = lipgloss.Color("#DC2626")
	PriorityHigh   = lipgloss.Color("#F97316")
	PriorityMedium = lipgloss.Color("#FBBF24")
	PriorityLow    = lipgloss.Color("#6B7280")

	// UI colors
	Background = lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#1A1A2E"}
	Foreground = lipgloss.AdaptiveColor{Light: "#1A1A2E", Dark: "#E4E4E7"}
	Border     = lipgloss.AdaptiveColor{Light: "#E4E4E7", Dark: "#3F3F46"}
	Muted      = lipgloss.AdaptiveColor{Light: "#9CA3AF", Dark: "#71717A"}

	// Git colors
	GitAdded    = lipgloss.Color("#22C55E")
	GitModified = lipgloss.Color("#EAB308")
	GitDeleted  = lipgloss.Color("#EF4444")

	// Pane-specific colors
	GraphPaneBorder  = lipgloss.AdaptiveColor{Light: "#5A56E0", Dark: "#7C78FF"}
	MainPaneBorder   = lipgloss.AdaptiveColor{Light: "#00D084", Dark: "#00E898"}
	DetailPaneBorder = lipgloss.AdaptiveColor{Light: "#6E6AE0", Dark: "#8E8AFF"}

	// Focus indicator
	FocusBorder = lipgloss.AdaptiveColor{Light: "#00D084", Dark: "#00E898"}

	// Status bar colors
	StatusBarBg = lipgloss.AdaptiveColor{Light: "#E4E4E7", Dark: "#27273A"}
	StatusBarFg = lipgloss.AdaptiveColor{Light: "#1A1A2E", Dark: "#A1A1AA"}
)

// StatusColor returns the appropriate color for a given status string.
func StatusColor(status string) lipgloss.Color {
	switch status {
	case "todo":
		return StatusTodo
	case "in_progress":
		return StatusInProgress
	case "done":
		return StatusDone
	case "canceled":
		return StatusCanceled
	case "blocked":
		return StatusBlocked
	default:
		return StatusTodo
	}
}

// PriorityColor returns the appropriate color for a given priority level.
// Priority: 1 = Urgent, 2 = High, 3 = Medium, 4+ = Low
func PriorityColor(priority int) lipgloss.Color {
	switch priority {
	case 1:
		return PriorityUrgent
	case 2:
		return PriorityHigh
	case 3:
		return PriorityMedium
	default:
		return PriorityLow
	}
}
