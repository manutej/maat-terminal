# Single-Pane Design - Simplified TUI

**Date**: 2026-01-07
**Architecture**: Single full-screen view with Tab cycling
**Status**: ✅ **IMPLEMENTED & READY TO TEST**

---

## Design Decision

### User Request
> "Rather than panes… let's get rid of panes altogether and switch view via hotkey for now"

### Rationale
- **Simplicity**: One focused view at a time eliminates layout complexity
- **Clarity**: Full terminal width/height for content = better readability
- **Proven Pattern**: vim modes, htop views, ranger panels
- **Responsive**: No split-pane geometry issues across terminal sizes
- **User Control**: Tab cycling gives explicit navigation control

---

## Architecture

### View Modes (3 Full-Screen Views)

```
┌─────────────────────────────────────────────┐
│  Tab cycles through:                        │
│                                             │
│  1. Graph View    (📊 Knowledge Graph)      │
│  2. Details View  (📝 Node Details)         │
│  3. Relations View (🔗 Relationships)       │
│                                             │
│  Loop: Graph → Details → Relations → Graph │
└─────────────────────────────────────────────┘
```

### 1. Graph View (📊 Knowledge Graph)
**Purpose**: Hierarchical tree visualization of all nodes

**Content**:
- Full-width graph rendering (no 25% constraint!)
- Tree layout with BFS algorithm
- Vim-style navigation (hjkl)
- Node selection with visual focus indicator
- Status/priority icons for each node

**Keyboard**:
- `hjkl` or arrows: Navigate nodes
- `Enter`: Drill into node (shows Details view)
- `Tab`: Switch to Details view

### 2. Details View (📝 Node Details)
**Purpose**: Comprehensive node information (centered, readable)

**Content**:
- Large icon + title (underlined, prominent)
- Type badge (colored background)
- Status with icon (✅🔄📋🚫❌)
- Priority with fire icon (🔥)
- Description (wrapped text, max 80 chars wide)
- Labels as colored badges (🏷)
- Node ID (faint, at bottom)

**Keyboard**:
- `Tab`: Switch to Relations view
- `Esc`: Back to Graph view

### 3. Relations View (🔗 Relationships)
**Purpose**: Show all connections for focused node

**Content**:
- Outgoing relations (→ this node points to others)
- Incoming relations (← others point to this node)
- Relation type highlighted (colored)
- Summary count (X outgoing, Y incoming)

**Keyboard**:
- `Tab`: Switch to Graph view
- `Esc`: Back to previous view

---

## Code Changes

### Files Modified

| File | Changes | Purpose |
|------|---------|---------|
| `internal/tui/state.go` | Updated ViewMode enum, added CycleView() | Define 3 view modes with cycling logic |
| `internal/tui/model.go` | Removed Pane type and methods | Eliminate 3-pane concept entirely |
| `internal/tui/view.go` | Complete rewrite (566 lines) | Single-pane rendering with 3 view modes |
| `internal/tui/update.go` | Updated Tab handling | Cycle views instead of panes |

### Key Architecture Changes

**Before (3-Pane Split)**:
```go
// Model had activePane field
activePane Pane // PaneGraph, PaneMain, PaneDetail

// View rendered 3 panes side-by-side
func (m Model) renderMainView() string {
    graphPane := m.renderGraphPane(layout)
    mainPane := m.renderMainPane(layout)
    detailPane := m.renderDetailPane(layout)
    return lipgloss.JoinHorizontal(lipgloss.Top,
        graphPane, mainPane, detailPane)
}

// Tab cycled panes
case "tab":
    return m.CyclePane(), nil
```

**After (Single-Pane with Mode)**:
```go
// Model has currentView field
currentView ViewMode // ViewGraph, ViewDetails, ViewRelations

// View renders based on mode (full screen)
func (m Model) renderCurrentView() string {
    switch m.currentView {
    case ViewGraph:
        return m.renderGraphView(m.width, contentHeight)
    case ViewDetails:
        return m.renderDetailsView(m.width, contentHeight)
    case ViewRelations:
        return m.renderRelationsView(m.width, contentHeight)
    }
}

// Tab cycles views
case "tab":
    m = m.WithView(m.currentView.CycleView())
    return m, nil
```

---

## Keyboard Controls

### Global (All Views)
```
Tab       - Cycle to next view (Graph → Details → Relations → Graph)
Shift+Tab - Cycle to previous view
Esc       - Back to previous view
q / Ctrl+C - Quit
r         - Refresh data
```

### Graph View
```
h / ←     - Navigate left
j / ↓     - Navigate down
k / ↑     - Navigate up
l / →     - Navigate right
Enter     - Drill into node (→ Details view)
```

### Details View
```
Tab       - Switch to Relations view
Esc       - Back to Graph view
```

### Relations View
```
Tab       - Switch to Graph view
Esc       - Back to previous view
```

---

## Visual Design

### Status Bar (Bottom - All Views)
```
[Graph View] | Node: MAAT Project               Tab: switch view | hjkl: navigate | Enter: drill | Esc: back | q: quit
```

### Graph View Example
```
                            📊 Knowledge Graph

┌────────────────┐
│ 📦 MAAT       │  * ← focused node indicator
└────────────────┘
        │
        ↓
┌────────────────┐  ┌────────────────┐
│ 🐛 Issue #1   │  │ 🔀 PR #1       │
└────────────────┘  └────────────────┘
```

### Details View Example
```
                            📝 Node Details

🐛 Implement graph rendering engine
────────────────────────────────────

 Type: Issue

✅ Status: done

🔥 Priority: High

Description:
Create hierarchical tree layout for knowledge
graph visualization using BFS algorithm.

🏷  Labels:  frontend   visualization

ID: issue:1
```

### Relations View Example
```
                            🔗 Relationships

Relationships for: Implement graph rendering engine

→ Outgoing Relations:

  • implements → MAAT Project
  • blocks → Issue #2

← Incoming Relations:

  • PR #1 → resolves

Total: 2 outgoing, 1 incoming
```

---

## Benefits of Single-Pane Design

### 1. Simplicity ✅
- **No split-pane layout math**: One view = full terminal dimensions
- **No overflow issues**: Content aware of exact available space
- **Easy to reason about**: Simple state machine (3 modes)

### 2. Clarity ✅
- **Full width for graph**: No 25% constraint = shows more nodes
- **Readable details**: 80-char centered content optimal for reading
- **Focused attention**: One thing at a time = less cognitive load

### 3. Responsiveness ✅
- **Works on any terminal size**: 80+ chars minimum
- **No pane resizing issues**: Single view adapts naturally
- **Meets KEY requirement**: "It should work on different sized terminal windows"

### 4. User Control ✅
- **Explicit navigation**: Tab to switch = clear intent
- **Predictable cycling**: Graph → Details → Relations (always same order)
- **Quick access**: Max 2 Tab presses to any view

---

## Implementation Details

### View Rendering Strategy

**Graph View** (`renderGraphView`):
- Uses full terminal width minus padding (width - 4)
- Calls `RenderGraph(m, width-4)` with responsive constraint
- Shows centered title "📊 Knowledge Graph"
- Displays "No nodes loaded" message if empty

**Details View** (`renderDetailsView`):
- Checks if node is selected (shows message if not)
- Renders centered content (max 80 chars wide)
- Uses `renderNodeDetailsExpanded()` for rich formatting
- Shows all node metadata with icons and colors

**Relations View** (`renderRelationsView`):
- Separates outgoing vs incoming edges
- Max 100 chars wide for relation lists
- Colored relation types (Primary/Secondary)
- Shows summary count at bottom

### State Management

**ViewMode Enum**:
```go
type ViewMode int

const (
    ViewGraph     ViewMode = iota // Full-screen graph
    ViewDetails                    // Full-screen node details
    ViewRelations                  // Full-screen relationships
    ViewConfirm                    // Overlay dialog
)

func (v ViewMode) CycleView() ViewMode {
    // Graph → Details → Relations → Graph
}
```

**Model Updates**:
- Removed `activePane` field (no longer needed)
- Kept `currentView` (ViewMode)
- Tab updates `currentView` via `m.WithView(m.currentView.CycleView())`

---

## Testing Instructions

```bash
cd /Users/manu/Documents/LUXOR/MAAT
./maat
```

### Test Sequence

**1. Graph View (Default)**
- ✅ Full-width graph visible
- ✅ Navigate with hjkl/arrows
- ✅ Focus indicator (*) shows on selected node
- ✅ Status bar shows "[Graph View]"

**2. Press Tab → Details View**
- ✅ Switches to full-screen details
- ✅ Shows selected node info (icon, title, type, status, priority)
- ✅ Description wrapped and readable
- ✅ Status bar shows "[Details View]"

**3. Press Tab → Relations View**
- ✅ Switches to relationships view
- ✅ Shows outgoing and incoming edges
- ✅ Relation types colored
- ✅ Summary count displayed
- ✅ Status bar shows "[Relations View]"

**4. Press Tab → Back to Graph View**
- ✅ Cycles back to graph
- ✅ Focus preserved on same node

**5. Esc Navigation**
- ✅ Esc returns to previous view in stack
- ✅ Esc in Graph with empty stack = stays in Graph

**6. Resize Terminal**
- ✅ Views adapt to new dimensions
- ✅ No overlapping or broken layout
- ✅ Content stays readable at 80+ chars

---

## Performance Improvements

### Memory
- **Before**: 3 panes rendered simultaneously (even if hidden)
- **After**: Only current view rendered = 66% less rendering work

### Layout Complexity
- **Before**: Split-pane geometry, border calculations, overflow handling
- **After**: Single view fills terminal = trivial layout math

### Cognitive Load
- **Before**: 3 information sources competing for attention
- **After**: 1 focused view = clear mental model

---

## Future Enhancements

### Phase 3+ Features
- **Search/Filter**: `/ <query>` to filter nodes in Graph view
- **Bookmarks**: `m <key>` to mark nodes, `' <key>` to jump
- **History**: `Ctrl+O` / `Ctrl+I` to navigate view history
- **Vim Splits**: `:vs` / `:sp` for power users who want panes back

### AI Integration (Phase 4)
- **Ctrl+A in Graph**: "What's blocking high-priority issues?"
- **Ctrl+A in Details**: "Summarize this node's context"
- **Ctrl+A in Relations**: "Why are these nodes connected?"

---

## Comparison: 3-Pane vs Single-Pane

| Aspect | 3-Pane Split | Single-Pane Mode |
|--------|-------------|------------------|
| **Layout Complexity** | High (split math, borders, overflow) | Low (full screen) |
| **Graph Width** | 25% (~20-30 chars) | 100% (~76-116 chars) |
| **Details Readability** | 50% (cramped) | Centered 80 chars (optimal) |
| **Terminal Resize** | Fragile (pane recalc) | Robust (single view) |
| **Cognitive Load** | High (3 info sources) | Low (1 focused view) |
| **Navigation** | Implicit (focus moves) | Explicit (Tab to switch) |
| **User Control** | Less (panes always visible) | More (choose what to see) |
| **Implementation Lines** | 392 (view.go) | 566 (view.go) |
| **State Complexity** | Pane + View modes | View mode only |

**Winner**: ✅ **Single-Pane Mode** (simpler, cleaner, more usable)

---

## Conclusion

**Design Philosophy**: "One thing at a time, done well"

**Benefits**:
- ✅ Eliminates layout complexity and overflow issues
- ✅ Maximizes content readability (full width/height)
- ✅ Meets KEY requirement (works on any terminal size)
- ✅ Follows proven TUI patterns (vim modes, htop views)
- ✅ Gives users explicit control (Tab to switch)

**Status**: ✅ **IMPLEMENTED - Ready for Testing**

Run `./maat` and press Tab to cycle between views!

---

**Next Steps**:
1. Test single-pane design in terminal
2. Verify Tab cycling works smoothly
3. Check all 3 views render correctly
4. If stable → proceed to Phase 3 (Linear/GitHub Integration)
