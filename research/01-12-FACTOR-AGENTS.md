# 12-Factor Agents: HumanLayer's Methodology for Production AI Agents

**Source**: HumanLayer (Dex Horthy)
**Research Date**: 2026-01-05
**Relevance to MAAT**: Critical - Defines agent architecture philosophy

---

## Overview

The 12-Factor Agents methodology provides principles for building reliable, production-ready LLM applications.

**Core Insight**: Successful production agents are not "magical autonomous beings" but well-engineered software with LLM capabilities strategically applied at key decision points.

---

## The 12 Factors

| # | Factor | Description |
|---|--------|-------------|
| 1 | **Natural Language to Tool Calls** | LLMs convert user input into structured JSON tool invocations |
| 2 | **Own Your Prompts** | Maintain direct control over prompts; avoid black-box framework abstractions |
| 3 | **Own Your Context Window** | Deliberately architect what information enters the LLM context |
| 4 | **Tools Are Just Structured Outputs** | Treat tools as structured output schemas, not special primitives |
| 5 | **Unify Execution State and Business State** | Synchronize agent execution with application business logic |
| 6 | **Launch/Pause/Resume with Simple APIs** | Design lifecycle management for interruptible workflows |
| 7 | **Contact Humans with Tool Calls** | Use tool-calling mechanisms to involve humans when needed |
| 8 | **Own Your Control Flow** | Explicitly manage decision logic; avoid opaque framework loops |
| 9 | **Compact Errors into Context Window** | Efficiently represent failures for LLM processing |
| 10 | **Small, Focused Agents** | Design narrow, specific responsibilities (3-10 steps max) |
| 11 | **Trigger from Anywhere** | Enable flexible invocation across platforms (Email, Slack, etc.) |
| 12 | **Make Your Agent a Stateless Reducer** | Structure agents as pure functions: input state -> output state |

---

## Design Philosophy

**Key Principle**: "The future of agent development isn't more magical frameworks -- it's better software engineering applied to LLM capabilities."

The methodology challenges the assumption that agents need complex autonomous reasoning:

- **Deterministic code paths** with LLM steps "sprinkled in at just the right points"
- **Engineering reliability** where models "almost succeed" (addressing the 70-80% functionality wall)
- **Incremental adoption** of agent patterns into existing products rather than wholesale framework adoption

---

## State Management: The Stateless Reducer Pattern

Factor 12 is central to the methodology:

```
state' = agent(state, event)
```

**Properties**:
- Agents receive complete context and return next state
- Execution state derives from business state (Factor 5)
- No hidden internal state -- full reproducibility and debuggability
- Enables **Launch/Pause/Resume** (Factor 6) for human approval workflows

---

## Tool Integration Principles

- **Tools are structured outputs** -- JSON schemas that trigger deterministic code
- **Human contact as a tool** (Factor 7) -- same mechanism for human-in-the-loop
- **Compact errors** into context for LLM reasoning about failures
- **Own the control flow** -- explicit loops, not framework magic

---

## Application to MAAT Architecture

| Factor | MAAT Application |
|--------|------------------|
| **Small, Focused Agents** | Separate modules for Linear sync, GitHub operations, Claude interaction |
| **Stateless Reducer** | Each command produces new state from current state + user action |
| **Launch/Pause/Resume** | Long-running operations (PR review, issue triage) can pause for approval |
| **Contact Humans** | Tool call pattern for approval gates before destructive actions |
| **Trigger from Anywhere** | TUI commands, keyboard shortcuts, external webhooks |
| **Own Your Context** | Curated context per view (only relevant Linear issues, focused PR diffs) |
| **Unify State** | Single state model representing Linear + GitHub + session context |

---

## Key Takeaway for MAAT

MAAT should be **mostly deterministic TUI code** with Claude Code invocations at strategic decision points, not an autonomous agent attempting end-to-end automation.

The TUI provides the control flow; AI augments specific decisions.

---

## Sources

- [12 Factor Agents - HumanLayer](https://www.humanlayer.dev/12-factor-agents)
- [GitHub - humanlayer/12-factor-agents](https://github.com/humanlayer/12-factor-agents)
- [HumanLayer Blog - 12 Factor Agents](https://www.humanlayer.dev/blog/12-factor-agents)
