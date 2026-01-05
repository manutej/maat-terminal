# Bubbles Widget Integration

**Source**: Context7 + Official Documentation
**Library**: `github.com/charmbracelet/bubbles`
**Version**: v0.18+
**Purpose**: Pre-built TUI components for MAAT

---

## DELTA FORCE: Exact Code for MAAT

### 1. Viewport (Scrollable Content)

**MAAT Requirement**: FR-003 Detail pane, FR-010 Commit details
**Use Case**: Scrollable issue descriptions, commit diffs

```go
// maat/internal/tui/components/detail_viewport.go
package components

import (
    "github.com/charmbracelet/bubbles/viewport"
    tea "github.com/charmbracelet/bubbletea"
    "maat/internal/tui/styles"
)

type DetailViewport struct {
    viewport viewport.Model
    title    string
    ready    bool
}

func NewDetailViewport() DetailViewport {
    return DetailViewport{}
}

func (d DetailViewport) Init() tea.Cmd {
    return nil
}

func (d DetailViewport) Update(msg tea.Msg) (DetailViewport, tea.Cmd) {
    var cmd tea.Cmd

    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        if !d.ready {
            d.viewport = viewport.New(msg.Width/4, msg.Height-6)
            d.viewport.HighPerformanceRendering = true
            d.ready = true
        } else {
            d.viewport.Width = msg.Width / 4
            d.viewport.Height = msg.Height - 6
        }
    }

    d.viewport, cmd = d.viewport.Update(msg)
    return d, cmd
}

func (d DetailViewport) View() string {
    if !d.ready {
        return "Loading..."
    }

    header := styles.DialogTitle.Render(d.title)
    content := d.viewport.View()
    scrollInfo := fmt.Sprintf("%3.f%%", d.viewport.ScrollPercent()*100)

    return styles.DetailPanel.Render(
        lipgloss.JoinVertical(
            lipgloss.Left,
            header,
            content,
            styles.Muted.Render(scrollInfo),
        ),
    )
}

func (d *DetailViewport) SetContent(title, content string) {
    d.title = title
    d.viewport.SetContent(content)
    d.viewport.GotoTop()
}

// Delegate key handling
func (d DetailViewport) HandleKey(key tea.KeyMsg) (DetailViewport, tea.Cmd) {
    switch key.String() {
    case "j", "down":
        d.viewport.LineDown(1)
    case "k", "up":
        d.viewport.LineUp(1)
    case "d":
        d.viewport.HalfViewDown()
    case "u":
        d.viewport.HalfViewUp()
    case "g":
        d.viewport.GotoTop()
    case "G":
        d.viewport.GotoBottom()
    }
    return d, nil
}
```

### 2. List (Issue/PR List)

**MAAT Requirement**: FR-001 Issue display, FR-007 Role-based views
**Use Case**: Issue list, PR list, commit list

```go
// maat/internal/tui/components/issue_list.go
package components

import (
    "github.com/charmbracelet/bubbles/list"
    tea "github.com/charmbracelet/bubbletea"
    "maat/internal/linear"
    "maat/internal/tui/styles"
)

// IssueItem implements list.Item interface
type IssueItem struct {
    issue linear.Issue
}

func (i IssueItem) Title() string {
    return fmt.Sprintf("%s %s", i.issue.Identifier, i.issue.Title)
}

func (i IssueItem) Description() string {
    return fmt.Sprintf("%s • %s", i.issue.State.Name, i.issue.Assignee.Name)
}

func (i IssueItem) FilterValue() string {
    return i.issue.Title + " " + i.issue.Identifier
}

// Custom delegate for issue rendering
type IssueDelegate struct {
    list.DefaultDelegate
}

func NewIssueDelegate() IssueDelegate {
    d := list.NewDefaultDelegate()

    // Customize styles
    d.Styles.SelectedTitle = d.Styles.SelectedTitle.
        Foreground(styles.Accent).
        Bold(true)

    d.Styles.SelectedDesc = d.Styles.SelectedDesc.
        Foreground(styles.Muted)

    return IssueDelegate{DefaultDelegate: d}
}

// IssueList component
type IssueList struct {
    list    list.Model
    focused bool
}

func NewIssueList(width, height int) IssueList {
    delegate := NewIssueDelegate()
    l := list.New([]list.Item{}, delegate, width, height)

    l.Title = "Issues"
    l.SetShowStatusBar(true)
    l.SetFilteringEnabled(true)
    l.SetShowHelp(false)

    // Custom key bindings
    l.KeyMap.CursorUp.SetKeys("k", "up")
    l.KeyMap.CursorDown.SetKeys("j", "down")

    return IssueList{list: l}
}

func (il IssueList) Update(msg tea.Msg) (IssueList, tea.Cmd) {
    var cmd tea.Cmd
    il.list, cmd = il.list.Update(msg)
    return il, cmd
}

func (il IssueList) View() string {
    style := styles.GraphPanel
    if il.focused {
        style = styles.PanelActive
    }
    return style.Render(il.list.View())
}

func (il *IssueList) SetIssues(issues []linear.Issue) {
    items := make([]list.Item, len(issues))
    for i, issue := range issues {
        items[i] = IssueItem{issue: issue}
    }
    il.list.SetItems(items)
}

func (il *IssueList) SetFocused(focused bool) {
    il.focused = focused
}

func (il IssueList) SelectedIssue() *linear.Issue {
    item := il.list.SelectedItem()
    if item == nil {
        return nil
    }
    issueItem := item.(IssueItem)
    return &issueItem.issue
}
```

### 3. Table (Structured Data)

**MAAT Requirement**: FR-010 Commit list, file changes
**Use Case**: Commit table, file change list

```go
// maat/internal/tui/components/commit_table.go
package components

import (
    "github.com/charmbracelet/bubbles/table"
    tea "github.com/charmbracelet/bubbletea"
    "maat/internal/tui/styles"
)

type CommitTable struct {
    table   table.Model
    commits []Commit
}

func NewCommitTable(width, height int) CommitTable {
    columns := []table.Column{
        {Title: "Hash", Width: 8},
        {Title: "Message", Width: width - 40},
        {Title: "Author", Width: 15},
        {Title: "Date", Width: 12},
    }

    t := table.New(
        table.WithColumns(columns),
        table.WithFocused(true),
        table.WithHeight(height),
    )

    // Style the table
    s := table.DefaultStyles()
    s.Header = s.Header.
        BorderStyle(lipgloss.NormalBorder()).
        BorderForeground(styles.Border).
        BorderBottom(true).
        Bold(true)

    s.Selected = s.Selected.
        Foreground(styles.Foreground).
        Background(lipgloss.Color("#2D2D3A")).
        Bold(true)

    t.SetStyles(s)

    return CommitTable{table: t}
}

func (ct CommitTable) Update(msg tea.Msg) (CommitTable, tea.Cmd) {
    var cmd tea.Cmd
    ct.table, cmd = ct.table.Update(msg)
    return ct, cmd
}

func (ct CommitTable) View() string {
    return ct.table.View()
}

func (ct *CommitTable) SetCommits(commits []Commit) {
    ct.commits = commits

    rows := make([]table.Row, len(commits))
    for i, c := range commits {
        rows[i] = table.Row{
            c.Hash[:7],
            truncate(c.Message, 50),
            c.Author,
            c.Date.Format("Jan 2"),
        }
    }
    ct.table.SetRows(rows)
}

func (ct CommitTable) SelectedCommit() *Commit {
    idx := ct.table.Cursor()
    if idx >= 0 && idx < len(ct.commits) {
        return &ct.commits[idx]
    }
    return nil
}

func truncate(s string, max int) string {
    if len(s) <= max {
        return s
    }
    return s[:max-3] + "..."
}
```

### 4. Text Input (Search, Commands)

**MAAT Requirement**: FR-008 Fuzzy search, FR-005 Claude input
**Use Case**: Search box, command palette, Claude prompt

```go
// maat/internal/tui/components/search_input.go
package components

import (
    "github.com/charmbracelet/bubbles/textinput"
    tea "github.com/charmbracelet/bubbletea"
    "maat/internal/tui/styles"
)

type SearchInput struct {
    input   textinput.Model
    results []SearchResult
}

func NewSearchInput() SearchInput {
    ti := textinput.New()
    ti.Placeholder = "Search issues, PRs, files..."
    ti.Prompt = "/ "
    ti.PromptStyle = lipgloss.NewStyle().Foreground(styles.Accent)
    ti.TextStyle = lipgloss.NewStyle().Foreground(styles.Foreground)
    ti.Cursor.Style = lipgloss.NewStyle().Foreground(styles.Accent)
    ti.CharLimit = 100
    ti.Width = 40

    return SearchInput{input: ti}
}

func (s SearchInput) Init() tea.Cmd {
    return textinput.Blink
}

func (s SearchInput) Update(msg tea.Msg) (SearchInput, tea.Cmd) {
    var cmd tea.Cmd

    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "esc":
            s.input.Blur()
            return s, nil
        case "enter":
            // Trigger search
            return s, s.doSearch()
        }
    }

    s.input, cmd = s.input.Update(msg)
    return s, cmd
}

func (s SearchInput) View() string {
    return styles.DialogBox.Render(
        lipgloss.JoinVertical(
            lipgloss.Left,
            s.input.View(),
            s.renderResults(),
        ),
    )
}

func (s *SearchInput) Focus() tea.Cmd {
    s.input.Focus()
    return textinput.Blink
}

func (s *SearchInput) Blur() {
    s.input.Blur()
}

func (s SearchInput) Value() string {
    return s.input.Value()
}

func (s *SearchInput) Reset() {
    s.input.Reset()
    s.results = nil
}

func (s SearchInput) doSearch() tea.Cmd {
    query := s.input.Value()
    return func() tea.Msg {
        // Return search results message
        return SearchResultsMsg{Query: query}
    }
}

func (s SearchInput) renderResults() string {
    if len(s.results) == 0 {
        return styles.Muted.Render("No results")
    }

    var lines []string
    for i, r := range s.results {
        line := fmt.Sprintf("%d. [%s] %s", i+1, r.Type, r.Title)
        lines = append(lines, line)
    }
    return strings.Join(lines, "\n")
}
```

### 5. Spinner (Loading States)

**MAAT Requirement**: Async operations feedback
**Use Case**: Data fetching, sync operations

```go
// maat/internal/tui/components/loading.go
package components

import (
    "github.com/charmbracelet/bubbles/spinner"
    tea "github.com/charmbracelet/bubbletea"
    "maat/internal/tui/styles"
)

type Loading struct {
    spinner spinner.Model
    message string
}

func NewLoading(message string) Loading {
    s := spinner.New()
    s.Spinner = spinner.Dot
    s.Style = lipgloss.NewStyle().Foreground(styles.Accent)

    return Loading{
        spinner: s,
        message: message,
    }
}

func (l Loading) Init() tea.Cmd {
    return l.spinner.Tick
}

func (l Loading) Update(msg tea.Msg) (Loading, tea.Cmd) {
    var cmd tea.Cmd

    switch msg := msg.(type) {
    case spinner.TickMsg:
        l.spinner, cmd = l.spinner.Update(msg)
    }

    return l, cmd
}

func (l Loading) View() string {
    return fmt.Sprintf("%s %s", l.spinner.View(), l.message)
}

func (l *Loading) SetMessage(message string) {
    l.message = message
}
```

---

## Component Composition Pattern

```go
// maat/internal/tui/model.go
type Model struct {
    // Embedded Bubbles components
    issueList    components.IssueList
    prList       components.IssueList  // Reuse for PRs
    commitTable  components.CommitTable
    detailView   components.DetailViewport
    searchInput  components.SearchInput
    loading      components.Loading

    // Active component tracking
    activePane   Pane
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmds []tea.Cmd

    // Update active component
    switch m.activePane {
    case PaneIssues:
        var cmd tea.Cmd
        m.issueList, cmd = m.issueList.Update(msg)
        cmds = append(cmds, cmd)
    case PaneDetail:
        var cmd tea.Cmd
        m.detailView, cmd = m.detailView.Update(msg)
        cmds = append(cmds, cmd)
    case PaneSearch:
        var cmd tea.Cmd
        m.searchInput, cmd = m.searchInput.Update(msg)
        cmds = append(cmds, cmd)
    }

    return m, tea.Batch(cmds...)
}

func (m Model) View() string {
    // Compose views
    left := m.issueList.View()
    middle := m.renderGraph()
    right := m.detailView.View()

    main := styles.JoinHorizontal(left, middle, right)

    if m.searchOpen {
        // Overlay search
        return styles.Center(m.searchInput.View(), m.width, m.height)
    }

    return main
}
```

---

## Widget Selection Matrix

| MAAT Feature | Bubbles Widget | Customization |
|--------------|----------------|---------------|
| Issue list | `list.Model` | Custom delegate, styles |
| PR list | `list.Model` | Different item type |
| Commit table | `table.Model` | Custom columns |
| Detail pane | `viewport.Model` | High-perf rendering |
| Search | `textinput.Model` | Custom prompt, fuzzy |
| Loading | `spinner.Model` | Dot spinner, accent color |
| Command palette | `textinput.Model` | Autocomplete |

---

## File Structure

```
maat/internal/tui/components/
├── issue_list.go      # Issue/PR list
├── commit_table.go    # Commit table
├── detail_viewport.go # Scrollable detail
├── search_input.go    # Search overlay
├── loading.go         # Loading spinner
├── confirm_dialog.go  # Confirmation modal
└── breadcrumb.go      # Navigation breadcrumb
```
