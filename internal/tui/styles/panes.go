// Package styles provides Lipgloss styling for the MAAT TUI.
package styles

import "github.com/charmbracelet/lipgloss"

// Layout holds calculated dimensions for the 3-pane layout.
type Layout struct {
	GraphWidth  int
	MainWidth   int
	DetailWidth int
	Height      int
	StatusHeight int
}

// CalculateLayout computes pane dimensions based on terminal size.
// Layout: Graph (25%) | Main (50%) | Detail (25%)
// Bottom: Status bar (1 line)
func CalculateLayout(width, height int) Layout {
	// Reserve for borders (2 chars per pane for left/right borders)
	// and gaps between panes (2 chars total)
	contentWidth := width - 8

	// 3-pane layout: 25% | 50% | 25%
	graphWidth := contentWidth * 25 / 100
	mainWidth := contentWidth * 50 / 100
	detailWidth := contentWidth - graphWidth - mainWidth

	// Ensure minimum widths
	if graphWidth < 15 {
		graphWidth = 15
	}
	if mainWidth < 20 {
		mainWidth = 20
	}
	if detailWidth < 15 {
		detailWidth = 15
	}

	// Reserve height for status bar
	statusHeight := 1
	contentHeight := height - statusHeight - 2 // -2 for top/bottom borders

	return Layout{
		GraphWidth:   graphWidth,
		MainWidth:    mainWidth,
		DetailWidth:  detailWidth,
		Height:       contentHeight,
		StatusHeight: statusHeight,
	}
}

// Pane styles for the 3-pane layout

var (
	// BasePaneStyle is the base style for all panes.
	BasePaneStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Border).
			Padding(0, 1)

	// GraphPaneStyle is the style for the left graph pane.
	GraphPaneStyle = BasePaneStyle.
			BorderForeground(GraphPaneBorder)

	// MainPaneStyle is the style for the middle main pane.
	MainPaneStyle = BasePaneStyle.
			BorderForeground(MainPaneBorder)

	// DetailPaneStyle is the style for the right detail pane.
	DetailPaneStyle = BasePaneStyle.
			BorderForeground(DetailPaneBorder)

	// FocusedPaneStyle is applied to the currently active pane.
	FocusedPaneStyle = lipgloss.NewStyle().
				Border(lipgloss.DoubleBorder()).
				BorderForeground(FocusBorder).
				Padding(0, 1)

	// PaneTitleStyle is the style for pane titles.
	PaneTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Accent).
			MarginBottom(1)

	// PaneContentStyle is the style for pane content.
	PaneContentStyle = lipgloss.NewStyle().
				Foreground(Foreground)
)

// StatusBar styles

var (
	// StatusBarStyle is the base style for the status bar.
	StatusBarStyle = lipgloss.NewStyle().
			Background(StatusBarBg).
			Foreground(StatusBarFg).
			Padding(0, 1)

	// StatusBarKeyStyle is the style for key hints in the status bar.
	StatusBarKeyStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(Accent)

	// StatusBarTextStyle is the style for descriptive text in the status bar.
	StatusBarTextStyle = lipgloss.NewStyle().
				Foreground(StatusBarFg)

	// StatusBarErrorStyle is the style for error messages in the status bar.
	StatusBarErrorStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#EF4444")).
				Bold(true)

	// StatusBarLoadingStyle is the style for loading indicator.
	StatusBarLoadingStyle = lipgloss.NewStyle().
				Foreground(StatusInProgress).
				Italic(true)
)

// Loading styles

var (
	// LoadingStyle is the style for the loading message.
	LoadingStyle = lipgloss.NewStyle().
			Foreground(Muted).
			Italic(true)

	// LoadingContainerStyle centers the loading message.
	LoadingContainerStyle = lipgloss.NewStyle().
				Align(lipgloss.Center)
)

// Node styles for the graph pane

var (
	// NodeStyle is the base style for nodes in the graph.
	NodeStyle = lipgloss.NewStyle().
			Foreground(Foreground).
			Padding(0, 1)

	// NodeSelectedStyle is the style for the selected node.
	NodeSelectedStyle = lipgloss.NewStyle().
				Background(lipgloss.AdaptiveColor{Light: "#E4E4E7", Dark: "#3F3F46"}).
				Foreground(Foreground).
				Bold(true).
				Padding(0, 1)

	// NodeTypeIssueStyle shows issue nodes with appropriate icon.
	NodeTypeIssueStyle = lipgloss.NewStyle().
				Foreground(Primary)

	// NodeTypePRStyle shows PR nodes with appropriate icon.
	NodeTypePRStyle = lipgloss.NewStyle().
			Foreground(Secondary)

	// NodeTypeCommitStyle shows commit nodes with appropriate icon.
	NodeTypeCommitStyle = lipgloss.NewStyle().
				Foreground(Muted)

	// NodeTypeFileStyle shows file nodes with appropriate icon.
	NodeTypeFileStyle = lipgloss.NewStyle().
				Foreground(Foreground)
)

// GetPaneStyle returns the appropriate style for a pane based on focus state.
func GetPaneStyle(isFocused bool, baseStyle lipgloss.Style) lipgloss.Style {
	if isFocused {
		return FocusedPaneStyle
	}
	return baseStyle
}

// RenderGraphPane creates a styled graph pane with dimensions.
func RenderGraphPane(content string, width, height int, isFocused bool) string {
	style := GetPaneStyle(isFocused, GraphPaneStyle)
	return style.
		Width(width).
		Height(height).
		Render(content)
}

// RenderMainPane creates a styled main pane with dimensions.
func RenderMainPane(content string, width, height int, isFocused bool) string {
	style := GetPaneStyle(isFocused, MainPaneStyle)
	return style.
		Width(width).
		Height(height).
		Render(content)
}

// RenderDetailPane creates a styled detail pane with dimensions.
func RenderDetailPane(content string, width, height int, isFocused bool) string {
	style := GetPaneStyle(isFocused, DetailPaneStyle)
	return style.
		Width(width).
		Height(height).
		Render(content)
}

// RenderStatusBar creates the styled status bar.
func RenderStatusBar(content string, width int) string {
	return StatusBarStyle.
		Width(width).
		Render(content)
}
