# ADR-006: IDP Self-Service Actions

**Status**: Accepted
**Date**: 2026-01-05
**Deciders**: MAAT Architecture Team
**Commandments**: #6 Human Contact, #10 Sovereignty Preservation

---

## Context

Gap analysis identified that MAAT should enable **creation and operation**, not just viewing:
- Users want to create issues, branches, PRs from within MAAT
- The tool should offer "golden paths" - standardized workflows
- This transforms MAAT from "viewer" to "Internal Developer Platform (IDP)"

However, this must respect:
- Human-in-loop principle (no autonomous actions)
- Sovereignty preservation (confirm before writing to external systems)
- Single responsibility (MAAT orchestrates, doesn't replace)

## Decision

Implement **command palette with golden paths** for self-service actions:

### Command Palette

```go
type CommandPalette struct {
    Input     string
    Filtered  []Command
    Selected  int
}

type Command struct {
    Name        string           // :new-issue
    Description string           // Create a new Linear issue
    Category    CommandCategory  // create | navigate | ai | system
    Handler     func(Model, []string) (Model, tea.Cmd)
    Confirm     bool             // Requires confirmation?
}

var commands = []Command{
    // Create commands (all require confirmation)
    {":new-issue", "Create Linear issue", CategoryCreate, handleNewIssue, true},
    {":new-branch", "Create git branch from issue", CategoryCreate, handleNewBranch, true},
    {":new-pr", "Open PR linked to issue", CategoryCreate, handleNewPR, true},
    {":new-comment", "Add comment to issue/PR", CategoryCreate, handleNewComment, true},

    // Navigate commands (no confirmation)
    {":goto", "Jump to node by ID", CategoryNavigate, handleGoto, false},
    {":search", "Fuzzy search workspace", CategoryNavigate, handleSearch, false},
    {":back", "Return to previous view", CategoryNavigate, handleBack, false},

    // AI commands (require confirmation for actions)
    {":ask-claude", "Ask Claude about current context", CategoryAI, handleAskClaude, false},
    {":explain", "Explain selected code/issue", CategoryAI, handleExplain, false},
    {":suggest-fix", "Get fix suggestions", CategoryAI, handleSuggestFix, true},

    // System commands
    {":sync", "Refresh from external systems", CategorySystem, handleSync, false},
    {":settings", "Open settings", CategorySystem, handleSettings, false},
    {":quit", "Exit MAAT", CategorySystem, handleQuit, false},
}
```

### Golden Paths

Pre-configured workflows that chain multiple actions:

```go
type GoldenPath struct {
    Name        string
    Description string
    Steps       []GoldenStep
}

type GoldenStep struct {
    Command     string
    AutoFill    map[string]string  // Values derived from context
    Confirm     bool
}

var goldenPaths = []GoldenPath{
    {
        Name: "Start Work on Issue",
        Description: "Create branch → Open editor → Update status",
        Steps: []GoldenStep{
            {Command: ":new-branch", AutoFill: map[string]string{
                "name": "{{issue.identifier}}-{{issue.title|slugify}}",
            }, Confirm: true},
            {Command: ":open-editor", AutoFill: map[string]string{
                "path": "{{repo.root}}",
            }, Confirm: false},
            {Command: ":update-status", AutoFill: map[string]string{
                "status": "In Progress",
            }, Confirm: true},
        },
    },
    {
        Name: "Submit for Review",
        Description: "Commit → Push → Create PR → Request review",
        Steps: []GoldenStep{
            {Command: ":git-commit", Confirm: true},
            {Command: ":git-push", Confirm: true},
            {Command: ":new-pr", AutoFill: map[string]string{
                "title": "{{issue.identifier}}: {{issue.title}}",
                "body": "Closes {{issue.url}}\n\n{{commit.messages}}",
            }, Confirm: true},
            {Command: ":request-review", Confirm: true},
        },
    },
    {
        Name: "Quick Bug Report",
        Description: "Create issue with template → Assign to self",
        Steps: []GoldenStep{
            {Command: ":new-issue", AutoFill: map[string]string{
                "template": "bug",
                "labels": "bug,triage",
            }, Confirm: true},
            {Command: ":assign-self", Confirm: true},
        },
    },
}
```

### Action Handlers

```go
// :new-issue handler
func handleNewIssue(m Model, args []string) (Model, tea.Cmd) {
    // Determine context
    project := m.currentProject()
    team := m.currentTeam()

    // Open issue creation form
    m.mode = ModeCreateIssue
    m.createForm = IssueForm{
        Team:    team.ID,
        Project: project.ID,
        Title:   "",
        Body:    "",
        Labels:  []string{},
    }

    return m, nil
}

// Form submission (requires confirmation)
func (m Model) SubmitIssueForm() (Model, tea.Cmd) {
    form := m.createForm

    // Build confirmation request
    confirm := ConfirmRequest{
        Action:      "Create Issue",
        Description: fmt.Sprintf("Create '%s' in %s", form.Title, form.Team),
        Preview:     form.Preview(),
        Execute: func() error {
            return m.linear.CreateIssue(form)
        },
    }

    return m.WithPendingConfirm(confirm), nil
}

// :new-branch handler
func handleNewBranch(m Model, args []string) (Model, tea.Cmd) {
    issue := m.focusedIssue()
    if issue == nil {
        return m.WithError("No issue selected"), nil
    }

    branchName := fmt.Sprintf("%s-%s",
        issue.Identifier,
        slugify(issue.Title),
    )

    confirm := ConfirmRequest{
        Action:      "Create Branch",
        Description: fmt.Sprintf("git checkout -b %s", branchName),
        Execute: func() error {
            return exec.Command("git", "checkout", "-b", branchName).Run()
        },
    }

    return m.WithPendingConfirm(confirm), nil
}
```

### UI: Command Palette

```
┌─────────────────────────────────────────────────┐
│  :                                              │
│  ───────────────────────────────────────────── │
│  CREATE                                         │
│  > :new-issue      Create Linear issue         │
│    :new-branch     Create git branch           │
│    :new-pr         Open PR linked to issue     │
│                                                 │
│  GOLDEN PATHS                                   │
│    :start-work     Branch → Editor → Status    │
│    :submit-review  Commit → Push → PR          │
│                                                 │
│  AI                                             │
│    :ask-claude     Ask about current context   │
│                                                 │
│  [Enter] Execute  [Tab] Autocomplete  [Esc] Cancel │
└─────────────────────────────────────────────────┘
```

### UI: Confirmation Dialog

```
┌─────────────────────────────────────────────────┐
│  CONFIRM ACTION                                 │
├─────────────────────────────────────────────────┤
│                                                 │
│  Action: Create Issue                           │
│  Target: Linear / Project Alpha                 │
│                                                 │
│  Preview:                                       │
│  ┌─────────────────────────────────────────┐   │
│  │ Title: Fix login timeout                │   │
│  │ Labels: bug, auth                       │   │
│  │ Assignee: You                           │   │
│  │ Project: Project Alpha                  │   │
│  └─────────────────────────────────────────┘   │
│                                                 │
│  [y] Confirm  [n] Cancel  [e] Edit             │
└─────────────────────────────────────────────────┘
```

## Consequences

### Positive
- **Productivity**: Create without leaving MAAT
- **Consistency**: Golden paths enforce team standards
- **Context Preservation**: Actions happen in context
- **IDP Capability**: MAAT becomes a platform, not just viewer

### Negative
- **Complexity**: More UI states and flows
- **Confirmation Fatigue**: Many confirmations for workflows
- **Scope Creep Risk**: Temptation to add more actions

### Mitigations
- Batch confirmations for golden paths
- "Trust mode" for experienced users (still logged)
- Strict action budget (max 20 commands total)

## Compliance

This ADR enforces:
- **Commandment #6**: All actions require explicit invocation
- **Commandment #10**: Confirm before writing to external systems

## References

- Internal Developer Platform Trends: [Software Development Trends 2025](https://www.codestringers.com/insights/software-development-trends/)
- MAAT Constitution: `/Users/manu/Documents/LUXOR/MAAT/docs/CONSTITUTION.md`
