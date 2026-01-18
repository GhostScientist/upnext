package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// Catppuccin Mocha color palette
var (
	Rosewater = lipgloss.Color("#f5e0dc")
	Flamingo  = lipgloss.Color("#f2cdcd")
	Pink      = lipgloss.Color("#f5c2e7")
	Mauve     = lipgloss.Color("#cba6f7")
	Red       = lipgloss.Color("#f38ba8")
	Maroon    = lipgloss.Color("#eba0ac")
	Peach     = lipgloss.Color("#fab387")
	Yellow    = lipgloss.Color("#f9e2af")
	Green     = lipgloss.Color("#a6e3a1")
	Teal      = lipgloss.Color("#94e2d5")
	Sky       = lipgloss.Color("#89dceb")
	Sapphire  = lipgloss.Color("#74c7ec")
	Blue      = lipgloss.Color("#89b4fa")
	Lavender  = lipgloss.Color("#b4befe")
	Text      = lipgloss.Color("#cdd6f4")
	Subtext1  = lipgloss.Color("#bac2de")
	Subtext0  = lipgloss.Color("#a6adc8")
	Overlay2  = lipgloss.Color("#9399b2")
	Overlay1  = lipgloss.Color("#7f849c")
	Overlay0  = lipgloss.Color("#6c7086")
	Surface2  = lipgloss.Color("#585b70")
	Surface1  = lipgloss.Color("#45475a")
	Surface0  = lipgloss.Color("#313244")
	Base      = lipgloss.Color("#1e1e2e")
	Mantle    = lipgloss.Color("#181825")
	Crust     = lipgloss.Color("#11111b")
)

// Styles for the TUI
var (
	// Header box style with rounded border
	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Mauve).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(Mauve).
			Padding(0, 1)

	// Title text inside header
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Text)

	// Subtitle/tagline
	SubtitleStyle = lipgloss.NewStyle().
			Foreground(Subtext0).
			Italic(true)

	// Selected item style
	SelectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Mauve).
			Background(Surface0).
			Padding(0, 1)

	// Normal item style
	ItemStyle = lipgloss.NewStyle().
			Foreground(Text).
			Padding(0, 1)

	// Cursor style
	CursorStyle = lipgloss.NewStyle().
			Foreground(Mauve).
			Bold(true)

	// Progress bar filled portion
	ProgressFilledStyle = lipgloss.NewStyle().
				Foreground(Green)

	// Progress bar empty portion
	ProgressEmptyStyle = lipgloss.NewStyle().
				Foreground(Surface1)

	// Status bar style
	StatusBarStyle = lipgloss.NewStyle().
			Foreground(Subtext0).
			Background(Surface0).
			Padding(0, 1)

	// Help key style
	HelpKeyStyle = lipgloss.NewStyle().
			Foreground(Mauve).
			Bold(true)

	// Help description style
	HelpDescStyle = lipgloss.NewStyle().
			Foreground(Subtext0)

	// Help overlay style
	HelpOverlayStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(Mauve).
				Padding(1, 2).
				Background(Base)

	// Empty state style
	EmptyStyle = lipgloss.NewStyle().
			Foreground(Overlay1).
			Italic(true)

	// Input prompt style
	InputPromptStyle = lipgloss.NewStyle().
				Foreground(Mauve).
				Bold(true)

	// Input text style
	InputTextStyle = lipgloss.NewStyle().
			Foreground(Text)

	// Celebration style
	CelebrationStyle = lipgloss.NewStyle().
				Foreground(Yellow).
				Bold(true)

	// Checkmark for completed items
	CheckmarkStyle = lipgloss.NewStyle().
			Foreground(Green).
			Bold(true)

	// Dimmed text for timestamps etc
	DimStyle = lipgloss.NewStyle().
			Foreground(Overlay0)

	// Table styles
	TableHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(Mauve).
				BorderStyle(lipgloss.NormalBorder()).
				BorderBottom(true).
				BorderForeground(Surface1)

	TableSelectedStyle = lipgloss.NewStyle().
				Background(Surface0).
				Foreground(Text).
				Bold(true)

	TableCellStyle = lipgloss.NewStyle().
			Foreground(Text).
			Padding(0, 1)

	// Priority styles
	PriorityHighStyle = lipgloss.NewStyle().
				Foreground(Red).
				Bold(true)

	PriorityMediumStyle = lipgloss.NewStyle().
				Foreground(Yellow)

	PriorityLowStyle = lipgloss.NewStyle().
				Foreground(Green)

	// Dialog/modal styles
	DialogStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(Mauve).
			Padding(1, 2).
			Background(Base)

	DialogTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(Text).
				MarginBottom(1)

	ButtonStyle = lipgloss.NewStyle().
			Foreground(Text).
			Background(Surface1).
			Padding(0, 2).
			MarginRight(1)

	ButtonActiveStyle = lipgloss.NewStyle().
				Foreground(Base).
				Background(Mauve).
				Padding(0, 2).
				MarginRight(1).
				Bold(true)

	// App container style
	AppStyle = lipgloss.NewStyle().
			Padding(1, 2)

	// Focus styles for form inputs
	FocusedStyle = lipgloss.NewStyle().
			Foreground(Mauve)

	BlurredStyle = lipgloss.NewStyle().
			Foreground(Subtext0)

	// Label style for form fields
	LabelStyle = lipgloss.NewStyle().
			Foreground(Subtext1).
			Width(12)
)

// Icons
const (
	IconCursor    = "❯"
	IconUnchecked = "○"
	IconChecked   = "●"
	IconStar      = "✦"
	IconCheckmark = "✓"
	IconHigh      = "▲"
	IconMedium    = "◆"
	IconLow       = "▽"
)

// RenderProgressBar creates a gradient progress bar
func RenderProgressBar(percent float64, width int) string {
	filled := int(float64(width) * percent)
	empty := width - filled

	// Gradient colors from green to yellow based on "age" or urgency
	var bar string
	for i := 0; i < filled; i++ {
		bar += ProgressFilledStyle.Render("█")
	}
	for i := 0; i < empty; i++ {
		bar += ProgressEmptyStyle.Render("░")
	}
	return bar
}
