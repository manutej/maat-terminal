# MAAT Execution Plan: Progressive Implementation Strategy

**Document Type**: Execution Roadmap
**Version**: 1.0
**Date**: 2026-01-06
**Status**: Ready to Execute
**Philosophy**: Ship simple & working, iterate to complex only if needed

---

## Executive Summary

This plan orchestrates **8 workflows** across **5 phases** using **20 atomic blocks** to complete MAAT implementation from current 70% foundation to production-ready v1.0.

**Current State**: 2,111 lines Go code (Phase 1: 70% complete)
**Target State**: Full MVP + Phase 2-5 features (est. 10,000+ lines)
**Execution Time**: 18-24 days (assumes 4-6 hrs/day work)
**Token Budget**: 35,400 tokens total

---

## Quick Start: Next 3 Actions

```bash
# 1. Complete Phase 1A: SQLite Store (TODAY)
/ois-compose --workflow-plan workflows/phase1-sqlite-store.yaml \
  --input-file internal/graph/schema.go \
  --output internal/graph/store_generated.go

# 2. Complete Phase 1B: Mock Data (TODAY)
/ois-compose --workflow-plan workflows/phase1-mock-data.yaml \
  --input-file "docs/CONSTITUTION.md specs/FUNCTIONAL-REQUIREMENTS.md" \
  --output internal/tui/mock_data.go

# 3. Test Phase 1 Complete (TOMORROW)
make test && ./maat  # Should render mock graph
```

**Estimated Completion**: 2-3 days for Phase 1 finish

---

## Phase 1: Foundation (CURRENT - 70% Complete)

### Status: ⚙️ In Progress

| Task | Status | Lines | Owner | ETA |
|------|--------|-------|-------|-----|
| Go project structure | ✅ Complete | 258 | Done | - |
| Bubble Tea model | ✅ Complete | 392 | Done | - |
| Graph schema | ✅ Complete | 170 | Done | - |
| 3-pane layout | ✅ Complete | 151 | Done | - |
| **SQLite store** | ⚙️ 30% | 150/500 | **YOU** | **Today** |
| **Mock data generator** | 🔜 Pending | 0/300 | **YOU** | **Today** |

**Total**: 2,111 → 3,400 lines (+1,289 lines)

---

### Workflow 1A: SQLite Store Implementation

**File**: `workflows/phase1-sqlite-store.yaml`

```yaml
metadata:
  task: "Implement SQLite store from schema.go"
  complexity: L1
  budget: 2000
  duration: 60-90 min

composition:
  pattern: sequential
  operators: [→]

  step_1_observe:
    block: structure-observer
    foundation: systems-thinking
    function: OBSERVE
    input:
      file: internal/graph/schema.go
      focus: [Node, Edge, NodeType, EdgeType, Store interface]
    output: Observable<Schema>
    budget: 800

  step_2_generate:
    block: type-generator
    foundation: category-theory
    function: GENERATE
    input: Observable<Schema>
    output: Generated<StoreImpl>
    budget: 1200
    validation:
      - Implements Store interface completely
      - Pure functions, no global state (#1)
      - SQLite via github.com/mattn/go-sqlite3
      - Table creation with schema migration

type_flow:
  - input: schema.go (types + interface)
  - intermediate: Observable<Schema> (structure extraction)
  - output: store.go (CRUD implementation)

output_file: internal/graph/store_generated.go

success_criteria:
  - AddNode/AddEdge/GetNode/GetEdges/GetNeighbors implemented
  - Tests pass (store_test.go)
  - Binary compiles
```

**Execution**:
```bash
/ois-compose --workflow-plan workflows/phase1-sqlite-store.yaml \
  --input-file internal/graph/schema.go \
  --output internal/graph/store_generated.go \
  --verbose
```

**Expected Output**: ~500 lines of SQLite store implementation

---

### Workflow 1B: Mock Data Generator

**File**: `workflows/phase1-mock-data.yaml`

```yaml
metadata:
  task: "Create realistic mock fixtures for TUI testing"
  complexity: L1
  budget: 1600
  duration: 45-60 min

composition:
  pattern: sequential
  operators: [→]

  step_1_gather:
    block: knowledge-gatherer
    foundation: knowledge-synthesis
    function: OBSERVE
    input:
      files: [docs/CONSTITUTION.md, specs/FUNCTIONAL-REQUIREMENTS.md]
      focus: [Issue, PR, Commit, File entities and relationships]
    output: Observable<DomainModel>
    budget: 600

  step_2_generate:
    block: mock-generator
    foundation: specification-driven
    function: GENERATE
    input: Observable<DomainModel>
    output: Generated<MockData>
    budget: 1000
    constraints:
      - 50-100 nodes (Issues, PRs, Commits, Files)
      - Realistic relationships (blocks, implements, mentions)
      - Covers all NodeTypes and EdgeTypes
      - Hierarchical structure (projects → issues → PRs → commits)

type_flow:
  - input: Spec documents (domain understanding)
  - intermediate: Observable<DomainModel> (entities + relationships)
  - output: mock_data.go (graph.Node[] + graph.Edge[])

output_file: internal/tui/mock_data.go

success_criteria:
  - GetMockGraph() returns populated graph
  - Nodes have realistic titles/descriptions
  - Edges reflect blocking/implementation relationships
  - TUI can render mock graph
```

**Execution**:
```bash
/ois-compose --workflow-plan workflows/phase1-mock-data.yaml \
  --input-file "docs/CONSTITUTION.md specs/FUNCTIONAL-REQUIREMENTS.md" \
  --output internal/tui/mock_data.go
```

**Expected Output**: ~300 lines of mock data generation

---

### Phase 1 Validation

```bash
# Run tests
make test

# Build binary
make build

# Test with mock data
./maat  # Should display mock graph in TUI

# Verify FR-001 acceptance criteria
# - Issues appear as nodes ✓
# - Edges show relationships ✓
# - Status reflected in styling ✓
```

**Gate**: All tests pass + TUI renders mock graph → Proceed to Phase 2

---

## Phase 2: Graph Navigation (Weeks 3-4)

### Status: 🔜 Ready to Start (after Phase 1)

| Task | Status | Lines | Owner | ETA |
|------|--------|-------|-------|-----|
| Graph rendering engine | 🔜 Pending | 0/800 | YOU | Week 3 |
| Keyboard navigation | 🔜 Pending | 0/600 | YOU | Week 3-4 |
| Detail pane viewport | 🔜 Pending | 0/400 | YOU | Week 4 |
| Drill-down/back-up | 🔜 Pending | 0/300 | YOU | Week 4 |

**Total**: 3,400 → 5,500 lines (+2,100 lines)

---

### Workflow 2A: Graph Rendering Engine

**File**: `workflows/phase2-graph-render.yaml`

```yaml
metadata:
  task: "Render knowledge graph in terminal with Bubble Tea"
  complexity: L2
  budget: 4500
  duration: 3-4 hours

composition:
  pattern: sequential
  operators: [→]

  step_1_observe:
    block: structure-observer
    foundation: systems-thinking
    function: OBSERVE
    input:
      file: internal/tui/model.go
      focus: [DisplayNode, DisplayEdge, Model]
    output: Observable<GraphStructure>
    budget: 1000

  step_2_reason:
    block: layout-reasoner
    foundation: abstraction-principles
    function: REASON
    input: Observable<GraphStructure>
    output: Reasoning<LayoutAlgorithm>
    budget: 1500
    constraints:
      - Hierarchical tree layout (issue → PR → commit)
      - Indented child nodes
      - Collapsible subtrees
      - < 100ms for 500 nodes (FR-001)

  step_3_generate:
    block: render-generator
    foundation: elm-architecture
    function: GENERATE
    input: Reasoning<LayoutAlgorithm>
    output: Generated<GraphView>
    budget: 2000
    includes:
      - Lipgloss styling (colors, borders, focus)
      - vim-style navigation integration
      - Status-based coloring (todo, in-progress, done)

type_flow:
  - input: DisplayNode[] + DisplayEdge[] (data model)
  - intermediate: Reasoning<LayoutAlgorithm> (tree layout)
  - output: views/graph.go (renderer)

output_file: internal/tui/views/graph.go

success_criteria:
  - Renders mock graph hierarchically
  - Focused node highlighted
  - Parent-child relationships visible
  - Performance < 100ms for 500 nodes
```

**Execution**:
```bash
/ois-compose --workflow-plan workflows/phase2-graph-render.yaml \
  --input-file internal/tui/model.go \
  --output internal/tui/views/graph.go \
  --verbose
```

---

### Workflow 2B: Keyboard Navigation Logic

**File**: `workflows/phase2-keyboard-nav.yaml`

```yaml
metadata:
  task: "Implement vim-style navigation (hjkl, Enter, Esc)"
  complexity: L2
  budget: 2600
  duration: 2-3 hours

composition:
  pattern: sequential
  operators: [→]

  step_1_observe:
    block: pattern-observer-spec
    foundation: unix-philosophy
    function: OBSERVE
    input:
      file: specs/FUNCTIONAL-REQUIREMENTS.md
      section: FR-002
    output: Observable<KeyBindings>
    budget: 800

  step_2_generate:
    block: navigation-generator
    foundation: elm-architecture
    function: GENERATE
    input: Observable<KeyBindings>
    output: Generated<KeyHandlers>
    budget: 1800
    constraints:
      - Pure functions (Commandment #1)
      - NavigationStack for drill-down/back-up
      - No goroutines (Commandment #5)
      - All FR-002 keybindings (hjkl, Enter, Esc, Tab, /, ?, q)

type_flow:
  - input: FR-002 specification
  - intermediate: Observable<KeyBindings> (key → action map)
  - output: internal/tui/keys_generated.go

output_file: internal/tui/keys_generated.go

success_criteria:
  - All FR-002 acceptance criteria met
  - hjkl navigates between nodes
  - Enter drills into focused node
  - Esc returns to parent context
  - Tab cycles panes
```

---

### Phase 2 Validation

```bash
# Test navigation
./maat
# Press j/k to navigate nodes
# Press Enter on issue to see PRs
# Press Esc to go back
# Press Tab to cycle panes

# Verify FR-002 acceptance criteria
# - hjkl navigation ✓
# - Enter drill-down ✓
# - Esc back-up ✓
# - Tab pane cycling ✓
```

**Gate**: All FR-002 criteria met → Proceed to Phase 3

---

## Phase 3: Integrations (Weeks 5-7)

### Status: 🔮 Future (requires Phase 2)

| Task | Status | Lines | Owner | ETA |
|------|--------|-------|-------|-----|
| Linear GraphQL client | 🔮 Future | 0/1200 | YOU | Week 5-6 |
| GitHub REST client | 🔮 Future | 0/1000 | YOU | Week 6 |
| Git commit history | 🔮 Future | 0/800 | YOU | Week 7 |

**Total**: 5,500 → 8,500 lines (+3,000 lines)

---

### Workflow 3A: Linear GraphQL Client

**Complexity**: L3 (API integration with error handling)

**File**: `workflows/phase3-linear-client.yaml`

```yaml
metadata:
  task: "Create thin Linear GraphQL client"
  complexity: L3
  budget: 6000
  duration: 4-5 hours

composition:
  pattern: sequential
  operators: [→]

  step_1_observe:
    block: api-observer
    foundation: systems-thinking
    function: OBSERVE
    input:
      source: Linear GraphQL API docs (Context7)
      operations: [fetchIssues, fetchProjects]
    output: Observable<APISchema>
    budget: 1500

  step_2_reason:
    block: schema-reasoner
    foundation: specification-driven
    function: REASON
    input: Observable<APISchema>
    output: Reasoning<ClientDesign>
    budget: 2000
    constraints:
      - Thin client only (Commandment #7)
      - Read operations priority
      - Write requires confirmation (Commandment #10)

  step_3_generate:
    block: client-generator
    foundation: category-theory
    function: GENERATE
    input: Reasoning<ClientDesign>
    output: Generated<LinearClient>
    budget: 2500
    includes:
      - genqlient GraphQL code generation
      - Exponential backoff retry
      - Rate limit handling
      - Adapter: Linear.Issue → graph.Node

output_files:
  - internal/linear/client.go
  - internal/linear/adapter.go
  - internal/linear/types.go

success_criteria:
  - FetchIssues returns []Issue
  - Converts to graph.Node correctly
  - Handles errors gracefully
  - Respects rate limits
```

---

### Workflow 3B: Git History Integration

**Complexity**: L2-L3 (local operations, parsing)

**File**: `workflows/phase3-git-history.yaml`

```yaml
metadata:
  task: "Parse git commits into knowledge graph nodes"
  complexity: L2
  budget: 2700
  duration: 2-3 hours

composition:
  pattern: sequential
  operators: [→]

  step_1_observe:
    block: git-observer
    foundation: unix-philosophy
    function: OBSERVE
    input:
      source: .git directory (go-git)
      limit: 1000 commits
    output: Observable<CommitHistory>
    budget: 1200

  step_2_generate:
    block: commit-generator
    foundation: systems-thinking
    function: GENERATE
    input: Observable<CommitHistory>
    output: Generated<CommitNodes>
    budget: 1500
    includes:
      - Parse commit messages for issue refs (#123)
      - Create edges: Commit → Issue (mentions)
      - Efficient (< 2s for 1000 commits)

output_files:
  - internal/git/adapter.go
  - internal/git/commits.go

success_criteria:
  - Parses commit history
  - Extracts issue references
  - Creates graph nodes + edges
  - Performance < 2s for 1000 commits
```

---

## Phase 4: AI & Polish (Weeks 8-10)

### Status: 🌌 Advanced (requires Phase 3)

| Task | Status | Lines | Owner | ETA |
|------|--------|-------|-------|-----|
| Claude MCP bridge | 🌌 Advanced | 0/1500 | YOU | Week 8-9 |
| Confirmation gates | 🌌 Advanced | 0/400 | YOU | Week 9 |
| Role-based views | 🌌 Advanced | 0/600 | YOU | Week 10 |

**Total**: 8,500 → 11,000 lines (+2,500 lines)

---

### Workflow 4A: Claude MCP Bridge

**Complexity**: L4 (human-in-loop coordination)

**File**: `workflows/phase4-claude-mcp.yaml`

```yaml
metadata:
  task: "Claude integration with human-in-loop gates"
  complexity: L4
  budget: 7500
  duration: 5-6 hours

composition:
  pattern: parallel_synthesis_refine
  operators: [||, →]

  stage_1_parallel:
    operator: ||
    agents:
      - context-observer:
          input: Focused node + edges
          output: Observable<NodeContext>
          budget: 1200

      - knowledge-gatherer:
          input: Navigation history + breadcrumb
          output: Observable<SessionContext>
          budget: 1000

  stage_2_synthesis:
    operator: →
    block: context-reasoner
    input: (Observable<NodeContext>, Observable<SessionContext>)
    output: Reasoning<ClaudePrompt>
    budget: 2500

  stage_3_generate:
    operator: →
    block: prompt-generator
    input: Reasoning<ClaudePrompt>
    output: Generated<MCPRequest>
    budget: 2000

  stage_4_refine:
    operator: →
    block: human-refiner
    input: Generated<MCPRequest>
    output: Refined<ConfirmedRequest>
    budget: 800

output_files:
  - internal/claude/bridge.go
  - internal/claude/context.go
  - internal/claude/audit.go

success_criteria:
  - Ctrl+A invokes Claude (Commandment #6)
  - Confirmation gate for actions (Commandment #10)
  - Full audit trail (ADR-004)
  - No ambient suggestions (AFR-005)
```

---

## Phase 5: Extensibility (Weeks 11-12)

### Status: 🚀 Future (requires Phase 4)

| Task | Status | Lines | Owner | ETA |
|------|--------|-------|-------|-----|
| Plugin interface | 🚀 Future | 0/800 | YOU | Week 11 |
| Example plugins | 🚀 Future | 0/600 | YOU | Week 11-12 |
| IDP self-service | 🚀 Future | 0/400 | YOU | Week 12 |

**Total**: 11,000 → 12,800 lines (+1,800 lines)

---

### Workflow 5A: Plugin Interface Design

**Complexity**: L5 (system-level architecture)

**File**: `workflows/phase5-plugin-interface.yaml`

```yaml
metadata:
  task: "Design plugin interface for extensibility"
  complexity: L5
  budget: 8500
  duration: 6-7 hours

composition:
  pattern: parallel_reasoning_optimize
  operators: [||, →]

  stage_1_parallel:
    operator: ||
    agents:
      - pattern-observer:
          input: [internal/linear, internal/github, internal/git]
          output: Observable<ClientPatterns>
          budget: 1800

      - architecture-reasoner:
          input: specs/ADR/ADR-007-plugin-architecture.md
          output: Reasoning<PluginDesign>
          budget: 2200

  stage_2_generate:
    operator: →
    block: interface-generator
    input: (Observable<ClientPatterns>, Reasoning<PluginDesign>)
    output: Generated<PluginInterface>
    budget: 2500

  stage_3_optimize:
    operator: →
    block: system-optimizer
    input: Generated<PluginInterface>
    output: Optimized<PluginAPI>
    budget: 2000

output_files:
  - internal/plugin/interface.go
  - internal/plugin/registry.go
  - internal/plugin/builtin/linear.go

success_criteria:
  - DataSourcePlugin interface defined
  - Type-safe registration
  - Lifecycle hooks (init, fetch, cleanup)
  - Example plugin works
```

---

## Execution Timeline

### Week-by-Week Breakdown

| Week | Phase | Tasks | Lines | Agent Budget |
|------|-------|-------|-------|--------------|
| **Week 1** | Phase 1 (finish) | SQLite + Mock data | +1,289 | 3,600 |
| **Week 2** | Phase 1 validation | Testing + bugfixes | +200 | 1,000 |
| **Week 3** | Phase 2 (start) | Graph render | +800 | 4,500 |
| **Week 4** | Phase 2 (finish) | Keyboard nav + detail | +1,300 | 2,600 |
| **Week 5-6** | Phase 3 (start) | Linear + GitHub | +2,200 | 7,000 |
| **Week 7** | Phase 3 (finish) | Git history | +800 | 2,700 |
| **Week 8-9** | Phase 4 (start) | Claude MCP | +1,900 | 7,500 |
| **Week 10** | Phase 4 (finish) | Role views | +600 | 1,500 |
| **Week 11-12** | Phase 5 | Plugin system | +1,800 | 8,500 |

**Total Duration**: 12 weeks (3 months)
**Total Lines**: 2,111 → 12,800 (+10,689 lines, 507% growth)
**Total Agent Budget**: 35,400 tokens

---

## Token Budget Allocation

| Phase | Workflows | Tokens | % of Total |
|-------|-----------|--------|------------|
| Phase 1 | 2 | 3,600 | 10.2% |
| Phase 2 | 2 | 7,100 | 20.1% |
| Phase 3 | 2 | 8,700 | 24.6% |
| Phase 4 | 1 | 7,500 | 21.2% |
| Phase 5 | 1 | 8,500 | 24.0% |
| **Total** | **8** | **35,400** | **100%** |

**Average per workflow**: 4,425 tokens

---

## Risk Mitigation

### Risk 1: API Integration Failures (Phase 3)

**Probability**: Medium
**Impact**: High (blocks Phase 4-5)

**Mitigation**:
- Start with mock data (Phase 1) → TUI works without APIs
- Implement git integration first (local, no API)
- Linear/GitHub clients independently testable
- Fallback: Use cached data from previous syncs

---

### Risk 2: Claude MCP Complexity (Phase 4)

**Probability**: Medium
**Impact**: Medium (AI features are P2, not P0)

**Mitigation**:
- Human-in-loop reduces AI unpredictability
- Audit trail allows debugging
- Can ship without AI initially (Phase 1-3 sufficient)
- Use `/think` for complex reasoning

---

### Risk 3: Plugin System Over-Engineering (Phase 5)

**Probability**: Medium
**Impact**: Low (extensibility is P3)

**Mitigation**:
- Unix philosophy: Keep interface minimal (5 methods)
- Built-in plugins as reference implementations
- Ship without plugins if needed (can add v1.1)
- Commandment #2: Single Responsibility per plugin

---

## Quality Gates

### Per-Phase Gates

| Phase | Quality Gate | Measurement |
|-------|--------------|-------------|
| Phase 1 | Store persists, mock renders | `make test && ./maat` |
| Phase 2 | Navigation works, < 100ms render | FR-001, FR-002 validation |
| Phase 3 | APIs fetch, adapt to graph | Integration tests |
| Phase 4 | Claude invokes, confirms | ADR-004 validation |
| Phase 5 | Plugin loads, fetches | Example plugin works |

### Failure Actions

```yaml
quality_below_threshold:
  - action: trigger_UNTIL_loop
  - max_iterations: 3
  - fallback: simplify_requirements

budget_exceeded:
  - action: reduce_complexity_level
  - options: [L5 → L4, L4 → L3, remove_feature]

time_overrun:
  - action: reassess_priorities
  - defer: [Phase 5, non-P0 features]
```

---

## Success Metrics

### MVP (Phase 1-2)

- ✅ SQLite backend functional
- ✅ Graph renders in TUI
- ✅ Keyboard navigation works
- ✅ Mock data populates graph
- ✅ Binary < 25MB
- ✅ Startup < 500ms

### V1.0 (Phase 1-4)

- ✅ Linear issues in graph
- ✅ GitHub PRs linked
- ✅ Git commits integrated
- ✅ Claude invokes via Ctrl+A
- ✅ Confirmation gates work
- ✅ All FR-001 to FR-007 met

### V1.1 (Phase 5)

- ✅ Plugin interface stable
- ✅ Example plugin works
- ✅ IDP commands functional
- ✅ FR-009 met

---

## Next Actions (Today)

```bash
# 1. Create workflows directory
mkdir -p workflows

# 2. Generate Phase 1A workflow file
cat > workflows/phase1-sqlite-store.yaml <<'EOF'
# [Copy Workflow 1A YAML from above]
EOF

# 3. Generate Phase 1B workflow file
cat > workflows/phase1-mock-data.yaml <<'EOF'
# [Copy Workflow 1B YAML from above]
EOF

# 4. Execute Phase 1A
/ois-compose --workflow-plan workflows/phase1-sqlite-store.yaml \
  --input-file internal/graph/schema.go \
  --output internal/graph/store_generated.go \
  --verbose

# 5. Execute Phase 1B
/ois-compose --workflow-plan workflows/phase1-mock-data.yaml \
  --input-file "docs/CONSTITUTION.md specs/FUNCTIONAL-REQUIREMENTS.md" \
  --output internal/tui/mock_data.go

# 6. Test Phase 1 complete
make test
./maat  # Should display mock graph
```

---

## Appendix: Command Reference

### Agent Composition

```bash
# Plan workflow
/ois-plan "implement feature X" --output workflow.yaml

# Execute workflow
/ois-compose --workflow-plan workflow.yaml --input-file code.go

# Dry run (preview)
/ois-compose --workflow-plan workflow.yaml --dry-run

# Verbose execution
/ois-compose --workflow-plan workflow.yaml --verbose
```

### Atomic Blocks

```bash
# List blocks
/blocks --list

# Show block info
/blocks --info structure-observer

# Compose blocks
/blocks "observer → generator" --input code.go --output impl.go
```

### Build & Test

```bash
# Build binary
make build

# Run tests
make test

# Run specific test
go test ./internal/graph -run TestStore

# Run binary
./maat
```

---

## Cross-References

- **Agent Composition Plan**: `docs/AGENT-COMPOSITION-PLAN.md`
- **Atomic Blocks**: `docs/ATOMIC-BLOCKS.md`
- **MAAT Spec**: `MAAT-SPEC.md`
- **Constitution**: `docs/CONSTITUTION.md`
- **Functional Requirements**: `specs/FUNCTIONAL-REQUIREMENTS.md`

---

**Status**: Ready to execute Phase 1 finish ✓
**Next Action**: Run Workflow 1A (SQLite Store)
**ETA**: Phase 1 complete in 2-3 days

