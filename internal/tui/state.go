package tui

import (
	"strings"

	"github.com/manutej/maat-terminal/internal/graph"
)

// ViewMode represents the current full-screen view in single-pane design.
// Tab key cycles through: Graph → Details → Relations → Graph...
type ViewMode int

const (
	ViewGraph     ViewMode = iota // Full-screen hierarchical graph
	ViewDetails                   // Full-screen node details
	ViewRelations                 // Full-screen relationship view
	ViewConfirm                   // Confirmation dialog (overlay)
)

// FilterMode controls which node types are displayed in the graph
type FilterMode int

const (
	FilterAll      FilterMode = iota // Show all nodes (overwhelming)
	FilterProjects                   // Projects + Issues + PRs only (useful default)
	FilterIssues                     // Issues only
	FilterPRs                        // PRs only
	FilterFiles                      // Files only
	FilterCommits                    // Commits only
)

// StatusFilter controls which statuses are displayed
type StatusFilter int

const (
	StatusAll        StatusFilter = iota // Show all statuses
	StatusActive                         // In Progress only (active work)
	StatusNotDone                        // In Progress + Backlog (hide completed)
	StatusDone                           // Done only (completed work)
)

// StatusFilterString returns the display name for the status filter
func (s StatusFilter) String() string {
	switch s {
	case StatusAll:
		return "All Status"
	case StatusActive:
		return "Active Only"
	case StatusNotDone:
		return "Not Done"
	case StatusDone:
		return "Done Only"
	default:
		return "Unknown"
	}
}

// CycleStatusFilter returns the next status filter
func (s StatusFilter) CycleStatusFilter() StatusFilter {
	switch s {
	case StatusAll:
		return StatusActive
	case StatusActive:
		return StatusNotDone
	case StatusNotDone:
		return StatusDone
	case StatusDone:
		return StatusAll
	default:
		return StatusAll
	}
}

// MatchesStatus returns true if the given status passes this filter
func (s StatusFilter) MatchesStatus(status string) bool {
	statusLower := strings.ToLower(status)

	switch s {
	case StatusAll:
		return true
	case StatusActive:
		// Only In Progress items
		return statusLower == "in progress" || statusLower == "in_progress" ||
			statusLower == "started" || statusLower == "in review"
	case StatusNotDone:
		// Everything except Done
		return statusLower != "done" && statusLower != "completed" &&
			statusLower != "merged" && statusLower != "closed"
	case StatusDone:
		// Only completed items
		return statusLower == "done" || statusLower == "completed" ||
			statusLower == "merged" || statusLower == "closed"
	default:
		return true
	}
}

// FilterModeTypes returns the node types to show for each filter mode
func (f FilterMode) Types() []graph.NodeType {
	switch f {
	case FilterAll:
		return nil // nil means show all
	case FilterProjects:
		return []graph.NodeType{graph.NodeTypeProject, graph.NodeTypeIssue, graph.NodeTypePR, graph.NodeTypeService}
	case FilterIssues:
		return []graph.NodeType{graph.NodeTypeIssue}
	case FilterPRs:
		return []graph.NodeType{graph.NodeTypePR}
	case FilterFiles:
		return []graph.NodeType{graph.NodeTypeFile}
	case FilterCommits:
		return []graph.NodeType{graph.NodeTypeCommit}
	default:
		return nil
	}
}

// String returns the display name for the filter mode
func (f FilterMode) String() string {
	switch f {
	case FilterAll:
		return "All"
	case FilterProjects:
		return "Projects"
	case FilterIssues:
		return "Issues"
	case FilterPRs:
		return "PRs"
	case FilterFiles:
		return "Files"
	case FilterCommits:
		return "Commits"
	default:
		return "Unknown"
	}
}

// CycleFilter returns the next filter mode
func (f FilterMode) CycleFilter() FilterMode {
	switch f {
	case FilterProjects:
		return FilterIssues
	case FilterIssues:
		return FilterPRs
	case FilterPRs:
		return FilterCommits
	case FilterCommits:
		return FilterFiles
	case FilterFiles:
		return FilterAll
	case FilterAll:
		return FilterProjects
	default:
		return FilterProjects
	}
}

// String returns the string representation of ViewMode
func (v ViewMode) String() string {
	switch v {
	case ViewGraph:
		return "Graph"
	case ViewDetails:
		return "Details"
	case ViewRelations:
		return "Relations"
	case ViewConfirm:
		return "Confirm"
	default:
		return "Unknown"
	}
}

// CycleView returns the next view in the Tab cycle sequence.
func (v ViewMode) CycleView() ViewMode {
	switch v {
	case ViewGraph:
		return ViewDetails
	case ViewDetails:
		return ViewRelations
	case ViewRelations:
		return ViewGraph
	default:
		return ViewGraph
	}
}

// NavigationStack maintains history for Esc navigation
type NavigationStack struct {
	stack []ViewMode
}

// NewNavigationStack creates an empty navigation stack
func NewNavigationStack() NavigationStack {
	return NavigationStack{
		stack: make([]ViewMode, 0),
	}
}

// Push adds a new view to the stack
func (n NavigationStack) Push(mode ViewMode) NavigationStack {
	newStack := make([]ViewMode, len(n.stack)+1)
	copy(newStack, n.stack)
	newStack[len(n.stack)] = mode
	return NavigationStack{stack: newStack}
}

// Pop removes the top view from the stack
func (n NavigationStack) Pop() (NavigationStack, ViewMode, bool) {
	if len(n.stack) == 0 {
		return n, ViewGraph, false
	}

	mode := n.stack[len(n.stack)-1]
	newStack := make([]ViewMode, len(n.stack)-1)
	copy(newStack, n.stack[:len(n.stack)-1])

	return NavigationStack{stack: newStack}, mode, true
}

// IsEmpty checks if the stack has no entries
func (n NavigationStack) IsEmpty() bool {
	return len(n.stack) == 0
}
