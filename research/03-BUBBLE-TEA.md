# Bubble Tea Framework: Go TUI Best Practices

**Source**: Charmbracelet
**Research Date**: 2026-01-05
**Relevance to MAAT**: Critical - Primary implementation framework

---

## The Elm Architecture in Go

Bubble Tea implements The Elm Architecture (TEA) with three core methods:

```go
type Model interface {
    Init() tea.Cmd
    Update(msg tea.Msg) (tea.Model, tea.Cmd)
    View() string
}
```

### Core Components

| Component | Purpose |
|-----------|---------|
| **Model** | Central state container holding all application data |
| **Init()** | Returns initial command for startup operations |
| **Update()** | Pure function handling events and state transitions |
| **View()** | Pure function rendering UI from model state |

### Key Benefits

- Eliminates race conditions by channeling all state changes through Update
- Unidirectional data flow: Message -> Update -> Model -> View
- Makes state predictable and debuggable
- Works naturally with Go's type system

---

## Component Composition Patterns

### Level 1: Top-Down Composition

Parent model embeds child models as fields:

```go
type AppModel struct {
    sidebar   SidebarModel
    viewport  ViewportModel
    statusbar StatusBarModel
}
```

### Level 2: Model Stack Architecture

Independent models composed into application hierarchy.

### Level 3: Hybrid Approach

Combines both for maximum flexibility.

### Key Techniques

- **Lazy initialization**: Defer setup until `WindowSizeMsg` arrives
- **Embedding Bubbles components**: viewport, spinner, text input
- **Layout composition**: Styled header/footer regions

---

## Async/Command Patterns

Commands (`tea.Cmd`) are functions returning messages, executed asynchronously:

```go
func fetchData(id int) tea.Cmd {
    return func() tea.Msg {
        data, _ := api.Get(id)
        return dataMsg{data}
    }
}
```

### Critical Rules

1. **Never spawn your own goroutines** - use commands exclusively
2. **`tea.Batch()`** - Run multiple commands concurrently
3. **`tea.Sequence()`** - Execute commands sequentially
4. **All I/O must go through commands** - network, disk, timers

---

## Styling with Lipgloss

Lipgloss provides CSS-like declarative styling:

```go
style := lipgloss.NewStyle().
    Bold(true).
    Foreground(lipgloss.Color("205")).
    Padding(0, 1)
```

### Features

- **Adaptive Colors**: `lipgloss.AdaptiveColor` for light/dark terminals
- **Layout Tools**: Padding, margins, borders, alignment
- **Automatic Degradation**: Colors degrade based on terminal capabilities
- **TTY Detection**: Strips colors when piping output

---

## Performance Optimizations

Built-in optimizations:

- **Framerate-based renderer**: Limits redraws to prevent flickering
- **Differential rendering**: Only changed lines are redrawn
- **High-performance viewport**: Specialized renderer for complex ANSI
- **Central message channel**: Safe goroutine communication

---

## Application to MAAT Design

| Concern | Bubble Tea Pattern |
|---------|-------------------|
| **Graph state** | Single model containing node/edge data, focus state |
| **Navigation stack** | Nested models representing drill-down levels |
| **Graph rendering** | Custom View() with Lipgloss-styled node boxes |
| **Async data loading** | `tea.Cmd` for fetching node details on expand |
| **Keyboard navigation** | Handle arrow keys, enter, escape in Update() |
| **Responsive layout** | Track `WindowSizeMsg`, calculate viewport bounds |
| **Component reuse** | Embed Bubbles viewport for scrollable regions |

---

## Recommended Architecture for MAAT

Use a **hybrid model stack**:

1. **Root Model**: Navigation state, global keybindings
2. **Graph Viewport Model**: Node rendering, focus management
3. **Detail Panel Model**: Issue/PR information display
4. **Search Model**: Fuzzy finder, ripgrep integration

Commands fetch data asynchronously while UI remains responsive through framerate-limited rendering.

---

## Sources

- [Bubble Tea GitHub](https://github.com/charmbracelet/bubbletea)
- [Lipgloss GitHub](https://github.com/charmbracelet/lipgloss)
- [Bubbles Components](https://github.com/charmbracelet/bubbles)
- [Commands Tutorial](https://charm.land/blog/commands-in-bubbletea/)
- [Managing Nested Models](https://donderom.com/posts/managing-nested-models-with-bubble-tea/)
