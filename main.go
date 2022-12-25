package main

import (
	"fmt"
	"os"

	key "github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
	"github.com/franklincm/bubbletea-template/components"
	frame "github.com/franklincm/bubbletea-template/components/frame"
	spinner "github.com/franklincm/bubbletea-template/components/spinner"
	text "github.com/franklincm/bubbletea-template/components/text"
)

type errMsg error

type frameId int

const (
	header frameId = iota
	body
	footer
)

var modelHeights = map[frameId]int{
	header: 2,
	footer: 2,
}

type Model struct {
	err        error
	quitting   bool
	styles     map[frameId]lipgloss.Style
	frames     map[frameId]tea.Model
	textModel  components.Component
	spinner    components.Component
	components map[string]components.Component
}

func New() Model {
	s := spinner.New()
	s.Spinner = spinner.Points

	m := Model{
		textModel: text.New().Content("text widget"),
		spinner:   s,

		styles: map[frameId]lipgloss.Style{
			header: headerStyle,
			body:   bodyStyle,
			footer: footerStyle,
		},

		components: map[string]components.Component{
			"spinner":   s,
			"textwidet": text.New().Content("text widget again"),
		},
	}

	frames := map[frameId]tea.Model{
		header: frame.New().Content(m.textModel),
		body:   frame.New().Content(m.spinner),
		footer: frame.New().Content(m.textModel),
	}

	for f := range frames {
		frames[f] = frames[f].(frame.Model).Style(m.styles[f])
	}

	m.frames = frames

	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		m.spinner.(spinner.Model).Tick,
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:

		// update header
		m.frames[header], cmd = m.frames[header].Update(tea.WindowSizeMsg{
			Width:  msg.Width,
			Height: modelHeights[header],
		})
		cmds = append(cmds, cmd)

		// update body
		m.frames[footer], cmd = m.frames[footer].Update(tea.WindowSizeMsg{
			Width:  msg.Width,
			Height: modelHeights[footer],
		})
		cmds = append(cmds, cmd)

		// update footer
		m.frames[body], cmd = m.frames[body].Update(tea.WindowSizeMsg{
			Width:  msg.Width,
			Height: msg.Height - modelHeights[header] - modelHeights[footer] - 3,
		})
		cmds = append(cmds, cmd)

		// redraw
		cmds = append(cmds, tea.ClearScreen)

	case tea.KeyMsg:
		if key.Matches(msg, quitKeys) {
			m.quitting = true
			return m, tea.Quit
		}

	case errMsg:
		m.err = msg
		return m, nil

	default:
		for f := range m.frames {
			m.frames[f], cmd = m.frames[f].Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.err != nil {
		return m.err.Error()
	}

	return fullscreenStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Center,
			m.frames[header].View(),
			m.frames[body].View(),
			m.frames[footer].View(),
		),
	)
}

func main() {
	p := tea.NewProgram(New(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
