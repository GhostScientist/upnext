package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"upnext/internal/model"
	"upnext/internal/ui"
)

// View implements tea.Model
func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	var sections []string

	// Header
	sections = append(sections, m.renderHeader())
	sections = append(sections, "")

	// Main content area
	switch m.mode {
	case ModeCelebration:
		sections = append(sections, m.renderCelebration())
	case ModeHelp:
		sections = append(sections, m.renderFullHelp())
	case ModeInput:
		sections = append(sections, m.renderTable())
		sections = append(sections, "")
		sections = append(sections, m.renderInputForm())
	default:
		if len(m.data.Items) == 0 {
			sections = append(sections, m.renderEmptyState())
		} else {
			sections = append(sections, m.renderTable())
		}
	}

	// Spacer to push status bar to bottom
	contentHeight := lipgloss.Height(strings.Join(sections, "\n"))
	spacerHeight := m.height - contentHeight - 3 // 3 for status bar + help
	if spacerHeight > 0 {
		sections = append(sections, strings.Repeat("\n", spacerHeight))
	}

	// Help bar (short version)
	if m.mode != ModeHelp && m.mode != ModeInput {
		sections = append(sections, m.renderHelpBar())
	}

	// Status bar
	sections = append(sections, m.renderStatusBar())

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m Model) renderHeader() string {
	title := ui.TitleStyle.Render("upnext")
	subtitle := ui.SubtitleStyle.Render(" - what's next?")
	content := title + subtitle

	return ui.HeaderStyle.Width(m.width - 4).Render(content)
}

func (m Model) renderTable() string {
	return m.table.View()
}

func (m Model) renderEmptyState() string {
	stars := `
      ✦  ·  ✦     ·    ✦
    ·    ✦    ·  ✦   ·
      ·     ✦  ·    ✦  ·
    ✦   ·  ✦    ·  ✦    ·
`
	message := "Nothing to do! Press 'a' to add a task."

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		ui.DimStyle.Render(stars),
		"",
		ui.EmptyStyle.Render(message),
	)

	// Center in available space
	return lipgloss.Place(
		m.width,
		m.height-10,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

func (m Model) renderInputForm() string {
	var b strings.Builder

	// Form title
	title := ui.DialogTitleStyle.Render("Add New Task")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Task title field
	titleLabel := ui.LabelStyle.Render("Task:")
	if m.inputFocus == 0 {
		titleLabel = ui.FocusedStyle.Render("Task:      ")
	}
	b.WriteString(titleLabel)
	b.WriteString(" ")
	b.WriteString(m.titleInput.View())
	b.WriteString("\n\n")

	// Description field
	descLabel := ui.LabelStyle.Render("Description:")
	if m.inputFocus == 1 {
		descLabel = ui.FocusedStyle.Render("Description:")
	}
	b.WriteString(descLabel)
	b.WriteString(" ")
	b.WriteString(m.descInput.View())
	b.WriteString("\n\n")

	// Priority selector
	priLabel := ui.LabelStyle.Render("Priority:")
	if m.inputFocus == 2 {
		priLabel = ui.FocusedStyle.Render("Priority:  ")
	}
	b.WriteString(priLabel)
	b.WriteString(" ")
	b.WriteString(m.renderPrioritySelector())
	b.WriteString("\n\n")

	// Help text
	helpText := ui.DimStyle.Render("tab: next field • enter: submit • esc: cancel")
	b.WriteString(helpText)

	// Wrap in dialog box
	return ui.DialogStyle.Width(m.width - 10).Render(b.String())
}

func (m Model) renderPrioritySelector() string {
	priorities := []struct {
		label string
		style lipgloss.Style
	}{
		{"Low", ui.PriorityLowStyle},
		{"Medium", ui.PriorityMediumStyle},
		{"High", ui.PriorityHighStyle},
	}

	var items []string
	for i, p := range priorities {
		label := p.label
		if i == m.priorityIndex {
			// Selected priority
			if m.inputFocus == 2 {
				// Focused on priority selector
				label = ui.ButtonActiveStyle.Render(label)
			} else {
				label = ui.ButtonStyle.Render(label)
			}
		} else {
			label = p.style.Render(label)
		}
		items = append(items, label)
	}

	return lipgloss.JoinHorizontal(lipgloss.Center, items...)
}

func (m Model) renderHelpBar() string {
	return m.help.View(m.keys)
}

func (m Model) renderStatusBar() string {
	// Left side: item count
	itemCount := fmt.Sprintf("%d items", len(m.data.Items))
	if len(m.data.Items) == 1 {
		itemCount = "1 item"
	}

	// Right side: total completed
	completed := fmt.Sprintf("%d completed", m.data.Stats.TotalCompleted)

	// Calculate padding
	leftContent := itemCount
	rightContent := completed
	padding := m.width - lipgloss.Width(leftContent) - lipgloss.Width(rightContent) - 4

	if padding < 1 {
		padding = 1
	}

	statusContent := leftContent + strings.Repeat(" ", padding) + rightContent

	return ui.StatusBarStyle.Width(m.width).Render(statusContent)
}

func (m Model) renderFullHelp() string {
	helpItems := []struct {
		key  string
		desc string
	}{
		{"↑/k", "Move up"},
		{"↓/j", "Move down"},
		{"enter/d", "Complete task"},
		{"a", "Add new task"},
		{"x", "Drop (delete) task"},
		{"b", "Bump task to top"},
		{"?", "Toggle help"},
		{"q/esc", "Quit"},
	}

	var lines []string
	lines = append(lines, ui.DialogTitleStyle.Render("Keyboard Shortcuts"))
	lines = append(lines, "")

	for _, item := range helpItems {
		key := ui.HelpKeyStyle.Render(fmt.Sprintf("%12s", item.key))
		desc := ui.HelpDescStyle.Render("  " + item.desc)
		lines = append(lines, key+desc)
	}

	lines = append(lines, "")
	lines = append(lines, ui.DimStyle.Render("Press any key to close"))

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)

	// Center the help dialog
	dialog := ui.HelpOverlayStyle.Render(content)
	return lipgloss.Place(
		m.width,
		m.height-6,
		lipgloss.Center,
		lipgloss.Center,
		dialog,
	)
}

func (m Model) renderCelebration() string {
	celebration := `
    ✨ ⭐ ✨ ⭐ ✨ ⭐ ✨

       🎉 MILESTONE! 🎉

       ` + m.celebrationMsg + `

    ✨ ⭐ ✨ ⭐ ✨ ⭐ ✨
`
	content := ui.CelebrationStyle.Render(celebration)

	return lipgloss.Place(
		m.width,
		m.height-6,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

// Unused but keeping for reference - shows selected task details
func (m Model) renderTaskDetails() string {
	if len(m.data.Items) == 0 {
		return ""
	}

	cursor := m.table.Cursor()
	if cursor >= len(m.data.Items) {
		return ""
	}

	item := m.data.Items[cursor]
	var lines []string

	lines = append(lines, ui.TitleStyle.Render(item.Text))
	if item.Description != "" {
		lines = append(lines, ui.DimStyle.Render(item.Description))
	}

	priStyle := ui.PriorityLowStyle
	switch item.Priority {
	case model.PriorityHigh:
		priStyle = ui.PriorityHighStyle
	case model.PriorityMedium:
		priStyle = ui.PriorityMediumStyle
	}
	lines = append(lines, priStyle.Render("Priority: "+item.Priority.String()))

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)
	return ui.DialogStyle.Width(m.width - 10).Render(content)
}
