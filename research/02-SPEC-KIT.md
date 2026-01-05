# GitHub Spec-Kit: Specification-Driven Development

**Source**: GitHub
**Research Date**: 2026-01-05
**Relevance to MAAT**: Critical - Defines development methodology

---

## Core Philosophy

Spec-Driven Development (SDD) inverts traditional software development:

- Traditional: Specifications are disposable scaffolding discarded once coding begins
- SDD: **Specifications are the source of truth that generates code**

The methodology shifts from "code is king" to "intent is king" - specifications become executable contracts that AI agents use to generate, test, and validate implementations.

**Key Insight**: Specifications are not documents that inform coding; they are precise definitions that **produce** code.

---

## Key Principles

1. **Specification as lingua franca**: Specs become the primary artifact; code merely expresses them
2. **Constitution-first governance**: Immutable architectural principles establish non-negotiable quality gates
3. **Four-phase workflow**: Specify -> Plan -> Tasks -> Implement
4. **Template-driven quality constraints**: Seven constraint mechanisms guide AI behavior
5. **Test-first imperative**: All implementation follows strict TDD
6. **Anti-speculation discipline**: No "might need" features; every feature traces to user stories
7. **Continuous refinement**: Consistency validation is ongoing, not one-time

---

## Directory Structure

```
.specify/
├── memory/constitution.md    # Immutable principles
├── specs/{feature-branch}/
│   ├── spec.md               # Requirements (WHAT/WHY)
│   ├── plan.md               # Architecture (HOW)
│   ├── tasks.md              # Executable breakdown
│   ├── data-model.md         # Entity definitions
│   ├── contracts/            # API specifications
│   └── research.md           # Technical context
└── templates/                # Reusable blueprints
```

---

## Workflow Commands

| Phase | Command | Output |
|-------|---------|--------|
| **Specify** | `/speckit.specify` | Requirements in `spec.md` |
| **Plan** | `/speckit.plan` | Architecture in `plan.md` |
| **Tasks** | `/speckit.tasks` | Breakdown in `tasks.md` |
| **Implement** | `/speckit.implement` | Code from tasks |

---

## Quality Gates

- **Simplicity Gate**: Maximum 3 projects in scope
- **Anti-Abstraction Gate**: Direct framework usage, no premature abstraction
- **Integration-First Gate**: Contracts before implementation

---

## Application to MAAT TUI Development

### Constitution-First

Define immutable principles upfront:
- Keyboard-driven interaction
- Composable panels
- Plugin architecture
- Performance budgets (<16ms frame time)

### Feature Specification Flow

Each MAAT capability gets its own spec branch:
- `specs/file-browser/spec.md`
- `specs/terminal-multiplexer/spec.md`
- `specs/editor-integration/spec.md`

### Contract-Driven Integration

Define TUI component interfaces in `contracts/` before implementation:
- Event systems
- Panel lifecycle
- Keybinding resolution

### Task Parallelization

Use `[P]` markers in `tasks.md` to identify concurrent work:
- `[P]` Status bar
- `[P]` Command palette
- `[P]` File tree

### Test-First TUI Validation

Write integration tests before widgets:
- Keyboard navigation
- Panel focus management
- State persistence

---

## Why This Matters for MAAT

Terminal interfaces have strict constraints:
- ANSI compliance
- Performance requirements
- Accessibility needs

These benefit from **upfront specification** rather than iterative discovery.

---

## Sources

- [GitHub Spec-Kit Repository](https://github.com/github/spec-kit)
- [GitHub Blog: Spec-Driven Development with AI](https://github.blog/ai-and-ml/generative-ai/spec-driven-development-with-ai-get-started-with-a-new-open-source-toolkit/)
- [Martin Fowler: Understanding Spec-Driven Development](https://martinfowler.com/articles/exploring-gen-ai/sdd-3-tools.html)
