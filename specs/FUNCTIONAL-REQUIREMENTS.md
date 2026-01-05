# MAAT Functional Requirements

**Document Type**: Functional Requirements Specification
**Version**: 1.0
**Date**: 2026-01-05
**Status**: Approved (RMP Quality: 0.85)

---

## Overview

This document defines the 10 core functional requirements for MAAT (Modular Agentic Architecture for Terminal). Each requirement includes:
- Priority (P0 = MVP, P1 = V1.0, P2 = V1.1, P3 = Future)
- Acceptance Criteria (testable conditions)
- Commandment Reference (constitutional alignment)
- Implementation Notes

---

## FR-001: View Linear Issues in Graph

**Priority**: P0 (MVP)
**Commandment**: #4 10x Differentiation, #7 Composition Monopoly

### Description
Display Linear issues as nodes in a navigable graph, with edges representing relationships (blocks, related, parent/child).

### Acceptance Criteria
- [ ] Issues from configured Linear team appear as nodes
- [ ] Blocking relationships shown as directed edges
- [ ] Related issues connected with bidirectional edges
- [ ] Parent-child (sub-issues) shown hierarchically
- [ ] Node color/style reflects issue status (todo, in-progress, done)
- [ ] Graph renders in < 100ms for up to 500 nodes

### Implementation Notes
```go
type IssueNode struct {
    LinearID   string
    Title      string
    Status     IssueStatus
    Priority   int
    Assignee   string
    BlockedBy  []string  // Edge: blocks
    Related    []string  // Edge: related
    ParentID   string    // Edge: parent
}
```

---

## FR-002: Keyboard Navigation (hjkl, Enter, Esc)

**Priority**: P0 (MVP)
**Commandment**: #4 10x Differentiation

### Description
Full keyboard-driven navigation without mouse dependency, using vim-style keybindings.

### Acceptance Criteria
- [ ] `h/j/k/l` or arrows navigate between nodes
- [ ] `Enter` drills into focused node (zoom)
- [ ] `Esc` or `Backspace` returns to parent context
- [ ] `Tab` cycles between panes (graph, detail, search)
- [ ] `/` opens search/fuzzy finder
- [ ] `?` shows help overlay
- [ ] `q` quits application
- [ ] All keybindings configurable via `~/.config/maat/keys.json`

### Implementation Notes
```go
type KeyBinding struct {
    Key     string `json:"key"`
    Action  string `json:"action"`
    Context string `json:"context"` // graph | detail | search | global
}
```

---

## FR-003: Drill Down: Issue → PR → File

**Priority**: P0 (MVP)
**Commandment**: #4 10x Differentiation, #7 Composition Monopoly

### Description
Navigate from a Linear issue to its linked GitHub PRs, and from PRs to modified files, preserving context via breadcrumb.

### Acceptance Criteria
- [ ] Enter on Issue shows linked PRs as child nodes
- [ ] Enter on PR shows modified files as child nodes
- [ ] Enter on File opens file preview in detail pane
- [ ] Breadcrumb always visible: `Project > Issue > PR > File`
- [ ] Breadcrumb items clickable for quick navigation
- [ ] Context panel shows relevant metadata at each level

### Implementation Notes
```go
type NavigationStack struct {
    Path      []NodeRef
    Current   int
    MaxDepth  int  // Default: 10
}

func (s *NavigationStack) Push(node NodeRef) {
    s.Path = append(s.Path[:s.Current+1], node)
    s.Current++
}

func (s *NavigationStack) Pop() NodeRef {
    if s.Current > 0 {
        s.Current--
    }
    return s.Path[s.Current]
}
```

---

## FR-004: Knowledge Graph Persistence

**Priority**: P0 (MVP)
**Commandment**: #1 Immutable Truth, #3 Text Interface

### Description
Persist fetched data in local SQLite database for offline access, fast queries, and relationship tracking.

### Acceptance Criteria
- [ ] All nodes stored in SQLite `nodes` table
- [ ] All edges stored in SQLite `edges` table
- [ ] Graph survives application restart
- [ ] Incremental sync (only fetch changes since last sync)
- [ ] Database location: `~/.local/share/maat/graph.db`
- [ ] Schema migrations handled automatically

### Implementation Notes
See ADR-003 for full schema.

---

## FR-005: Invoke Claude with Ctrl+A

**Priority**: P1 (V1.0)
**Commandment**: #6 Human Contact, #10 Sovereignty Preservation

### Description
Explicit AI invocation that assembles context from current focus and presents suggestions for human review.

### Acceptance Criteria
- [ ] `Ctrl+A` opens Claude panel
- [ ] Context auto-assembled from focused node + linked content
- [ ] User can type question before submission
- [ ] Response shown in Claude pane
- [ ] Suggestions listed with approve/reject options
- [ ] All interactions logged to audit trail
- [ ] NO ambient or proactive AI suggestions

### Implementation Notes
See ADR-004 for full specification.

---

## FR-006: Update Issue Status (with Confirm)

**Priority**: P1 (V1.0)
**Commandment**: #10 Sovereignty Preservation

### Description
Allow users to update Linear issue status from within MAAT, with mandatory confirmation before write.

### Acceptance Criteria
- [ ] `s` on focused issue opens status picker
- [ ] Status options fetched from Linear workflow
- [ ] Confirmation dialog shows: action, target, new status
- [ ] `y` confirms, `n` cancels
- [ ] Successful update reflects immediately in graph
- [ ] Action logged to audit trail

### Implementation Notes
```go
type ConfirmDialog struct {
    Action      string
    Target      string
    Description string
    OnConfirm   tea.Cmd
    OnCancel    tea.Cmd
}
```

---

## FR-007: Role-Based Views (exec/lead/ic)

**Priority**: P1 (V1.0)
**Commandment**: #2 Single Responsibility

### Description
Filter visible nodes based on user role, providing appropriate abstraction levels.

### Acceptance Criteria
- [ ] Role configurable in `~/.config/maat/config.yaml`
- [ ] `exec` sees: Initiatives, Projects, Risks, Milestones
- [ ] `lead` sees: Projects, Cycles, Services, PRs, Issues (P0-P1)
- [ ] `ic` sees: All nodes (full access)
- [ ] Role switch via `:role <name>` command
- [ ] Node visibility filter applied after fetch

### Implementation Notes
See ADR-003 `NodeMetadata.AccessLevel` field.

---

## FR-008: Fuzzy Search Across Workspace

**Priority**: P2 (V1.1)
**Commandment**: #4 10x Differentiation

### Description
Fast fuzzy search across all nodes (issues, PRs, files, commits) with instant results.

### Acceptance Criteria
- [ ] `/` opens search overlay
- [ ] Fuzzy matching on title, ID, and description
- [ ] Results ranked by relevance
- [ ] Results filterable by type (issue, PR, file)
- [ ] `Enter` navigates to selected result
- [ ] Search < 50ms for 10,000 nodes

### Implementation Notes
```go
// Use Bubble Tea's built-in textinput + custom fuzzy matcher
type SearchModel struct {
    Input     textinput.Model
    Results   []SearchResult
    Selected  int
    Filter    NodeType  // nil = all types
}
```

---

## FR-009: AI Action Audit Trail

**Priority**: P2 (V1.1)
**Commandment**: #6 Human Contact, #10 Sovereignty Preservation

### Description
Log all AI interactions with full context, response, and user decision.

### Acceptance Criteria
- [ ] All Claude invocations logged to `ai_audit` table
- [ ] Logged data: session_id, timestamp, context, response, user_action
- [ ] For edits: diff included in log
- [ ] Audit viewable via `:audit` command
- [ ] Export to JSON via `:audit --export`
- [ ] Retention configurable (default: 90 days)

### Implementation Notes
See ADR-004 for schema and logging implementation.

---

## FR-010: Git Commit History Graph with Code Snippets

**Priority**: P1 (V1.0)
**Commandment**: #4 10x Differentiation, #7 Composition Monopoly

### Description
Display git commit history as navigable graph with commit message previews and code diff snippets, enabling quick context about code evolution.

### Acceptance Criteria
- [ ] Commit history displayed as linear graph (branching optional V1.1)
- [ ] Each commit node shows: hash (short), message (first line), author, date
- [ ] Enter on commit shows full message + diff stats
- [ ] `d` on commit shows code diff in detail pane
- [ ] Branch/merge visualization with connecting edges
- [ ] Commits linked to PRs/issues via commit message parsing
- [ ] Window into last N commits (default: 50, configurable)
- [ ] Commits filterable by path, author, date range

### Implementation Notes
```go
type CommitNode struct {
    Hash      string    // Short hash (7 chars)
    FullHash  string    // Full SHA
    Message   string    // First line
    FullMsg   string    // Complete message
    Author    string
    Date      time.Time
    Parents   []string  // Parent commit hashes
    Branch    string    // Branch name if head
    Tags      []string  // Associated tags
    Stats     DiffStats // Files changed, insertions, deletions
}

type DiffStats struct {
    FilesChanged int
    Insertions   int
    Deletions    int
    Files        []FileChange
}

type FileChange struct {
    Path       string
    Status     string  // A (added), M (modified), D (deleted)
    Insertions int
    Deletions  int
    Snippet    string  // First 5 lines of diff for preview
}

// Git integration via go-git
func FetchCommitGraph(repo *git.Repository, limit int) ([]CommitNode, error) {
    iter, _ := repo.Log(&git.LogOptions{All: true})
    var commits []CommitNode

    iter.ForEach(func(c *object.Commit) error {
        if len(commits) >= limit {
            return storer.ErrStop
        }

        node := CommitNode{
            Hash:    c.Hash.String()[:7],
            FullHash: c.Hash.String(),
            Message: strings.Split(c.Message, "\n")[0],
            FullMsg: c.Message,
            Author:  c.Author.Name,
            Date:    c.Author.When,
        }

        for _, parent := range c.ParentHashes {
            node.Parents = append(node.Parents, parent.String()[:7])
        }

        // Parse PR/issue links from commit message
        node.LinkedIssues = parseIssueRefs(c.Message)

        commits = append(commits, node)
        return nil
    })

    return commits, nil
}

// Commit message parsing for issue links
func parseIssueRefs(msg string) []string {
    // Matches: LIN-123, #456, fixes #789, closes LINEAR-101
    re := regexp.MustCompile(`(?i)(LIN-\d+|#\d+|closes?\s+#?\w+-?\d+|fixes?\s+#?\w+-?\d+)`)
    return re.FindAllString(msg, -1)
}
```

### UI: Commit History View
```
┌─────────────────────────────────────────────────┐
│  GIT HISTORY: feature/context-pack (50 commits) │
├─────────────────────────────────────────────────┤
│                                                 │
│  ● a3f2c9e  Add compression algorithm    (HEAD)│
│  │          Jan 5, 2026 • You                  │
│  │          +142 -23 • 3 files                 │
│  │                                              │
│  ● 8b1d4f2  Fix edge case in CRDT merge        │
│  │          Jan 4, 2026 • Alice                │
│  │          Fixes LIN-234                      │
│  │                                              │
│  ●─┬ 2e5a9c1  Merge branch 'main'              │
│  │ │         Jan 3, 2026 • You                 │
│  │ │                                            │
│  ○ │ 7f3b2d4  Update dependencies              │
│    │         Jan 3, 2026 • Bot                 │
│    │                                            │
│  ──┘                                            │
│                                                 │
│  [Enter] Details  [d] Diff  [c] Checkout  [/] Search │
└─────────────────────────────────────────────────┘
```

### UI: Commit Detail Pane
```
┌─────────────────────────────────────────────────┐
│  COMMIT: a3f2c9e                                │
├─────────────────────────────────────────────────┤
│  Author: You <you@example.com>                  │
│  Date:   2026-01-05 14:32:01 -0800              │
│  Branch: feature/context-pack                   │
│                                                 │
│  Add compression algorithm                      │
│                                                 │
│  Implements delta-state compression for CRDT    │
│  context packing. Uses the approach described   │
│  in docs/crdts.md.                              │
│                                                 │
│  Refs: LIN-234                                  │
│                                                 │
│  ─────────────────────────────────────────────  │
│  CHANGES (3 files, +142 -23)                    │
│                                                 │
│  M src/compress.rs        +98  -12              │
│  │ + pub fn compress_delta(state: &State) {     │
│  │ +     let baseline = state.baseline();       │
│  │ +     ...                                    │
│  │                                              │
│  M src/lib.rs             +12  -3               │
│  A tests/compress_test.rs +32  -8               │
│                                                 │
│  [d] Full diff  [f] Jump to file  [Esc] Back   │
└─────────────────────────────────────────────────┘
```

---

## Summary Table

| FR | Requirement | Priority | Commandments |
|----|-------------|----------|--------------|
| FR-001 | View Linear issues in graph | P0 | #4, #7 |
| FR-002 | Keyboard navigation (hjkl) | P0 | #4 |
| FR-003 | Drill down: Issue → PR → File | P0 | #4, #7 |
| FR-004 | Knowledge graph persistence | P0 | #1, #3 |
| FR-005 | Invoke Claude with Ctrl+A | P1 | #6, #10 |
| FR-006 | Update issue status (confirm) | P1 | #10 |
| FR-007 | Role-based views | P1 | #2 |
| FR-008 | Fuzzy search | P2 | #4 |
| FR-009 | AI action audit trail | P2 | #6, #10 |
| FR-010 | Git commit history graph | P1 | #4, #7 |

---

## Dependency Graph

```
FR-004 (Graph Persistence)
    ↓
FR-001 (Linear Issues) ───→ FR-003 (Drill Down)
    ↓                           ↓
FR-002 (Navigation) ←──────────┘
    ↓
FR-010 (Git History) ───→ FR-003 (links commits to PRs)
    ↓
FR-005 (Claude) ───→ FR-009 (Audit Trail)
    ↓
FR-006 (Update Status)
    ↓
FR-007 (Role Views) ───→ FR-008 (Search)
```

---

## MVP Scope (P0)

For initial release, implement:
1. **FR-001**: Linear issues as graph
2. **FR-002**: Keyboard navigation
3. **FR-003**: Drill-down navigation
4. **FR-004**: SQLite persistence

Estimated effort: 4-6 weeks

---

## References

- MAAT Constitution: `/Users/manu/Documents/LUXOR/MAAT/docs/CONSTITUTION.md`
- ADRs: `/Users/manu/Documents/LUXOR/MAAT/specs/ADR/`
- PRE-RMP-SPEC: `/Users/manu/Documents/LUXOR/MAAT/specs/PRE-RMP-SPEC.md`
