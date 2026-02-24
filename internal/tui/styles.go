package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	colorPrimary   = lipgloss.Color("#5faf5f")
	colorHighlight = lipgloss.Color("#ffaf5f")
	colorMuted     = lipgloss.Color("#666666")
	colorURL       = lipgloss.Color("#888888")
	colorSnippet   = lipgloss.Color("#aaaaaa")
	colorError     = lipgloss.Color("#ff5f5f")
	colorBorder    = lipgloss.Color("#3a4a3a")
	colorActiveBg  = lipgloss.Color("#0a1a0a")

	// Result block styles
	resultBlock = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder).
			Padding(0, 1).
			MarginBottom(1)

	selectedBlock = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorHighlight).
			Padding(0, 1).
			MarginBottom(1).
			Background(colorActiveBg)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorPrimary)

	selectedTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(colorHighlight)

	urlStyle = lipgloss.NewStyle().
			Foreground(colorURL).
			Italic(true)

	snippetStyle = lipgloss.NewStyle().
			Foreground(colorSnippet)

	// Status bar
	statusBar = lipgloss.NewStyle().
			Foreground(colorMuted).
			Padding(0, 1)

	// Prompt
	promptStyle = lipgloss.NewStyle().
			Foreground(colorPrimary).
			Bold(true)

	// Error
	errorStyle = lipgloss.NewStyle().
			Foreground(colorError).
			Bold(true).
			Padding(1, 2)

	// Spinner
	spinnerStyle = lipgloss.NewStyle().
			Foreground(colorPrimary)
)
