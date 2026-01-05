# MAAT Constitution: The 10 Design Commandments

**Document Type**: Constitutional Design Principles
**Created**: 2026-01-05
**Methodology**: JUPITER Cross-Domain Pattern Exchange
**Source Domains**: Elm Architecture, 12-Factor Agents, Unix Philosophy, Peter Thiel, Functional Programming, Spec-Driven Development

---

## Preamble

MAAT (Modular Agentic Architecture for Terminal) is a unified terminal workspace that integrates Linear, GitHub, and Claude Code into a single navigable context. These commandments establish the constitutional principles that govern all architectural decisions.

**The Central Insight**: MAAT's value is in composition, not competition. We do not build a better Linear, a better GitHub, or a better Claude Code. We build the only place where all three compose into unified context.

---

## The 10 Commandments

### Commandment 1: The Immutable Truth

**Source Exchange**: Elm Architecture ⇄ Functional Programming

| Domain | Contribution |
|--------|--------------|
| Elm | Model-View-Update - state transitions are explicit, traceable, debuggable |
| FP | Immutability - data flows forward, never mutates in place |

**Principle**: State shall never be mutated; it shall only be transformed through pure functions that return new state.

**Design Decision**: All MAAT state lives in a single immutable Model struct. The Update function receives (Model, Msg) and returns (Model, Cmd). No pointer receivers that mutate. No global state. No side effects in reducers.

```go
// YES: Pure transformation
func Update(model Model, msg Msg) (Model, tea.Cmd) {
    switch msg := msg.(type) {
    case IssuesFetched:
        return model.WithIssues(msg.Issues), nil
    }
    return model, nil
}

// NO: Mutation (FORBIDDEN)
func (m *Model) Update(msg Msg) {
    m.issues = msg.Issues
}
```

**Anti-Pattern**: Scattered state mutations. Multiple sources of truth. Pointer receivers modifying state. Global variables.

---

### Commandment 2: The Single Responsibility Sovereignty

**Source Exchange**: Unix Philosophy ⇄ 12-Factor Agents

| Domain | Contribution |
|--------|--------------|
| Unix | "Do one thing well" - singular, clear purpose |
| 12-Factor | Small focused agents - complexity from composition |

**Principle**: Each component shall have exactly one reason to change; composition of simple parts yields complex capability.

**Design Decision**: MAAT decomposes into sovereign domains:
- `linear/` - Linear API only
- `github/` - GitHub API only
- `claude/` - Claude interface only
- `graph/` - Graph rendering only
- `tui/` - Terminal UI only

No component knows about others' internals. Communication through central Model.

**Anti-Pattern**: God objects. Components importing each other. Business logic in UI code.

---

### Commandment 3: The Text Interface Covenant

**Source Exchange**: Unix Philosophy ⇄ Spec-Driven Development

| Domain | Contribution |
|--------|--------------|
| Unix | Text as universal interface - human-readable, debuggable |
| Spec-Driven | Specification as source of truth - explicit, verifiable |

**Principle**: All inter-component communication shall be through explicit, serializable data structures.

**Design Decision**: Messages (Msg types) are the universal interface. Every state change is a message. Messages can be serialized to JSON for debugging.

```go
type IssuesFetched struct {
    Issues []linear.Issue `json:"issues"`
    Error  error          `json:"error,omitempty"`
}
```

**Anti-Pattern**: Callbacks mutating state. Channels bypassing message flow. Implicit state changes.

---

### Commandment 4: The 10x Differentiation

**Source Exchange**: Peter Thiel ⇄ Elm Architecture

| Domain | Contribution |
|--------|--------------|
| Thiel | 10x better, not 10% - paradigm shift beats incremental |
| Elm | State machine UI modes - discrete states, explicit transitions |

**Principle**: MAAT shall not be "a slightly better terminal"; it shall be a fundamentally different paradigm - a spatial, navigable workspace where context is never lost.

**Design Decision**: The graph IS the primary navigation paradigm. Every item is a node. Relationships are edges. Drill-down is zoom.

```go
type ViewMode int

const (
    GraphView ViewMode = iota  // Primary: spatial navigation
    ListFallback               // Secondary: when graph overwhelms
    DetailPane                 // Contextual: never loses graph
)
```

**Anti-Pattern**: Building "Linear CLI with colors." Adding features to match competitors. Graph as optional visualization.

---

### Commandment 5: The Controlled Effect Boundary

**Source Exchange**: Functional Programming ⇄ 12-Factor Agents

| Domain | Contribution |
|--------|--------------|
| FP | Explicit effect handling - isolated, tracked, controlled |
| 12-Factor | Own your control flow - never cede to frameworks |

**Principle**: Side effects (API calls, file I/O, Claude invocations) shall occur only at explicit boundaries, initiated by Commands.

**Design Decision**: Update is pure. It returns Commands describing effects. Bubble Tea runtime executes effects. MAAT owns when and how effects occur.

```go
// Effect is isolated, testable, replaceable
func fetchIssuesCmd(teamID string) tea.Cmd {
    return func() tea.Msg {
        issues, err := linearClient.FetchIssues(teamID)
        return IssuesFetched{Issues: issues, Error: err}
    }
}
```

**Anti-Pattern**: API calls inside Update. Goroutines spawned without commands. Unmockable effects.

---

### Commandment 6: The Human Contact Protocol

**Source Exchange**: 12-Factor Agents ⇄ Spec-Driven Development

| Domain | Contribution |
|--------|--------------|
| 12-Factor | Human contact as tool call - bounded, auditable |
| Spec-Driven | Anti-speculation discipline - never assume |

**Principle**: Claude integration shall be an explicit tool invocation with clear I/O boundaries, not ambient presence.

**Design Decision**: Claude invoked through explicit commands (Ctrl+A). Input context explicitly assembled. Output presented for review. No auto-suggestions. Human remains in control.

```go
type ClaudeRequest struct {
    Context     string   `json:"context"`
    Question    string   `json:"question"`
    Constraints []string `json:"constraints"`
}

// Human reviews before execution
func (m Model) HandleClaudeResponse(resp ClaudeResponse) (Model, tea.Cmd) {
    return m.WithPendingActions(resp.Actions), nil  // Show, don't execute
}
```

**Anti-Pattern**: Auto-completing without request. Executing suggestions without confirmation. Ambient AI.

---

### Commandment 7: The Composition Monopoly

**Source Exchange**: Peter Thiel ⇄ Unix Philosophy

| Domain | Contribution |
|--------|--------------|
| Thiel | Competition is for losers - find monopoly |
| Unix | Composability - power from combining simple tools |

**Principle**: MAAT's monopoly is unified context - the only place where Linear, GitHub, Claude, and visualization compose.

**Design Decision**: MAAT does not compete with Linear's UI, GitHub's interface, or Claude Code. It composes them. Each integration is thin. Value emerges from composition.

```go
type WorkspaceGraph struct {
    Nodes map[string]Node  // Issues, PRs, Files unified
    Edges []Edge           // Relationships across systems
}
```

**Anti-Pattern**: Building "better Linear." Feature parity as goal. Recreating existing capabilities.

---

### Commandment 8: The Async Purity Discipline

**Source Exchange**: Elm Architecture ⇄ Functional Programming

| Domain | Contribution |
|--------|--------------|
| Elm | Async via commands - returns to pure update cycle |
| FP | Pure functions for testability - deterministic I/O |

**Principle**: Async operations shall be described, not performed, within Update; runtime executes, results return as messages.

**Design Decision**: tea.Cmd is a description. Update returns commands. Runtime executes. Results are messages. Update is 100% testable.

```go
func TestIssueRefresh(t *testing.T) {
    model := NewModel()
    newModel, cmd := model.Update(RefreshRequested{})
    assert.True(t, newModel.Loading)
    assert.NotNil(t, cmd)  // Described, not executed
}
```

**Anti-Pattern**: Async inside Update. Testing requiring mocks. Non-deterministic tests.

---

### Commandment 9: The Specification Constitution

**Source Exchange**: Spec-Driven Development ⇄ Peter Thiel

| Domain | Contribution |
|--------|--------------|
| Spec-Driven | Constitution-first governance - spec precedes code |
| Thiel | Zero-to-one - design for end state |

**Principle**: Architecture decisions shall be documented before implementation; specification governs code.

**Design Decision**: Every significant choice becomes an ADR (Architecture Decision Record). Written BEFORE implementation. Code violating ADRs is unconstitutional.

```
docs/adr/
├── 001-elm-architecture.md
├── 002-graph-first-navigation.md
├── 003-explicit-effects.md
├── 004-human-in-loop-ai.md
└── 005-thin-integrations.md
```

**Anti-Pattern**: Implementing first. "Code is documentation." Specifications drifting from reality.

---

### Commandment 10: The Sovereignty Preservation

**Source Exchange**: 12-Factor Agents ⇄ Unix Philosophy

| Domain | Contribution |
|--------|--------------|
| 12-Factor | Stateless reducer - explicit state, pure transitions |
| Unix | Small sharp tools - domain sovereignty preserved |

**Principle**: Each integrated system shall retain its sovereignty; MAAT orchestrates without colonizing.

**Design Decision**: MAAT never writes without explicit user action. Read is free; write requires confirmation. Each system's model is respected.

```go
// Write requires explicit confirmation
func (c *LinearClient) UpdateIssue(id string, update IssueUpdate) tea.Cmd {
    return func() tea.Msg {
        return ConfirmationRequired{
            Action:      "Update Issue",
            Description: fmt.Sprintf("Set %s status to %s", id, update.Status),
            OnConfirm:   actualUpdateCmd(id, update),
        }
    }
}
```

**Anti-Pattern**: Auto-syncing changes. Modifying without confirmation. Flattening system-specific concepts.

---

## Summary Table

| # | Commandment | Core Principle |
|---|-------------|----------------|
| 1 | Immutable Truth | State transforms, never mutates |
| 2 | Single Responsibility | One component, one purpose |
| 3 | Text Interface | Messages are the specification |
| 4 | 10x Differentiation | Graph navigation, not lists |
| 5 | Controlled Effects | Commands describe, runtime executes |
| 6 | Human Contact | Explicit invocation, human review |
| 7 | Composition Monopoly | Value from integration, not features |
| 8 | Async Purity | Deterministic updates, testable code |
| 9 | Specification Constitution | ADRs before implementation |
| 10 | Sovereignty Preservation | Orchestrate, never colonize |

---

## Domain Contributions

Each domain contributed unique wisdom through JUPITER's exchange principle:

| Domain | Key Contribution |
|--------|------------------|
| **Elm Architecture** | State machines, pure updates, async commands |
| **12-Factor Agents** | Explicit control, human contact, stateless reducers |
| **Unix Philosophy** | Composition, sovereignty, text interfaces |
| **Peter Thiel** | 10x thinking, monopoly positioning, zero-to-one |
| **Functional Programming** | Immutability, effect isolation, testability |
| **Spec-Driven Development** | Constitution-first, anti-speculation, verification |

---

**Status**: RATIFIED
**Applies To**: All MAAT development
**Amendment Process**: ADR required, reviewed by maintainers
