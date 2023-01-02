package main

import (
	"fmt"
	"os"

	key "github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
	commandprompt "github.com/franklincm/bubbles/commandPrompt"
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

var frameHeights = map[frameId]int{
	header: 2,
	footer: 2,
}

type Model struct {
	width  int
	height int

	err         error
	quitting    bool
	spinner     tea.Model
	headerModel tea.Model
	footerModel tea.Model
	prompt      tea.Model
	input       string
	showprompt  bool

	frames map[frameId]*frame.Model
}

func New() Model {
	s := spinner.New()
	s.Spinner = spinner.Dot

	p := commandprompt.New(": ")
	p.InputShow = key.NewBinding(key.WithKeys(":"))

	m := Model{
		headerModel: text.New().Content("header"),
		footerModel: text.New().Content("footer"),
		spinner:     s,
		prompt:      p,
	}

	frames := map[frameId]*frame.Model{
		header: frame.
			New().
			Style(headerStyle).
			Content(m.headerModel),
		body: frame.
			New().
			Style(bodyStyle).
			Content(m.spinner),
		footer: frame.
			New().
			Style(footerStyle).
			Content(m.footerModel),
	}

	m.frames = frames

	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		m.spinner.(spinner.Model).Tick,
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
		m.width = msg.Width
		m.height = msg.Height

		// update header
		m.frames[header], cmd = m.frames[header].Update(tea.WindowSizeMsg{
			Width:  msg.Width,
			Height: frameHeights[header],
		})
		cmds = append(cmds, cmd)

		// update footer
		m.frames[footer], cmd = m.frames[footer].Update(tea.WindowSizeMsg{
			Width:  msg.Width,
			Height: frameHeights[footer],
		})
		cmds = append(cmds, cmd)

		// update body
		m.frames[body], cmd = m.frames[body].Update(tea.WindowSizeMsg{
			Width:  msg.Width,
			Height: msg.Height - frameHeights[header] - frameHeights[footer] - len(frameHeights),
		})
		cmds = append(cmds, cmd)

		// redraw
		cmds = append(cmds, tea.ClearScreen)

	case commandprompt.PromptInput:
		if msg == "quit" || msg == "q" {
			m.quitting = true
			return m, tea.Quit

		} else if msg == "c" {
			tmp := m.frames[header].GetContent()
			m.frames[header] = m.frames[header].Content(m.frames[body].GetContent())
			m.frames[body] = m.frames[body].Content(tmp)

		}

	case commandprompt.PromptEditing:
		m.showprompt = bool(msg)
		return m, nil

	case tea.KeyMsg:
		if key.Matches(msg, quitKeys) && !m.showprompt {
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

	m.prompt, cmd = m.prompt.Update(msg)
	cmds = append(cmds, cmd)

	if m.showprompt {
		m.frames[footer] = m.frames[footer].Content(m.prompt)
	} else {
		m.frames[footer] = m.frames[footer].Content(m.footerModel)
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
