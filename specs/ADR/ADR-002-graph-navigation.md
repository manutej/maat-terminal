# ADR-002: Graph-First Navigation Paradigm

**Status**: Accepted
**Date**: 2026-01-05
**Deciders**: MAAT Architecture Team
**Commandments**: #4 10x Differentiation, #7 Composition Monopoly

---

## Context

Traditional developer tools present work as:
- **Lists**: Linear issues, GitHub PRs, file trees
- **Separate contexts**: Switch between apps to see relationships
- **Lost context**: Navigate away and lose mental model

The 10x opportunity is **spatial, graph-based navigation** where:
- Relationships are visible (blocks, implements, calls)
- Drill-down preserves context (breadcrumb)
- Composition of Linear + GitHub + Claude creates unique value

## Decision

Adopt **graph as the primary UI paradigm**:

```go
type WorkspaceGraph struct {
    Nodes map[string]Node
    Edges []Edge
}

type Node struct {
    ID       string
    Type     NodeType  // Issue | PR | File | Commit | Service | Document
    Source   string    // linear | github | local | plugin:<name>
    Data     any
    Position Position  // x, y coordinates for layout
}

type Edge struct {
    From     string
    To       string
    Relation EdgeType  // blocks | related | implements | calls | owns
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

### Interaction Patterns

| Action | Key | Behavior |
|--------|-----|----------|
| Navigate | hjkl/arrows | Move focus between nodes |
| Drill Down | Enter | Zoom into node's children |
| Back Up | Esc/Backspace | Return to parent context |
| Expand | Tab | Toggle node expansion |
| Search | / | Fuzzy find across graph |

### View Modes

```go
const (
    GraphView ViewMode = iota  // Primary: spatial navigation
    ListView                    // Fallback: when graph overwhelms
    DetailPane                  // Contextual: right pane for focused node
    SearchMode                  // Overlay: fuzzy finder
    ClaudeMode                  // AI: Claude interaction panel
)
```

### Layout Algorithm

Force-directed layout with constraints:
- Parent-child vertical hierarchy
- Sibling horizontal spread
- Edge bundling for readability
- Fit-to-viewport with zoom

```go
func (g *WorkspaceGraph) Layout(width, height int) {
    // 1. Assign layers (topological sort)
    layers := g.topologicalLayers()

    // 2. Position nodes within layers
    for i, layer := range layers {
        y := i * (height / len(layers))
        for j, node := range layer {
            x := j * (width / len(layer))
            g.Nodes[node.ID].Position = Position{X: x, Y: y}
        }
    }

    // 3. Minimize edge crossings (barycenter method)
    g.minimizeCrossings()
}
```

## Consequences

### Positive
- **10x Differentiation**: No other tool shows work as navigable graph
- **Context Preservation**: Breadcrumb + zoom maintains mental model
- **Relationship Discovery**: See blocking issues, linked PRs at a glance
- **Composition Value**: Linear + GitHub unified in single spatial view

### Negative
- **Complexity**: Graph rendering is harder than lists
- **Screen Real Estate**: Graphs need space; small terminals challenging
- **Learning Curve**: Users expect lists; graph is novel

### Mitigations
- ListView as fallback for narrow terminals
- Progressive disclosure (collapsed by default)
- Keyboard shortcuts familiar from vim/tmux

## Compliance

This ADR enforces:
- **Commandment #4**: Graph navigation, not lists (10x differentiation)
- **Commandment #7**: Value from integration, not features (composition monopoly)

## References

- LUMINA Patterns: `/Users/manu/Documents/LUXOR/MAAT/research/04-LUMINA-PATTERNS.md`
- MAAT Constitution: `/Users/manu/Documents/LUXOR/MAAT/docs/CONSTITUTION.md`
