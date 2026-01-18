package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// KeyMap defines the key bindings using bubbles key.Binding
type KeyMap struct {
	Up         key.Binding
	Down       key.Binding
	PageUp     key.Binding
	PageDown   key.Binding
	GotoTop    key.Binding
	GotoBottom key.Binding
	Done       key.Binding
	Add        key.Binding
	Drop       key.Binding
	Bump       key.Binding
	Help       key.Binding
	Quit       key.Binding
	Confirm    key.Binding
	Cancel     key.Binding
	Tab        key.Binding
	Left       key.Binding
	Right      key.Binding
}

// DefaultKeyMap returns the default key bindings
var DefaultKeyMap = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("↓/j", "down"),
	),
	PageUp: key.NewBinding(
		key.WithKeys("pgup", "ctrl+u"),
		key.WithHelp("pgup/^u", "page up"),
	),
	PageDown: key.NewBinding(
		key.WithKeys("pgdown", "ctrl+d"),
		key.WithHelp("pgdn/^d", "page down"),
	),
	GotoTop: key.NewBinding(
		key.WithKeys("home", "g"),
		key.WithHelp("g/home", "go to top"),
	),
	GotoBottom: key.NewBinding(
		key.WithKeys("end", "G"),
		key.WithHelp("G/end", "go to bottom"),
	),
	Done: key.NewBinding(
		key.WithKeys("enter", "d"),
		key.WithHelp("enter/d", "complete"),
	),
	Add: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add"),
	),
	Drop: key.NewBinding(
		key.WithKeys("x"),
		key.WithHelp("x", "drop"),
	),
	Bump: key.NewBinding(
		key.WithKeys("b"),
		key.WithHelp("b", "bump"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc"),
		key.WithHelp("q", "quit"),
	),
	Confirm: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "confirm"),
	),
	Cancel: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab", "shift+tab"),
		key.WithHelp("tab", "next field"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "right"),
	),
}

// matchesKey checks if a key message matches a key binding
func matchesKey(msg tea.KeyMsg, binding key.Binding) bool {
	return key.Matches(msg, binding)
}
