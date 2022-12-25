package main

import (
	"fmt"
	"os"

	key "github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
	frame "github.com/franklincm/bubbletea-template/components/frame"
)

type errMsg error

type componentId int

const (
	header componentId = iota
	body
	footer
)

var modelHeights = map[componentId]int{
	header: 2,
	footer: 2,
}

type Model struct {
	border   bool
	err      error
	quitting bool
	styles   map[componentId]lipgloss.Style
	frames   map[componentId]tea.Model
}

func New() Model {
	styles := map[componentId]lipgloss.Style{
		header: headerStyle,
		body:   bodyStyle,
		footer: footerStyle,
	}

	frames := map[componentId]tea.Model{
		header: frame.New().Content("header"),
		body:   frame.New().Content("body"),
		footer: frame.New().Content("footer"),
	}

	for component := range frames {
		frames[component] = frames[component].(frame.Model).Style(styles[component])
	}

	return Model{
		border: false,
		frames: frames,
		styles: styles,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:

		m.frames[header], cmd = m.frames[header].Update(tea.WindowSizeMsg{
			Width:  msg.Width,
			Height: modelHeights[header],
		})
		cmds = append(cmds, cmd)

		m.frames[footer], cmd = m.frames[footer].Update(tea.WindowSizeMsg{
			Width:  msg.Width,
			Height: modelHeights[footer],
		})
		cmds = append(cmds, cmd)

		m.frames[body], cmd = m.frames[body].Update(tea.WindowSizeMsg{
			Width:  msg.Width,
			Height: msg.Height - modelHeights[header] - modelHeights[footer] - 3,
		})
		cmds = append(cmds, cmd)
		cmds = append(cmds, tea.ClearScreen)

	case tea.KeyMsg:
		if key.Matches(msg, quitKeys) {
			m.quitting = true
			return m, tea.Quit
		}

	case errMsg:
		m.err = msg
		return m, nil
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
