# ADR-005: Thin Integration Philosophy

**Status**: Accepted
**Date**: 2026-01-05
**Deciders**: MAAT Architecture Team
**Commandments**: #5 Controlled Effects, #7 Composition Monopoly

---

## Context

MAAT integrates multiple external systems:
- **Linear**: Project management, issues, cycles
- **GitHub**: Repositories, PRs, code
- **Claude**: AI assistance
- **Future**: Jira, Azure DevOps, etc.

The temptation is to build rich clients that replicate functionality. This leads to:
- Feature parity competition (unwinnable)
- Maintenance burden (API changes break us)
- Bloated codebase (each integration grows)
- Lost focus (building Linear features instead of MAAT value)

## Decision

Implement **thin integrations** - API clients only, no business logic:

### Integration Architecture

```
maat/
├── internal/
│   ├── linear/          # ~300-500 LOC
│   │   ├── client.go    # API client
│   │   ├── types.go     # Data types
│   │   └── adapter.go   # Node/Edge conversion
│   ├── github/          # ~300-500 LOC
│   │   ├── client.go
│   │   ├── types.go
│   │   └── adapter.go
│   ├── claude/          # ~200-300 LOC
│   │   ├── client.go
│   │   └── context.go
│   └── graph/           # Core MAAT logic HERE
│       ├── knowledge.go
│       ├── layout.go
│       └── query.go
```

### Linear Client (Example)

```go
// internal/linear/client.go
package linear

import (
    "context"
    "github.com/machinebox/graphql"
)

type Client struct {
    gql    *graphql.Client
    teamID string
}

func NewClient(apiKey, teamID string) *Client {
    return &Client{
        gql:    graphql.NewClient("https://api.linear.app/graphql"),
        teamID: teamID,
    }
}

// Thin: Just fetches data, no business logic
func (c *Client) FetchIssues(ctx context.Context) ([]Issue, error) {
    var resp struct {
        Team struct {
            Issues struct {
                Nodes []Issue
            }
        }
    }

    req := graphql.NewRequest(`
        query($teamId: String!) {
            team(id: $teamId) {
                issues {
                    nodes {
                        id
                        title
                        state { name }
                        assignee { name }
                        labels { nodes { name } }
                    }
                }
            }
        }
    `)
    req.Var("teamId", c.teamID)

    if err := c.gql.Run(ctx, req, &resp); err != nil {
        return nil, err
    }
    return resp.Team.Issues.Nodes, nil
}

// Write operations return confirmation requests
func (c *Client) UpdateIssueStatus(ctx context.Context, id, status string) ConfirmRequest {
    return ConfirmRequest{
        Action:      "Update Issue Status",
        Description: fmt.Sprintf("Set %s to %s", id, status),
        Execute: func() error {
            // Actual mutation
            req := graphql.NewRequest(`
                mutation($id: String!, $status: String!) {
                    issueUpdate(id: $id, input: { stateId: $status }) {
                        success
                    }
                }
            `)
            req.Var("id", id)
            req.Var("status", status)
            return c.gql.Run(ctx, req, nil)
        },
    }
}
```

### Adapter Pattern

```go
// internal/linear/adapter.go
package linear

import "maat/internal/graph"

// Convert Linear types to graph nodes
func (i Issue) ToNode() graph.Node {
    return graph.Node{
        ID:     "linear:" + i.ID,
        Type:   graph.NodeTypeIssue,
        Source: "linear",
        Data:   mustMarshal(i),
        Metadata: graph.NodeMetadata{
            CreatedAt:   i.CreatedAt,
            UpdatedAt:   i.UpdatedAt,
            AccessLevel: graph.RoleIC,  // Issues visible to all
        },
    }
}

// Convert Linear relationships to edges
func (i Issue) ToEdges() []graph.Edge {
    var edges []graph.Edge

    for _, blocked := range i.BlockedBy {
        edges = append(edges, graph.Edge{
            FromID:   "linear:" + blocked.ID,
            ToID:     "linear:" + i.ID,
            Relation: graph.EdgeBlocks,
        })
    }

    for _, related := range i.Related {
        edges = append(edges, graph.Edge{
            FromID:   "linear:" + i.ID,
            ToID:     "linear:" + related.ID,
            Relation: graph.EdgeRelated,
        })
    }

    return edges
}
```

### Read vs Write Sovereignty

| Operation | Confirmation | Rationale |
|-----------|--------------|-----------|
| Fetch issues | None | No side effects |
| Fetch PRs | None | No side effects |
| Fetch files | None | No side effects |
| Update issue status | **Required** | Writes to Linear |
| Create comment | **Required** | Writes to Linear/GitHub |
| Merge PR | **Required** | Critical write to GitHub |
| Apply Claude edit | **Required** | Modifies local files |

```go
type ConfirmRequest struct {
    Action      string
    Description string
    Target      string      // Node ID or file path
    Preview     string      // Diff or preview
    Execute     func() error
}

// In TUI: show confirmation dialog
func (m Model) HandleConfirmRequest(req ConfirmRequest) (Model, tea.Cmd) {
    m.pendingConfirm = &req
    m.showConfirmDialog = true
    return m, nil
}

func (m Model) ConfirmAction() (Model, tea.Cmd) {
    if m.pendingConfirm == nil {
        return m, nil
    }

    // Log to audit before execution
    logAudit(m.pendingConfirm)

    // Execute the action
    return m.WithLoading(true), executeConfirmedCmd(m.pendingConfirm)
}
```

### Size Constraints

Each integration should be:
- **~500 LOC max** for client + types + adapter
- **No business logic** - just data transformation
- **No caching** - graph backend handles that
- **No UI** - TUI is separate concern

## Consequences

### Positive
- **Maintainability**: Small surface area for API changes
- **Focus**: MAAT value is in composition, not features
- **Testability**: Adapters are pure functions
- **Extensibility**: New integrations follow same pattern

### Negative
- **Feature Limits**: Can't do everything Linear/GitHub can
- **API Dependence**: Offline mode limited to cached data
- **Lowest Common Denominator**: Must abstract across systems

### Mitigations
- Accept limits as feature (composition is the value)
- Robust caching in graph backend
- System-specific extensions via plugins (ADR-007)

## Compliance

This ADR enforces:
- **Commandment #5**: Commands describe, runtime executes
- **Commandment #7**: Value from integration, not features

## References

- MAAT Constitution: `/Users/manu/Documents/LUXOR/MAAT/docs/CONSTITUTION.md`
- Spec-Kit Methodology: `/Users/manu/Documents/LUXOR/MAAT/research/02-SPEC-KIT.md`
