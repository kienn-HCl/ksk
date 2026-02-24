package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/frort/ksk/internal/browser"
	"github.com/frort/ksk/internal/search"
)

type state int

const (
	stateInput state = iota
	stateLoading
	stateResults
)

// Messages
type searchResultMsg struct {
	page *search.Page
	err  error
}

type Model struct {
	state   state
	input   inputModel
	results resultsModel
	spinner spinner.Model
	query   string
	page    *search.Page
	errMsg  string
	width   int
	height  int
	backend search.Backend
}

func NewModel(initialQuery string, backend search.Backend) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle

	m := Model{
		state:   stateInput,
		input:   newInputModel("search the web..."),
		results: newResultsModel(),
		spinner: s,
		backend: backend,
	}

	if initialQuery != "" {
		m.query = initialQuery
		m.input.SetValue(initialQuery)
		m.state = stateLoading
	}

	return m
}

func (m Model) Init() tea.Cmd {
	cmds := []tea.Cmd{m.input.textInput.Cursor.BlinkCmd()}
	if m.state == stateLoading {
		cmds = append(cmds, m.spinner.Tick, m.doSearch(m.query))
	}
	return tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Reserve space for input(1) + status(1) + padding(2)
		m.results.SetSize(msg.Width, msg.Height-4)
		return m, nil

	case searchResultMsg:
		if msg.err != nil {
			m.errMsg = msg.err.Error()
			m.state = stateInput
			return m, m.input.Focus()
		}
		m.page = msg.page
		m.errMsg = ""
		m.results.SetResults(msg.page.Results, msg.page.PageNum, msg.page.HasMore)
		m.state = stateResults
		m.input.Blur()
		return m, nil

	case spinner.TickMsg:
		if m.state == stateLoading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
		return m, nil
	}

	switch m.state {
	case stateInput:
		return m.updateInput(msg)
	case stateResults:
		return m.updateResults(msg)
	}

	return m, nil
}

func (m Model) updateInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEnter:
			q := m.input.Value()
			if q == "" {
				return m, nil
			}
			m.query = q
			m.state = stateLoading
			m.errMsg = ""
			m.input.Blur()
			return m, tea.Batch(m.spinner.Tick, m.doSearch(q))
		case tea.KeyEsc:
			if len(m.results.results) > 0 {
				m.state = stateResults
				m.input.Blur()
				return m, nil
			}
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m Model) updateResults(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "j", "down":
			m.results.CursorDown()
		case "k", "up":
			m.results.CursorUp()
		case "g":
			m.results.CursorTop()
		case "G":
			m.results.CursorBottom()
		case "enter", "o":
			if r := m.results.SelectedResult(); r != nil {
				_ = browser.Open(r.URL)
			}
		case "/":
			m.state = stateInput
			return m, m.input.Focus()
		case "l", "right":
			if m.page != nil && m.page.HasMore {
				m.state = stateLoading
				return m, tea.Batch(m.spinner.Tick, m.doNextPage())
			}
		case "h", "left":
			if m.page != nil && m.page.PageNum > 1 {
				m.state = stateLoading
				return m, tea.Batch(m.spinner.Tick, m.doPrevPage())
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	var sections []string

	// Input bar
	sections = append(sections, m.input.View())

	// Error message
	if m.errMsg != "" {
		sections = append(sections, errorStyle.Render("Error: "+m.errMsg))
	}

	switch m.state {
	case stateLoading:
		sections = append(sections, fmt.Sprintf("\n  %s Searching...\n", m.spinner.View()))

	case stateResults:
		sections = append(sections, m.results.View())

	case stateInput:
		if len(m.results.results) > 0 {
			sections = append(sections, m.results.View())
		}
	}

	// Status bar
	if m.state == stateResults || (m.state == stateInput && len(m.results.results) > 0) {
		sections = append(sections, statusBar.Render(m.results.StatusView(m.backend.Name())))
	}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m Model) doSearch(query string) tea.Cmd {
	backend := m.backend
	return func() tea.Msg {
		page, err := backend.Search(query)
		return searchResultMsg{page: page, err: err}
	}
}

func (m Model) doNextPage() tea.Cmd {
	page := m.page
	query := m.query
	backend := m.backend
	return func() tea.Msg {
		next, err := backend.NextPage(page, query)
		return searchResultMsg{page: next, err: err}
	}
}

func (m Model) doPrevPage() tea.Cmd {
	query := m.query
	pageNum := m.page.PageNum - 1
	backend := m.backend
	return func() tea.Msg {
		prev, err := backend.PrevPage(query, pageNum)
		return searchResultMsg{page: prev, err: err}
	}
}
