# MAAT Anti-Requirements

**Document Type**: Anti-Requirements Specification (What We Refuse)
**Version**: 1.0
**Date**: 2026-01-05
**Status**: Approved (RMP Quality: 0.85)

---

## Overview

Anti-requirements define what MAAT explicitly **refuses** to do. These are not "nice-to-haves we'll skip" but **constitutional prohibitions** that protect the system's integrity, focus, and trust.

Each anti-requirement:
- References the violated Commandment(s)
- Explains the rationale
- Describes the anti-pattern to avoid
- Provides the correct alternative

---

## A-001: NO Autonomous AI Actions

**Commandment**: #6 Human Contact

### Prohibition
MAAT shall **never** execute AI-suggested actions without explicit human approval.

### Anti-Pattern
```go
// ❌ FORBIDDEN: Auto-applying AI suggestions
func (m Model) HandleClaudeResponse(resp ClaudeResponse) (Model, tea.Cmd) {
    for _, suggestion := range resp.Suggestions {
        applyEdit(suggestion)  // NO! Auto-execution
    }
    return m, nil
}
```

### Correct Pattern
```go
// ✅ CORRECT: Present for human review
func (m Model) HandleClaudeResponse(resp ClaudeResponse) (Model, tea.Cmd) {
    m.pendingSuggestions = resp.Suggestions
    // User must explicitly [a]pprove each
    return m, nil
}
```

### Rationale
AI should augment human judgment, not replace it. Every AI action must be auditable with clear human decision points. This builds trust and prevents runaway automation.

---

## A-002: NO Feature Parity with Linear

**Commandment**: #7 Composition Monopoly

### Prohibition
MAAT shall **never** attempt to replicate Linear's full feature set (views, workflows, custom fields, integrations).

### Anti-Pattern
- Building a "Linear TUI" with identical functionality
- Implementing Linear's filtering/sorting in MAAT
- Recreating Linear's project templates
- Adding Linear-specific workflow editors

### Correct Pattern
- Fetch issues via thin API client
- Display in MAAT's graph paradigm
- Link to Linear for advanced operations
- Value from composition, not features

### Rationale
Competition with Linear is unwinnable and wastes resources. MAAT's value is **unified context across systems**, not being a better Linear. If a user needs full Linear functionality, they should use Linear.

---

## A-003: NO Write Without Confirmation

**Commandment**: #10 Sovereignty Preservation

### Prohibition
MAAT shall **never** write to external systems (Linear, GitHub) without explicit user confirmation.

### Anti-Pattern
```go
// ❌ FORBIDDEN: Silent writes
func (c *LinearClient) UpdateIssue(id string, update Update) error {
    return c.gql.Mutate(updateMutation, update)  // No confirmation!
}
```

### Correct Pattern
```go
// ✅ CORRECT: Return confirmation request
func (c *LinearClient) UpdateIssue(id string, update Update) ConfirmRequest {
    return ConfirmRequest{
        Action:      "Update Issue",
        Description: fmt.Sprintf("Set %s to %s", id, update.Status),
        Execute:     func() error { return c.gql.Mutate(updateMutation, update) },
    }
}
```

### Rationale
External systems are sovereign. MAAT orchestrates but does not colonize. Users must maintain explicit control over what data flows where.

---

## A-004: NO Global Mutable State

**Commandment**: #1 Immutable Truth

### Prohibition
MAAT shall **never** use global variables or mutable singletons for application state.

### Anti-Pattern
```go
// ❌ FORBIDDEN: Global mutable state
var (
    currentUser *User
    issueCache  map[string]Issue
    isLoading   bool
)

func UpdateIssue(id string) {
    issueCache[id].Status = "done"  // Mutation!
    isLoading = true
}
```

### Correct Pattern
```go
// ✅ CORRECT: All state in Model, transforms via Update
type Model struct {
    currentUser User
    issueCache  map[string]Issue
    isLoading   bool
}

func (m Model) Update(msg Msg) (Model, tea.Cmd) {
    switch msg := msg.(type) {
    case IssueUpdated:
        m.issueCache[msg.ID] = msg.Issue  // New map entry, not mutation
        return m, nil
    }
    return m, nil
}
```

### Rationale
Global state creates race conditions, makes testing impossible without mocks, and destroys debuggability. The Elm Architecture's power comes from centralized, immutable state.

---

## A-005: NO Ambient AI Suggestions

**Commandment**: #6 Human Contact

### Prohibition
MAAT shall **never** proactively show AI suggestions, auto-complete, or "helpful tips" without explicit user invocation.

### Anti-Pattern
- Pop-up: "Claude noticed you're looking at a bug. Want suggestions?"
- Auto-complete: Typing in comment field triggers AI completion
- Ambient: "AI analyzed this PR and found 3 issues" (unsolicited)

### Correct Pattern
- User presses `Ctrl+A` to invoke Claude
- User types question explicitly
- AI responds only to direct queries
- No "AI is typing..." unless invoked

### Rationale
Ambient AI is distracting, breaks flow, and erodes user agency. The user should feel in complete control of when AI is active. "Human contact is a tool call" - explicit, bounded, auditable.

---

## A-006: NO Competing with GitHub UI

**Commandment**: #7 Composition Monopoly

### Prohibition
MAAT shall **never** attempt to replicate GitHub's UI functionality (code review, CI/CD, Actions, Discussions).

### Anti-Pattern
- Building inline code review with line comments
- Implementing CI/CD pipeline visualization
- Creating GitHub Actions workflow editor
- Adding repository settings management

### Correct Pattern
- Show PR status and linked issues
- Display commit history with diffs
- Link to GitHub for full review
- Focus on cross-system relationships

### Rationale
GitHub's UI is mature and comprehensive. Competing wastes resources and confuses users. MAAT adds value by **connecting** GitHub to Linear and Claude, not by replacing GitHub.

---

## A-007: NO Auto-Sync Without User Action

**Commandment**: #10 Sovereignty Preservation

### Prohibition
MAAT shall **never** push local changes to external systems automatically via background sync.

### Anti-Pattern
```go
// ❌ FORBIDDEN: Auto-sync in background
func (sm *SyncManager) Start() {
    ticker := time.NewTicker(5 * time.Minute)
    for range ticker.C {
        sm.PushLocalChanges()  // NO! Silent push
    }
}
```

### Correct Pattern
```go
// ✅ CORRECT: Explicit sync command
func (m Model) Update(msg Msg) (Model, tea.Cmd) {
    switch msg.(type) {
    case SyncRequested:  // User pressed 'S' or ran :sync
        return m.WithSyncing(true), syncCmd()
    }
}
```

### Exceptions
- **Pull sync** (fetching new data) is allowed in background
- **Push sync** (writing data) requires explicit user action

### Rationale
Users must maintain control over what leaves their system. Background pushes can cause unexpected side effects (triggering webhooks, notifications, CI runs).

---

## A-008: NO Pointer Mutations in Update

**Commandment**: #8 Async Purity

### Prohibition
MAAT shall **never** use pointer receivers that mutate state within the Update function.

### Anti-Pattern
```go
// ❌ FORBIDDEN: Pointer receiver mutation
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case IssuesFetched:
        m.issues = msg.Issues  // Mutation via pointer!
        m.loading = false
    }
    return m, nil
}
```

### Correct Pattern
```go
// ✅ CORRECT: Value receiver with new model returned
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case IssuesFetched:
        return m.WithIssues(msg.Issues), nil
    }
    return m, nil
}

func (m Model) WithIssues(issues []Issue) Model {
    m.issues = issues
    m.loading = false
    return m  // Return copy
}
```

### Rationale
Pointer mutations break referential transparency, making it impossible to:
- Compare old/new state for debugging
- Implement undo/redo
- Test Update deterministically
- Reason about state transitions

---

## Summary Table

| Anti-FR | Prohibition | Commandment | Severity |
|---------|-------------|-------------|----------|
| A-001 | Autonomous AI actions | #6 | Critical |
| A-002 | Feature parity with Linear | #7 | High |
| A-003 | Write without confirmation | #10 | Critical |
| A-004 | Global mutable state | #1 | Critical |
| A-005 | Ambient AI suggestions | #6 | High |
| A-006 | Competing with GitHub | #7 | High |
| A-007 | Auto-sync without action | #10 | High |
| A-008 | Pointer mutations in Update | #8 | Critical |

---

## Enforcement

### Code Review Checklist
- [ ] No `var` declarations at package level (except constants)
- [ ] No `*Model` receivers on Update
- [ ] No API writes without ConfirmRequest return
- [ ] No `go` statements outside tea.Cmd
- [ ] No AI invocations without user trigger

### Static Analysis
```go
// Linter rules (golangci-lint custom)
// - gochecknoglobals: Forbid global variables
// - gocritic: Detect pointer receiver patterns
// - custom: Detect direct API mutations
```

---

## References

- MAAT Constitution: `/Users/manu/Documents/LUXOR/MAAT/docs/CONSTITUTION.md`
- ADR-001: Elm Architecture
- ADR-004: Human-in-Loop AI
- ADR-005: Thin Integrations
