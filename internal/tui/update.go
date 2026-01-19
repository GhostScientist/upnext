package tui

import (
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"

	"upnext/internal/model"
	"upnext/internal/ui"
)

// Update implements tea.Model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Update table dimensions
		tableHeight := m.height - 14 // Reserve space for header, tabs, status bar, help, details
		if tableHeight < 5 {
			tableHeight = 5
		}
		m.table.SetHeight(tableHeight)

		// Update column widths based on available space
		// Reserve space: status(3) + pri(5) + age(10) + padding(~17) = 35
		availableWidth := m.width - 35
		if availableWidth < 60 {
			availableWidth = 60
		}
		// Split available space between task, description, and context
		taskWidth := availableWidth * 45 / 100
		descWidth := availableWidth * 30 / 100
		ctxWidth := availableWidth * 25 / 100
		if taskWidth < 15 {
			taskWidth = 15
		}
		if descWidth < 10 {
			descWidth = 10
		}
		if ctxWidth < 10 {
			ctxWidth = 10
		}
		m.table.SetColumns([]table.Column{
			{Title: "", Width: 3},
			{Title: "Pri", Width: 5},
			{Title: "Task", Width: taskWidth},
			{Title: "Description", Width: descWidth},
			{Title: "Context", Width: ctxWidth},
			{Title: "Age", Width: 10},
		})
		m.refreshTable()

		// Update help width
		m.help.Width = m.width

		// Update input widths
		inputWidth := m.width - 20
		if inputWidth < 30 {
			inputWidth = 30
		}
		m.titleInput.Width = inputWidth
		m.descInput.Width = inputWidth

		return m, nil

	case celebrationTickMsg:
		m.mode = ModeNormal
		m.celebrationMsg = ""
		return m, nil

	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	}

	// Handle component updates based on mode
	var cmd tea.Cmd
	switch m.mode {
	case ModeInput:
		return m.updateInputMode(msg)
	case ModeNormal:
		m.table, cmd = m.table.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Handle celebration mode - any key dismisses
	if m.mode == ModeCelebration {
		m.mode = ModeNormal
		m.celebrationMsg = ""
		return m, nil
	}

	// Handle help mode
	if m.mode == ModeHelp {
		if key.Matches(msg, m.keys.Help) || key.Matches(msg, m.keys.Quit) || msg.String() == "enter" {
			m.mode = ModeNormal
			return m, nil
		}
		return m, nil
	}

	// Handle input mode
	if m.mode == ModeInput {
		return m.handleInputKeyPress(msg)
	}

	// Normal mode key handling
	switch {
	case key.Matches(msg, m.keys.Quit):
		return m, tea.Quit

	case key.Matches(msg, m.keys.Tab):
		// Switch between Active and Completed tabs
		m.SwitchTab()
		return m, nil

	case msg.String() == "1":
		// Go to Active tab
		if m.tab != TabActive {
			m.tab = TabActive
			m.refreshTable()
			m.table.SetCursor(0)
		}
		return m, nil

	case msg.String() == "2":
		// Go to Completed tab
		if m.tab != TabCompleted {
			m.tab = TabCompleted
			m.refreshTable()
			m.table.SetCursor(0)
		}
		return m, nil

	case key.Matches(msg, m.keys.ToggleAll):
		m.ToggleShowAll()
		return m, nil

	case key.Matches(msg, m.keys.Uncomplete):
		if m.UncompleteTodo() {
			if err := m.Save(); err != nil {
				m.err = err
			}
		}
		return m, nil

	case key.Matches(msg, m.keys.Done):
		if m.tab == TabActive {
			if m.CompleteTodo() {
				if err := m.Save(); err != nil {
					m.err = err
				}
				// Check for celebration milestone
				if m.IsCelebrationMilestone() {
					m.mode = ModeCelebration
					m.celebrationMsg = getCelebrationMessage(m.data.Stats.TotalCompleted)
					return m, tea.Tick(time.Second*3, func(time.Time) tea.Msg {
						return celebrationTickMsg{}
					})
				}
			}
		}
		return m, nil

	case key.Matches(msg, m.keys.Add):
		// Only allow adding tasks in Active tab
		if m.tab == TabActive {
			m.mode = ModeInput
			m.inputFocus = 0
			m.priorityIndex = 1 // Default to Medium
			m.titleInput.SetValue("")
			m.descInput.SetValue("")
			m.titleInput.Focus()
			m.titleInput.PromptStyle = ui.FocusedStyle
			m.descInput.PromptStyle = ui.BlurredStyle
			return m, m.titleInput.Focus()
		}
		return m, nil

	case key.Matches(msg, m.keys.Drop):
		m.DropTodo()
		if err := m.Save(); err != nil {
			m.err = err
		}
		return m, nil

	case key.Matches(msg, m.keys.Bump):
		m.BumpTodo()
		if err := m.Save(); err != nil {
			m.err = err
		}
		return m, nil

	case key.Matches(msg, m.keys.Help):
		m.mode = ModeHelp
		return m, nil

	default:
		// Pass through to table for navigation
		var cmd tea.Cmd
		m.table, cmd = m.table.Update(msg)
		return m, cmd
	}
}

func (m Model) handleInputKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Cancel):
		m.mode = ModeNormal
		m.titleInput.Blur()
		m.descInput.Blur()
		return m, nil

	case msg.String() == "enter":
		if m.inputFocus < 2 {
			// Move to next field
			m.inputFocus++
			m.updateInputFocus()
			return m, nil
		}
		// Submit the form
		text := m.titleInput.Value()
		if text != "" {
			priority := model.Priority(m.priorityIndex)
			m.AddTodo(text, m.descInput.Value(), priority)
			if err := m.Save(); err != nil {
				m.err = err
			}
		}
		m.mode = ModeNormal
		m.titleInput.Blur()
		m.descInput.Blur()
		return m, nil

	case msg.String() == "tab" || msg.String() == "shift+tab":
		if msg.String() == "shift+tab" {
			m.inputFocus--
			if m.inputFocus < 0 {
				m.inputFocus = 2
			}
		} else {
			m.inputFocus++
			if m.inputFocus > 2 {
				m.inputFocus = 0
			}
		}
		m.updateInputFocus()
		return m, nil

	case m.inputFocus == 2: // Priority selection
		if key.Matches(msg, m.keys.Left) || msg.String() == "h" {
			m.priorityIndex--
			if m.priorityIndex < 0 {
				m.priorityIndex = 2
			}
		} else if key.Matches(msg, m.keys.Right) || msg.String() == "l" {
			m.priorityIndex++
			if m.priorityIndex > 2 {
				m.priorityIndex = 0
			}
		}
		return m, nil
	}

	return m.updateInputMode(msg)
}

func (m Model) updateInputMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.inputFocus {
	case 0:
		m.titleInput, cmd = m.titleInput.Update(msg)
	case 1:
		m.descInput, cmd = m.descInput.Update(msg)
	}

	return m, cmd
}

func (m *Model) updateInputFocus() {
	switch m.inputFocus {
	case 0:
		m.titleInput.Focus()
		m.titleInput.PromptStyle = ui.FocusedStyle
		m.descInput.Blur()
		m.descInput.PromptStyle = ui.BlurredStyle
	case 1:
		m.titleInput.Blur()
		m.titleInput.PromptStyle = ui.BlurredStyle
		m.descInput.Focus()
		m.descInput.PromptStyle = ui.FocusedStyle
	case 2:
		m.titleInput.Blur()
		m.titleInput.PromptStyle = ui.BlurredStyle
		m.descInput.Blur()
		m.descInput.PromptStyle = ui.BlurredStyle
	}
}

func getCelebrationMessage(total int) string {
	messages := []string{
		"Amazing! You're on fire!",
		"Incredible progress!",
		"You're crushing it!",
		"Productivity champion!",
		"Keep up the great work!",
	}
	return messages[(total/10-1)%len(messages)]
}
