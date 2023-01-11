package main

import (
	"fmt"
	"math"
	"os"

	key "github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
	commandprompt "github.com/franklincm/bubbles/commandPrompt"
	tabs "github.com/franklincm/bubbles/tabs"
	spinner "github.com/franklincm/bubbletea-template/components/spinner"
	table "github.com/franklincm/bubbletea-template/components/table"
	text "github.com/franklincm/bubbletea-template/components/text"
	vframe "github.com/franklincm/bubbletea-template/components/vframe"
	config "github.com/franklincm/bubbletea-template/config"
)

type errMsg error

type frameId int

const (
	header frameId = iota
	body
	footer
)

var conf = config.New()

var frameHeights = map[frameId]int{
	header: 9,
	footer: 2,
}

var dataTable table.Model

type Model struct {
	err      error
	quitting bool

	width  int
	height int

	cursor     int
	input      string
	showprompt bool

	headerModel tea.Model
	footerModel tea.Model
	prompt      tea.Model
	tabs        tea.Model

	frames    map[frameId]*vframe.Model
	models    map[string]tea.Model
	viewOrder map[int]string
	viewPos   map[string]int
}

func New() Model {
	s1 := spinner.New()
	s1.Spinner = spinner.Dot

	s2 := spinner.New()
	s2.Spinner = spinner.Points

	s3 := spinner.New()
	s3.Spinner = spinner.Pulse

	prompt := commandprompt.New(":")
	prompt.InputShow = key.NewBinding(key.WithKeys(":"))

	dataTable = table.New(
		table.WithColumns(cityColumns),
		table.WithRows(cityRows),
		table.WithFocused(true),
	)

	headings := []string{
		"data",
		"two",
		"three",
		"four",
	}
	var viewOrder = map[int]string{}
	var viewPos = map[string]int{}
	for i, label := range headings {
		viewOrder[i] = label
		viewPos[label] = i
	}

	cursor := 2
	tabs := tabs.New(headings)
	tabs = tabs.FocusedStyle(tabFocusedStyle)
	tabs = tabs.BlurredStyle(tabBlurredStyle)
	tabs = tabs.SetFocused(cursor)

	m := Model{
		headerModel: text.New().Content(logo),
		footerModel: text.New().Content("footer"),
		prompt:      prompt,
		tabs:        tabs,
		cursor:      cursor,
		models: map[string]tea.Model{
			"data":  dataTable,
			"two":   s1,
			"three": s2,
			"four":  s3,
		},
		viewOrder: viewOrder,
		viewPos:   viewPos,
	}

	frames := map[frameId]*vframe.Model{
		header: vframe.
			New().
			Style(headerStyle).
			Content(
				[]tea.Model{
					m.headerModel,
					m.tabs,
				},
			),
		body: vframe.
			New().
			Style(bodyStyle).
			Content(
				[]tea.Model{m.models["three"]},
			),
		footer: vframe.
			New().
			Style(footerStyle).
			Content(
				[]tea.Model{m.footerModel},
			),
	}

	m.frames = frames

	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		m.models["two"].(spinner.Model).Tick,
		m.models["three"].(spinner.Model).Tick,
		m.models["four"].(spinner.Model).Tick,
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
			Width:  msg.Width - 2,
			Height: msg.Height - frameHeights[header] - frameHeights[footer] - len(frameHeights),
		})
		cmds = append(cmds, cmd)

		// redraw
		cmds = append(cmds, tea.ClearScreen)

	case commandprompt.PromptInput:
		if msg == "quit" || msg == "q" {
			m.quitting = true
			return m, tea.Quit

		} else if msg == "dnd" {

			dataTable.SetRows(charRows)
			dataTable.SetColumns(charColumns)
			m.models["one"] = dataTable

			if m.cursor == 0 {
				m.SetContent(m.models["one"])
			}

		} else if msg == "city" {
			dataTable.SetColumns(cityColumns)
			dataTable.SetRows(cityRows)
			m.models["one"] = dataTable

			if m.cursor == 0 {
				m.SetContent(m.models["one"])
			}

		} else {
			model, ok := m.models[string(msg)]
			if ok {
				m.setCursor(m.viewPos[string(msg)])
				m.SetContent(model)
			}
		}

	case commandprompt.PromptEditing:
		m.showprompt = bool(msg)
		return m, nil

	case tea.KeyMsg:
		if key.Matches(msg, quitKeys) && !m.showprompt {
			m.quitting = true
			return m, tea.Quit

		} else if key.Matches(msg, key.NewBinding(key.WithKeys(conf.Keys["global"]["left"]))) && !m.showprompt {
			m.cursorPrev()
			m.SetContent(m.models[m.viewOrder[m.cursor]])

			return m, cmd

		} else if key.Matches(msg, key.NewBinding(key.WithKeys(conf.Keys["global"]["right"]))) && !m.showprompt {
			m.cursorNext()
			m.SetContent(m.models[m.viewOrder[m.cursor]])

		} else if key.Matches(msg, key.NewBinding(key.WithKeys(conf.Keys["global"]["down"]))) && !m.showprompt {
			m.frames[body], cmd = m.frames[body].Update(msg)
			return m, cmd

		} else if key.Matches(msg, key.NewBinding(key.WithKeys(conf.Keys["global"]["up"]))) && !m.showprompt {
			m.frames[body], cmd = m.frames[body].Update(msg)
			return m, cmd

		}

		m.prompt, cmd = m.prompt.Update(msg)
		cmds = append(cmds, cmd)

		return m, tea.Batch(cmds...)

	case errMsg:
		m.err = msg
		return m, nil

	default:
		for f := range m.frames {
			m.frames[f], cmd = m.frames[f].Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	m.tabs, cmd = m.tabs.Update(msg)
	cmds = append(cmds, cmd)

	// update main models
	for model := range m.models {
		m.models[model], cmd = m.models[model].Update(msg)
		cmds = append(cmds, cmd)
	}

	if m.showprompt {
		m.frames[footer] = m.frames[footer].Content(
			[]tea.Model{m.prompt},
		)
	} else {
		m.frames[footer] = m.frames[footer].Content(
			[]tea.Model{m.footerModel},
		)
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

func (m Model) SetContent(tm tea.Model) {
	m.frames[header] = m.frames[header].Content(
		[]tea.Model{
			m.headerModel,
			m.tabs,
		},
	)

	m.frames[body] = m.frames[body].Content(
		[]tea.Model{tm},
	)

	m.frames[body], _ = m.frames[body].Update(tea.WindowSizeMsg{
		Width:  bodyStyle.GetWidth(),
		Height: bodyStyle.GetHeight(),
	})
}

func (m *Model) setCursor(index int) {
	if index >= 0 && index <= len(m.tabs.(tabs.Model).GetHeadings()) {
		m.cursor = index
		m.tabs = m.tabs.(tabs.Model).SetFocused(m.cursor)
	}
}

func (m *Model) cursorNext() {
	numHeadings := len(m.tabs.(tabs.Model).GetHeadings())
	m.cursor = int(math.Min(float64(m.cursor+1), float64(numHeadings-1)))
	m.tabs = m.tabs.(tabs.Model).SetFocused(m.cursor)
}

func (m *Model) cursorPrev() {
	m.cursor = int(math.Max(float64(m.cursor-1), 0))
	m.tabs = m.tabs.(tabs.Model).SetFocused(m.cursor)
}

func main() {
	p := tea.NewProgram(New(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
