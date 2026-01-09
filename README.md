# MAAT: Modular Agentic Architecture for Terminal

> "One graph to rule them all — where issues become nodes, relationships become edges, and the developer becomes the navigator."

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

## What is MAAT?

MAAT is a unified terminal workspace that integrates **Linear issues**, **GitHub PRs**, and **Claude Code** into a navigable knowledge graph. Built with Go using the [Bubble Tea](https://github.com/charmbracelet/bubbletea) framework.

```
┌──────────────┬──────────────────────┬──────────────────────┐
│ GRAPH PANE   │     MAIN PANE        │   DETAIL PANE        │
│              │                      │                       │
│ • Node tree  │ • Issue details      │ • Metadata           │
│ • Edges      │ • PR content         │ • Git history        │
│ • Hierarchy  │ • Code diff          │ • AI suggestions     │
└──────────────┴──────────────────────┴──────────────────────┘
```

## Key Features

- **Knowledge Graph**: Issues, PRs, commits, and files as connected nodes
- **Keyboard-First**: vim-style navigation (h/j/k/l), Enter to drill down, Esc to back up
- **Human-in-Loop AI**: Claude integration with explicit invocation and confirmation gates
- **Thin Integrations**: API clients only — no feature competition with Linear or GitHub
- **Single Binary**: One `go build`, instant startup, zero runtime dependencies

## Design Philosophy

MAAT follows **10 Design Commandments** that govern all decisions:

| Priority | Commandment | Summary |
|----------|-------------|---------|
| Tier 1 | #1 Immutable Truth | All state in Model, Update is pure |
| Tier 1 | #6 Human Contact | AI requires explicit invocation |
| Tier 1 | #10 Sovereignty | External writes need confirmation |
| Tier 2 | #2 Graph Supremacy | Every entity is a node or edge |
| Tier 2 | #4 Navigation Monopoly | Enter drills, Esc backs |
| Tier 2 | #7 Composition | Value from weaving, not competing |

See [docs/CONSTITUTION.md](docs/CONSTITUTION.md) for full details.

## Quick Start

```bash
# Clone
git clone https://github.com/manutej/maat-terminal.git
cd maat

# Build
go build -o maat ./cmd/maat

# Run
./maat
```

## Configuration

```yaml
# configs/default.yaml
linear:
  api_key: ${LINEAR_API_KEY}
  team_id: "your-team-id"

github:
  token: ${GITHUB_TOKEN}
  owner: "your-org"
  repo: "your-repo"

claude:
  enabled: true
  require_confirmation: true
```

## Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `h/j/k/l` | Navigate graph |
| `Enter` | Drill down into node |
| `Esc` | Navigate back |
| `Tab` | Cycle panes |
| `/` | Search |
| `Ctrl+A` | Invoke Claude |
| `?` | Help |
| `q` | Quit |

## Documentation

| Document | Purpose |
|----------|---------|
| [MAAT-SPEC.md](docs/specs/MAAT-SPEC.md) | Unified specification |
| [docs/CONSTITUTION.md](docs/CONSTITUTION.md) | 10 Commandments |
| [specs/FUNCTIONAL-REQUIREMENTS.md](specs/FUNCTIONAL-REQUIREMENTS.md) | FR-001 to FR-010 |
| [specs/ANTI-REQUIREMENTS.md](specs/ANTI-REQUIREMENTS.md) | What we refuse to build |
| [specs/ADR/](specs/ADR/) | Architecture decisions |
| [docs/components/](docs/components/) | Implementation guides |

## Architecture

```
┌─────────────────────────────────────────────┐
│              Bubble Tea TUI                 │
│         (Elm Architecture / TEA)            │
└─────────────────────────────────────────────┘
                     │
     ┌───────────────┼───────────────┐
     ▼               ▼               ▼
┌─────────┐    ┌─────────┐    ┌─────────┐
│ Linear  │    │ GitHub  │    │ Claude  │
│ Client  │    │ Client  │    │ Bridge  │
└─────────┘    └─────────┘    └─────────┘
     │               │               │
     └───────────────┼───────────────┘
                     ▼
           ┌─────────────────┐
           │ Knowledge Graph │
           │    (SQLite)     │
           └─────────────────┘
```

## Project Status

🚧 **Phase 1: Foundation** (In Progress)

- [ ] Go project structure
- [ ] Basic Bubble Tea model
- [ ] SQLite knowledge graph schema
- [ ] 3-pane layout

See [MAAT-SPEC.md](docs/specs/MAAT-SPEC.md) for full roadmap.

## Contributing

1. Read the [Constitution](docs/CONSTITUTION.md)
2. Check [Anti-Requirements](specs/ANTI-REQUIREMENTS.md) for what NOT to build
3. Follow the Elm Architecture pattern
4. Ensure all PRs trace to a Functional Requirement

## Name

**MAAT** (מַעַת) honors the Egyptian goddess of truth, justice, and cosmic order — appropriate for a system that brings order to developer chaos through truthful state management and balanced composition.

## License

MIT License - See [LICENSE](LICENSE) for details.
