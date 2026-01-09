# MAAT Atomic Blocks: Compositional Spec Development

**Document Type**: Atomic Block Definitions
**Version**: 1.0
**Date**: 2026-01-06
**Philosophy**: `/blocks` - Compose atomic primitives into workflows

---

## Purpose

This document defines **atomic blocks** (indivisible units of agent behavior) for MAAT spec development. Each block is a **pure function** with explicit input/output types, enabling type-safe composition via `/blocks` and `/ois-compose`.

**Atomic = Cannot be decomposed further without losing semantic meaning**

---

## Block Taxonomy

### 7 Foundation Skills × 7 Functions = 49 Potential Blocks

**Active Blocks for MAAT**: 20 blocks (subset of 49 based on actual usage)

---

## Foundation: Abstraction Principles

### AP-OBSERVE: Pattern Observer

```yaml
block_id: pattern-observer
foundation: abstraction-principles
function: OBSERVE
description: "Extract recurring patterns, abstractions, and design structures"

input_type: Observable<Code>
output_type: Observable<Patterns>

signature: |
  patternObserver: Code → Observable<Patterns>
  where Patterns = { syntactic, semantic, architectural }

example_input: |
  internal/tui/model.go (258 lines)

example_output: |
  Observable<Patterns>:
    - syntactic: WithX transformation methods (12 found)
    - semantic: Immutable state pattern (Commandment #1)
    - architectural: Navigation stack pattern

budget: 1200-1800 tokens
complexity: L1 (OBSERVE)
```

---

### AP-REASON: Layout Reasoner

```yaml
block_id: layout-reasoner
foundation: abstraction-principles
function: REASON
description: "Reason about optimal layout algorithms for graph visualization"

input_type: Observable<GraphStructure>
output_type: Reasoning<LayoutAlgorithm>

signature: |
  layoutReasoner: Observable<GraphStructure> → Reasoning<LayoutAlgorithm>
  where LayoutAlgorithm = { hierarchical, force-directed, tree }

example_input: |
  Observable<GraphStructure>:
    nodes: 250 (DisplayNode with parent references)
    edges: 380 (DisplayEdge with relationship types)
    structure: hierarchical (issue → PR → commit)

example_output: |
  Reasoning<LayoutAlgorithm>:
    chosen: hierarchical_tree
    rationale: "Parent-child relationships dominate (78%), force-directed overkill"
    properties:
      - Top-down layout (root at top)
      - Indent child nodes
      - Collapse/expand for large subtrees

budget: 1500-2000 tokens
complexity: L2 (REASON)
```

---

## Foundation: Systems Thinking

### ST-OBSERVE: Structure Observer

```yaml
block_id: structure-observer
foundation: systems-thinking
function: OBSERVE
description: "Inspect code structure, extract types, interfaces, relationships"

input_type: Observable<Code>
output_type: Observable<Structure>

signature: |
  structureObserver: Code → Observable<Structure>
  where Structure = { types, interfaces, dependencies, patterns }

example_input: |
  internal/graph/schema.go (170 lines)

example_output: |
  Observable<Structure>:
    types:
      - NodeType (enum: Issue, PR, Commit, File, Project, Service)
      - EdgeType (enum: blocks, related, implements, calls, owns, modifies)
      - Node (struct with ID, Type, Source, Data, Metadata)
      - Edge (struct with FromID, ToID, Relation, Metadata)
    interfaces:
      - Store (CRUD operations for nodes/edges)
    dependencies:
      - encoding/json (for Data field)
      - time (for Metadata timestamps)

budget: 800-1200 tokens
complexity: L1 (OBSERVE)
```

---

### ST-REASON: Architecture Reasoner

```yaml
block_id: architecture-reasoner
foundation: systems-thinking
function: REASON
description: "Reason about system architecture, emergence, feedback loops"

input_type: Observable<Specification>
output_type: Reasoning<Architecture>

signature: |
  architectureReasoner: Observable<Specification> → Reasoning<Architecture>
  where Architecture = { components, boundaries, interactions }

example_input: |
  specs/ADR/ADR-007-plugin-architecture.md

example_output: |
  Reasoning<Architecture>:
    pattern: plugin_registry
    components:
      - Registry (central registry with lifecycle management)
      - PluginInterface (DataSourcePlugin with 5 methods)
      - BuiltinPlugins (linear, github, git as reference impls)
    boundaries:
      - Plugins never access each other directly
      - All communication through graph.Node/Edge types
    feedback_loops:
      - Plugin registration validates interface compliance
      - Failed plugins logged but don't crash system

budget: 2000-2500 tokens
complexity: L2-L3 (REASON with architecture)
```

---

### ST-GENERATE: Commit Generator

```yaml
block_id: commit-generator
foundation: systems-thinking
function: GENERATE
description: "Generate graph nodes from git commit history"

input_type: Observable<CommitHistory>
output_type: Generated<CommitNodes>

signature: |
  commitGenerator: Observable<CommitHistory> → Generated<CommitNodes>
  where CommitNodes = graph.Node[] + graph.Edge[]

example_input: |
  Observable<CommitHistory>:
    commits: [
      { hash: "abc123", message: "Fix #42: Update graph schema", author: "dev" },
      { hash: "def456", message: "Add Linear integration", author: "dev" }
    ]

example_output: |
  Generated<CommitNodes>:
    nodes: [
      { id: "commit:abc123", type: Commit, data: { message: "Fix #42: ...", author: "dev" } },
      { id: "commit:def456", type: Commit, data: { message: "Add Linear...", author: "dev" } }
    ]
    edges: [
      { from: "commit:abc123", to: "issue:42", relation: mentions }
    ]

budget: 1500-2000 tokens
complexity: L2-L3 (GENERATE with parsing)
```

---

## Foundation: Knowledge Synthesis

### KS-OBSERVE: Knowledge Gatherer

```yaml
block_id: knowledge-gatherer
foundation: knowledge-synthesis
function: OBSERVE
description: "Collect and catalog domain knowledge from documentation"

input_type: Observable<Documents>
output_type: Observable<DomainModel>

signature: |
  knowledgeGatherer: Documents → Observable<DomainModel>
  where DomainModel = { entities, relationships, constraints }

example_input: |
  docs/CONSTITUTION.md + specs/FUNCTIONAL-REQUIREMENTS.md

example_output: |
  Observable<DomainModel>:
    entities:
      - Issue (from Linear, with status/priority/labels)
      - PR (from GitHub, with number/author/status)
      - Commit (from Git, with hash/message/author)
      - File (from repository, with path/language/lines)
    relationships:
      - Issue blocks Issue
      - Issue implemented_by PR
      - PR modifies File
      - Commit mentions Issue
    constraints:
      - No mutable state (Commandment #1)
      - Human confirmation for writes (Commandment #10)

budget: 600-1000 tokens
complexity: L1 (OBSERVE)
```

---

### KS-REASON: Context Reasoner

```yaml
block_id: context-reasoner
foundation: knowledge-synthesis
function: REASON
description: "Synthesize context from multiple observations for Claude prompts"

input_type: (Observable<NodeContext>, Observable<SessionContext>)
output_type: Reasoning<ClaudePrompt>

signature: |
  contextReasoner: (NodeContext, SessionContext) → Reasoning<ClaudePrompt>
  where ClaudePrompt = { context, question, constraints }

example_input: |
  NodeContext: {
    focused: "issue:LIN-42 (Implement graph rendering)",
    edges: ["blocks issue:LIN-41", "implemented_by pr:123"]
  }
  SessionContext: {
    breadcrumb: ["Project > Backend > Issues > LIN-42"],
    recent_actions: ["viewed PR #123", "checked commit abc123"]
  }

example_output: |
  Reasoning<ClaudePrompt>:
    context: |
      User is viewing Linear issue LIN-42: "Implement graph rendering"
      This issue blocks LIN-41 and is implemented in PR #123.
      Recent navigation: viewed PR #123, checked commit abc123.
    question: "What should I work on next for this issue?"
    constraints:
      - Respect Commandment #6 (explicit invocation)
      - Include graph relationships in suggestions
      - Reference recent commits for context

budget: 2000-2500 tokens
complexity: L3-L4 (REASON with synthesis)
```

---

## Foundation: Specification-Driven

### SD-OBSERVE: Pattern Observer (Spec)

```yaml
block_id: pattern-observer-spec
foundation: specification-driven
function: OBSERVE
description: "Extract patterns and requirements from specifications"

input_type: Observable<Specification>
output_type: Observable<Requirements>

signature: |
  patternObserverSpec: Specification → Observable<Requirements>
  where Requirements = { functional, constraints, anti_patterns }

example_input: |
  specs/FUNCTIONAL-REQUIREMENTS.md (FR-002: Keyboard Navigation)

example_output: |
  Observable<Requirements>:
    functional:
      - hjkl or arrows navigate between nodes
      - Enter drills into focused node
      - Esc/Backspace returns to parent context
      - Tab cycles between panes
      - / opens search
      - ? shows help
      - q quits
    constraints:
      - All keybindings configurable via ~/.config/maat/keys.json
      - Response time < 50ms (Commandment #9)
    anti_patterns:
      - NO mouse dependency
      - NO blocking operations in key handlers

budget: 800-1200 tokens
complexity: L1 (OBSERVE)
```

---

### SD-GENERATE: Mock Generator

```yaml
block_id: mock-generator
foundation: specification-driven
function: GENERATE
description: "Generate realistic mock data from domain model specifications"

input_type: Observable<DomainModel>
output_type: Generated<MockData>

signature: |
  mockGenerator: Observable<DomainModel> → Generated<MockData>
  where MockData = { nodes: graph.Node[], edges: graph.Edge[] }

example_input: |
  Observable<DomainModel>:
    entities: [Issue, PR, Commit, File]
    relationships: [blocks, related, implements, modifies, mentions]

example_output: |
  Generated<MockData>:
    nodes: [
      { id: "issue:1", type: Issue, data: { title: "Setup project", status: "done" } },
      { id: "issue:2", type: Issue, data: { title: "Implement graph", status: "in_progress" } },
      { id: "pr:101", type: PR, data: { title: "Add graph rendering", number: 101 } },
      { id: "commit:abc", type: Commit, data: { message: "Initial commit" } }
    ]
    edges: [
      { from: "issue:2", to: "issue:1", relation: blocks },
      { from: "pr:101", to: "issue:2", relation: implements }
    ]
    count: 50 nodes, 75 edges

budget: 1000-1500 tokens
complexity: L2 (GENERATE with realistic data)
```

---

### SD-REASON: Schema Reasoner

```yaml
block_id: schema-reasoner
foundation: specification-driven
function: REASON
description: "Reason about API schemas and design thin client interfaces"

input_type: Observable<APISchema>
output_type: Reasoning<ClientDesign>

signature: |
  schemaReasoner: Observable<APISchema> → Reasoning<ClientDesign>
  where ClientDesign = { operations, error_handling, adapters }

example_input: |
  Observable<APISchema>:
    api: Linear GraphQL
    operations: [fetchIssues, fetchProjects, updateIssueStatus]
    types: [Issue, Project, WorkflowState]

example_output: |
  Reasoning<ClientDesign>:
    pattern: thin_client
    operations:
      - FetchIssues(teamID) → []Issue (read-only, respects #7 Composition)
      - FetchProjects(teamID) → []Project (read-only)
      - UpdateIssueStatus(id, status) → Confirmation (write, requires confirm per #10)
    error_handling:
      - Rate limiting: exponential backoff
      - Network errors: retry with timeout
      - Auth errors: surface to user
    adapters:
      - IssueToNode: Linear.Issue → graph.Node
      - ProjectToNode: Linear.Project → graph.Node

budget: 2000-2500 tokens
complexity: L3 (REASON with API design)
```

---

## Foundation: Category Theory

### CT-GENERATE: Type Generator

```yaml
block_id: type-generator
foundation: category-theory
function: GENERATE
description: "Generate type-safe implementations from schema specifications"

input_type: Observable<Schema>
output_type: Generated<Implementation>

signature: |
  typeGenerator: Observable<Schema> → Generated<Implementation>
  where Implementation = { types, functions, tests }

example_input: |
  Observable<Schema>:
    source: internal/graph/schema.go
    types: [Node, Edge, NodeType, EdgeType]
    operations: [AddNode, AddEdge, GetNode, GetEdges, GetNeighbors]

example_output: |
  Generated<Implementation>:
    file: internal/graph/store.go
    types: [Store interface, sqliteStore struct]
    functions:
      - AddNode(node Node) error
      - AddEdge(edge Edge) error
      - GetNode(id string) (Node, error)
      - GetEdges(fromID string) ([]Edge, error)
      - GetNeighbors(nodeID string) ([]Node, error)
    tests: store_test.go (table-driven tests for all operations)

budget: 1200-2000 tokens
complexity: L2-L3 (GENERATE with type safety)
```

---

### CT-GENERATE: Interface Generator

```yaml
block_id: interface-generator
foundation: category-theory
function: GENERATE
description: "Generate plugin interfaces with type-safe composition"

input_type: (Observable<ClientPatterns>, Reasoning<PluginDesign>)
output_type: Generated<PluginInterface>

signature: |
  interfaceGenerator: (ClientPatterns, PluginDesign) → Generated<PluginInterface>
  where PluginInterface = { interface, lifecycle, validation }

example_input: |
  ClientPatterns: [linear.Client, github.Client, git.Repository]
  PluginDesign: { pattern: registry, lifecycle: [init, fetch, cleanup] }

example_output: |
  Generated<PluginInterface>:
    interface: |
      type DataSourcePlugin interface {
        Name() string
        NodeTypes() []graph.NodeType
        EdgeTypes() []graph.EdgeType
        Fetch(ctx context.Context, query Query) ([]graph.Node, []graph.Edge, error)
        Transform(raw interface{}) (graph.Node, error)
      }
    lifecycle:
      - Register(plugin DataSourcePlugin) → Validates interface compliance
      - Init() → Setup (config, auth)
      - Fetch() → Data retrieval
      - Cleanup() → Resource release
    validation:
      - Type-safe registration
      - Panic recovery
      - Resource limits

budget: 2500-3000 tokens
complexity: L4-L5 (GENERATE with system design)
```

---

### CT-GENERATE: Prompt Generator

```yaml
block_id: prompt-generator
foundation: category-theory
function: GENERATE
description: "Generate type-safe MCP protocol messages for Claude"

input_type: Reasoning<ClaudePrompt>
output_type: Generated<MCPRequest>

signature: |
  promptGenerator: Reasoning<ClaudePrompt> → Generated<MCPRequest>
  where MCPRequest = { protocol, context, tools, validation }

example_input: |
  Reasoning<ClaudePrompt>:
    context: "Issue LIN-42 blocks LIN-41, implemented in PR #123"
    question: "What should I work on next?"
    constraints: ["explicit invocation", "graph relationships", "recent commits"]

example_output: |
  Generated<MCPRequest>:
    protocol: "mcp/v1"
    request:
      context: |
        Current Focus: issue:LIN-42 "Implement graph rendering"
        Blocking: issue:LIN-41
        Implementation: pr:123 (reviewed, needs tests)
        Recent: commit:abc123 "Add basic graph structure"
      question: "What should I work on next for this issue?"
      tools: ["code_search", "git_history", "issue_lookup"]
    validation:
      - Context < 4000 tokens
      - Tools available in MCP
      - Request respects Commandment #6

budget: 2000-2500 tokens
complexity: L3 (GENERATE with protocol)
```

---

## Foundation: Elm Architecture

### EA-GENERATE: Render Generator

```yaml
block_id: render-generator
foundation: elm-architecture
function: GENERATE
description: "Generate Bubble Tea View functions with Lipgloss styling"

input_type: Reasoning<LayoutAlgorithm>
output_type: Generated<ViewFunction>

signature: |
  renderGenerator: Reasoning<LayoutAlgorithm> → Generated<ViewFunction>
  where ViewFunction = { view, styles, helpers }

example_input: |
  Reasoning<LayoutAlgorithm>:
    chosen: hierarchical_tree
    properties: [top-down, indented, collapsible]

example_output: |
  Generated<ViewFunction>:
    view: |
      func (m Model) renderGraph() string {
        var lines []string
        for _, node := range m.nodes {
          indent := strings.Repeat("  ", node.Depth)
          icon := nodeIcon(node.Type)
          title := styles.NodeTitle.Render(node.Title)
          lines = append(lines, indent + icon + " " + title)
        }
        return lipgloss.JoinVertical(lipgloss.Left, lines...)
      }
    styles: |
      NodeTitle = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
      FocusedNode = lipgloss.NewStyle().Background(lipgloss.Color("8"))
    helpers: |
      func nodeIcon(t graph.NodeType) string { ... }

budget: 2000-2500 tokens
complexity: L3 (GENERATE with UI)
```

---

### EA-GENERATE: Navigation Generator

```yaml
block_id: navigation-generator
foundation: elm-architecture
function: GENERATE
description: "Generate pure key handler functions (Model, Msg) → (Model, Cmd)"

input_type: Observable<KeyBindings>
output_type: Generated<KeyHandlers>

signature: |
  navigationGenerator: Observable<KeyBindings> → Generated<KeyHandlers>
  where KeyHandlers = { update, commands, stack }

example_input: |
  Observable<KeyBindings>:
    - { key: "j", action: "move_down" }
    - { key: "k", action: "move_up" }
    - { key: "Enter", action: "drill_down" }
    - { key: "Esc", action: "back_up" }

example_output: |
  Generated<KeyHandlers>:
    update: |
      func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
        switch msg := msg.(type) {
        case tea.KeyMsg:
          switch msg.String() {
          case "j": return m.WithFocusedNode(m.NextNode()), nil
          case "k": return m.WithFocusedNode(m.PrevNode()), nil
          case "enter": return m.PushView(DetailView), nil
          case "esc": return m.PopView(), nil
          }
        }
        return m, nil
      }
    stack: NavigationStack (for drill-down/back-up)
    pure: All functions return new Model, never mutate

budget: 1800-2200 tokens
complexity: L2-L3 (GENERATE with pure functions)
```

---

## Foundation: Unix Philosophy

### UP-OBSERVE: Git Observer

```yaml
block_id: git-observer
foundation: unix-philosophy
function: OBSERVE
description: "Observe git repository history using go-git library"

input_type: Observable<Repository>
output_type: Observable<CommitHistory>

signature: |
  gitObserver: Repository → Observable<CommitHistory>
  where CommitHistory = { commits: Commit[], refs: Ref[] }

example_input: |
  .git directory (go-git opens repository)

example_output: |
  Observable<CommitHistory>:
    commits: [
      { hash: "abc123", message: "Fix #42: Update schema", author: "dev", date: "2026-01-05" },
      { hash: "def456", message: "Add Linear client", author: "dev", date: "2026-01-04" }
    ]
    refs: [
      { name: "main", commit: "abc123" },
      { name: "feature/graph", commit: "def456" }
    ]

budget: 1200-1500 tokens
complexity: L1-L2 (OBSERVE with library)
```

---

### UP-REFINE: Human Refiner

```yaml
block_id: human-refiner
foundation: unix-philosophy
function: REFINE
description: "Present generated artifacts to human for approval/modification"

input_type: Generated<MCPRequest>
output_type: Refined<ConfirmedRequest>

signature: |
  humanRefiner: Generated<MCPRequest> → Refined<ConfirmedRequest>
  where ConfirmedRequest = { request, approved, modifications }

example_input: |
  Generated<MCPRequest>:
    action: "Update issue LIN-42 status to 'Done'"
    context: "All tests passing, PR #123 merged"

example_output: |
  Refined<ConfirmedRequest>:
    action: "Update issue LIN-42 status to 'Done'"
    approved: true
    modifications: []
    confirmation_dialog: |
      ┌─────────────────────────────────────┐
      │ Confirm External Write              │
      ├─────────────────────────────────────┤
      │ Action: Update Issue Status         │
      │ Issue: LIN-42                       │
      │ Status: Todo → Done                 │
      │                                     │
      │ Press 'y' to confirm, 'n' to cancel│
      └─────────────────────────────────────┘

budget: 800-1200 tokens
complexity: L4 (REFINE with human-in-loop)
```

---

### UP-OPTIMIZE: System Optimizer

```yaml
block_id: system-optimizer
foundation: unix-philosophy
function: OPTIMIZE
description: "Optimize interfaces for minimal API surface (small, sharp tools)"

input_type: Generated<PluginInterface>
output_type: Optimized<PluginAPI>

signature: |
  systemOptimizer: Generated<PluginInterface> → Optimized<PluginAPI>
  where PluginAPI = { interface, minimal, composable }

example_input: |
  Generated<PluginInterface>:
    methods: [Name, NodeTypes, EdgeTypes, Fetch, Transform, Init, Cleanup, Validate]
    complexity: 8 methods

example_output: |
  Optimized<PluginAPI>:
    methods: [Name, NodeTypes, EdgeTypes, Fetch, Transform]
    removed: [Init, Cleanup, Validate]
    rationale: |
      - Init/Cleanup → Registry responsibility (Commandment #2)
      - Validate → Type system enforces (compile-time check)
    complexity: 5 methods (37.5% reduction)
    unix_principle: "Do one thing well: fetch and transform data"

budget: 2000-2500 tokens
complexity: L5 (OPTIMIZE with system design)
```

---

## Block Composition Examples

### Example 1: Sequential (Phase 1A)

```yaml
workflow: sqlite-store-implementation
composition: StructureObserver → TypeGenerator

blocks:
  - structure-observer:
      input: internal/graph/schema.go
      output: Observable<Schema>

  - type-generator:
      input: Observable<Schema>
      output: Generated<StoreImpl>

type_flow:
  Code → Observable<Schema> → Generated<StoreImpl>

operators: [→]
```

---

### Example 2: Parallel Synthesis (Phase 4A)

```yaml
workflow: claude-mcp-bridge
composition: (ContextObserver || KnowledgeGatherer) → ContextReasoner → PromptGenerator → HumanRefiner

blocks:
  - context-observer:
      input: FocusedNode + Edges
      output: Observable<NodeContext>

  - knowledge-gatherer:
      input: NavigationHistory + Breadcrumb
      output: Observable<SessionContext>

  - context-reasoner:
      input: (Observable<NodeContext>, Observable<SessionContext>)
      output: Reasoning<ClaudePrompt>

  - prompt-generator:
      input: Reasoning<ClaudePrompt>
      output: Generated<MCPRequest>

  - human-refiner:
      input: Generated<MCPRequest>
      output: Refined<ConfirmedRequest>

type_flow:
  (Observable<NodeContext> || Observable<SessionContext>)
    → Reasoning<ClaudePrompt>
    → Generated<MCPRequest>
    → Refined<ConfirmedRequest>

operators: [||, →]
```

---

## Type System

### Input Types

- `Observable<Code>`: Source code files
- `Observable<Specification>`: Requirement documents
- `Observable<APISchema>`: API documentation
- `Observable<Repository>`: Git repository
- `Observable<GraphStructure>`: Node/Edge data
- `Observable<Documents>`: Markdown files

### Output Types

- `Observable<T>`: Raw observations
- `Reasoning<T>`: Analysis + insights
- `Generated<T>`: New artifacts
- `Refined<T>`: Improved versions
- `Optimized<T>`: Performance-tuned

### Type Conversion Rules

```
Observable<T> → Reasoning<T>      (Analysis)
Reasoning<T> → Generated<T>       (Generation)
Generated<T> → Refined<T>         (Improvement)
Refined<T> → Optimized<T>         (Optimization)
Any<T> → Observable<T>            (Re-observation)
```

---

## Usage with /blocks Command

```bash
# List all available blocks
/blocks --list

# Show block details
/blocks --info pattern-observer

# Compose blocks into workflow
/blocks "structure-observer → type-generator" \
  --input internal/graph/schema.go \
  --output internal/graph/store.go

# Parallel composition
/blocks "(context-observer || knowledge-gatherer) → context-reasoner" \
  --input "focused-node.json navigation-history.json" \
  --output claude-prompt.json
```

---

## Summary Statistics

| Category | Count |
|----------|-------|
| **Total Blocks Defined** | 20 |
| **Foundations Used** | 7 |
| **Functions Used** | 5 (OBSERVE, REASON, GENERATE, REFINE, OPTIMIZE) |
| **Complexity Range** | L1-L5 |
| **Token Budget Range** | 600-3000 tokens |

**Blocks per Foundation**:
- Abstraction Principles: 2
- Systems Thinking: 4
- Knowledge Synthesis: 2
- Specification-Driven: 4
- Category Theory: 3
- Elm Architecture: 2
- Unix Philosophy: 3

**Blocks per Function**:
- OBSERVE: 7 blocks
- REASON: 4 blocks
- GENERATE: 7 blocks
- REFINE: 1 block
- OPTIMIZE: 1 block

---

## Cross-References

- **Agent Composition Plan**: `docs/AGENT-COMPOSITION-PLAN.md`
- **OIS Taxonomy**: `~/.claude/skills/ois-taxonomy.md`
- **Functional Requirements**: `specs/FUNCTIONAL-REQUIREMENTS.md`
- **Constitution**: `docs/CONSTITUTION.md`

---

**Status**: 20 atomic blocks defined ✓
**Philosophy**: Indivisible, composable, type-safe
**Next Action**: Execute compositions via `/blocks` or `/ois-compose`

