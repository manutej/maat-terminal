# ADR-004: Human-in-Loop AI Integration

**Status**: Accepted
**Date**: 2026-01-05
**Deciders**: MAAT Architecture Team
**Commandments**: #6 Human Contact, #10 Sovereignty Preservation

---

## Context

MAAT integrates Claude for AI assistance. The risk is:
- **Autonomous actions**: AI making changes without human review
- **Context pollution**: AI suggestions appearing unbidden
- **Audit gaps**: No record of AI-driven changes
- **Trust erosion**: Users losing confidence in system predictability

The 12-Factor Agents principle states: "Human contact is a tool call" - explicit, bounded, auditable.

## Decision

Implement **Human-in-Loop AI** with explicit invocation and confirmation gates:

### Invocation Model

```go
// AI is ONLY invoked via explicit user action
type ClaudeInvocation struct {
    Trigger     InvokeTrigger  // Ctrl+A | :ask-claude | command palette
    Context     ClaudeContext  // Assembled from current focus
    Constraints []string       // User-specified limits
}

type InvokeTrigger string
const (
    TriggerHotkey  InvokeTrigger = "ctrl+a"
    TriggerCommand InvokeTrigger = ":ask-claude"
    TriggerPalette InvokeTrigger = "palette"
)

// NO ambient suggestions, NO auto-complete, NO proactive AI
```

### Context Assembly

```go
type ClaudeContext struct {
    FocusedNode    Node              // Currently selected node
    VisibleGraph   []Node            // Nodes in view
    RecentHistory  []Node            // Last 10 visited nodes
    LinkedContent  map[string]string // File contents, issue descriptions
    UserQuestion   string            // What the user asked
}

func (m Model) AssembleContext() ClaudeContext {
    ctx := ClaudeContext{
        FocusedNode:  m.graph.Nodes[m.focusedNode],
        UserQuestion: m.claudeInput,
    }

    // Add linked content based on node type
    switch ctx.FocusedNode.Type {
    case NodeTypeIssue:
        ctx.LinkedContent["issue"] = ctx.FocusedNode.Data
        ctx.LinkedContent["comments"] = m.fetchComments(ctx.FocusedNode.ID)
    case NodeTypePR:
        ctx.LinkedContent["pr"] = ctx.FocusedNode.Data
        ctx.LinkedContent["diff"] = m.fetchDiff(ctx.FocusedNode.ID)
    case NodeTypeFile:
        ctx.LinkedContent["file"] = m.readFile(ctx.FocusedNode.ID)
    }

    return ctx
}
```

### Confirmation Gates

```go
type ClaudeResponse struct {
    Answer      string           // Text response
    Suggestions []Suggestion     // Proposed actions
    Confidence  float64          // 0.0-1.0
}

type Suggestion struct {
    Action      SuggestionAction  // Edit | Create | Update | Comment
    Target      string            // File path or node ID
    Content     string            // Proposed content
    Diff        string            // For edits, show diff
    Explanation string            // Why Claude suggests this
}

type SuggestionAction string
const (
    ActionEdit    SuggestionAction = "edit"
    ActionCreate  SuggestionAction = "create"
    ActionUpdate  SuggestionAction = "update"   // Linear/GitHub
    ActionComment SuggestionAction = "comment"
)

// ALL suggestions require explicit approval
func (m Model) HandleClaudeResponse(resp ClaudeResponse) (Model, tea.Cmd) {
    // Show response in Claude pane
    m.claudeResponse = resp

    // Suggestions are SHOWN, not EXECUTED
    m.pendingSuggestions = resp.Suggestions

    // User must explicitly approve each suggestion
    // [a] Apply  [s] Skip  [e] Edit  [r] Reject All
    return m, nil
}

func (m Model) ApplySuggestion(idx int) (Model, tea.Cmd) {
    suggestion := m.pendingSuggestions[idx]

    // Log to audit trail BEFORE execution
    audit := AIAuditEntry{
        SessionID:  m.claude.SessionID,
        Timestamp:  time.Now(),
        Action:     suggestion.Action,
        Context:    m.claudeContext.UserQuestion,
        Response:   suggestion.Content,
        UserAction: "approved",
        Diff:       suggestion.Diff,
    }

    // Execute with confirmation if write operation
    switch suggestion.Action {
    case ActionEdit:
        return m, editFileCmd(suggestion.Target, suggestion.Content, audit)
    case ActionUpdate:
        return m, confirmAndUpdateCmd(suggestion.Target, suggestion.Content, audit)
    }

    return m, nil
}
```

### Audit Trail

```go
type AIAuditEntry struct {
    SessionID   string          `json:"session_id"`
    Timestamp   time.Time       `json:"timestamp"`
    Action      SuggestionAction `json:"action"`
    Context     string          `json:"context"`      // User's question
    Response    string          `json:"response"`     // Claude's suggestion
    UserAction  string          `json:"user_action"`  // approved | rejected | modified
    Diff        string          `json:"diff"`         // If edit, what changed
    NodeID      string          `json:"node_id"`      // Affected node
}

// SQLite storage
const createAuditTable = `
CREATE TABLE IF NOT EXISTS ai_audit (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id TEXT NOT NULL,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    action TEXT NOT NULL,
    context TEXT,
    response TEXT,
    user_action TEXT,
    diff TEXT,
    node_id TEXT,
    FOREIGN KEY (node_id) REFERENCES nodes(id)
);

CREATE INDEX idx_audit_session ON ai_audit(session_id);
CREATE INDEX idx_audit_timestamp ON ai_audit(timestamp);
`

func (kg *KnowledgeGraph) LogAIAction(entry AIAuditEntry) error {
    _, err := kg.db.Exec(`
        INSERT INTO ai_audit
        (session_id, action, context, response, user_action, diff, node_id)
        VALUES (?, ?, ?, ?, ?, ?, ?)
    `, entry.SessionID, entry.Action, entry.Context,
       entry.Response, entry.UserAction, entry.Diff, entry.NodeID)
    return err
}
```

### UI Presentation

```
┌─────────────────────────────────────────────────┐
│  CLAUDE ASSISTANT                    [Ctrl+A]   │
├─────────────────────────────────────────────────┤
│  Context: Issue LIN-234 "Implement Context"     │
│  Files: src/compress.rs, docs/crdts.md          │
│                                                 │
│  > How can I optimize the compression algo?     │
│                                                 │
│  ─────────────────────────────────────────────  │
│  Claude:                                        │
│  Based on the CRDT approach in docs/crdts.md,   │
│  I suggest using delta-state compression...     │
│                                                 │
│  Suggestions:                                   │
│  [1] Edit src/compress.rs (show diff)           │
│  [2] Add comment to LIN-234                     │
│                                                 │
│  [a] Apply  [s] Skip  [v] View diff  [r] Reject │
└─────────────────────────────────────────────────┘
```

## Consequences

### Positive
- **Human Control**: AI never acts without approval
- **Auditability**: Full trail of AI interactions
- **Trust**: Users know exactly what AI can/cannot do
- **Compliance**: Enterprise-ready audit logs

### Negative
- **Friction**: Extra steps to apply suggestions
- **Speed**: Can't auto-apply obvious fixes
- **UX Complexity**: More UI states to manage

### Mitigations
- Batch approval for low-risk suggestions
- Keyboard shortcuts for quick approve/reject
- Session history for quick re-application

## Compliance

This ADR enforces:
- **Commandment #6**: Explicit invocation, human review
- **Commandment #10**: Orchestrate, never colonize

## References

- 12-Factor Agents: `/Users/manu/Documents/LUXOR/MAAT/research/01-12-FACTOR-AGENTS.md`
- MAAT Constitution: `/Users/manu/Documents/LUXOR/MAAT/docs/CONSTITUTION.md`
