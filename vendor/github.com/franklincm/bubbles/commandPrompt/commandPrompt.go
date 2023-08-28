package commandprompt

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	TextInput textinput.Model
	style     lipgloss.Style
	editing   bool
	input     string

	InputAbort  key.Binding
	InputAccept key.Binding
	InputShow   key.Binding

	Prompt string
}

type PromptInput string
type PromptEditing bool

func (m Model) PromptInput() tea.Msg {
	return PromptInput(m.input)
}

func (m Model) PromptEditing() tea.Msg {
	return PromptEditing(m.editing)
}

func New(prefix string) Model {
	ti := textinput.New()
	ti.Prompt = prefix

	return Model{
		TextInput: ti,
		style:     lipgloss.NewStyle(),

		InputAbort:  key.NewBinding(key.WithKeys("esc")),
		InputAccept: key.NewBinding(key.WithKeys("enter")),
		InputShow:   key.NewBinding(key.WithKeys(":")),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.PromptInput)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, m.InputShow) && !m.editing {
			m.editing = true
			m.TextInput.Focus()
			m.TextInput, cmd = m.TextInput.Update(msg)
			m.TextInput.SetValue("")
			cmd = m.PromptEditing
			return m, cmd

		} else if key.Matches(msg, m.InputAbort) {
			m.TextInput.Reset()
			m.TextInput.Blur()
			m.editing = false
			cmd = m.PromptEditing
			return m, cmd
		} else if key.Matches(msg, m.InputAccept) {
			m.input = m.TextInput.Value()
			m.TextInput.Reset()
			m.TextInput.Blur()
			m.editing = false
			cmds = append(cmds, m.PromptInput)
			cmds = append(cmds, m.PromptEditing)
		} else {
			m.TextInput, cmd = m.TextInput.Update(msg)
			return m, cmd
		}

	}
	m.TextInput, cmd = m.TextInput.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.editing {
		return m.style.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				m.TextInput.View(),
			),
		)
	}

	return ""
}

func (m Model) SetValue(str string) Model {
	m.TextInput.SetValue(str)
	m.TextInput.SetCursor(len(str))
	return m
}
