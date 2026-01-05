# MAAT Project Instructions

## Project Overview

MAAT (Modular Agentic Architecture for Terminal) is a unified terminal workspace integrating Linear, GitHub, and Claude Code into a navigable knowledge graph.

**Repository**: https://github.com/manutej/maat-terminal
**Language**: Go 1.21+
**Framework**: Bubble Tea (Charmbracelet)

## 10 Design Commandments

All code must comply with these constitutional principles:

| Tier | # | Commandment | Enforcement |
|------|---|-------------|-------------|
| 1 | #1 | **Immutable Truth** | No pointer receivers on Update, no global state |
| 1 | #6 | **Human Contact** | AI requires explicit `Ctrl+A` invocation |
| 1 | #10 | **Sovereignty** | External writes require ConfirmRequest |
| 2 | #2 | **Graph Supremacy** | All entities as nodes/edges |
| 2 | #4 | **Navigation Monopoly** | Enter drills, Esc backs |
| 2 | #7 | **Composition** | Thin API clients only |
| 3 | #3 | **Text Interface** | Msg types are the API |
| 3 | #5 | **Controlled Effects** | tea.Cmd only, no goroutines |
| 3 | #8 | **Async Purity** | Commands describe, runtime executes |
| 3 | #9 | **Terminal Citizenship** | Dark theme, Ctrl+C exits |

## Code Patterns

### Model-View-Update (REQUIRED)

```go
// ✅ CORRECT: Value receiver, returns new Model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case IssuesFetched:
        return m.WithIssues(msg.Issues), nil
    }
    return m, nil
}

// ❌ FORBIDDEN: Pointer receiver mutation
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    m.issues = msg.Issues  // NEVER DO THIS
    return m, nil
}
```

### Async Operations (REQUIRED)

```go
// ✅ CORRECT: tea.Cmd for async
func (m Model) fetchData() tea.Cmd {
    return func() tea.Msg {
        data, err := api.Fetch()
        if err != nil {
            return ErrorOccurred{Error: err}
        }
        return DataFetched{Data: data}
    }
}

// ❌ FORBIDDEN: Direct goroutines
go func() {
    data := api.Fetch()  // NEVER DO THIS
}()
```

### External Writes (REQUIRED)

```go
// ✅ CORRECT: Return ConfirmRequest
func (c *Client) UpdateIssue(id string, update Update) ConfirmRequest {
    return ConfirmRequest{
        Action:  "Update Issue",
        Execute: func() error { return c.api.Mutate(update) },
    }
}

// ❌ FORBIDDEN: Silent writes
func (c *Client) UpdateIssue(id string, update Update) error {
    return c.api.Mutate(update)  // NEVER without confirmation
}
```

## Project Structure

```
maat/
├── cmd/maat/main.go          # Entry point
├── internal/
│   ├── tui/                  # Bubble Tea UI
│   │   ├── model.go          # Central Model
│   │   ├── update.go         # Update handlers
│   │   ├── view.go           # View rendering
│   │   ├── commands.go       # tea.Cmd definitions
│   │   └── messages.go       # Msg types
│   ├── graph/                # Knowledge graph (SQLite)
│   ├── linear/               # Linear API client
│   ├── github/               # GitHub API client
│   ├── git/                  # go-git operations
│   └── claude/               # MCP bridge
├── configs/                  # YAML configuration
└── specs/                    # Specifications
```

## Anti-Requirements (What NOT to build)

- **A-001**: NO autonomous AI actions
- **A-002**: NO feature parity with Linear
- **A-003**: NO writes without confirmation
- **A-004**: NO global mutable state
- **A-005**: NO ambient AI suggestions
- **A-006**: NO competing with GitHub UI
- **A-007**: NO auto-sync without user action
- **A-008**: NO pointer mutations in Update

## Key Documents

| Document | Purpose |
|----------|---------|
| `MAAT-SPEC.md` | Unified specification |
| `docs/CONSTITUTION.md` | 10 Commandments |
| `specs/FUNCTIONAL-REQUIREMENTS.md` | FR-001 to FR-010 |
| `specs/ANTI-REQUIREMENTS.md` | What we refuse |
| `docs/components/` | DELTA FORCE implementation guides |

## Development Commands

```bash
# Build
go build -o maat ./cmd/maat

# Test
go test ./...

# Format
gofmt -w .

# Run
./maat
```

## Commit Convention

All commits should reference the relevant FR or ADR:

```
feat(FR-001): Add knowledge graph display

Implements node rendering per ADR-002 graph navigation pattern.
Uses SQLite backend per ADR-003.
```
