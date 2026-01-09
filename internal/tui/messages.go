package tui

// Message types define the TUI API (Commandment #3: Text Interface)
// All async operations communicate via these message types

// WindowSizeMsg is sent when the terminal is resized
type WindowSizeMsg struct {
	Width  int
	Height int
}

// ErrorOccurred is sent when an operation fails
type ErrorOccurred struct {
	Err error
}

// DataLoadedMsg is sent when async data fetch completes
type DataLoadedMsg struct {
	Data interface{}
}

// GraphDataLoadedMsg is sent when graph data is loaded
type GraphDataLoadedMsg struct {
	Nodes []DisplayNode
	Edges []DisplayEdge
}

// RefreshRequested is sent when user presses 'r'
type RefreshRequested struct{}

// AIInvoked is sent when user presses Ctrl+A (Commandment #6: Human Contact)
type AIInvoked struct{}

// ConfirmationRequested is sent when an external write is attempted (Commandment #10: Sovereignty)
type ConfirmationRequested struct {
	Action  string
	Execute func() error
}

// ConfirmationAccepted is sent when user confirms an action
type ConfirmationAccepted struct{}

// ConfirmationRejected is sent when user rejects an action
type ConfirmationRejected struct{}

// NavigateDown is sent when user presses Enter (Commandment #4: Navigation Monopoly)
type NavigateDown struct{}

// NavigateUp is sent when user presses Esc
type NavigateUp struct{}
