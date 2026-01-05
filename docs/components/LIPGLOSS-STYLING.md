# Lipgloss Styling Guide

**Source**: Context7 + Official Documentation
**Library**: `github.com/charmbracelet/lipgloss`
**Version**: v0.9+ (stable), v2.0.0-beta available
**Purpose**: CSS-like terminal styling for MAAT UI

---

## DELTA FORCE: Exact Code for MAAT

### 1. Color Palette Definition

**MAAT Requirement**: Consistent visual hierarchy
**Pattern**: Define colors once, use everywhere

```go
// maat/internal/tui/styles/colors.go
package styles

import "github.com/charmbracelet/lipgloss"

// Adaptive colors for light/dark terminals
var (
    // Primary palette
    Primary   = lipgloss.AdaptiveColor{Light: "#5A56E0", Dark: "#7C78FF"}
    Secondary = lipgloss.AdaptiveColor{Light: "#6E6AE0", Dark: "#8E8AFF"}
    Accent    = lipgloss.AdaptiveColor{Light: "#00D084", Dark: "#00E898"}

    // Status colors (Linear-inspired)
    StatusTodo       = lipgloss.Color("#6B7280") // Gray
    StatusInProgress = lipgloss.Color("#F59E0B") // Amber
    StatusDone       = lipgloss.Color("#10B981") // Green
    StatusCanceled   = lipgloss.Color("#EF4444") // Red
    StatusBlocked    = lipgloss.Color("#8B5CF6") // Purple

    // Priority colors
    PriorityUrgent = lipgloss.Color("#DC2626")
    PriorityHigh   = lipgloss.Color("#F97316")
    PriorityMedium = lipgloss.Color("#FBBF24")
    PriorityLow    = lipgloss.Color("#6B7280")

    // UI colors
    Background = lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#1A1A2E"}
    Foreground = lipgloss.AdaptiveColor{Light: "#1A1A2E", Dark: "#E4E4E7"}
    Border     = lipgloss.AdaptiveColor{Light: "#E4E4E7", Dark: "#3F3F46"}
    Muted      = lipgloss.AdaptiveColor{Light: "#9CA3AF", Dark: "#71717A"}

    // Git colors
    GitAdded    = lipgloss.Color("#22C55E")
    GitModified = lipgloss.Color("#EAB308")
    GitDeleted  = lipgloss.Color("#EF4444")
)
```

### 2. Component Styles

**MAAT Requirement**: Visual consistency across views

```go
// maat/internal/tui/styles/components.go
package styles

import "github.com/charmbracelet/lipgloss"

// Node styles for graph view
var (
    // Issue node
    IssueNode = lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(Primary).
        Padding(0, 1).
        Width(30)

    IssueNodeFocused = IssueNode.Copy().
        BorderForeground(Accent).
        Bold(true)

    // PR node
    PRNode = lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(Secondary).
        Padding(0, 1).
        Width(30)

    // File node
    FileNode = lipgloss.NewStyle().
        Border(lipgloss.NormalBorder()).
        BorderForeground(Muted).
        Padding(0, 1).
        Width(25)

    // Commit node (for git history graph)
    CommitNode = lipgloss.NewStyle().
        Foreground(Foreground).
        Padding(0, 1)

    CommitNodeFocused = CommitNode.Copy().
        Background(lipgloss.Color("#2D2D3A")).
        Bold(true)
)

// Edge styles
var (
    EdgeBlocks = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#EF4444"))

    EdgeRelated = lipgloss.NewStyle().
        Foreground(Muted)

    EdgeImplements = lipgloss.NewStyle().
        Foreground(Accent)
)

// Panel styles
var (
    // Main panels
    GraphPanel = lipgloss.NewStyle().
        Border(lipgloss.NormalBorder()).
        BorderForeground(Border).
        Padding(1)

    DetailPanel = lipgloss.NewStyle().
        Border(lipgloss.NormalBorder()).
        BorderForeground(Border).
        Padding(1)

    // Active panel indicator
    PanelActive = lipgloss.NewStyle().
        Border(lipgloss.DoubleBorder()).
        BorderForeground(Accent).
        Padding(1)
)

// Status indicators
var (
    StatusBadge = func(status string) lipgloss.Style {
        var color lipgloss.Color
        switch status {
        case "todo":
            color = StatusTodo
        case "in_progress":
            color = StatusInProgress
        case "done":
            color = StatusDone
        case "canceled":
            color = StatusCanceled
        default:
            color = Muted
        }
        return lipgloss.NewStyle().
            Background(color).
            Foreground(lipgloss.Color("#FFFFFF")).
            Padding(0, 1).
            Bold(true)
    }

    PriorityBadge = func(priority int) lipgloss.Style {
        var color lipgloss.Color
        switch priority {
        case 1:
            color = PriorityUrgent
        case 2:
            color = PriorityHigh
        case 3:
            color = PriorityMedium
        default:
            color = PriorityLow
        }
        return lipgloss.NewStyle().
            Foreground(color).
            Bold(true)
    }
)
```

### 3. Layout Helpers

**MAAT Requirement**: Responsive 3-pane layout

```go
// maat/internal/tui/styles/layout.go
package styles

import "github.com/charmbracelet/lipgloss"

// Layout calculations
func CalculateLayout(width, height int) Layout {
    // Reserve for borders and padding
    contentWidth := width - 4
    contentHeight := height - 4

    // 3-pane layout: 25% | 50% | 25%
    graphWidth := contentWidth * 25 / 100
    mainWidth := contentWidth * 50 / 100
    detailWidth := contentWidth - graphWidth - mainWidth

    return Layout{
        GraphWidth:  graphWidth,
        MainWidth:   mainWidth,
        DetailWidth: detailWidth,
        Height:      contentHeight,
    }
}

type Layout struct {
    GraphWidth  int
    MainWidth   int
    DetailWidth int
    Height      int
}

// Horizontal join with spacing
func JoinHorizontal(left, middle, right string) string {
    return lipgloss.JoinHorizontal(
        lipgloss.Top,
        left,
        middle,
        right,
    )
}

// Vertical join for stacked elements
func JoinVertical(elements ...string) string {
    return lipgloss.JoinVertical(lipgloss.Left, elements...)
}

// Center content in available space
func Center(content string, width, height int) string {
    return lipgloss.Place(
        width, height,
        lipgloss.Center, lipgloss.Center,
        content,
    )
}
```

### 4. Breadcrumb Style

**MAAT Requirement**: FR-003 Navigation context

```go
// maat/internal/tui/styles/breadcrumb.go
package styles

import (
    "strings"
    "github.com/charmbracelet/lipgloss"
)

var (
    BreadcrumbContainer = lipgloss.NewStyle().
        Foreground(Muted).
        Padding(0, 1)

    BreadcrumbItem = lipgloss.NewStyle().
        Foreground(Foreground)

    BreadcrumbCurrent = lipgloss.NewStyle().
        Foreground(Accent).
        Bold(true)

    BreadcrumbSeparator = lipgloss.NewStyle().
        Foreground(Muted).
        SetString(" › ")
)

func RenderBreadcrumb(path []string) string {
    if len(path) == 0 {
        return ""
    }

    var parts []string
    for i, item := range path {
        if i == len(path)-1 {
            parts = append(parts, BreadcrumbCurrent.Render(item))
        } else {
            parts = append(parts, BreadcrumbItem.Render(item))
        }
    }

    return BreadcrumbContainer.Render(
        strings.Join(parts, BreadcrumbSeparator.String()),
    )
}
```

### 5. Confirmation Dialog Style

**MAAT Requirement**: A-003 Write confirmation

```go
// maat/internal/tui/styles/dialog.go
package styles

import "github.com/charmbracelet/lipgloss"

var (
    DialogBox = lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(Accent).
        Padding(1, 2).
        Width(50)

    DialogTitle = lipgloss.NewStyle().
        Bold(true).
        Foreground(Foreground).
        MarginBottom(1)

    DialogContent = lipgloss.NewStyle().
        Foreground(Foreground)

    DialogPreview = lipgloss.NewStyle().
        Border(lipgloss.NormalBorder()).
        BorderForeground(Border).
        Padding(0, 1).
        MarginTop(1).
        MarginBottom(1)

    DialogButtons = lipgloss.NewStyle().
        MarginTop(1)

    ButtonConfirm = lipgloss.NewStyle().
        Background(Accent).
        Foreground(lipgloss.Color("#000000")).
        Padding(0, 2).
        Bold(true)

    ButtonCancel = lipgloss.NewStyle().
        Background(Muted).
        Foreground(lipgloss.Color("#FFFFFF")).
        Padding(0, 2)

    ButtonEdit = lipgloss.NewStyle().
        Background(Primary).
        Foreground(lipgloss.Color("#FFFFFF")).
        Padding(0, 2)
)

func RenderDialog(title, content, preview string) string {
    return DialogBox.Render(
        lipgloss.JoinVertical(
            lipgloss.Left,
            DialogTitle.Render(title),
            DialogContent.Render(content),
            DialogPreview.Render(preview),
            DialogButtons.Render(
                lipgloss.JoinHorizontal(
                    lipgloss.Top,
                    ButtonConfirm.Render("[y] Confirm"),
                    "  ",
                    ButtonCancel.Render("[n] Cancel"),
                    "  ",
                    ButtonEdit.Render("[e] Edit"),
                ),
            ),
        ),
    )
}
```

### 6. Git History Styling

**MAAT Requirement**: FR-010 Git commit graph

```go
// maat/internal/tui/styles/git.go
package styles

import "github.com/charmbracelet/lipgloss"

var (
    // Commit graph characters
    CommitDot     = "●"
    CommitLine    = "│"
    CommitMerge   = "─┬"
    CommitBranch  = "├"
    CommitEnd     = "└"

    // Commit line styles
    CommitHash = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#F59E0B")).
        Bold(true)

    CommitMessage = lipgloss.NewStyle().
        Foreground(Foreground)

    CommitAuthor = lipgloss.NewStyle().
        Foreground(Muted).
        Italic(true)

    CommitDate = lipgloss.NewStyle().
        Foreground(Muted)

    CommitRef = lipgloss.NewStyle().
        Background(Primary).
        Foreground(lipgloss.Color("#FFFFFF")).
        Padding(0, 1)

    CommitTag = lipgloss.NewStyle().
        Background(lipgloss.Color("#F59E0B")).
        Foreground(lipgloss.Color("#000000")).
        Padding(0, 1)

    // Diff stats
    DiffAdditions = lipgloss.NewStyle().
        Foreground(GitAdded)

    DiffDeletions = lipgloss.NewStyle().
        Foreground(GitDeleted)

    DiffFile = lipgloss.NewStyle().
        Foreground(Foreground)

    DiffFileAdded = lipgloss.NewStyle().
        Foreground(GitAdded).
        Bold(true)

    DiffFileModified = lipgloss.NewStyle().
        Foreground(GitModified)

    DiffFileDeleted = lipgloss.NewStyle().
        Foreground(GitDeleted).
        Strikethrough(true)
)

// Render a commit line
func RenderCommitLine(hash, message, author, date string, isHead bool) string {
    hashStyled := CommitHash.Render(hash)
    msgStyled := CommitMessage.Render(message)
    authorStyled := CommitAuthor.Render(author)
    dateStyled := CommitDate.Render(date)

    prefix := CommitDot
    if isHead {
        prefix = CommitRef.Render("HEAD") + " " + CommitDot
    }

    return lipgloss.JoinHorizontal(
        lipgloss.Top,
        prefix,
        " ",
        hashStyled,
        "  ",
        msgStyled,
        "  ",
        authorStyled,
        " • ",
        dateStyled,
    )
}

// Render diff stats
func RenderDiffStats(additions, deletions int) string {
    return lipgloss.JoinHorizontal(
        lipgloss.Top,
        DiffAdditions.Render(fmt.Sprintf("+%d", additions)),
        " ",
        DiffDeletions.Render(fmt.Sprintf("-%d", deletions)),
    )
}
```

---

## Style Hierarchy

```
styles/
├── colors.go      # Color definitions
├── components.go  # Component styles
├── layout.go      # Layout calculations
├── breadcrumb.go  # Breadcrumb styling
├── dialog.go      # Dialog/modal styles
├── git.go         # Git history styles
└── theme.go       # Theme switching (future)
```

---

## Performance Tips

1. **Reuse styles**: Define styles at package level, not in render loops
2. **Copy for variants**: Use `Style.Copy()` to create variants
3. **Width constraints**: Always set `Width()` for layout stability
4. **Adaptive colors**: Use `AdaptiveColor` for light/dark terminal support
