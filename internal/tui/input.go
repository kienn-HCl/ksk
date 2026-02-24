package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type inputModel struct {
	textInput textinput.Model
}

func newInputModel(placeholder string) inputModel {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Prompt = "ksk> "
	ti.PromptStyle = promptStyle
	ti.Focus()
	ti.CharLimit = 256
	return inputModel{textInput: ti}
}

func (m inputModel) Update(msg tea.Msg) (inputModel, tea.Cmd) {
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m inputModel) View() string {
	return m.textInput.View()
}

func (m inputModel) Value() string {
	return m.textInput.Value()
}

func (m *inputModel) SetValue(s string) {
	m.textInput.SetValue(s)
}

func (m *inputModel) Focus() tea.Cmd {
	return m.textInput.Focus()
}

func (m *inputModel) Blur() {
	m.textInput.Blur()
}
