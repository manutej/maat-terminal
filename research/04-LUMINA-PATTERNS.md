# LUMINA Project: Markdown TUI Implementation Patterns

**Source**: Local codebase (LUXOR/PROJECTS/LUMINA/ccn)
**Research Date**: 2026-01-05
**Relevance to MAAT**: Critical - Production-validated patterns

---

## Project Overview

LUMINA (Claude Code Navigator) is a production-ready TUI application:

- **Language**: Go (~6,500 lines across 23 files)
- **Framework**: Charm Ecosystem (Bubble Tea, Glamour, Lip Gloss, Bubbles)
- **Layout**: 3-pane design for documentation browsing
- **Binary Size**: 14MB
- **Test Coverage**: 55+ tests

```
┌──────────────┬──────────────────────┬──────────────┐
│  File Tree   │    Viewer (Markdown) │   Preview    │
│    (1/4)     │        (2/4)         │    (1/4)     │
└──────────────┴──────────────────────┴──────────────┘
```

---

## State Machine Architecture

LUMINA uses 4 primary modes:

```go
const (
    NormalMode UIMode = iota  // Default interaction
    FinderMode                 // Fuzzy file search (/)
    SearchMode                 // Content search (Ctrl+F)
    HelpMode                   // Help overlay
)
```

**Transition Pattern**: Clean state transitions with cleanup of previous mode.

---

## Markdown Rendering Approach

### Glamour Integration

```go
type GlamourRendererImpl struct {
    renderer *glamour.TermRenderer  // Lazy-initialized
    theme    string                 // Auto-detected
    width    int                    // Terminal width
    themeMu, widthMu sync.RWMutex  // Thread safety
}
```

### Key Features

1. **Theme Detection**: Auto-detects light/dark from `COLORFGBG` environment
2. **Lazy Initialization**: Renderer created on first use (performance)
3. **Thread-Safe Rendering**: RWMutex protects width/theme during concurrent access
4. **Width Adaptation**: Adjusts to terminal width via WindowSizeMsg

---

## Vim-Style Keybindings

Fully configurable JSON-based system (`~/.config/lumina/keybindings.json`):

| Key | Action | Description |
|-----|--------|-------------|
| j/k | Navigate | Move down/up in lists |
| h/Backspace | Back | Return to parent |
| Enter | Open | Select item |
| d/u | Page | Half-page down/up |
| g/G | Jump | Top/bottom of document |
| / | Fuzzy Find | Search files by name |
| Ctrl+F | Search | Search file contents |
| ? | Help | Toggle help overlay |
| y | Copy | Copy to clipboard |
| Tab | Cycle | Cycle between panes |
| m | Mode | Cycle context panel modes |

### View-Specific Binding

Each binding targets specific views:

```go
type KeyBinding struct {
    Key    string  // "j"
    Action string  // "cursor_down"
    View   string  // "filetree", "viewer", "any"
}
```

---

## Component Structure

### AppModel (Single Source of Truth)

```go
type AppModel struct {
    // Window dimensions
    width, height int
    ready         bool

    // Navigation
    rootPath, currentPath string
    currentView           ViewMode  // FileTreeView, ViewerView, PreviewView
    currentMode           UIMode    // State machine state

    // UI Components (Bubbles)
    fileList    list.Model
    viewer      viewport.Model
    preview     viewport.Model

    // Rendering
    markdownRenderer *utils.MarkdownRenderer
    viewerContent    string
    renderedContent  string

    // Search (Flat State)
    finderInput    string
    finderFiltered []string
    searchQuery    string
    searchResults  []RipgrepResult

    // Context Panel
    contextPanel *ContextPanel
}
```

### Context Panel Multi-Mode

Right pane cycles through 4 modes:

1. **Table of Contents**: Navigate headings, jump with Enter
2. **File Info**: Metadata (size, word count, reading time)
3. **Document Stats**: Heading counts, code blocks, links
4. **Quick Actions**: Future features foundation

---

## Key Architectural Decisions

| Decision | Rationale |
|----------|-----------|
| **Value Receivers** | Used for Bubble Tea Update() compatibility |
| **Pointer Receivers** | For internal helpers that modify state |
| **Flat State Structure** | Finder/Search state directly in AppModel |
| **tea.Cmd for Async** | File watching, ripgrep via commands |
| **Thread Safety** | RWMutex for concurrent width/theme access |
| **Context Preservation** | Scroll position saved on file reload |

---

## Reusable Patterns for MAAT

### 1. Three-Pane Layout

Adapt LUMINA's proportional width system:
- 25% navigation
- 50% main content
- 25% context/details

### 2. State Machine Router

Use UIMode pattern for mode management:
- NormalMode (graph navigation)
- SearchMode (issue/PR search)
- DetailMode (drill-down view)

### 3. Keybindings Registry

JSON-based customization without code changes:
- User preferences persist
- Category-based organization
- View-specific bindings

### 4. TOC Jump Navigation

Implement for MAAT tree nodes:
- Store line numbers for context switching
- Heading extraction for structure navigation

### 5. Context Panel

Right pane cycles between:
- Node details
- Tree statistics
- Quick actions

### 6. Async Search

Use tea.Cmd + channels for background operations:
- Linear API calls
- GitHub API calls
- File watching

### 7. Viewport Management

Bubble Tea's viewport handles:
- Scrolling
- Selection
- Rendering

---

## File Structure Reference

```
ccn/
├── main.go            # Entry point, UI, message handling (~1200 lines)
├── model.go           # AppModel state definition (~650 lines)
├── keybindings.go     # Configurable keybindings (175 lines)
├── context_panel.go   # Multi-mode right pane (385 lines)
├── toc.go             # Table of Contents (237 lines)
├── clipboard.go       # Copy functionality (221 lines)
├── ripgrep.go         # Ripgrep search manager (~200 lines)
├── colors.go          # Color management
├── help.go            # Help overlay
├── glamour_impl.go    # Markdown rendering
└── *_test.go          # Test files (55+ tests)
```

---

## Sources

- Local codebase: `/Users/manu/Documents/LUXOR/PROJECTS/LUMINA/ccn/`
- CLAUDE.md project documentation
- Test files for implementation patterns
