package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"upnext/internal/model"
	"upnext/internal/store"
	"upnext/internal/ui"
)

// Mode represents the current UI mode
type Mode int

const (
	ModeNormal Mode = iota
	ModeInput
	ModeHelp
	ModeCelebration
	ModeConfirm
)

// Model is the main Bubble Tea model
type Model struct {
	data           *model.Data
	store          store.Store
	table          table.Model
	help           help.Model
	keys           KeyMap
	width          int
	height         int
	mode           Mode
	titleInput     textinput.Model
	descInput      textinput.Model
	priorityIndex  int
	inputFocus     int // 0 = title, 1 = description, 2 = priority
	celebrationMsg string
	confirmAction  string
	err            error
}

// celebrationTickMsg is sent to end the celebration animation
type celebrationTickMsg struct{}

// New creates a new TUI model
func New(s store.Store) (Model, error) {
	data, err := s.Load()
	if err != nil {
		return Model{}, err
	}

	// Create table with expanded columns for fuller view
	columns := []table.Column{
		{Title: "", Width: 3},              // Status icon
		{Title: "Pri", Width: 5},           // Priority
		{Title: "Task", Width: 35},         // Task text
		{Title: "Description", Width: 25},  // Description preview
		{Title: "Age", Width: 10},          // Age
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	// Style the table
	s2 := table.DefaultStyles()
	s2.Header = s2.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(ui.Surface1).
		BorderBottom(true).
		Bold(true).
		Foreground(ui.Mauve)
	s2.Selected = s2.Selected.
		Foreground(ui.Text).
		Background(ui.Surface0).
		Bold(true)
	s2.Cell = s2.Cell.
		Foreground(ui.Text)
	t.SetStyles(s2)

	// Create help
	h := help.New()
	h.Styles.ShortKey = lipgloss.NewStyle().Foreground(ui.Mauve).Bold(true)
	h.Styles.ShortDesc = lipgloss.NewStyle().Foreground(ui.Subtext0)
	h.Styles.ShortSeparator = lipgloss.NewStyle().Foreground(ui.Surface1)

	// Create inputs
	titleInput := textinput.New()
	titleInput.Placeholder = "What needs to be done?"
	titleInput.CharLimit = 100
	titleInput.Width = 40
	titleInput.PromptStyle = ui.FocusedStyle
	titleInput.TextStyle = lipgloss.NewStyle().Foreground(ui.Text)

	descInput := textinput.New()
	descInput.Placeholder = "Optional description..."
	descInput.CharLimit = 200
	descInput.Width = 40
	descInput.PromptStyle = ui.BlurredStyle
	descInput.TextStyle = lipgloss.NewStyle().Foreground(ui.Text)

	m := Model{
		data:          data,
		store:         s,
		table:         t,
		help:          h,
		keys:          DefaultKeyMap,
		width:         80,
		height:        24,
		mode:          ModeNormal,
		titleInput:    titleInput,
		descInput:     descInput,
		priorityIndex: 1, // Default to Medium
	}

	m.refreshTable()
	return m, nil
}

// refreshTable updates the table rows from the data
func (m *Model) refreshTable() {
	rows := make([]table.Row, len(m.data.Items))

	// Calculate column widths dynamically
	taskWidth := 33
	descWidth := 23
	if m.width > 0 {
		availableWidth := m.width - 30
		if availableWidth > 40 {
			taskWidth = availableWidth * 60 / 100 - 2
			descWidth = availableWidth - taskWidth - 4
		}
	}

	for i, item := range m.data.Items {
		desc := item.Description
		if desc == "" {
			desc = "-"
		}
		rows[i] = table.Row{
			ui.IconUnchecked,
			m.priorityIcon(item.Priority),
			truncateText(item.Text, taskWidth),
			truncateText(desc, descWidth),
			formatAge(item.Created),
		}
	}
	m.table.SetRows(rows)
}

func (m *Model) priorityIcon(p model.Priority) string {
	switch p {
	case model.PriorityHigh:
		return ui.PriorityHighStyle.Render(ui.IconHigh)
	case model.PriorityMedium:
		return ui.PriorityMediumStyle.Render(ui.IconMedium)
	default:
		return ui.PriorityLowStyle.Render(ui.IconLow)
	}
}

func truncateText(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}

func formatAge(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	case d < 7*24*time.Hour:
		return fmt.Sprintf("%dd ago", int(d.Hours()/24))
	default:
		return fmt.Sprintf("%dw ago", int(d.Hours()/(24*7)))
	}
}

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	return nil
}

// AddTodo adds a new todo item
func (m *Model) AddTodo(text, description string, priority model.Priority) {
	if text == "" {
		return
	}

	todo := model.Todo{
		ID:          model.GenerateID(),
		Text:        text,
		Description: description,
		Priority:    priority,
		Created:     time.Now(),
		Position:    0,
	}

	// Shift all positions down
	for i := range m.data.Items {
		m.data.Items[i].Position++
	}

	// Insert at the beginning
	m.data.Items = append([]model.Todo{todo}, m.data.Items...)
	m.refreshTable()
	m.table.SetCursor(0)
}

// CompleteTodo marks the current todo as done
func (m *Model) CompleteTodo() bool {
	cursor := m.table.Cursor()
	if len(m.data.Items) == 0 || cursor >= len(m.data.Items) {
		return false
	}

	item := m.data.Items[cursor]

	// Add to archive
	archived := model.ArchivedTodo{
		ID:        item.ID,
		Text:      item.Text,
		Created:   item.Created,
		Completed: time.Now(),
	}
	m.data.Archive = append(m.data.Archive, archived)

	// Remove from items
	m.data.Items = append(m.data.Items[:cursor], m.data.Items[cursor+1:]...)

	// Update stats
	m.data.Stats.TotalCompleted++

	m.refreshTable()
	return true
}

// DropTodo removes the current todo without archiving
func (m *Model) DropTodo() {
	cursor := m.table.Cursor()
	if len(m.data.Items) == 0 || cursor >= len(m.data.Items) {
		return
	}

	m.data.Items = append(m.data.Items[:cursor], m.data.Items[cursor+1:]...)
	m.refreshTable()
}

// BumpTodo moves the current todo to the top
func (m *Model) BumpTodo() {
	cursor := m.table.Cursor()
	if len(m.data.Items) <= 1 || cursor == 0 || cursor >= len(m.data.Items) {
		return
	}

	item := m.data.Items[cursor]
	m.data.Items = append(m.data.Items[:cursor], m.data.Items[cursor+1:]...)
	m.data.Items = append([]model.Todo{item}, m.data.Items...)
	m.refreshTable()
	m.table.SetCursor(0)
}

// Save persists the current state
func (m *Model) Save() error {
	return m.store.Save(m.data)
}

// IsCelebrationMilestone checks if we hit a celebration milestone
func (m *Model) IsCelebrationMilestone() bool {
	return m.data.Stats.TotalCompleted > 0 && m.data.Stats.TotalCompleted%10 == 0
}

// Run starts the TUI application
func Run(s store.Store) error {
	m, err := New(s)
	if err != nil {
		return err
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err = p.Run()
	return err
}

// ShortHelp returns keybindings to be shown in the mini help view
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Done, k.Add, k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Done, k.Add},
		{k.Drop, k.Bump, k.Help, k.Quit},
	}
}
