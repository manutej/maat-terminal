# MAAT Pre-RMP Specification

**Document Type**: Pre-Refinement Specification
**Created**: 2026-01-05
**Next Step**: `/rmp @quality:0.85` to generate ADRs and FRs
**Status**: Ready for recursive refinement

---

## 1. Project Vision

**MAAT** (Modular Agentic Architecture for Terminal) is a unified terminal workspace that integrates:
- **Linear** (project/issue management)
- **GitHub** (repositories, PRs, code)
- **Claude Code** (AI assistance)
- **Graph Visualization** (spatial navigation with drill-down)

### The 10x Insight (Thiel)

MAAT is NOT:
- A better Linear CLI
- A better GitHub CLI
- A better Claude Code interface

MAAT IS:
- The **only place** where all three compose into unified navigable context
- A **paradigm shift** from list-based to graph-based developer workflow
- **10x differentiation** through composition, not competition

---

## 2. The 10 Commandments (Constitution)

| # | Commandment | Principle | Anti-Pattern |
|---|-------------|-----------|--------------|
| 1 | Immutable Truth | State transforms, never mutates | Pointer mutations |
| 2 | Single Responsibility | One component, one purpose | God objects |
| 3 | Text Interface | Messages are specification | Callbacks |
| 4 | 10x Differentiation | Graph navigation, not lists | Feature parity |
| 5 | Controlled Effects | Commands describe, runtime executes | Side effects in Update |
| 6 | Human Contact | Explicit invocation, human review | Ambient AI |
| 7 | Composition Monopoly | Value from integration | Competing with sources |
| 8 | Async Purity | Deterministic updates, testable | Goroutines in Update |
| 9 | Specification Constitution | ADRs before implementation | Code-first design |
| 10 | Sovereignty Preservation | Orchestrate, never colonize | Auto-sync writes |

**Full Document**: `/Users/manu/Documents/LUXOR/MAAT/docs/CONSTITUTION.md`

---

## 3. Technical Foundation

### Framework Choice: Go + Bubble Tea

**Rationale**:
- Elm Architecture (TEA) enforces immutability
- Single binary distribution
- Proven by LUMINA, LUMOS, Crush codebases
- Charmbracelet ecosystem (Lipgloss, Bubbles, Glamour)

### Architecture Pattern: Model-View-Update

```go
type Model struct {
    // Single source of truth
    graph       WorkspaceGraph
    viewMode    ViewMode
    focusedNode string
    // Integrations
    linear      LinearState
    github      GitHubState
    claude      ClaudeState
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Pure function - no side effects
    // Returns new model + commands
}

func (m Model) View() string {
    // Pure rendering from state
}
```

### State Machine (UI Modes)

```go
const (
    GraphView ViewMode = iota  // Primary: spatial navigation
    ListView                    // Fallback: when graph overwhelms
    DetailPane                  // Contextual: issue/PR details
    SearchMode                  // Fuzzy finder
    ClaudeMode                  // AI assistance panel
)
```

---

## 4. Integration Philosophy

### Thin Integrations

Each integration is a **thin API client only**:

```
maat/
├── internal/
│   ├── linear/      # Linear API client (read issues, write with confirm)
│   ├── github/      # GitHub API client (repos, PRs, files)
│   ├── claude/      # Claude API client (context-aware assistance)
│   └── graph/       # Graph data structure and rendering
├── tui/
│   ├── model.go     # Central state
│   ├── update.go    # Pure update functions
│   ├── view.go      # Rendering
│   └── commands.go  # Async effect descriptions
└── main.go          # Composition root
```

### Read vs Write Sovereignty

| Operation | Confirmation | Rationale |
|-----------|--------------|-----------|
| Read issues | None | No side effects |
| Read PRs | None | No side effects |
| Update issue status | Required | Writes to Linear |
| Create comment | Required | Writes to Linear/GitHub |
| Claude suggestion | Show, don't execute | Human reviews first |

---

## 5. Graph-First Navigation

### Core Data Structure

```go
type WorkspaceGraph struct {
    Nodes map[string]Node
    Edges []Edge
}

type Node struct {
    ID     string
    Type   NodeType  // Issue | PR | File | Commit
    Source string    // linear | github | local
    Data   any       // Source-specific payload
}

type Edge struct {
    From     string
    To       string
    Relation EdgeType  // blocks | related | links | mentions
}
```

### Navigation Model

```
Initiative
    └── Project
        └── Cycle
            └── Issue ←──────────┐
                ├── PR ──────────┤ (edges)
                │   └── File     │
                └── Comment      │
                    └── Claude ──┘
```

- **Drill Down**: Enter on node → zoom into children
- **Breadcrumb**: Always visible path to current location
- **Context Panel**: Right pane shows selected node details

---

## 6. RMP Input: What to Generate

### ADRs Needed (Top 5)

| ADR | Title | Decision |
|-----|-------|----------|
| ADR-001 | Elm Architecture | Adopt TEA with pure Update |
| ADR-002 | Graph-First Navigation | Graph as primary UI paradigm |
| ADR-003 | Explicit Effect Boundaries | tea.Cmd for all async |
| ADR-004 | Human-in-Loop AI | Claude with confirmation gates |
| ADR-005 | Thin Integrations | API clients only, no business logic |

### Functional Requirements (Max 10)

| FR | Requirement | Priority |
|----|-------------|----------|
| FR-001 | View Linear issues in graph | P0 |
| FR-002 | Navigate with keyboard (hjkl, Enter, Esc) | P0 |
| FR-003 | Drill down into issue → PR → file | P0 |
| FR-004 | View GitHub PRs linked to issues | P1 |
| FR-005 | Invoke Claude with Ctrl+A | P1 |
| FR-006 | Update issue status (with confirm) | P1 |
| FR-007 | Fuzzy search across workspace | P2 |
| FR-008 | Breadcrumb navigation | P2 |
| FR-009 | Theme support (dark/light) | P3 |
| FR-010 | Configurable keybindings | P3 |

### Anti-Requirements (What We Refuse)

| Anti-FR | Refusal | Rationale |
|---------|---------|-----------|
| A-001 | NO autonomous AI actions | Human contact protocol |
| A-002 | NO feature parity with Linear | Composition monopoly |
| A-003 | NO write without confirmation | Sovereignty preservation |
| A-004 | NO global mutable state | Immutable truth |
| A-005 | NO ambient AI suggestions | Explicit invocation |
| A-006 | NO competing with GitHub UI | Thin integrations |
| A-007 | NO sync without user action | Orchestrate, don't colonize |
| A-008 | NO pointer mutations in Update | Async purity |

---

## 7. Quality Gates for /rmp

### Dimensions

| Dimension | Weight | Criteria |
|-----------|--------|----------|
| Correctness | 40% | Follows 10 Commandments |
| Clarity | 30% | Unambiguous requirements |
| Completeness | 20% | Covers core user journey |
| Efficiency | 10% | Minimal viable scope |

### Target Quality: 0.85

```bash
/rmp @quality:0.85 @max_iterations:5 "Generate MAAT ADRs and FRs"
```

### Convergence Criteria

1. All ADRs reference specific Commandments
2. All FRs have clear acceptance criteria
3. All Anti-FRs explain which Commandment they enforce
4. No redundancy between requirements
5. Core user journey is complete (view → navigate → drill → act)

---

## 8. Reference Materials

### Research Documents
- `/Users/manu/Documents/LUXOR/MAAT/research/01-12-FACTOR-AGENTS.md`
- `/Users/manu/Documents/LUXOR/MAAT/research/02-SPEC-KIT.md`
- `/Users/manu/Documents/LUXOR/MAAT/research/03-BUBBLE-TEA.md`
- `/Users/manu/Documents/LUXOR/MAAT/research/04-LUMINA-PATTERNS.md`

### Constitution
- `/Users/manu/Documents/LUXOR/MAAT/docs/CONSTITUTION.md`

### Reference Implementations
- `/Users/manu/Documents/LUXOR/PROJECTS/LUMINA/ccn/` (Go TUI, 6.5K lines)
- `/Users/manu/Documents/LUXOR/PROJECTS/LUMOS/` (PDF TUI)
- `/Users/manu/Documents/LUXOR/crush/` (Production Charm TUI)

---

## 9. Post-RMP Deliverables

After `/rmp` completes, generate:

1. **`/Users/manu/Documents/LUXOR/MAAT/specs/ADR/`**
   - ADR-001-elm-architecture.md
   - ADR-002-graph-navigation.md
   - ADR-003-explicit-effects.md
   - ADR-004-human-in-loop.md
   - ADR-005-thin-integrations.md

2. **`/Users/manu/Documents/LUXOR/MAAT/specs/FUNCTIONAL-REQUIREMENTS.md`**
   - FR-001 through FR-010 with acceptance criteria

3. **`/Users/manu/Documents/LUXOR/MAAT/specs/ANTI-REQUIREMENTS.md`**
   - A-001 through A-008 with Commandment references

4. **`/Users/manu/Documents/LUXOR/MAAT/MAAT-SPEC.md`**
   - Final unified specification document

---

## 10. Summary

**MAAT = Unified Context Workspace**

- Linear + GitHub + Claude → Single navigable graph
- 10 Commandments → Constitutional governance
- Elm Architecture → Immutable, testable state
- Graph-first → 10x differentiation from list UIs
- Human-in-loop → AI assists, never acts autonomously

**Ready for**: `/rmp @quality:0.85`
