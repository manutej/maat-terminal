# ADR-001: Elm Architecture with Immutable State

**Status**: Accepted
**Date**: 2026-01-05
**Deciders**: MAAT Architecture Team
**Commandments**: #1 Immutable Truth, #8 Async Purity

---

## Context

MAAT requires a state management approach that is:
- Predictable and debuggable
- Testable without mocks
- Safe for concurrent access
- Compatible with TUI frameworks

Traditional approaches using mutable state and callbacks lead to:
- Race conditions in UI updates
- Difficult-to-trace state changes
- Non-deterministic test behavior
- Spaghetti code from callback chains

## Decision

Adopt **The Elm Architecture (TEA)** implemented via Bubble Tea framework:

```go
type Model struct {
    // Single source of truth - all state here
    graph       WorkspaceGraph
    viewMode    ViewMode
    focusedNode string
    currentRole Role

    // Integration state
    linear      LinearState
    github      GitHubState
    claude      ClaudeState

    // UI state
    width, height int
    ready         bool
}

// Pure function - no side effects
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case IssuesFetched:
        return m.WithIssues(msg.Issues), nil
    case RefreshRequested:
        return m.WithLoading(true), fetchIssuesCmd(m.linear.TeamID)
    }
    return m, nil
}

// Pure rendering from state
func (m Model) View() string {
    if !m.ready {
        return "Loading..."
    }
    return m.renderGraph()
}
```

### Key Principles

1. **Immutable State**: Model is never mutated; Update returns new Model
2. **Pure Update**: No side effects in Update function
3. **Commands for Effects**: All I/O via tea.Cmd
4. **Messages as Events**: All state changes triggered by messages

### Value Receivers (Not Pointer)

```go
// YES: Value receiver returns new model
func (m Model) WithIssues(issues []Issue) Model {
    m.linear.Issues = issues
    m.loading = false
    return m
}

// NO: Pointer receiver mutates (FORBIDDEN)
func (m *Model) SetIssues(issues []Issue) {
    m.linear.Issues = issues  // Mutation!
}
```

## Consequences

### Positive
- **Testability**: Update is 100% testable without mocks
- **Debuggability**: State transitions are explicit and traceable
- **Concurrency Safety**: No shared mutable state
- **Time Travel**: Can replay message sequences for debugging

### Negative
- **Learning Curve**: Developers unfamiliar with FP patterns
- **Memory Allocation**: New model per update (Go GC handles well)
- **Verbosity**: More code than direct mutation

### Mitigations
- Document patterns in CONTRIBUTING.md
- Provide helper methods (WithX pattern)
- Profile memory if issues arise (unlikely for TUI scale)

## Compliance

This ADR enforces:
- **Commandment #1**: State transforms, never mutates
- **Commandment #8**: Deterministic updates, testable code

## References

- [Bubble Tea Framework](https://github.com/charmbracelet/bubbletea)
- [The Elm Architecture](https://guide.elm-lang.org/architecture/)
- MAAT Constitution: `/Users/manu/Documents/LUXOR/MAAT/docs/CONSTITUTION.md`
