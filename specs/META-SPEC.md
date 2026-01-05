# MAAT Meta-Specification

**Document Type**: Meta-Specification (Philosophy, Process, Creed)
**Version**: 1.0
**Date**: 2026-01-05
**Status**: Canonical

---

## Purpose

This META-SPEC governs how MAAT specifications are created, prioritized, and validated. It is the specification that specifies how we specify.

---

## Part I: Philosophy Hierarchy

### The Commandment Priority Stack

When commandments conflict, higher-ranked commandments take precedence:

```
┌────────────────────────────────────────────┐
│  TIER 1: INVIOLABLE (Never Compromise)     │
├────────────────────────────────────────────┤
│  #1  Immutable Truth                       │
│      └─ All state in Model, Update pure    │
│  #6  Human Contact                         │
│      └─ AI requires explicit invocation    │
│  #10 Sovereignty Preservation              │
│      └─ External writes require confirm    │
└────────────────────────────────────────────┘
                    │
                    ▼
┌────────────────────────────────────────────┐
│  TIER 2: ARCHITECTURAL (Shape the System)  │
├────────────────────────────────────────────┤
│  #2  Graph Supremacy                       │
│      └─ All entities as nodes/edges        │
│  #4  Navigation Monopoly                   │
│      └─ Enter drills, Esc backs            │
│  #7  Composition Monopoly                  │
│      └─ Value from composition, not parity │
└────────────────────────────────────────────┘
                    │
                    ▼
┌────────────────────────────────────────────┐
│  TIER 3: TACTICAL (Guide Implementation)   │
├────────────────────────────────────────────┤
│  #3  Text Interface                        │
│      └─ Msg types are the language         │
│  #5  Controlled Effects                    │
│      └─ tea.Cmd only, no goroutines        │
│  #8  Async Purity                          │
│      └─ Commands describe, runtime executes│
│  #9  Terminal Citizenship                  │
│      └─ Pure stdout, dark theme, ^C exits  │
└────────────────────────────────────────────┘
```

### Conflict Resolution Examples

**Scenario**: User wants faster updates via auto-sync
- **#10 Sovereignty** (Tier 1) vs **Performance** (unlisted)
- **Resolution**: No auto-push sync. User must explicitly trigger.

**Scenario**: Beautiful inline editing vs state purity
- **#1 Immutable Truth** (Tier 1) vs **UX Polish** (unlisted)
- **Resolution**: Edit via Msg→Update→View cycle, never direct mutation.

**Scenario**: Full Linear feature parity requested
- **#7 Composition Monopoly** (Tier 2) says NO
- **Resolution**: Thin client only. Link to Linear for advanced ops.

---

## Part II: Spec-Driven Development Methodology

### The Specification Lifecycle

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   IDEATE    │────▶│   SPECIFY   │────▶│  VALIDATE   │
│             │     │             │     │             │
│ • User need │     │ • FR/AFR    │     │ • RMP loop  │
│ • Gap found │     │ • ADR       │     │ • Quality   │
│ • Research  │     │ • Anti-req  │     │   gates     │
└─────────────┘     └─────────────┘     └─────────────┘
                           │
                           ▼
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   REFINE    │◀────│  IMPLEMENT  │◀────│   APPROVE   │
│             │     │             │     │             │
│ • Feedback  │     │ • Code      │     │ • Review    │
│ • Version   │     │ • Tests     │     │ • Sign-off  │
│ • Archive   │     │ • Docs      │     │ • Merge     │
└─────────────┘     └─────────────┘     └─────────────┘
```

### Specification Types

| Type | Purpose | Location | Governance |
|------|---------|----------|------------|
| **Constitution** | Foundational principles | `docs/CONSTITUTION.md` | Change requires ADR |
| **ADR** | Architecture decisions | `specs/ADR/` | Append-only, versioned |
| **FR** | Functional requirements | `specs/FUNCTIONAL-REQUIREMENTS.md` | Numbered, traceable |
| **AFR** | Anti-requirements | `specs/ANTI-REQUIREMENTS.md` | Prohibitions |
| **Component Spec** | Implementation detail | `docs/components/` | Per-library guides |
| **META-SPEC** | This document | `specs/META-SPEC.md` | Self-referential |

### Quality Gates

Every specification must pass:

1. **Commandment Alignment**: Does it comply with Tier 1-3 priorities?
2. **Anti-Pattern Check**: Does it avoid all A-001 through A-008?
3. **Traceable Reference**: Does it cite parent documents?
4. **Implementation Path**: Is there a clear code path?
5. **Test Strategy**: How do we verify compliance?

---

## Part III: The MAAT Creed

### Core Beliefs

```
We believe in...

THE GRAPH AS TRUTH
  Knowledge lives in relationships, not tables.
  Every entity connects. Isolation is debt.

THE HUMAN AS SOVEREIGN
  AI advises, humans decide.
  Every action must be authorized.
  No surprises. No magic. No autonomy.

THE TERMINAL AS HOME
  We honor the shell that raised us.
  Dark backgrounds. Keyboard primacy.
  Ctrl+C is sacred. Response is instant.

COMPOSITION OVER COMPETITION
  We integrate, not replicate.
  Linear is Linear. GitHub is GitHub.
  Our value is the weaving, not the threads.

PURITY AS STRENGTH
  State flows one way.
  Effects are commands.
  The Update function tells no lies.

SPECIFICATION AS CODE
  If it's not written, it doesn't exist.
  ADRs are permanent records.
  Anti-requirements are laws.
```

### The Developer Oath

Before implementing any feature, the developer affirms:

> "I will not mutate state outside Update.
> I will not spawn goroutines outside tea.Cmd.
> I will not write to external systems without confirmation.
> I will not invoke AI without explicit user request.
> I will not compete with integrated platforms.
> I will trace every feature to a specification.
> I will honor the Commandments in priority order."

---

## Part IV: Requirements Governance

### FR Numbering Convention

```
FR-{DDD}: {Title}

DDD = Domain + Sequence
  0xx: Core UI/Graph (001-009)
  1xx: Integration (010-019)
  2xx: AI Features (020-029)
  3xx: Self-Service (030-039)
  4xx: Plugin System (040-049)
```

### AFR Numbering Convention

```
A-{NNN}: NO {Prohibition}

NNN = Sequential
  001-010: Critical (code review blocks)
  011-020: High (PR review flags)
  021+: Medium (linter warnings)
```

### ADR Lifecycle

```
Status Flow:
  PROPOSED → ACCEPTED → DEPRECATED → SUPERSEDED
              │
              └─ REJECTED (terminal)

ADR Rules:
  • Never delete an ADR
  • Superseded ADRs reference successor
  • Context section explains the "why"
  • Decision section is definitive
```

### Specification Versioning

```
Major.Minor.Patch

Major: Constitution change, new Commandment
Minor: New FR, new ADR
Patch: Clarification, typo fix

Examples:
  1.0.0 - Initial specification
  1.1.0 - Added FR-010 (Git History)
  1.1.1 - Clarified A-003 language
  2.0.0 - New Commandment added
```

---

## Part V: Specification Templates

### FR Template

```markdown
## FR-{NNN}: {Title}

**Commandment**: #{N} {Name}
**Priority**: P{1-3}
**Status**: Draft | Approved | Implemented

### Description
{What this feature does}

### Acceptance Criteria
- [ ] Criterion 1
- [ ] Criterion 2

### Anti-Patterns (must avoid)
- {Specific thing NOT to do}

### Implementation Notes
{Technical guidance}
```

### ADR Template

```markdown
# ADR-{NNN}: {Title}

**Status**: Proposed | Accepted | Deprecated | Superseded
**Date**: YYYY-MM-DD
**Commandments**: #{N}, #{M}

## Context
{Why this decision is needed}

## Decision
{What we decided}

## Consequences
### Positive
- {Benefit}

### Negative
- {Tradeoff}

### Neutral
- {Observation}
```

### AFR Template

```markdown
## A-{NNN}: NO {Prohibition}

**Commandment**: #{N} {Name}
**Severity**: Critical | High | Medium

### Prohibition
{What is forbidden}

### Anti-Pattern
```go
// ❌ FORBIDDEN
{bad code}
```

### Correct Pattern
```go
// ✅ CORRECT
{good code}
```

### Rationale
{Why this is forbidden}
```

---

## Part VI: Recursive Refinement Protocol

### RMP Application to Specifications

MAAT specifications are themselves subject to /rmp:

```
Initial Spec → Assess Quality → If < 0.85 → Refine → Repeat
```

### Quality Dimensions for Specs

| Dimension | Weight | Criteria |
|-----------|--------|----------|
| **Completeness** | 30% | All FRs have acceptance criteria |
| **Consistency** | 25% | No contradictions between docs |
| **Traceability** | 20% | Every FR traces to Commandment |
| **Clarity** | 15% | Unambiguous language |
| **Testability** | 10% | Verifiable acceptance criteria |

### Convergence Criteria

A specification is considered converged when:

1. Quality score ≥ 0.85 across all dimensions
2. All Tier 1 Commandments are respected
3. No unresolved anti-pattern conflicts
4. Implementation path is clear
5. Human reviewer approves

---

## Part VII: Living Document Protocol

### Update Triggers

This META-SPEC must be updated when:

1. New Commandment is added to Constitution
2. Specification methodology changes
3. New document type is introduced
4. Governance process is modified
5. Creed is amended

### Amendment Process

```
1. Propose change via ADR
2. Review against existing Commandments
3. If Tier 1 conflict: REJECT
4. If approved: Update META-SPEC
5. Increment version number
6. Archive previous version
```

---

## Cross-References

- **Constitution**: `/Users/manu/Documents/LUXOR/MAAT/docs/CONSTITUTION.md`
- **ADRs**: `/Users/manu/Documents/LUXOR/MAAT/specs/ADR/`
- **Functional Requirements**: `/Users/manu/Documents/LUXOR/MAAT/specs/FUNCTIONAL-REQUIREMENTS.md`
- **Anti-Requirements**: `/Users/manu/Documents/LUXOR/MAAT/specs/ANTI-REQUIREMENTS.md`
- **Component Specs**: `/Users/manu/Documents/LUXOR/MAAT/docs/components/`

---

## Appendix: Philosophical Influences

MAAT draws from:

1. **Elm Architecture**: Pure functional UI, Msg-based state
2. **12-Factor Agents**: Declarative config, stateless processes
3. **Spec-Kit**: Specification-driven development
4. **Unix Philosophy**: Do one thing well, compose via pipes
5. **Ma'at (Egyptian)**: Balance, truth, cosmic order

The name MAAT honors the Egyptian goddess of truth and order - appropriate for a system that brings order to developer chaos through truthful state management and balanced composition.

---

**Document Hash**: META-SPEC-v1.0-2026-01-05
**Specification of specifications**: Complete ✓
