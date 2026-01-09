package tui

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines all keybindings
type KeyMap struct {
	Quit        key.Binding
	Enter       key.Binding
	Back        key.Binding
	Up          key.Binding
	Down        key.Binding
	Left        key.Binding
	Right       key.Binding
	Help        key.Binding
	Refresh     key.Binding
	AI          key.Binding
	OpenBrowser key.Binding
	CopyURL     key.Binding
}

// DefaultKeyMap returns the default keybindings
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c", "q"),
			key.WithHelp("ctrl+c/q", "quit"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "drill down"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "left"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "right"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		Refresh: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "refresh"),
		),
		AI: key.NewBinding(
			key.WithKeys("ctrl+a"),
			key.WithHelp("ctrl+a", "invoke AI"),
		),
		OpenBrowser: key.NewBinding(
			key.WithKeys("o"),
			key.WithHelp("o", "open in browser"),
		),
		CopyURL: key.NewBinding(
			key.WithKeys("y"),
			key.WithHelp("y", "copy URL"),
		),
	}
}

// ShortHelp returns a slice of key bindings for the short help view
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Enter, k.Back, k.Quit, k.Help}
}

// FullHelp returns a slice of key bindings for the full help view
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Enter, k.Back, k.Refresh, k.AI},
		{k.OpenBrowser, k.CopyURL, k.Help, k.Quit},
	}
}
