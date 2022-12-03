package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	quitting  bool
	err       error
	textStyle lipgloss.Style
	width     int
	height    int
	ready     bool
	stopwatch stopwatch.Model
	border    bool
}

func NewModel() Model {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("black")).
		Background(lipgloss.Color("240")).
		Align(lipgloss.Center)

	stopwatch := stopwatch.NewWithInterval(time.Millisecond)

	return Model{
		textStyle: style,
		stopwatch: stopwatch,
		border:    false,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.stopwatch.Init(), tea.EnterAltScreen)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		fullscreenStyle.Width(msg.Width)
		fullscreenStyle.Height(msg.Height)
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		if key.Matches(msg, quitKeys) {
			m.quitting = true
			return m, tea.Quit

		} else if key.Matches(msg, borderKeys) {
			m.border = !m.border
		}
		return m, nil
	case errMsg:
		m.err = msg
		return m, nil

	default:
		var cmd tea.Cmd
		m.stopwatch, cmd = m.stopwatch.Update(msg)
		return m, cmd
	}
}

func (m Model) View() string {
	if m.err != nil {
		return m.err.Error()
	}

	str := fmt.Sprintf(
		"%-*s %*s %*s",
		20,
		"c9s",
		(m.width-20-12)/2,
		fmt.Sprintf("width:%d", (m.width-20-15)/2),
		(m.width-20-12)/2,
		quitKeys.Help().Desc,
	)

	if m.quitting {
		return str + "\n"
	}

	testStr := `
	oh
	hai
	wut
	up
	`
	return fullscreenStyle.Height(m.height).Render(
		lipgloss.JoinVertical(
			lipgloss.Top,
			headerStyle.Width(m.width).Render(str),
			testStyle.Height(m.height-5).Width(m.width).Margin(0).Padding(0).BorderStyle(lipgloss.NormalBorder()).BorderTop(false).BorderBottom(true).Render(testStr),
			testStyle.Width(m.width).Margin(0).Padding(0).BorderStyle(lipgloss.NormalBorder()).BorderTop(false).BorderBottom(false).Render(fmt.Sprintf("%s:%d", "footer, height", m.height)),
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
