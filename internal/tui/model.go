package tui

import (
	"fmt"
	"os"
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

// Tab represents which tab is active
type Tab int

const (
	TabActive Tab = iota
	TabCompleted
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
	tab            Tab
	titleInput     textinput.Model
	descInput      textinput.Model
	priorityIndex  int
	inputFocus     int // 0 = title, 1 = description, 2 = priority
	celebrationMsg string
	confirmAction  string
	err            error
	cwd            string // Current working directory for context filtering
	showAllTasks   bool   // If true, show all tasks regardless of context
	filteredItems  []model.Todo
	filteredArchive []model.ArchivedTodo
}

// celebrationTickMsg is sent to end the celebration animation
type celebrationTickMsg struct{}

// NewWithContext creates a new TUI model with context awareness
func NewWithContext(s store.Store, cwd string, showAll bool) (Model, error) {
	data, err := s.Load()
	if err != nil {
		return Model{}, err
	}

	// Create table with expanded columns for fuller view
	columns := []table.Column{
		{Title: "", Width: 3},             // Status icon
		{Title: "Pri", Width: 5},          // Priority
		{Title: "Task", Width: 30},        // Task text
		{Title: "Description", Width: 20}, // Description preview
		{Title: "Context", Width: 15},     // Context/path
		{Title: "Age", Width: 10},         // Age
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(10),
		table.WithKeyMap(table.KeyMap{
			LineUp: key.NewBinding(
				key.WithKeys("up", "k"),
			),
			LineDown: key.NewBinding(
				key.WithKeys("down", "j"),
			),
			PageUp: key.NewBinding(
				key.WithKeys("pgup", "ctrl+u"),
			),
			PageDown: key.NewBinding(
				key.WithKeys("pgdown", "ctrl+d"),
			),
			HalfPageUp: key.NewBinding(
				key.WithKeys("ctrl+b"),
			),
			HalfPageDown: key.NewBinding(
				key.WithKeys("ctrl+f"),
			),
			GotoTop: key.NewBinding(
				key.WithKeys("home", "g"),
			),
			GotoBottom: key.NewBinding(
				key.WithKeys("end", "G"),
			),
		}),
	)

	// Style the table with vibrant colors
	s2 := table.DefaultStyles()
	s2.Header = s2.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(ui.BrightViolet).
		BorderBottom(true).
		Bold(true).
		Foreground(ui.NeonPurple)
	s2.Selected = s2.Selected.
		Foreground(ui.ElectricBlue).
		Background(ui.Surface0).
		Bold(true)
	s2.Cell = s2.Cell.
		Foreground(ui.Text)
	t.SetStyles(s2)

	// Create help
	h := help.New()
	h.Styles.ShortKey = lipgloss.NewStyle().Foreground(ui.ElectricBlue).Bold(true)
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
		tab:           TabActive,
		titleInput:    titleInput,
		descInput:     descInput,
		priorityIndex: 1, // Default to Medium
		cwd:           cwd,
		showAllTasks:  showAll,
	}

	m.refreshFiltered()
	m.refreshTable()
	return m, nil
}

// New creates a new TUI model (legacy, no context)
func New(s store.Store) (Model, error) {
	cwd, _ := os.Getwd()
	return NewWithContext(s, cwd, false)
}

// refreshFiltered updates the filtered items based on context
func (m *Model) refreshFiltered() {
	if m.showAllTasks || m.cwd == "" {
		m.filteredItems = m.data.Items
		m.filteredArchive = m.data.Archive
	} else {
		m.filteredItems = m.data.FilterByContext(m.cwd)
		m.filteredArchive = m.data.FilterArchiveByContext(m.cwd)
	}
}

// refreshTable updates the table rows from the data
func (m *Model) refreshTable() {
	m.refreshFiltered()

	// Calculate column widths dynamically
	taskWidth := 28
	descWidth := 18
	ctxWidth := 13
	if m.width > 0 {
		availableWidth := m.width - 35 // Reserve for status, pri, age, padding
		if availableWidth > 60 {
			taskWidth = availableWidth * 45 / 100
			descWidth = availableWidth * 30 / 100
			ctxWidth = availableWidth * 25 / 100
		}
	}

	if m.tab == TabActive {
		rows := make([]table.Row, len(m.filteredItems))
		for i, item := range m.filteredItems {
			desc := item.Description
			if desc == "" {
				desc = "-"
			}
			ctx := model.GetContextDisplay(item.Context, m.cwd)
			rows[i] = table.Row{
				ui.IconUnchecked,
				m.priorityIcon(item.Priority),
				truncateText(item.Text, taskWidth),
				truncateText(desc, descWidth),
				truncateText(ctx, ctxWidth),
				formatAge(item.Created),
			}
		}
		m.table.SetRows(rows)
	} else {
		// Completed tab - show archived items (most recent first)
		rows := make([]table.Row, len(m.filteredArchive))
		for i := range m.filteredArchive {
			// Reverse order - most recent first
			item := m.filteredArchive[len(m.filteredArchive)-1-i]
			desc := item.Description
			if desc == "" {
				desc = "-"
			}
			ctx := model.GetContextDisplay(item.Context, m.cwd)
			rows[i] = table.Row{
				ui.IconChecked,
				m.priorityIcon(item.Priority),
				truncateText(item.Text, taskWidth),
				truncateText(desc, descWidth),
				truncateText(ctx, ctxWidth),
				formatAge(item.Completed),
			}
		}
		m.table.SetRows(rows)
	}
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

// AddTodo adds a new todo item with context
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
		Context:     m.cwd, // Set context to current working directory
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
	if m.tab != TabActive {
		return false
	}

	cursor := m.table.Cursor()
	if len(m.filteredItems) == 0 || cursor >= len(m.filteredItems) {
		return false
	}

	// Get the actual item from filtered list
	item := m.filteredItems[cursor]

	// Find and remove from actual data.Items
	for i, dataItem := range m.data.Items {
		if dataItem.ID == item.ID {
			// Add to archive with all fields preserved
			archived := model.ArchivedTodo{
				ID:          item.ID,
				Text:        item.Text,
				Description: item.Description,
				Priority:    item.Priority,
				Created:     item.Created,
				Completed:   time.Now(),
				Context:     item.Context,
			}
			m.data.Archive = append(m.data.Archive, archived)

			// Remove from items
			m.data.Items = append(m.data.Items[:i], m.data.Items[i+1:]...)
			break
		}
	}

	// Update stats
	m.data.Stats.TotalCompleted++

	m.refreshTable()
	return true
}

// UncompleteTodo moves a completed task back to active
func (m *Model) UncompleteTodo() bool {
	if m.tab != TabCompleted {
		return false
	}

	cursor := m.table.Cursor()
	if len(m.filteredArchive) == 0 || cursor >= len(m.filteredArchive) {
		return false
	}

	// Get the actual item from filtered list (reversed)
	item := m.filteredArchive[len(m.filteredArchive)-1-cursor]

	// Find and remove from actual data.Archive
	for i, archiveItem := range m.data.Archive {
		if archiveItem.ID == item.ID {
			// Create active todo from archived
			todo := model.Todo{
				ID:          item.ID,
				Text:        item.Text,
				Description: item.Description,
				Priority:    item.Priority,
				Created:     item.Created,
				Position:    0,
				Context:     item.Context,
			}

			// Shift all positions down
			for j := range m.data.Items {
				m.data.Items[j].Position++
			}

			// Insert at the beginning
			m.data.Items = append([]model.Todo{todo}, m.data.Items...)

			// Remove from archive
			m.data.Archive = append(m.data.Archive[:i], m.data.Archive[i+1:]...)
			break
		}
	}

	m.refreshTable()
	return true
}

// DropTodo removes the current todo without archiving
func (m *Model) DropTodo() {
	if m.tab == TabActive {
		cursor := m.table.Cursor()
		if len(m.filteredItems) == 0 || cursor >= len(m.filteredItems) {
			return
		}

		item := m.filteredItems[cursor]
		for i, dataItem := range m.data.Items {
			if dataItem.ID == item.ID {
				m.data.Items = append(m.data.Items[:i], m.data.Items[i+1:]...)
				break
			}
		}
	} else {
		// Drop from completed (permanently delete)
		cursor := m.table.Cursor()
		if len(m.filteredArchive) == 0 || cursor >= len(m.filteredArchive) {
			return
		}

		item := m.filteredArchive[len(m.filteredArchive)-1-cursor]
		for i, archiveItem := range m.data.Archive {
			if archiveItem.ID == item.ID {
				m.data.Archive = append(m.data.Archive[:i], m.data.Archive[i+1:]...)
				break
			}
		}
	}
	m.refreshTable()
}

// BumpTodo moves the current todo to the top
func (m *Model) BumpTodo() {
	if m.tab != TabActive {
		return
	}

	cursor := m.table.Cursor()
	if len(m.filteredItems) <= 1 || cursor == 0 || cursor >= len(m.filteredItems) {
		return
	}

	item := m.filteredItems[cursor]

	// Find in actual data and move to top
	for i, dataItem := range m.data.Items {
		if dataItem.ID == item.ID {
			m.data.Items = append(m.data.Items[:i], m.data.Items[i+1:]...)
			m.data.Items = append([]model.Todo{dataItem}, m.data.Items...)
			break
		}
	}

	m.refreshTable()
	m.table.SetCursor(0)
}

// SwitchTab switches between Active and Completed tabs
func (m *Model) SwitchTab() {
	if m.tab == TabActive {
		m.tab = TabCompleted
	} else {
		m.tab = TabActive
	}
	m.refreshTable()
	m.table.SetCursor(0)
}

// ToggleShowAll toggles between showing all tasks and context-filtered tasks
func (m *Model) ToggleShowAll() {
	m.showAllTasks = !m.showAllTasks
	m.refreshTable()
}

// Save persists the current state
func (m *Model) Save() error {
	return m.store.Save(m.data)
}

// IsCelebrationMilestone checks if we hit a celebration milestone
func (m *Model) IsCelebrationMilestone() bool {
	return m.data.Stats.TotalCompleted > 0 && m.data.Stats.TotalCompleted%10 == 0
}

// GetCurrentItems returns the currently displayed items based on tab
func (m *Model) GetCurrentItems() int {
	if m.tab == TabActive {
		return len(m.filteredItems)
	}
	return len(m.filteredArchive)
}

// Run starts the TUI application (legacy)
func Run(s store.Store) error {
	cwd, _ := os.Getwd()
	return RunWithContext(s, cwd, false)
}

// RunWithContext starts the TUI application with context awareness
func RunWithContext(s store.Store, cwd string, showAll bool) error {
	m, err := NewWithContext(s, cwd, showAll)
	if err != nil {
		return err
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err = p.Run()
	return err
}

// ShortHelp returns keybindings to be shown in the mini help view
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Done, k.Add, k.Tab, k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.PageUp, k.PageDown},
		{k.GotoTop, k.GotoBottom, k.Tab, k.ToggleAll},
		{k.Done, k.Add, k.Drop, k.Bump},
		{k.Help, k.Quit},
	}
}
