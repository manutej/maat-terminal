# Phase 2 Kickoff: Graph Rendering + Navigation

**Date**: 2026-01-07T06:47:00Z
**Status**: Ready to Begin
**Complexity**: L2 (Reasoning layer added)

---

## Phase 1 Review

### Completed ✅

- ✅ SQLite store (479 lines, already existed)
- ✅ Mock data generation (1,893 lines)
- ✅ Mock data integration with TUI
- ✅ Binary builds and runs (8.1 MB)
- ✅ 95 nodes + 67 edges loaded on startup

### Key Metrics

| Metric | Value |
|--------|-------|
| Total Code | 4,483 lines (+112% from baseline) |
| Token Usage | ~500 tokens (7.2x efficiency) |
| Build Time | ~2 seconds |
| Binary Size | 8.1 MB (67% under limit) |

---

## Phase 2 Objectives

### Goal
Implement **graph visualization** and **keyboard navigation** to make MAAT's TUI functional for exploring the knowledge graph.

### Deliverables

1. **Graph Rendering Engine** (Phase 2A)
   - Hierarchical tree layout algorithm
   - ASCII/Unicode box drawing for graph visualization
   - Focused node highlighting
   - Relationship line rendering

2. **Keyboard Navigation** (Phase 2B)
   - Vim-style hjkl navigation
   - Tab for pane cycling
   - Enter to drill down, Esc to back out
   - Graph traversal following edges

---

## Workflows Overview

### Phase 2A: Graph Rendering Engine

**Complexity**: L2 (OBSERVE → REASON → GENERATE)

**Blocks**:
1. **StructureObserver** (L1): Analyze TUI layout requirements
2. **AlgorithmReasoner** (L2): Select optimal layout algorithm
3. **RenderGenerator** (L2): Generate graph rendering code

**Token Budget**: ~4,500 tokens
**Duration**: 3-5 days
**Output**: `internal/tui/render_graph.go` (~600 lines)

**Success Criteria**:
- Hierarchical tree layout renders correctly
- Focused node highlighted visually
- Edges displayed with ASCII/Unicode lines
- Graph viewport updates on model changes

### Phase 2B: Keyboard Navigation Logic

**Complexity**: L2 (OBSERVE → REASON → GENERATE)

**Blocks**:
1. **BehaviorObserver** (L1): Analyze navigation requirements
2. **PatternReasoner** (L2): Design navigation state machine
3. **NavigatorGenerator** (L2): Implement navigation handlers

**Token Budget**: ~2,600 tokens
**Duration**: 2-3 days
**Output**: `internal/tui/navigation.go` (~350 lines)

**Success Criteria**:
- hjkl moves focus through graph
- Tab cycles between panes
- Enter drills down into focused node
- Esc backs up navigation stack
- Navigation respects graph relationships

---

## Complexity Increase: L1 → L2

### What Changes?

**L1 (Phase 1)**:
- OBSERVE → GENERATE (direct transformation)
- No reasoning needed (schema → store, domain → mock data)
- Single-pass generation

**L2 (Phase 2)**:
- OBSERVE → **REASON** → GENERATE (adds reasoning step)
- Algorithm selection required (tree layout vs force-directed vs hierarchical)
- Pattern design needed (navigation state machine)
- Multiple valid solutions exist

### Why L2?

**Graph Rendering**:
- Multiple layout algorithms available (tree, force-directed, circular)
- Must reason about: performance, readability, space efficiency
- Must select optimal algorithm for terminal constraints

**Keyboard Navigation**:
- Multiple navigation patterns possible (vim, emacs, custom)
- Must reason about: ergonomics, consistency, discoverability
- Must design state machine for navigation stack

---

## Technical Approach

### Graph Rendering Strategy

**Layout Algorithm Decision Tree**:
```
IF (graph is hierarchical AND < 100 nodes)
  → Use Tree Layout (simple, fast, readable)
ELSE IF (graph has cycles AND < 50 nodes)
  → Use Force-Directed (handles cycles, slower)
ELSE
  → Use Hierarchical Layout (scales better, consistent)
```

**For Phase 2**: Use **Tree Layout** (graph is hierarchical, 95 nodes)

**Rendering Approach**:
1. Calculate node positions (tree layout algorithm)
2. Render nodes (boxes with titles)
3. Draw edges (ASCII lines: │ ─ ┌ ┐ └ ┘)
4. Highlight focused node (bold + color)
5. Update viewport on navigation

### Navigation Strategy

**State Machine**:
```
States: {Idle, Navigating, Drilling, Backing}

Events:
- hjkl → Navigating (move focus)
- Enter → Drilling (push view)
- Esc → Backing (pop view)
- Tab → Pane cycling

Transitions:
- Idle + hjkl → Navigating → update focusedNode
- Idle + Enter → Drilling → push current node to stack
- Idle + Esc → Backing → pop stack if not empty
```

**Navigation Rules**:
- h/j/k/l moves focus to adjacent nodes in graph
- Following edges (left/right/up/down based on relationship)
- Boundary checking (don't move beyond visible nodes)
- Focus wrapping (optional: wrap at edges)

---

## Updated Code Metrics Projection

| Component | Phase 1 | Phase 2 (Projected) | Growth |
|-----------|---------|---------------------|--------|
| internal/graph | 649 lines | 649 lines | 0% |
| internal/tui | 2,303 lines | 3,253 lines | +41% |
| cmd/maat | 24 lines | 24 lines | 0% |
| **Total** | **4,483 lines** | **5,433 lines** | **+21%** |

**New Files**:
- `internal/tui/render_graph.go` (~600 lines)
- `internal/tui/navigation.go` (~350 lines)

**Binary Size Projection**: 8.1 MB → 9.5 MB (still < 25 MB ✅)

---

## Risk Assessment

### Technical Risks

**Rendering Performance** (Medium Risk)
- **Risk**: Rendering 95 nodes + 67 edges may be slow
- **Mitigation**: Viewport culling (only render visible nodes)
- **Fallback**: Simple list view if tree layout too slow

**Navigation Complexity** (Low Risk)
- **Risk**: Graph cycles may confuse navigation
- **Mitigation**: Maintain visited set, detect cycles
- **Fallback**: Breadth-first navigation (simpler)

**Terminal Compatibility** (Low Risk)
- **Risk**: Unicode box drawing may not render on all terminals
- **Mitigation**: Fallback to ASCII-only mode
- **Detection**: Check TERM environment variable

### Schedule Risks

**Workflow Dependency** (Low Risk)
- Phase 2A must complete before 2B (navigation needs rendering)
- **Mitigation**: Can prototype navigation logic with placeholder rendering

**Integration Complexity** (Medium Risk)
- Rendering + navigation must work together seamlessly
- **Mitigation**: Phase 2B tests against Phase 2A rendering
- **Validation**: Manual testing with all navigation patterns

---

## Constitutional Alignment

### Commandment Checklist

**Phase 2A: Graph Rendering**
- ✅ #1 Immutable Truth: Rendering pure function (Model → View)
- ✅ #2 Graph Supremacy: Render from graph.Node/Edge types
- ✅ #3 Text Interface: ASCII/Unicode rendering
- ✅ #8 Async Purity: No goroutines, just view functions

**Phase 2B: Keyboard Navigation**
- ✅ #1 Immutable Truth: Navigation returns new Model
- ✅ #4 Navigation Monopoly: Enter drills down, Esc backs out
- ✅ #7 Composition Monopoly: Build complex navigation from simple moves
- ✅ #9 Terminal Citizenship: Vim-style hjkl (standard in terminal UIs)

---

## Execution Plan

### Week 1: Phase 2A (Graph Rendering)

**Day 1-2**: Algorithm Research + Selection
- Analyze tree layout algorithms
- Prototype positioning logic
- Validate against 95-node mock graph

**Day 3-4**: Implementation
- Generate `render_graph.go`
- Implement tree layout algorithm
- ASCII/Unicode line drawing

**Day 5**: Integration + Testing
- Integrate with view.go
- Test with all node types
- Verify performance < 100ms render time

### Week 2: Phase 2B (Keyboard Navigation)

**Day 1-2**: Navigation Design
- Design state machine
- Map hjkl to graph movements
- Plan edge-following logic

**Day 3-4**: Implementation
- Generate `navigation.go`
- Implement movement handlers
- Integrate with update.go

**Day 5**: Integration + Testing
- End-to-end navigation testing
- Edge case handling (cycles, boundaries)
- Performance validation

---

## Success Metrics

### Functional

- ✅ Graph renders with all 95 nodes visible
- ✅ Focused node highlighted clearly
- ✅ Edges displayed with proper connections
- ✅ hjkl navigation works in all directions
- ✅ Enter/Esc drill down/back up correctly
- ✅ Tab cycles through panes

### Non-Functional

- ✅ Render time < 100ms (60 FPS target)
- ✅ Navigation latency < 16ms (responsive feel)
- ✅ Terminal compatibility: xterm, alacritty, iTerm2
- ✅ Binary size < 25 MB
- ✅ No memory leaks during extended navigation

### Constitutional

- ✅ Pure functions (no pointer mutations)
- ✅ Value receivers throughout
- ✅ No global state
- ✅ All effects via tea.Cmd
- ✅ Commandments #1, #2, #3, #4, #7, #8, #9 followed

---

## Token Budget

| Phase | Workflow | Budget | Actual | Status |
|-------|----------|--------|--------|--------|
| Phase 1A | SQLite Store | 2,000 | 0* | ✅ Pre-existing |
| Phase 1B | Mock Data | 1,600 | ~500 | ✅ 3.2x efficient |
| **Phase 1 Total** | | **3,600** | **~500** | **✅ 7.2x efficient** |
| Phase 2A | Graph Rendering | 4,500 | TBD | ⏳ Pending |
| Phase 2B | Keyboard Navigation | 2,600 | TBD | ⏳ Pending |
| **Phase 2 Total** | | **7,100** | **TBD** | **⏳ Pending** |

**Phase 1 Efficiency**: 7.2x better than planned
**Phase 2 Budget**: 7,100 tokens (conservative estimate)
**Total Budget**: 10,700 tokens (Phase 1 + 2)

---

## Dependencies

### External

- ✅ Bubble Tea (charmbracelet/bubbletea) - TUI framework
- ✅ Lipgloss (charmbracelet/lipgloss) - styling
- ⏳ go-graph (optional, if complex layout needed)

### Internal

- ✅ graph.Node, graph.Edge (data structures)
- ✅ GetMockGraph() (test data source)
- ✅ Model.WithNodes(), Model.WithEdges() (state setters)
- ⏳ ViewGraph rendering (to be implemented)

---

## Next Actions

1. **Generate Phase 2A YAML workflows** (graph rendering)
2. **Execute Phase 2A** (StructureObserver → AlgorithmReasoner → RenderGenerator)
3. **Test rendering** with mock graph data
4. **Generate Phase 2B YAML workflows** (keyboard navigation)
5. **Execute Phase 2B** (BehaviorObserver → PatternReasoner → NavigatorGenerator)
6. **Integration testing** (rendering + navigation together)
7. **Phase 2 completion report**

---

## Questions to Answer

### Phase 2A (Rendering)

- Q: Which layout algorithm performs best for 95 nodes?
- A: TBD (will test tree layout first, measure performance)

- Q: ASCII-only or Unicode box drawing?
- A: TBD (check TERM variable, fallback to ASCII if needed)

- Q: Should viewport cull off-screen nodes?
- A: TBD (likely yes if > 50 visible nodes)

### Phase 2B (Navigation)

- Q: How to handle graph cycles in navigation?
- A: TBD (maintain visited set, allow revisiting but prevent loops)

- Q: Should focus wrap at boundaries?
- A: TBD (probably no, vim-style hard boundaries)

- Q: Edge-following direction mapping?
- A: TBD (h/l = horizontal edges, j/k = vertical hierarchy)

---

## Philosophy Validation

### "Simple & Working" Checklist

**Phase 2A**:
- ✅ Start with simplest algorithm (tree layout)
- ✅ ASCII rendering before Unicode (fallback ready)
- ✅ No premature optimization (profile first)
- ✅ Manual testing before automated tests

**Phase 2B**:
- ✅ Basic movement before advanced features
- ✅ Single pane navigation before multi-pane
- ✅ Hard boundaries before smart wrapping
- ✅ Simple state machine before complex logic

**Complexity Gates**:
1. ✅ L1 insufficient? (Yes - need layout algorithm selection)
2. ✅ L2 necessary? (Yes - multiple valid approaches exist)
3. ✅ Risk acceptable? (Yes - can fallback to simpler layouts)
4. ✅ Value justifiable? (Yes - core UX feature)
5. ✅ Budget available? (Yes - 7,100 tokens allocated)

---

## Conclusion

Phase 2 is **ready to begin**. All prerequisites complete:

1. ✅ Phase 1 delivered working foundation
2. ✅ Mock data integrated and tested
3. ✅ TUI initializes successfully
4. ✅ Complexity increase justified (L1 → L2)
5. ✅ Token budget allocated
6. ✅ Success criteria defined

**Proceed to Phase 2A workflow generation.**

---

**Status**: ✅ Ready to Execute
**Risk Level**: Low (clear requirements, tested foundation)
**Next Action**: Generate `workflows/phase2-graph-rendering.yaml`
