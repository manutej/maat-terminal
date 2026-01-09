# 🎯 Test Single-Pane Design - Quick Guide

**Date**: 2026-01-07
**Status**: ✅ **READY TO TEST**

---

## What Changed

### Before ❌
- 3-pane split layout (Graph | Main | Detail)
- Overlapping text, broken layout
- Fixed 25% width for graph = unusable
- Complex pane cycling with Tab

### After ✅
- **Single full-screen view** with Tab switching
- **3 modes**: Graph → Details → Relations
- **Full terminal width** for all content
- **Simple**: One thing at a time, centered and readable

---

## Quick Test

```bash
cd /Users/manu/Documents/LUXOR/MAAT
./maat
```

### Test Flow (30 seconds)

**1. Launch → Graph View** (default)
```
Expected:
✅ Full-width hierarchical graph
✅ 97 nodes visible in tree layout
✅ Status bar: "[Graph View]"
✅ Navigation hints visible
```

**2. Press `Tab` → Details View**
```
Expected:
✅ Full-screen node details (centered)
✅ Large icon + title
✅ Type, status, priority visible
✅ Description readable
✅ Status bar: "[Details View]"
```

**3. Press `Tab` → Relations View**
```
Expected:
✅ Relationship list (outgoing/incoming)
✅ Relation types colored
✅ Summary count shown
✅ Status bar: "[Relations View]"
```

**4. Press `Tab` → Back to Graph View**
```
Expected:
✅ Cycles back to graph
✅ Focus preserved on same node
```

**5. Test Navigation in Graph**
```
hjkl or arrows: Move between nodes
Expected:
✅ Focus indicator (*) moves
✅ No visual glitches
✅ Smooth navigation
```

**6. Test Enter → Details**
```
Press Enter on focused node
Expected:
✅ Switches to Details view
✅ Shows selected node's info
```

**7. Test Esc → Back**
```
Press Esc in Details view
Expected:
✅ Returns to Graph view
✅ Focus preserved
```

**8. Resize Terminal**
```
Make terminal bigger/smaller
Expected:
✅ Layout adapts automatically
✅ No overlapping
✅ Content readable at any size (80+ chars)
```

---

## Keyboard Cheat Sheet

```
┌─────────────────────────────────────────┐
│  Navigation                             │
├─────────────────────────────────────────┤
│  Tab          Switch view (cycle)       │
│  Shift+Tab    Switch view (reverse)     │
│  Esc          Back to previous view     │
│  q / Ctrl+C   Quit                      │
│  r            Refresh data              │
├─────────────────────────────────────────┤
│  Graph View                             │
├─────────────────────────────────────────┤
│  hjkl / ←↓↑→  Navigate nodes            │
│  Enter        Drill into node (Details) │
└─────────────────────────────────────────┘
```

---

## Success Criteria

If you see this, it's working ✅:

- [ ] **Graph View** shows full-width tree (no truncation at 25%)
- [ ] **Tab key** smoothly cycles: Graph → Details → Relations → Graph
- [ ] **Details View** shows centered, readable node info with icons
- [ ] **Relations View** shows all connections with colored types
- [ ] **Status bar** updates to show current view mode
- [ ] **hjkl navigation** works in Graph view
- [ ] **Enter key** switches from Graph to Details
- [ ] **Esc key** returns to previous view
- [ ] **Resize terminal** doesn't break layout

---

## If Something Looks Wrong

**Symptom**: Graph still truncated at ~25 chars
- **Fix**: Restart terminal, ensure latest binary: `go build -o maat cmd/maat/main.go`

**Symptom**: Tab doesn't cycle views
- **Issue**: Old binary running
- **Fix**: Kill maat, rebuild, relaunch

**Symptom**: Blank screen or garbled text
- **Debug**: Take screenshot, save to `errors/` with timestamp
- **Report**: Terminal size (`echo $COLUMNS x $LINES`) + what you pressed

**Symptom**: Status bar shows wrong view mode
- **Expected**: This shouldn't happen (pure function rendering)
- **Debug**: Press `r` to refresh, see if it fixes itself

---

## What to Look For

### Graph View ✅
- Tree hierarchy visible with proper indentation
- Node icons: 📦 🐛 🔀 💾 📄
- Status icons: ✓ ◐ ○
- Focus indicator: * character next to selected node
- No horizontal scrolling needed

### Details View ✅
- Centered content (not edge-to-edge)
- Large icon + underlined title
- Colored type badge
- Status with icon (✅🔄📋)
- Priority with 🔥 icon
- Wrapped description (max 80 chars)
- Colored label badges
- Faint ID at bottom

### Relations View ✅
- Clear sections: "→ Outgoing" and "← Incoming"
- Colored relation types
- Arrow indicators (→)
- Summary count: "Total: X outgoing, Y incoming"

---

## Performance Check

**Memory**: Single view rendered (not 3 panes) = 66% less work

**Responsiveness**:
- Tab switching should be instant (< 50ms)
- hjkl navigation should feel immediate
- No lag or stuttering

**Smooth Cycling**:
```
Tab → [Graph View]
Tab → [Details View]
Tab → [Relations View]
Tab → [Graph View]  (loops back)
```

---

## Next Steps

### If It Works ✅
1. Mark design as stable
2. Proceed to Phase 3: Linear/GitHub Integration
3. Replace mock data with real API calls

### If Issues Found ❌
1. Take screenshot of broken state
2. Note terminal size and what you did
3. Report back with:
   - What looks wrong?
   - What did you press?
   - What's your terminal size?

---

## Quick Debug Commands

```bash
# Check terminal size
echo "$COLUMNS x $LINES"

# Rebuild if needed
cd /Users/manu/Documents/LUXOR/MAAT
go build -o maat cmd/maat/main.go

# Run with clean state
./maat
```

---

**Status**: ✅ **Build successful, ready to test**

**Architecture**: Single-pane with 3 view modes (Graph/Details/Relations)

**Key Improvement**: Full terminal width for all views = maximum readability

Run `./maat` and press Tab to experience the simplified design! 🚀
