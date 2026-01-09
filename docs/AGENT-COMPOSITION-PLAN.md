# MAAT Agent Composition Plan: Progressive Spec Development

**Document Type**: Agent Orchestration Strategy
**Version**: 1.0
**Date**: 2026-01-06
**Status**: Active
**Philosophy**: Simple & Working → Complex & Broken

---

## Executive Summary

This document defines the **progressive agent composition strategy** for building MAAT specifications, following the OIS (Operator-based Integration System) taxonomy with 49 primitive agents and 6 composition operators.

**Core Principle**: Start with **simple, working implementations** (L1-L3) and progressively increase complexity only when needed, avoiding "complex and broken" outcomes.

---

## Current State Analysis

### Existing Specifications (3,051 lines)

| Document | Lines | Status | Quality |
|----------|-------|--------|---------|
| **MAAT-SPEC.md** | 497 | ✅ Complete | 0.85 (RMP validated) |
| **CONSTITUTION.md** | 326 | ✅ Complete | Constitutional |
| **FUNCTIONAL-REQUIREMENTS.md** | 465 | ✅ Complete | FR-001 to FR-010 |
| **ANTI-REQUIREMENTS.md** | ~300 | ✅ Complete | A-001 to A-008 |
| **META-SPEC.md** | 395 | ✅ Complete | Self-referential |
| **ADRs** | ~1,068 | ✅ Complete | 7 decisions |

**Total**: 3,051 lines of high-quality, RMP-validated specification

**Implementation Progress**: ~2,111 lines of Go code (Foundation Phase 1)

---

## Progressive Complexity Strategy

### Complexity Levels (L1-L7)

```
L1 (OBSERVE)  → Extract what exists
L2 (REASON)   → Analyze relationships
L3 (GENERATE) → Create new artifacts
L4 (REFINE)   → Improve quality
L5 (OPTIMIZE) → System-level tuning
L6 (INTEGRATE)→ Multi-system synthesis
L7 (REFLECT)  → Meta-level learning
```

**MAAT Philosophy**: Start at L1-L3, move to L4+ only when **demonstrably necessary**.

### The "Simple & Working" Gates

Before increasing complexity level, validate:

| Gate | Question | Evidence Required |
|------|----------|-------------------|
| **G1: Current Level Insufficient** | Does L(n) fail to meet requirements? | Failed test cases, user feedback |
| **G2: Next Level Necessary** | Will L(n+1) actually solve the problem? | Proof of concept, research |
| **G3: Risk Acceptable** | Can we afford L(n+1) breaking? | Rollback plan, budget |
| **G4: Value Justifiable** | Does benefit exceed complexity cost? | Cost-benefit analysis |

**Rule**: If ANY gate fails → Stay at current level, improve within constraints

---

## OIS Agent Taxonomy for MAAT

### 7 Foundation Skills × 7 Functions = 49 Agents

#### Foundation Skills

1. **abstraction-principles**: Generalization, pattern extraction
2. **systems-thinking**: Holistic view, emergence, feedback loops
3. **knowledge-synthesis**: Integration across domains
4. **specification-driven**: Constitutional governance, ADRs
5. **category-theory**: Type-safe composition, morphisms
6. **elm-architecture**: Pure functional UI, MVU pattern
7. **unix-philosophy**: Composition, simplicity, text interfaces

#### 7 Generic Functions

```
OBSERVE   → Data collection, inspection
REASON    → Analysis, inference
GENERATE  → Artifact creation
REFINE    → Quality improvement
OPTIMIZE  → Performance tuning
INTEGRATE → System synthesis
REFLECT   → Meta-learning
```

---

## Phase-Based Composition Workflows

### Phase 1: Foundation (L1-L2, Weeks 1-2)

**Goal**: SQLite graph backend, basic TUI structure

**Status**: ✅ 70% complete (2,111 lines Go code)

**Remaining Tasks**:
- SQLite store implementation (schema exists)
- Mock data injection for testing
- Basic viewport rendering

#### Workflow 1A: SQLite Store Implementation

**Complexity**: L1 (OBSERVE existing schema → GENERATE implementation)

**Composition**: `StructureObserver → TypeGenerator`

```yaml
task: "Implement SQLite store from schema.go"
complexity: L1
pattern: sequential

composition:
  step_1:
    agent: structure-observer
    foundation: systems-thinking
    function: OBSERVE
    input: internal/graph/schema.go
    output: Observable<Schema>
    budget: 800 tokens

  step_2:
    agent: type-generator
    foundation: category-theory
    function: GENERATE
    input: Observable<Schema>
    output: Generated<StoreImpl>
    budget: 1200 tokens

type_flow:
  - input: schema.go (existing types)
  - intermediate: Observable<Schema> (structure + operations)
  - output: store.go (CRUD implementation)

validation:
  - Implements all Node/Edge operations
  - Follows graph.Store interface
  - Pure functions, no global state (#1 Immutable Truth)

total_budget: 2000 tokens
expected_quality: 0.80+ (L1 baseline)
```

**Operators**: Sequential (`→`)

**Rationale**: Simple task, well-defined schema exists, straightforward generation.

---

#### Workflow 1B: Mock Data Generator

**Complexity**: L1 (GENERATE test fixtures)

**Composition**: `KnowledgeGatherer → MockGenerator`

```yaml
task: "Create mock issues/PRs/commits for TUI testing"
complexity: L1
pattern: sequential

composition:
  step_1:
    agent: knowledge-gatherer
    foundation: knowledge-synthesis
    function: OBSERVE
    input: docs/CONSTITUTION.md + specs/FUNCTIONAL-REQUIREMENTS.md
    output: Observable<DomainModel>
    budget: 600 tokens

  step_2:
    agent: mock-generator
    foundation: specification-driven
    function: GENERATE
    input: Observable<DomainModel>
    output: Generated<MockData>
    budget: 1000 tokens

type_flow:
  - input: Spec documents (domain understanding)
  - intermediate: Observable<DomainModel> (entity types)
  - output: mock_data.go (realistic fixtures)

validation:
  - Covers all NodeTypes (Issue, PR, Commit, File)
  - Realistic relationships (blocks, related, parent)
  - 50-100 nodes for testing

total_budget: 1600 tokens
expected_quality: 0.75+ (mock data quality)
```

**Operators**: Sequential (`→`)

**Rationale**: No real APIs needed yet, mock data enables rapid TUI iteration.

---

### Phase 2: Graph Navigation (L2-L3, Weeks 3-4)

**Goal**: Node/edge rendering, keyboard navigation, drill-down/back-up

**Status**: ⚙️ Ready to start (models exist, view.go partially implemented)

#### Workflow 2A: Graph Rendering Engine

**Complexity**: L2 (OBSERVE graph → REASON layout → GENERATE renderer)

**Composition**: `StructureObserver → LayoutReasoner → RenderGenerator`

```yaml
task: "Render knowledge graph in terminal with Bubble Tea"
complexity: L2
pattern: sequential

composition:
  step_1:
    agent: structure-observer
    foundation: systems-thinking
    function: OBSERVE
    input: internal/tui/model.go (DisplayNode, DisplayEdge)
    output: Observable<GraphStructure>
    budget: 1000 tokens

  step_2:
    agent: layout-reasoner
    foundation: abstraction-principles
    function: REASON
    input: Observable<GraphStructure>
    output: Reasoning<LayoutAlgorithm>
    budget: 1500 tokens

  step_3:
    agent: render-generator
    foundation: elm-architecture
    function: GENERATE
    input: Reasoning<LayoutAlgorithm>
    output: Generated<GraphView>
    budget: 2000 tokens

type_flow:
  - input: DisplayNode[] + DisplayEdge[] (data model)
  - intermediate: Reasoning<LayoutAlgorithm> (hierarchical tree layout)
  - output: views/graph.go (Lipgloss-styled renderer)

validation:
  - Renders up to 500 nodes in < 100ms (FR-001)
  - Hierarchical tree layout (parent-child)
  - vim-style navigation (h/j/k/l)
  - Respects Commandment #4 (10x Navigation)

total_budget: 4500 tokens
expected_quality: 0.82+ (core functionality)
```

**Operators**: Sequential (`→`)

**Rationale**: L2 reasoning required for layout algorithm (hierarchical vs force-directed), but implementation is straightforward generation.

---

#### Workflow 2B: Keyboard Navigation Logic

**Complexity**: L2 (OBSERVE key patterns → GENERATE handlers)

**Composition**: `PatternObserver → NavigationGenerator`

```yaml
task: "Implement vim-style navigation (hjkl, Enter, Esc)"
complexity: L2
pattern: sequential

composition:
  step_1:
    agent: pattern-observer
    foundation: unix-philosophy
    function: OBSERVE
    input: specs/FUNCTIONAL-REQUIREMENTS.md (FR-002)
    output: Observable<KeyBindings>
    budget: 800 tokens

  step_2:
    agent: navigation-generator
    foundation: elm-architecture
    function: GENERATE
    input: Observable<KeyBindings>
    output: Generated<KeyHandlers>
    budget: 1800 tokens

type_flow:
  - input: FR-002 specification (keyboard requirements)
  - intermediate: Observable<KeyBindings> (key → action mapping)
  - output: internal/tui/keys.go (pure key handlers)

validation:
  - All FR-002 acceptance criteria met
  - Pure functions (Commandment #1)
  - NavigationStack for drill-down/back-up
  - No goroutines (Commandment #5)

total_budget: 2600 tokens
expected_quality: 0.85+ (straightforward implementation)
```

**Operators**: Sequential (`→`)

**Rationale**: Pattern observation (L2) needed to extract all keybindings from spec, then simple generation.

---

### Phase 3: Integrations (L3-L4, Weeks 5-7)

**Goal**: Linear GraphQL, GitHub REST, Git commit history

**Status**: 🔜 Depends on Phase 2 completion

**Complexity Increase Rationale**: External APIs introduce **error handling, rate limiting, auth** → L3-L4 required

#### Workflow 3A: Linear GraphQL Client (L3)

**Complexity**: L3 (OBSERVE API → REASON schema → GENERATE client)

**Composition**: `APIObserver → SchemaReasoner → ClientGenerator`

```yaml
task: "Create thin Linear GraphQL client"
complexity: L3
pattern: sequential

composition:
  step_1:
    agent: api-observer
    foundation: systems-thinking
    function: OBSERVE
    input: Linear GraphQL API docs (Context7)
    output: Observable<APISchema>
    budget: 1500 tokens

  step_2:
    agent: schema-reasoner
    foundation: category-theory
    function: REASON
    input: Observable<APISchema>
    output: Reasoning<ClientDesign>
    budget: 2000 tokens

  step_3:
    agent: client-generator
    foundation: specification-driven
    function: GENERATE
    input: Reasoning<ClientDesign>
    output: Generated<LinearClient>
    budget: 2500 tokens

type_flow:
  - input: Linear API documentation
  - intermediate: Reasoning<ClientDesign> (thin client, error handling)
  - output: internal/linear/client.go + adapter.go

validation:
  - Thin client only (Commandment #7 Composition)
  - Read operations only (Commandment #10 Sovereignty)
  - Converts to graph.Node (adapter pattern)
  - Respects rate limits

total_budget: 6000 tokens
expected_quality: 0.80+ (API integration complexity)
```

**Operators**: Sequential (`→`)

**Rationale**: L3 generation for API clients (error handling, retries, auth).

---

#### Workflow 3B: Git History Integration (L2-L3)

**Complexity**: L2-L3 (OBSERVE git history → GENERATE graph nodes)

**Composition**: `GitObserver → CommitGenerator`

```yaml
task: "Parse git commits into knowledge graph nodes"
complexity: L2
pattern: sequential

composition:
  step_1:
    agent: git-observer
    foundation: unix-philosophy
    function: OBSERVE
    input: .git directory (go-git library)
    output: Observable<CommitHistory>
    budget: 1200 tokens

  step_2:
    agent: commit-generator
    foundation: systems-thinking
    function: GENERATE
    input: Observable<CommitHistory>
    output: Generated<CommitNodes>
    budget: 1500 tokens

type_flow:
  - input: .git/objects (git repository)
  - intermediate: Observable<CommitHistory> (parsed commits)
  - output: internal/git/adapter.go (graph.Node conversions)

validation:
  - Parses commit messages for issue references (#123)
  - Creates edges: Commit → Issue (mentions)
  - Efficient (< 2s for 1000 commits)

total_budget: 2700 tokens
expected_quality: 0.82+ (local git, no API complexity)
```

**Operators**: Sequential (`→`)

**Rationale**: L2 sufficient for local git operations (no external API).

---

### Phase 4: AI & Polish (L4-L5, Weeks 8-10)

**Goal**: Claude MCP integration, role-based views, fuzzy search

**Status**: 🔮 Future (depends on Phase 3)

**Complexity Increase Rationale**: Human-in-loop AI requires **explicit invocation, confirmation gates, audit trails** → L4-L5

#### Workflow 4A: Claude MCP Bridge (L4)

**Complexity**: L4 (OBSERVE context → REASON assembly → GENERATE → REFINE)

**Composition**: `ContextObserver || KnowledgeGatherer → ContextReasoner → PromptGenerator → HumanRefiner`

```yaml
task: "Claude integration with human-in-loop gates"
complexity: L4
pattern: parallel_synthesis_refine

composition:
  stage_1_parallel:
    operator: ||
    agents:
      - context-observer:
          foundation: systems-thinking
          function: OBSERVE
          input: Current focused node + edges
          output: Observable<NodeContext>
          budget: 1200

      - knowledge-gatherer:
          foundation: knowledge-synthesis
          function: OBSERVE
          input: Navigation history + breadcrumb
          output: Observable<SessionContext>
          budget: 1000

  stage_2_synthesis:
    operator: →
    agent: context-reasoner
    foundation: specification-driven
    function: REASON
    input: (Observable<NodeContext>, Observable<SessionContext>)
    output: Reasoning<ClaudePrompt>
    budget: 2500

  stage_3_generate:
    operator: →
    agent: prompt-generator
    foundation: category-theory
    function: GENERATE
    input: Reasoning<ClaudePrompt>
    output: Generated<MCPRequest>
    budget: 2000

  stage_4_refine:
    operator: →
    agent: human-refiner
    foundation: unix-philosophy
    function: REFINE
    input: Generated<MCPRequest>
    output: Refined<ConfirmedRequest>
    budget: 800

type_flow:
  - parallel: Observable<NodeContext> || Observable<SessionContext>
  - synthesis: Reasoning<ClaudePrompt> (combined context)
  - generation: Generated<MCPRequest> (MCP protocol message)
  - refinement: Refined<ConfirmedRequest> (human approval)

validation:
  - Explicit invocation only (Commandment #6 Human Contact)
  - Confirmation gate for all actions (Commandment #10 Sovereignty)
  - Full audit trail (ADR-004)
  - No ambient suggestions (AFR-005)

total_budget: 7500 tokens
expected_quality: 0.80+ (L4 human-AI coordination)
```

**Operators**: Parallel (`||`), Sequential (`→`)

**Rationale**: L4 required for human-in-loop refinement, parallel observation for efficiency.

**Key Insight**: This is the FIRST workflow using parallel composition (`||`) because context gathering (node + session) is independent.

---

### Phase 5: Extensibility (L5-L6, Weeks 11-12)

**Goal**: Plugin architecture, IDP self-service

**Status**: 🌌 Advanced (requires all previous phases)

**Complexity Increase Rationale**: Plugin system requires **type-safe composition, interface validation, lifecycle management** → L5-L6

#### Workflow 5A: Plugin Interface Design (L5)

**Complexity**: L5 (OBSERVE patterns → REASON architecture → GENERATE → OPTIMIZE)

**Composition**: `PatternObserver || ArchitectureReasoner → InterfaceGenerator → SystemOptimizer`

```yaml
task: "Design plugin interface for extensibility"
complexity: L5
pattern: parallel_reasoning_optimize

composition:
  stage_1_parallel:
    operator: ||
    agents:
      - pattern-observer:
          foundation: abstraction-principles
          function: OBSERVE
          input: internal/linear + internal/github + internal/git
          output: Observable<ClientPatterns>
          budget: 1800

      - architecture-reasoner:
          foundation: systems-thinking
          function: REASON
          input: specs/ADR/ADR-007-plugin-architecture.md
          output: Reasoning<PluginDesign>
          budget: 2200

  stage_2_generate:
    operator: →
    agent: interface-generator
    foundation: category-theory
    function: GENERATE
    input: (Observable<ClientPatterns>, Reasoning<PluginDesign>)
    output: Generated<PluginInterface>
    budget: 2500

  stage_3_optimize:
    operator: →
    agent: system-optimizer
    foundation: unix-philosophy
    function: OPTIMIZE
    input: Generated<PluginInterface>
    output: Optimized<PluginAPI>
    budget: 2000

type_flow:
  - parallel: Observable<ClientPatterns> || Reasoning<PluginDesign>
  - generation: Generated<PluginInterface> (Go interfaces)
  - optimization: Optimized<PluginAPI> (minimized API surface)

validation:
  - Type-safe plugin registration
  - Lifecycle hooks (init, fetch, transform, cleanup)
  - Respects Single Responsibility (Commandment #2)
  - Thin integration boundary (Commandment #7)

total_budget: 8500 tokens
expected_quality: 0.78+ (L5 system-level design)
```

**Operators**: Parallel (`||`), Sequential (`→`)

**Rationale**: L5 optimization needed for minimal API surface (Unix philosophy: small, sharp interfaces).

---

## Composition Operator Usage Matrix

| Operator | Phase 1 | Phase 2 | Phase 3 | Phase 4 | Phase 5 |
|----------|---------|---------|---------|---------|---------|
| **Sequential (→)** | ✅✅✅ | ✅✅✅ | ✅✅✅ | ✅✅ | ✅✅ |
| **Parallel (\|\|)** | ❌ | ❌ | ❌ | ✅ | ✅ |
| **Product (×)** | ❌ | ❌ | ❌ | ❌ | ❌ |
| **Tensor (⊗)** | ❌ | ❌ | ❌ | ❌ | ❌ |
| **Conditional (IF)** | ❌ | ❌ | ❌ | ❌ | ❌ |
| **Recursive (UNTIL)** | ❌ | ❌ | ❌ | ❌ | ❌ |

**Observation**: MAAT spec development uses **primarily sequential composition** through Phase 3, introducing **parallel** only in Phase 4 (AI context gathering).

**Philosophy**: Simple compositions (`→`) work until proven insufficient. Complex operators (`⊗`, `IF`, `UNTIL`) reserved for future needs.

---

## Token Budget Analysis

### Phase-Level Budgets

| Phase | Workflows | Total Agents | Est. Tokens | Duration |
|-------|-----------|--------------|-------------|----------|
| **Phase 1** | 2 (1A, 1B) | 4 | 3,600 | 2-3 days |
| **Phase 2** | 2 (2A, 2B) | 5 | 7,100 | 4-5 days |
| **Phase 3** | 2 (3A, 3B) | 5 | 8,700 | 5-7 days |
| **Phase 4** | 1 (4A) | 4 | 7,500 | 3-4 days |
| **Phase 5** | 1 (5A) | 4 | 8,500 | 4-5 days |
| **Total** | 8 workflows | 22 agents | 35,400 | 18-24 days |

**Budget per Workflow**: 3,600 - 8,700 tokens (avg: 4,425 tokens)

**Rationale**: Token budgets increase with complexity level:
- L1-L2: 2,000-4,500 tokens (observation + generation)
- L3-L4: 6,000-7,500 tokens (API integration + refinement)
- L5-L6: 8,000-10,000 tokens (system optimization)

---

## Quality Thresholds by Phase

| Phase | Min Quality | Target Quality | Rationale |
|-------|-------------|----------------|-----------|
| Phase 1 | 0.75 | 0.80 | Foundation must be solid |
| Phase 2 | 0.80 | 0.85 | Core UX, user-facing |
| Phase 3 | 0.75 | 0.82 | API integration (external dependency risk) |
| Phase 4 | 0.78 | 0.85 | Human-in-loop critical |
| Phase 5 | 0.75 | 0.80 | Extensibility, advanced |

**Note**: Quality < 0.75 triggers UNTIL loop (iterative refinement)

---

## Anti-Patterns to Avoid

### ❌ Complex & Broken

**Symptom**: Jumping to L5-L7 workflows without validating L1-L3

**Example**:
```yaml
# ❌ BAD: Parallel + UNTIL loop for simple mock data generation
composition:
  stage_1_parallel:
    operator: ||
    agents: [MockGeneratorA, MockGeneratorB, MockGeneratorC]

  stage_2_refine:
    operator: UNTIL
    quality_threshold: 0.90
    max_iterations: 10
    agent: MockRefiner

# Token overhead: 15,000+ tokens for 1,000 token task
```

**Fix**: Use simple sequential generation (Workflow 1B: 1,600 tokens)

---

### ❌ Premature Optimization

**Symptom**: Using L5 OPTIMIZE before L3 GENERATE works

**Example**:
```yaml
# ❌ BAD: Optimizing graph renderer before it renders anything
composition:
  - graph-renderer-optimizer  # L5 OPTIMIZE
  - graph-renderer-generator  # L3 GENERATE (should be first!)
```

**Fix**: Generate working implementation first, optimize ONLY if performance fails validation.

---

### ❌ Parallel Overkill

**Symptom**: Using parallel (`||`) for dependent operations

**Example**:
```yaml
# ❌ BAD: Schema observation and store implementation in parallel
composition:
  operator: ||
  agents:
    - structure-observer   # Must complete first!
    - type-generator       # Depends on observer output
```

**Fix**: Use sequential (`→`) for dependent operations.

---

## Progressive Disclosure Strategy

### When to Increase Complexity

**Trigger Checklist**:

1. ✅ **Current level implemented and tested**
2. ✅ **Functional requirements unmet despite iteration**
3. ✅ **Higher complexity provably necessary** (not speculative)
4. ✅ **Budget available** (tokens + time)
5. ✅ **Rollback plan exists**

**If ANY checkbox fails** → Stay at current level

---

### Complexity Increase Decision Tree

```
Current level working?
├─ YES → Ship it, move to next feature at same level
└─ NO
   └─ Have you tried 3 iterations at current level?
      ├─ NO → Iterate (UNTIL loop)
      └─ YES
         └─ Is higher complexity necessary (not just desirable)?
            ├─ NO → Simplify requirements or accept limitation
            └─ YES
               └─ Can you afford the risk?
                  ├─ NO → Stay at current level
                  └─ YES → Increase complexity with caution
```

---

## Execution Commands

### Phase 1: Foundation

```bash
# Workflow 1A: SQLite Store
/ois-compose --workflow-plan phase1-sqlite-store.yaml \
  --input-file internal/graph/schema.go \
  --output internal/graph/store_generated.go

# Workflow 1B: Mock Data
/ois-compose --workflow-plan phase1-mock-data.yaml \
  --input-file "docs/CONSTITUTION.md specs/FUNCTIONAL-REQUIREMENTS.md" \
  --output internal/tui/mock_data.go
```

### Phase 2: Graph Navigation

```bash
# Workflow 2A: Graph Rendering
/ois-compose --workflow-plan phase2-graph-render.yaml \
  --input-file internal/tui/model.go \
  --output internal/tui/views/graph.go

# Workflow 2B: Keyboard Navigation
/ois-compose --workflow-plan phase2-keyboard-nav.yaml \
  --input-file specs/FUNCTIONAL-REQUIREMENTS.md \
  --output internal/tui/keys_generated.go
```

---

## Success Metrics

### Per-Phase Validation

| Phase | Success Criteria |
|-------|------------------|
| **Phase 1** | • SQLite stores nodes/edges<br>• Mock data renders in TUI<br>• Binary < 25MB |
| **Phase 2** | • Graph renders 500 nodes < 100ms<br>• hjkl navigation works<br>• Enter/Esc drill-down functional |
| **Phase 3** | • Linear issues populate graph<br>• GitHub PRs linked<br>• Git commits integrated |
| **Phase 4** | • Claude invokes via Ctrl+A<br>• Confirmation gates work<br>• Audit trail complete |
| **Phase 5** | • Plugin interface stable<br>• Example plugin works<br>• IDP commands functional |

### Quality Gates

```yaml
minimum_quality_per_phase:
  phase_1: 0.75
  phase_2: 0.80
  phase_3: 0.75
  phase_4: 0.78
  phase_5: 0.75

failure_action: "Trigger UNTIL loop (max 3 iterations)"
```

---

## Appendix: Complete Agent List

### 22 Agents Used Across All Phases

| Agent | Foundation | Function | Phases |
|-------|------------|----------|--------|
| structure-observer | systems-thinking | OBSERVE | 1, 2, 5 |
| type-generator | category-theory | GENERATE | 1 |
| knowledge-gatherer | knowledge-synthesis | OBSERVE | 1, 4 |
| mock-generator | specification-driven | GENERATE | 1 |
| layout-reasoner | abstraction-principles | REASON | 2 |
| render-generator | elm-architecture | GENERATE | 2 |
| pattern-observer | unix-philosophy | OBSERVE | 2, 5 |
| navigation-generator | elm-architecture | GENERATE | 2 |
| api-observer | systems-thinking | OBSERVE | 3 |
| schema-reasoner | category-theory | REASON | 3 |
| client-generator | specification-driven | GENERATE | 3 |
| git-observer | unix-philosophy | OBSERVE | 3 |
| commit-generator | systems-thinking | GENERATE | 3 |
| context-observer | systems-thinking | OBSERVE | 4 |
| context-reasoner | specification-driven | REASON | 4 |
| prompt-generator | category-theory | GENERATE | 4 |
| human-refiner | unix-philosophy | REFINE | 4 |
| architecture-reasoner | systems-thinking | REASON | 5 |
| interface-generator | category-theory | GENERATE | 5 |
| system-optimizer | unix-philosophy | OPTIMIZE | 5 |

**Total**: 20 unique agents (some reused across phases)

---

## Cross-References

- **OIS Taxonomy**: `~/.claude/skills/ois-taxonomy.md`
- **OIS Plan**: `/ois-plan` command
- **OIS Compose**: `/ois-compose` command
- **MAAT Constitution**: `docs/CONSTITUTION.md`
- **Functional Requirements**: `specs/FUNCTIONAL-REQUIREMENTS.md`
- **ADRs**: `specs/ADR/`

---

**Status**: Ready for Phase 1 execution ✓
**Philosophy**: Simple & Working beats Complex & Broken
**Next Action**: Execute Workflow 1A (SQLite Store)

