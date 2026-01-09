# MAAT Testing Guide

## Prerequisites

- Terminal with TTY support (iTerm2, Terminal.app, etc.)
- MAAT binary built (`go build ./cmd/maat`)

## Quick Start Test

```bash
cd /Users/manu/Documents/LUXOR/MAAT
./maat
```

## Phase 2 Feature Tests

### Test 1: Basic Navigation (hjkl)

**Test**: Verify vim-style navigation works

1. Launch MAAT: `./maat`
2. Observe the graph pane (left side)
3. Note the focused node (marked with `*`)
4. Press `j` (down) - focus should move to a child node or node below
5. Press `k` (up) - focus should move to parent node or node above
6. Press `h` (left) - focus should move to node on the left
7. Press `l` (right) - focus should move to node on the right

**Expected Result**:
- `*` marker moves to different nodes
- Main pane (center) updates with focused node details
- Navigation feels responsive (< 100ms latency)

### Test 2: Pane Cycling (Tab)

**Test**: Verify Tab cycles through panes

1. Press `Tab` - active pane should switch from Graph to Main
2. Press `Tab` again - should switch from Main to Detail
3. Press `Tab` again - should cycle back to Graph

**Expected Result**:
- Active pane indicator changes
- Border highlighting shows which pane is active
- Cycle completes in 3 Tab presses

### Test 3: Drill Down/Back (Enter/Esc)

**Test**: Verify hierarchical navigation

1. Navigate to a Project node (use hjkl)
2. Press `Enter` - should drill into Project (filter to show only related nodes)
3. Press `Esc` - should return to full graph view

**Expected Result**:
- Enter: Graph filters to subgraph
- Esc: Graph restores to full view
- Navigation stack works correctly

### Test 4: Edge Following

**Test**: Verify j/k follow parent-child edges

1. Navigate to a node with children (e.g., Project)
2. Press `j` - should move to first child (Issue)
3. Press `j` again - should move to child's child (PR or Commit)
4. Press `k` - should move back to parent

**Expected Result**:
- j/k follow edge relationships (not just spatial)
- Hierarchical traversal works correctly
- Parent-child relationships visible in movement

### Test 5: Spatial Fallback

**Test**: Verify spatial navigation works for orphan nodes

1. Navigate to an orphan node (node with no edges)
2. Press `j` - should move to nearest node below (spatial)
3. Press `k` - should move to nearest node above (spatial)

**Expected Result**:
- Navigation works even without edges
- Selects nearest node in direction
- No crashes or hangs

### Test 6: Boundary Handling

**Test**: Verify vim-style hard boundaries

1. Navigate to the topmost node
2. Press `k` repeatedly - focus should stay at top (not wrap)
3. Navigate to the leftmost node
4. Press `h` repeatedly - focus should stay at left (not wrap)

**Expected Result**:
- No focus wrapping
- Focus stays at boundary (vim behavior)
- No errors or crashes

### Test 7: Quit Commands

**Test**: Verify exit mechanisms

1. Press `q` - should quit immediately
2. Relaunch, press `Ctrl+C` - should quit immediately

**Expected Result**:
- Clean exit, no errors
- Terminal returns to normal state

## Visual Inspection Checklist

When MAAT is running, verify:

- [ ] Graph pane renders ASCII boxes with node titles
- [ ] Edges rendered with vertical lines (│) and arrows (↓)
- [ ] Focused node marked with `*` prefix
- [ ] Main pane shows focused node details (type, status, title)
- [ ] Detail pane shows relationships (edges to/from focused node)
- [ ] Status bar shows help text (hjkl: navigate | Tab: switch | etc.)
- [ ] No visual glitches or flickering
- [ ] Terminal colors look correct (if using colors)

## Mock Data Validation

The TUI loads 95 mock nodes:

| Node Type | Count |
|-----------|-------|
| Project | 5 |
| Issue | 20 |
| PR | 15 |
| Commit | 30 |
| File | 25 |
| **Total** | **95** |

**Edges**: 67 relationships between nodes

Verify the graph shows all node types with proper icons:
- Project: No icon (root nodes)
- Issue: `○` or `✓` (based on status)
- PR: `◐` or `●` (based on status)
- Commit: `·` (small dot)
- File: No special icon

## Performance Checks

**Navigation Latency**:
- Press hjkl rapidly - should feel responsive
- No visible lag between keypress and focus change
- Target: < 16ms (60 FPS feel)

**Rendering Performance**:
- Initial render should be instant (< 100ms)
- Graph should not flicker on navigation
- Smooth redraw on pane switching

## Known Limitations (Not Bugs)

1. **Multiple Children**: `j` selects first child only (not nearest)
2. **No Scrolling**: Graph may overflow terminal on very large displays
3. **No Undo**: Navigation has no history (except Enter/Esc stack)
4. **ASCII Only**: No fancy graphics (by design - Commandment #3)

## Troubleshooting

### "No graph data loaded"

**Cause**: Mock data not loading
**Fix**: Check `internal/tui/commands.go` calls `GetMockGraph()`

### "Cannot open TTY"

**Cause**: Running in non-interactive environment
**Fix**: Run in real terminal (iTerm2, Terminal.app, etc.)

### Navigation doesn't work

**Cause**: Key bindings not integrated
**Fix**: Check `internal/tui/update.go` has hjkl cases

### Focus stays stuck

**Cause**: Navigation algorithm issue
**Fix**: Check `internal/tui/navigation.go` moveLeft/Right/Up/Down

### Binary won't run

**Cause**: Not executable or missing dependencies
**Fix**:
```bash
chmod +x maat
go build ./cmd/maat  # Rebuild
```

## Advanced Testing (Optional)

### Performance Profiling

```bash
# CPU profiling
go build -o maat ./cmd/maat
./maat --cpuprofile=cpu.prof  # If implemented
go tool pprof cpu.prof

# Memory profiling
go tool pprof -alloc_space maat mem.prof
```

### Unit Tests (Phase 3)

```bash
# Run tests (when added in Phase 3)
go test ./internal/tui -v
go test ./internal/tui -bench=. -benchmem
```

## Success Criteria

Phase 2 is successful if:

- [x] Binary compiles without errors
- [x] TUI launches in terminal
- [x] All 95 nodes render correctly
- [x] hjkl navigation works in all directions
- [x] Tab cycles through panes
- [x] Enter/Esc drill down/back
- [x] No crashes or panics
- [x] Navigation feels responsive (< 100ms)
- [x] Focus indicator updates correctly
- [x] Main/Detail panes update with focused node

## Reporting Issues

If you find bugs, document:

1. **Steps to reproduce**: Exact key sequence
2. **Expected behavior**: What should happen
3. **Actual behavior**: What actually happened
4. **Error messages**: Any terminal output
5. **Environment**: macOS version, terminal app, MAAT version

Example bug report:
```
Steps:
1. Launch MAAT
2. Navigate to Project node (3rd from top)
3. Press 'j'

Expected: Focus moves to first child Issue
Actual: Focus stays on Project
Error: None

Environment: macOS Ventura 13.1, iTerm2 3.4.19, MAAT Phase 2B
```

---

**Quick Test Command**:
```bash
cd /Users/manu/Documents/LUXOR/MAAT && ./maat
```

Press `q` to quit.
