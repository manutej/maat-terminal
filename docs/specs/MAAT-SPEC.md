# MAAT: Modular Agentic Architecture for Terminal

**Unified Specification Document**
**Version**: 1.0.0
**Date**: 2026-01-05
**Status**: Approved (RMP Quality: 0.85)

---

## Executive Summary

MAAT is a unified terminal workspace that integrates Linear issues, GitHub PRs, and Claude Code into a navigable knowledge graph. Built with Go using the Bubble Tea framework, MAAT provides keyboard-first navigation, human-in-loop AI assistance, and a plugin architecture for extensibility.

### Vision Statement

> "One graph to rule them all — where issues become nodes, relationships become edges, and the developer becomes the navigator."

### Key Differentiators

| Traditional Tools | MAAT |
|-------------------|------|
| Isolated platforms | Unified knowledge graph |
| Mouse-driven UI | Keyboard-first navigation |
| AI as magic | AI as explicit tool |
| Feature competition | Composition over parity |
| Mutable state chaos | Immutable Elm Architecture |

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                         MAAT TUI                                │
│  ┌──────────────┬──────────────────────┬──────────────────────┐│
│  │ GRAPH PANE   │     MAIN PANE        │   DETAIL PANE       ││
│  │              │                      │                      ││
│  │ • Node tree  │ • Issue details      │ • Metadata          ││
│  │ • Edges      │ • PR content         │ • Git history       ││
│  │ • Hierarchy  │ • Code diff          │ • AI suggestions    ││
│  │              │                      │                      ││
│  └──────────────┴──────────────────────┴──────────────────────┘│
│  ┌────────────────────────────────────────────────────────────┐│
│  │ STATUS BAR: [Mode] [Role] [Focus] [Sync] [Help: ?]        ││
│  └────────────────────────────────────────────────────────────┘│
└─────────────────────────────────────────────────────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        ▼                     ▼                     ▼
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│   LINEAR     │     │   GITHUB     │     │   CLAUDE     │
│   CLIENT     │     │   CLIENT     │     │   BRIDGE     │
│              │     │              │     │              │
│ • GraphQL    │     │ • REST API   │     │ • MCP        │
│ • Issues     │     │ • PRs        │     │ • Human-gate │
│ • Projects   │     │ • Commits    │     │ • Audit log  │
└──────────────┘     └──────────────┘     └──────────────┘
        │                     │                     │
        └─────────────────────┼─────────────────────┘
                              ▼
                    ┌──────────────────┐
                    │ KNOWLEDGE GRAPH  │
                    │    (SQLite)      │
                    │                  │
                    │ • Nodes table    │
                    │ • Edges table    │
                    │ • Graph views    │
                    └──────────────────┘
```

---

## 10 Design Commandments

These are the constitutional principles that govern all MAAT decisions:

| # | Commandment | Summary |
|---|-------------|---------|
| 1 | **Immutable Truth** | All state in Model, Update is pure |
| 2 | **Graph Supremacy** | Every entity is a node or edge |
| 3 | **Text Interface** | Msg types are the API |
| 4 | **Navigation Monopoly** | Enter drills down, Esc backs up |
| 5 | **Controlled Effects** | tea.Cmd only, no goroutines |
| 6 | **Human Contact** | AI requires explicit invocation |
| 7 | **Composition Monopoly** | Value from weaving, not competing |
| 8 | **Async Purity** | Commands describe, runtime executes |
| 9 | **Terminal Citizenship** | Dark theme, Ctrl+C sacred |
| 10 | **Sovereignty Preservation** | External writes need confirmation |

**Priority**: Tier 1 (#1, #6, #10) > Tier 2 (#2, #4, #7) > Tier 3 (#3, #5, #8, #9)

---

## Functional Requirements Summary

| FR | Title | Priority | Commandment |
|----|-------|----------|-------------|
| FR-001 | Knowledge Graph Display | P1 | #2 Graph Supremacy |
| FR-002 | Keyboard Navigation | P1 | #4 Navigation Monopoly |
| FR-003 | Detail Pane | P1 | #2 Graph Supremacy |
| FR-004 | State Management | P1 | #1 Immutable Truth |
| FR-005 | Claude Integration | P2 | #6 Human Contact |
| FR-006 | Linear Integration | P2 | #7 Composition Monopoly |
| FR-007 | Role-Based Views | P2 | #6 Human Contact |
| FR-008 | Fuzzy Search | P2 | #4 Navigation Monopoly |
| FR-009 | IDP Self-Service | P3 | #7 Composition Monopoly |
| FR-010 | Git History Graph | P2 | #2 Graph Supremacy |

**Full Details**: See `specs/FUNCTIONAL-REQUIREMENTS.md`

---

## Anti-Requirements Summary

| AFR | Prohibition | Severity |
|-----|-------------|----------|
| A-001 | NO Autonomous AI Actions | Critical |
| A-002 | NO Feature Parity with Linear | High |
| A-003 | NO Write Without Confirmation | Critical |
| A-004 | NO Global Mutable State | Critical |
| A-005 | NO Ambient AI Suggestions | High |
| A-006 | NO Competing with GitHub UI | High |
| A-007 | NO Auto-Sync Without User Action | High |
| A-008 | NO Pointer Mutations in Update | Critical |

**Full Details**: See `specs/ANTI-REQUIREMENTS.md`

---

## Architecture Decision Records

| ADR | Decision | Status |
|-----|----------|--------|
| ADR-001 | Adopt Elm Architecture (TEA) | Accepted |
| ADR-002 | Graph as Primary Navigation | Accepted |
| ADR-003 | SQLite Knowledge Graph Backend | Accepted |
| ADR-004 | Human-in-Loop AI Pattern | Accepted |
| ADR-005 | Thin Integration Clients | Accepted |
| ADR-006 | IDP Self-Service Layer | Accepted |
| ADR-007 | Plugin Architecture | Accepted |

**Full Details**: See `specs/ADR/`

---

## Technology Stack

### Core Framework

| Layer | Technology | Purpose |
|-------|------------|---------|
| TUI Framework | Bubble Tea v0.25+ | Elm Architecture for terminals |
| Styling | Lipgloss v0.9+ | CSS-like terminal styling |
| Widgets | Bubbles v0.18+ | Pre-built components |
| Git | go-git v5.11+ | Pure Go git operations |
| Database | SQLite 3.40+ | Knowledge graph storage |

### Integration Clients

| Service | Library | Pattern |
|---------|---------|---------|
| Linear | `github.com/Khan/genqlient` | GraphQL code generation |
| GitHub | `github.com/google/go-github/v57` | REST client |
| Claude | Custom MCP Bridge | MCP protocol |

### Build & Distribution

| Aspect | Choice |
|--------|--------|
| Language | Go 1.21+ |
| Build | `go build` → single binary |
| Config | YAML + env vars |
| Distribution | Homebrew, direct download |

---

## File Structure

```
maat/
├── cmd/
│   └── maat/
│       └── main.go           # Entry point
├── internal/
│   ├── tui/
│   │   ├── model.go          # Central Model struct
│   │   ├── update.go         # Update handlers
│   │   ├── view.go           # View rendering
│   │   ├── commands.go       # tea.Cmd definitions
│   │   ├── messages.go       # Msg type definitions
│   │   ├── keys.go           # Key handling
│   │   ├── state.go          # WithX helpers
│   │   ├── components/       # Bubbles wrappers
│   │   │   ├── issue_list.go
│   │   │   ├── commit_table.go
│   │   │   ├── detail_viewport.go
│   │   │   └── search_input.go
│   │   ├── styles/           # Lipgloss styles
│   │   │   ├── colors.go
│   │   │   ├── components.go
│   │   │   ├── layout.go
│   │   │   └── git.go
│   │   └── views/            # View renderers
│   │       ├── graph.go
│   │       ├── detail.go
│   │       └── confirm.go
│   ├── graph/
│   │   ├── schema.go         # Node/Edge types
│   │   ├── store.go          # SQLite operations
│   │   └── query.go          # Graph queries
│   ├── linear/
│   │   ├── client.go         # Thin API client
│   │   ├── types.go          # GraphQL types
│   │   └── adapter.go        # Graph node conversion
│   ├── github/
│   │   ├── client.go         # Thin API client
│   │   ├── types.go          # REST types
│   │   └── adapter.go        # Graph node conversion
│   ├── git/
│   │   ├── repository.go     # Repo access
│   │   ├── commits.go        # Commit fetching
│   │   ├── diff.go           # Diff extraction
│   │   ├── refs.go           # Branch/tag resolution
│   │   ├── links.go          # Issue link parsing
│   │   └── adapter.go        # Graph node conversion
│   ├── claude/
│   │   ├── bridge.go         # MCP protocol
│   │   ├── context.go        # Context assembly
│   │   └── audit.go          # Audit logging
│   └── plugin/
│       ├── interface.go      # DataSourcePlugin
│       ├── registry.go       # Plugin management
│       └── builtin/          # Built-in plugins
│           ├── linear.go
│           ├── github.go
│           └── git.go
├── configs/
│   └── default.yaml          # Default configuration
├── docs/
│   ├── CONSTITUTION.md       # 10 Commandments
│   └── components/           # DELTA FORCE guides
│       ├── BUBBLE-TEA-PATTERNS.md
│       ├── LIPGLOSS-STYLING.md
│       ├── BUBBLES-WIDGETS.md
│       └── GO-GIT-INTEGRATION.md
├── specs/
│   ├── FUNCTIONAL-REQUIREMENTS.md
│   ├── ANTI-REQUIREMENTS.md
│   ├── META-SPEC.md
│   └── ADR/
│       ├── ADR-001-elm-architecture.md
│       ├── ADR-002-graph-navigation.md
│       ├── ADR-003-knowledge-graph-backend.md
│       ├── ADR-004-human-in-loop-ai.md
│       ├── ADR-005-thin-integrations.md
│       ├── ADR-006-idp-self-service.md
│       └── ADR-007-plugin-architecture.md
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

---

## Key Code Patterns

### Model-View-Update (Elm Architecture)

```go
type Model struct {
    graph       *graph.KnowledgeGraph
    focusedNode string
    navStack    NavigationStack
    viewMode    ViewMode
}

// Update is PURE - no side effects
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case IssuesFetched:
        return m.WithIssues(msg.Issues), nil
    }
    return m, nil
}

// WithX pattern for state transformation
func (m Model) WithIssues(issues []Issue) Model {
    m.linear.Issues = issues
    return m  // Return copy
}
```

### tea.Cmd for Async Operations

```go
func (m Model) fetchLinearIssues() tea.Cmd {
    return func() tea.Msg {
        issues, err := m.linear.Client.FetchIssues(context.Background())
        if err != nil {
            return ErrorOccurred{Source: "linear", Error: err}
        }
        return IssuesFetched{Issues: issues}
    }
}

// Parallel execution via tea.Batch
func (m Model) refreshAll() tea.Cmd {
    return tea.Batch(
        m.fetchLinearIssues(),
        m.fetchGitHubPRs(),
        m.fetchCommits(),
    )
}
```

### Human-in-Loop AI Pattern

```go
// Claude invocation requires explicit user action
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    switch msg.String() {
    case "ctrl+a":  // Explicit invocation
        return m.WithViewMode(ClaudeMode), nil
    }
    return m, nil
}

// AI actions require confirmation
type ClaudeResponse struct {
    Suggestions []Suggestion
}

func (m Model) HandleClaudeResponse(resp ClaudeResponse) (Model, tea.Cmd) {
    // Present for human review - NEVER auto-apply
    m.pendingSuggestions = resp.Suggestions
    return m, nil
}
```

### Knowledge Graph Schema

```sql
CREATE TABLE nodes (
    id TEXT PRIMARY KEY,
    type TEXT NOT NULL,
    source TEXT NOT NULL,
    data JSON NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE edges (
    id TEXT PRIMARY KEY,
    from_id TEXT NOT NULL REFERENCES nodes(id),
    to_id TEXT NOT NULL REFERENCES nodes(id),
    relation TEXT NOT NULL,
    metadata JSON
);

-- Graph traversal view
CREATE VIEW graph_neighbors AS
SELECT
    n.id AS node_id,
    e.relation,
    n2.id AS neighbor_id,
    n2.type AS neighbor_type
FROM nodes n
JOIN edges e ON n.id = e.from_id OR n.id = e.to_id
JOIN nodes n2 ON (e.from_id = n2.id OR e.to_id = n2.id) AND n2.id != n.id;
```

### Plugin Interface

```go
type DataSourcePlugin interface {
    Name() string
    NodeTypes() []graph.NodeType
    EdgeTypes() []graph.EdgeType
    Fetch(ctx context.Context, query Query) ([]graph.Node, []graph.Edge, error)
    Transform(raw interface{}) (graph.Node, error)
}

// Registration
func (r *Registry) Register(plugin DataSourcePlugin) error {
    if _, exists := r.plugins[plugin.Name()]; exists {
        return fmt.Errorf("plugin %s already registered", plugin.Name())
    }
    r.plugins[plugin.Name()] = plugin
    return nil
}
```

---

## Implementation Roadmap

### Phase 1: Foundation (Weeks 1-2)

- [ ] Set up Go project structure
- [ ] Implement basic Bubble Tea model
- [ ] Create SQLite knowledge graph schema
- [ ] Build basic 3-pane layout

### Phase 2: Graph Navigation (Weeks 3-4)

- [ ] Implement node/edge rendering
- [ ] Add keyboard navigation (h/j/k/l)
- [ ] Build drill-down/back-up navigation
- [ ] Create detail pane viewport

### Phase 3: Integrations (Weeks 5-7)

- [ ] Linear GraphQL client
- [ ] GitHub REST client
- [ ] Git commit history (go-git)
- [ ] Graph adapter layer

### Phase 4: AI & Polish (Weeks 8-10)

- [ ] Claude MCP bridge
- [ ] Human confirmation dialogs
- [ ] Role-based views
- [ ] Fuzzy search

### Phase 5: Extensibility (Weeks 11-12)

- [ ] Plugin interface
- [ ] IDP self-service commands
- [ ] Documentation
- [ ] Beta release

---

## Success Metrics

| Metric | Target | Measurement |
|--------|--------|-------------|
| Startup Time | < 500ms | Time to interactive |
| Navigation Latency | < 50ms | Key press to render |
| Sync Refresh | < 2s | Full graph update |
| Memory Usage | < 100MB | Steady state |
| Binary Size | < 20MB | Single executable |
| Test Coverage | > 80% | Unit + integration |

---

## Document Index

| Document | Purpose | Location |
|----------|---------|----------|
| **MAAT-SPEC** | This unified spec | `MAAT-SPEC.md` |
| **Constitution** | 10 Commandments | `docs/CONSTITUTION.md` |
| **META-SPEC** | Philosophy & process | `specs/META-SPEC.md` |
| **Functional Requirements** | FR-001 to FR-010 | `specs/FUNCTIONAL-REQUIREMENTS.md` |
| **Anti-Requirements** | A-001 to A-008 | `specs/ANTI-REQUIREMENTS.md` |
| **ADRs** | Architecture decisions | `specs/ADR/` |
| **Bubble Tea Patterns** | TUI framework guide | `docs/components/BUBBLE-TEA-PATTERNS.md` |
| **Lipgloss Styling** | Styling reference | `docs/components/LIPGLOSS-STYLING.md` |
| **Bubbles Widgets** | Widget integration | `docs/components/BUBBLES-WIDGETS.md` |
| **Go-Git Integration** | Git implementation | `docs/components/GO-GIT-INTEGRATION.md` |

---

## Appendix: Research Sources

### Context7 Documentation Queries

| Library | Query | Context |
|---------|-------|---------|
| Bubble Tea | Model-View-Update, tea.Cmd | TUI architecture |
| Lipgloss | Styling, colors, layout | Visual design |
| Bubbles | viewport, list, table, textinput | Component selection |

### Methodological Influences

| Source | Contribution |
|--------|--------------|
| 12-Factor Agents | Declarative config, stateless processes |
| GitHub Spec-Kit | Specification-driven development |
| LUMINA Project | Markdown TUI patterns |
| Elm Architecture | Pure functional UI |

---

## Revision History

| Version | Date | Changes |
|---------|------|---------|
| 1.0.0 | 2026-01-05 | Initial unified specification |

---

**Specification Complete** ✓

This document represents the converged MAAT specification at RMP quality 0.85. All functional requirements, anti-requirements, and architecture decisions have been validated against the 10 Design Commandments.

**Next Steps**: Begin Phase 1 implementation per roadmap.
