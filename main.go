package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	commandprompt "github.com/franklincm/bubbles/commandPrompt"
)

var borderKeys = key.NewBinding(
	key.WithKeys("b", "ctrl+b"),
	key.WithHelp("", "toggle border"),
)

var headerStyle = lipgloss.NewStyle().
	Margin(0).
	Padding(0).
	BorderBottom(true).
	BorderTop(false).
	BorderLeft(false).
	BorderRight(false).
	BorderStyle(lipgloss.NormalBorder())

var testStyle = lipgloss.NewStyle()
var fullscreenStyle = lipgloss.NewStyle().
	Padding(0).
	Margin(0).
	Align(lipgloss.Center)

var quitKeys = key.NewBinding(
	key.WithKeys("q", "esc", "ctrl+c"),
	key.WithHelp("", "press q to quit"),
)

type errMsg error

type Model struct {
	border     bool
	err        error
	height     int
	input      []string
	prompt     tea.Model
	quitting   bool
	ready      bool
	showprompt bool
	textStyle  lipgloss.Style
	width      int
}

func NewModel() Model {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("black")).
		Background(lipgloss.Color("240")).
		Align(lipgloss.Center)

	prompt := commandprompt.New()

	return Model{
		textStyle: style,
		border:    false,
		prompt:    prompt,
	}
}

func (m Model) Init() tea.Cmd {
	m.input = make([]string, 3)
	return tea.Batch(
		tea.EnterAltScreen,
		m.prompt.Init(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		fullscreenStyle.Width(msg.Width)
		fullscreenStyle.Height(msg.Height)
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case commandprompt.PromptInput:
		if msg == "quit" || msg == "q" {
			m.quitting = true
			return m, tea.Quit
		}

		m.input = append(m.input, string(msg))
		return m, nil

	case commandprompt.PromptEditing:
		m.showprompt = bool(msg)
		return m, nil

	case tea.KeyMsg:
		if key.Matches(msg, quitKeys) && !m.showprompt {
			m.quitting = true
			return m, tea.Quit

		} else if key.Matches(msg, borderKeys) && !m.showprompt {
			m.border = !m.border
		}

		// return m, nil

	case errMsg:
		m.err = msg
		return m, nil
	}

	m.prompt, cmd = m.prompt.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.err != nil {
		return m.err.Error()
	}

	str := fmt.Sprintf(
		"%-*s %*s %*s",
		20,
		"test",
		(m.width-20-12)/2,
		fmt.Sprintf("width:%d", (m.width-20-15)/2),
		(m.width-20-12)/2,
		quitKeys.Help().Desc,
	)

	if m.quitting {
		return str + "\n"
	}

	testStr := `
some
text
here
i
guess
	`
	return fullscreenStyle.Height(m.height).Render(
		lipgloss.JoinVertical(
			lipgloss.Top,
			headerStyle.
				Width(m.width).
				Render(str),

			testStyle.
				AlignHorizontal(lipgloss.Center).
				Height(m.height-6).
				Width(m.width).
				Margin(0).
				Padding(0).
				BorderStyle(lipgloss.NormalBorder()).
				BorderTop(false).
				BorderBottom(true).
				Render(testStr),

			testStyle.
				Width(m.width).
				Margin(0).
				Padding(0).
				BorderStyle(lipgloss.NormalBorder()).
				BorderTop(false).
				BorderBottom(false).
				Render(fmt.Sprintf("%s:%d", "footer, height", m.height)),

			testStyle.Render(m.prompt.View()+"\n"),
		),
	)
}

func main() {
	p := tea.NewProgram(NewModel(), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
