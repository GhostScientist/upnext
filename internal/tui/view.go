package tui

import (
	"fmt"
	"path/filepath"
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

	// Tabs
	sections = append(sections, m.renderTabs())
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
		if m.GetCurrentItems() == 0 {
			sections = append(sections, m.renderEmptyState())
		} else {
			sections = append(sections, m.renderTable())
			// Show task details panel for selected task
			details := m.renderTaskDetails()
			if details != "" {
				sections = append(sections, "")
				sections = append(sections, details)
			}
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
	title := ui.TitleStyle.Render("⚡ upnext")
	subtitle := ui.SubtitleStyle.Render(" - what's next?")
	content := title + subtitle

	return ui.HeaderStyle.Width(m.width - 4).Render(content)
}

func (m Model) renderTabs() string {
	// Tab labels
	activeLabel := fmt.Sprintf(" Active (%d) ", len(m.filteredItems))
	completedLabel := fmt.Sprintf(" Completed (%d) ", len(m.filteredArchive))

	var activeTab, completedTab string
	if m.tab == TabActive {
		activeTab = ui.TabActiveStyle.Render(activeLabel)
		completedTab = ui.TabInactiveStyle.Render(completedLabel)
	} else {
		activeTab = ui.TabInactiveStyle.Render(activeLabel)
		completedTab = ui.TabActiveStyle.Render(completedLabel)
	}

	tabs := lipgloss.JoinHorizontal(lipgloss.Bottom, activeTab, " ", completedTab)

	// Context indicator
	var contextInfo string
	if m.showAllTasks {
		contextInfo = ui.ContextStyle.Render("  " + ui.IconGlobal + " showing all tasks")
	} else if m.cwd != "" {
		shortPath := filepath.Base(m.cwd)
		contextInfo = ui.ContextStyle.Render("  " + ui.IconFolder + " " + shortPath)
	}

	return tabs + contextInfo
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
	var message string
	if m.tab == TabActive {
		message = "Nothing to do! Press 'a' to add a task."
	} else {
		message = "No completed tasks yet. Complete some tasks to see them here!"
	}

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		ui.DimStyle.Render(stars),
		"",
		ui.EmptyStyle.Render(message),
	)

	// Center in available space
	return lipgloss.Place(
		m.width,
		m.height-14,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

func (m Model) renderInputForm() string {
	var b strings.Builder

	// Form title
	title := ui.DialogTitleStyle.Render("✨ Add New Task")
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
	// Left side: item count based on current tab
	var itemCount string
	if m.tab == TabActive {
		count := len(m.filteredItems)
		if count == 1 {
			itemCount = "1 active task"
		} else {
			itemCount = fmt.Sprintf("%d active tasks", count)
		}
	} else {
		count := len(m.filteredArchive)
		if count == 1 {
			itemCount = "1 completed task"
		} else {
			itemCount = fmt.Sprintf("%d completed tasks", count)
		}
	}

	// Right side: total completed
	completed := fmt.Sprintf("🏆 %d total completed", m.data.Stats.TotalCompleted)

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
		{"↑/k, ↓/j", "Navigate tasks"},
		{"pgup/^u, pgdn/^d", "Page up/down"},
		{"g/G", "Go to top/bottom"},
		{"1/2", "Switch to Active/Completed tab"},
		{"enter/d", "Complete task (Active) / View (Completed)"},
		{"u", "Uncomplete task (Completed tab)"},
		{"a", "Add new task"},
		{"x", "Drop (delete) task"},
		{"b", "Bump task to top"},
		{"A", "Toggle show all tasks"},
		{"?", "Toggle help"},
		{"q/esc", "Quit"},
	}

	var lines []string
	lines = append(lines, ui.DialogTitleStyle.Render("⌨️  Keyboard Shortcuts"))
	lines = append(lines, "")

	for _, item := range helpItems {
		key := ui.HelpKeyStyle.Render(fmt.Sprintf("%18s", item.key))
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

// renderTaskDetails shows expanded details for the selected task
func (m Model) renderTaskDetails() string {
	if m.tab == TabActive {
		return m.renderActiveTaskDetails()
	}
	return m.renderCompletedTaskDetails()
}

func (m Model) renderActiveTaskDetails() string {
	if len(m.filteredItems) == 0 {
		return ""
	}

	cursor := m.table.Cursor()
	if cursor >= len(m.filteredItems) {
		return ""
	}

	item := m.filteredItems[cursor]
	var lines []string

	// Title with task number indicator
	taskNum := fmt.Sprintf("Task %d of %d", cursor+1, len(m.filteredItems))
	headerLine := ui.DimStyle.Render(taskNum)
	lines = append(lines, headerLine)
	lines = append(lines, "")

	// Full task text
	lines = append(lines, ui.TitleStyle.Render(item.Text))

	// Description (if available)
	if item.Description != "" {
		lines = append(lines, "")
		descLabel := ui.LabelStyle.Render("Description:")
		lines = append(lines, descLabel)
		descText := ui.SubtitleStyle.Render(item.Description)
		lines = append(lines, descText)
	}

	lines = append(lines, "")

	// Priority and age on same line
	priStyle := ui.PriorityLowStyle
	priIcon := ui.IconLow
	switch item.Priority {
	case model.PriorityHigh:
		priStyle = ui.PriorityHighStyle
		priIcon = ui.IconHigh
	case model.PriorityMedium:
		priStyle = ui.PriorityMediumStyle
		priIcon = ui.IconMedium
	}
	priText := priStyle.Render(priIcon + " " + item.Priority.String() + " priority")
	ageText := ui.DimStyle.Render("Created: " + formatAge(item.Created))
	infoLine := priText + "  " + ui.DimStyle.Render("•") + "  " + ageText

	// Context info
	ctx := model.GetContextDisplay(item.Context, m.cwd)
	if ctx != "" {
		infoLine += "  " + ui.DimStyle.Render("•") + "  " + ui.ContextStyle.Render(ctx)
	}

	lines = append(lines, infoLine)

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)

	// Create a styled panel
	panelStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(ui.BrightViolet).
		Padding(0, 1).
		Width(m.width - 6)

	return panelStyle.Render(content)
}

func (m Model) renderCompletedTaskDetails() string {
	if len(m.filteredArchive) == 0 {
		return ""
	}

	cursor := m.table.Cursor()
	if cursor >= len(m.filteredArchive) {
		return ""
	}

	// Get the item (reversed order)
	item := m.filteredArchive[len(m.filteredArchive)-1-cursor]
	var lines []string

	// Title with task number indicator
	taskNum := fmt.Sprintf("Completed task %d of %d", cursor+1, len(m.filteredArchive))
	headerLine := ui.DimStyle.Render(taskNum)
	lines = append(lines, headerLine)
	lines = append(lines, "")

	// Full task text with checkmark
	lines = append(lines, ui.CheckmarkStyle.Render("✓ ")+ui.TitleStyle.Render(item.Text))

	// Description (if available)
	if item.Description != "" {
		lines = append(lines, "")
		descLabel := ui.LabelStyle.Render("Description:")
		lines = append(lines, descLabel)
		descText := ui.SubtitleStyle.Render(item.Description)
		lines = append(lines, descText)
	}

	lines = append(lines, "")

	// Completion info
	completedText := ui.CheckmarkStyle.Render("Completed: " + formatAge(item.Completed))
	createdText := ui.DimStyle.Render("Created: " + formatAge(item.Created))
	infoLine := completedText + "  " + ui.DimStyle.Render("•") + "  " + createdText

	// Context info
	ctx := model.GetContextDisplay(item.Context, m.cwd)
	if ctx != "" {
		infoLine += "  " + ui.DimStyle.Render("•") + "  " + ui.ContextStyle.Render(ctx)
	}

	lines = append(lines, infoLine)

	// Hint about uncomplete
	lines = append(lines, "")
	lines = append(lines, ui.DimStyle.Render("Press 'u' to move back to active tasks"))

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)

	// Create a styled panel
	panelStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(ui.BrightViolet).
		Padding(0, 1).
		Width(m.width - 6)

	return panelStyle.Render(content)
}
