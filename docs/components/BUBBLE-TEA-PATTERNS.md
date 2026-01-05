# Bubble Tea Component Mapping

**Source**: Context7 + Official Documentation
**Library**: `github.com/charmbracelet/bubbletea`
**Version**: v0.25+ (stable), v2.0.0-beta available
**Purpose**: Core TUI framework implementing Elm Architecture

---

## DELTA FORCE: Exact Code for MAAT

### 1. Core Model Interface

**MAAT Requirement**: FR-001, FR-002, FR-003, FR-004
**Commandment**: #1 Immutable Truth, #8 Async Purity

```go
// maat/internal/tui/model.go
package tui

import (
    "github.com/charmbracelet/bubbletea"
    "maat/internal/graph"
    "maat/internal/linear"
    "maat/internal/github"
)

// Model is the single source of truth for MAAT state
type Model struct {
    // Core graph state
    graph       *graph.KnowledgeGraph
    focusedNode string
    navStack    NavigationStack

    // View state
    viewMode    ViewMode
    currentRole Role
    ready       bool

    // Window dimensions
    width, height int

    // Integration state (thin clients)
    linear LinearState
    github GitHubState
    claude ClaudeState

    // UI components (Bubbles)
    viewport viewport.Model
    list     list.Model
    search   textinput.Model

    // Pending operations
    pendingConfirm *ConfirmRequest
    loading        bool
    err            error
}

// Init returns initial command - fetch data
func (m Model) Init() tea.Cmd {
    return tea.Batch(
        m.fetchLinearIssues(),
        m.fetchGitHubPRs(),
    )
}

// Update handles all messages - PURE FUNCTION, no side effects
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        return m.handleWindowSize(msg)
    case tea.KeyMsg:
        return m.handleKeyPress(msg)
    case IssuesFetched:
        return m.handleIssuesFetched(msg)
    case PRsFetched:
        return m.handlePRsFetched(msg)
    case ErrorOccurred:
        return m.handleError(msg)
    }
    return m, nil
}

// View renders UI - PURE FUNCTION from state
func (m Model) View() string {
    if !m.ready {
        return m.renderLoading()
    }
    if m.pendingConfirm != nil {
        return m.renderConfirmDialog()
    }
    return m.renderMainView()
}
```

### 2. Message Types

**MAAT Requirement**: #3 Text Interface
**Pattern**: All state changes via explicit messages

```go
// maat/internal/tui/messages.go
package tui

// Data fetching messages
type IssuesFetched struct {
    Issues []linear.Issue
    Error  error
}

type PRsFetched struct {
    PRs   []github.PullRequest
    Error error
}

type CommitsFetched struct {
    Commits []git.Commit
    Error   error
}

// Navigation messages
type NodeSelected struct {
    NodeID string
}

type DrillDown struct {
    NodeID string
}

type NavigateBack struct{}

// Action messages
type ConfirmAction struct {
    Confirmed bool
}

type SyncRequested struct{}

type ClaudeInvoked struct {
    Question string
    Context  ClaudeContext
}

type ClaudeResponded struct {
    Response ClaudeResponse
    Error    error
}

// Error wrapper
type ErrorOccurred struct {
    Source string
    Error  error
}
```

### 3. Command Patterns

**MAAT Requirement**: #5 Controlled Effects
**Pattern**: Commands describe effects, runtime executes

```go
// maat/internal/tui/commands.go
package tui

import (
    tea "github.com/charmbracelet/bubbletea"
    "maat/internal/linear"
)

// Fetch issues - async operation via tea.Cmd
func (m Model) fetchLinearIssues() tea.Cmd {
    return func() tea.Msg {
        issues, err := m.linear.Client.FetchIssues(context.Background())
        if err != nil {
            return ErrorOccurred{Source: "linear", Error: err}
        }
        return IssuesFetched{Issues: issues}
    }
}

// Fetch PRs - can run in parallel with issues
func (m Model) fetchGitHubPRs() tea.Cmd {
    return func() tea.Msg {
        prs, err := m.github.Client.FetchPRs(context.Background())
        if err != nil {
            return ErrorOccurred{Source: "github", Error: err}
        }
        return PRsFetched{PRs: prs}
    }
}

// Batch multiple commands - run in parallel
func (m Model) refreshAll() tea.Cmd {
    return tea.Batch(
        m.fetchLinearIssues(),
        m.fetchGitHubPRs(),
        m.fetchCommits(),
    )
}

// Sequence commands - run in order
func (m Model) createAndSync() tea.Cmd {
    return tea.Sequence(
        m.createIssue(),
        m.syncGraph(),
    )
}
```

### 4. Key Handling Pattern

**MAAT Requirement**: FR-002 Keyboard Navigation

```go
// maat/internal/tui/keys.go
package tui

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    // Global keys first
    switch msg.String() {
    case "ctrl+c", "q":
        return m, tea.Quit
    case "?":
        return m.toggleHelp(), nil
    case "/":
        return m.enterSearchMode(), nil
    case "ctrl+a":
        return m.invokeClaudePanel(), nil
    }

    // Mode-specific handling
    switch m.viewMode {
    case GraphView:
        return m.handleGraphKeys(msg)
    case SearchMode:
        return m.handleSearchKeys(msg)
    case ClaudeMode:
        return m.handleClaudeKeys(msg)
    case ConfirmMode:
        return m.handleConfirmKeys(msg)
    }

    return m, nil
}

func (m Model) handleGraphKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    switch msg.String() {
    case "h", "left":
        return m.navigateLeft(), nil
    case "j", "down":
        return m.navigateDown(), nil
    case "k", "up":
        return m.navigateUp(), nil
    case "l", "right":
        return m.navigateRight(), nil
    case "enter":
        return m.drillDown(), nil
    case "esc", "backspace":
        return m.navigateBack(), nil
    case "tab":
        return m.cyclePanes(), nil
    case "s":
        return m.openStatusPicker(), nil
    case "r":
        return m, m.refreshAll()
    }
    return m, nil
}
```

### 5. State Transformation Helpers

**MAAT Requirement**: #1 Immutable Truth
**Pattern**: WithX methods return new model

```go
// maat/internal/tui/state.go
package tui

// WithX pattern - returns copy with modification
func (m Model) WithIssues(issues []linear.Issue) Model {
    m.linear.Issues = issues
    m.loading = false
    return m
}

func (m Model) WithPRs(prs []github.PullRequest) Model {
    m.github.PRs = prs
    return m
}

func (m Model) WithFocus(nodeID string) Model {
    m.focusedNode = nodeID
    return m
}

func (m Model) WithError(err error) Model {
    m.err = err
    m.loading = false
    return m
}

func (m Model) WithLoading(loading bool) Model {
    m.loading = loading
    return m
}

func (m Model) WithConfirm(req *ConfirmRequest) Model {
    m.pendingConfirm = req
    m.viewMode = ConfirmMode
    return m
}

func (m Model) WithViewMode(mode ViewMode) Model {
    m.viewMode = mode
    return m
}
```

---

## Critical Rules

1. **NEVER use pointer receivers** on Model in Update
2. **NEVER spawn goroutines** - use tea.Cmd exclusively
3. **NEVER mutate state** - always return new Model
4. **ALWAYS handle WindowSizeMsg** for responsive layout
5. **ALWAYS use tea.Batch** for parallel operations
6. **ALWAYS use tea.Sequence** for dependent operations

---

## Performance Patterns

```go
// Lazy initialization - wait for window size
func (m Model) handleWindowSize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
    m.width = msg.Width
    m.height = msg.Height

    if !m.ready {
        // Initialize components with correct dimensions
        m.viewport = viewport.New(msg.Width, msg.Height-4)
        m.list = list.New([]list.Item{}, list.NewDefaultDelegate(), msg.Width, msg.Height-4)
        m.ready = true

        // Now safe to fetch data
        return m, m.refreshAll()
    }

    // Resize existing components
    m.viewport.Width = msg.Width
    m.viewport.Height = msg.Height - 4
    m.list.SetSize(msg.Width, msg.Height-4)

    return m, nil
}

// Alt screen for full TUI experience
func main() {
    p := tea.NewProgram(
        NewModel(),
        tea.WithAltScreen(),      // Full screen mode
        tea.WithMouseCellMotion(), // Mouse support
    )
    if _, err := p.Run(); err != nil {
        log.Fatal(err)
    }
}
```

---

## Integration Points

| MAAT Component | Bubble Tea Element | Location |
|----------------|-------------------|----------|
| Graph View | Custom View() | `tui/views/graph.go` |
| Issue List | bubbles/list | `tui/views/list.go` |
| Detail Pane | bubbles/viewport | `tui/views/detail.go` |
| Search | bubbles/textinput | `tui/views/search.go` |
| Confirm Dialog | Custom View() | `tui/views/confirm.go` |
| Loading | bubbles/spinner | `tui/views/loading.go` |

---

## File Structure

```
maat/internal/tui/
├── model.go       # Model struct definition
├── update.go      # Update function handlers
├── view.go        # View rendering
├── commands.go    # tea.Cmd definitions
├── messages.go    # Msg type definitions
├── keys.go        # Key handling
├── state.go       # WithX helpers
└── views/
    ├── graph.go   # Graph rendering
    ├── list.go    # List view
    ├── detail.go  # Detail pane
    ├── search.go  # Search overlay
    └── confirm.go # Confirmation dialog
```
