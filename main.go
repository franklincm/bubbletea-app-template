package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strings"

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
	nav
	body
	footer
)

var (
	command   string
	activeTab int
)

var (
	conf   = config.New()
	styles = NewStyle(conf)
)

var frameHeights = map[frameId]int{
	header: 8,
	nav:    2,
	footer: 2,
}

var dataTable table.Model

type Model struct {
	err      error
	quitting bool

	width  int
	height int

	inputSuggestionCounter int
	inputHint              string
	suggestions            []string
	showprompt             bool

	headerModel tea.Model
	blank       tea.Model
	info        tea.Model
	footerModel tea.Model
	prompt      tea.Model
	tabs        tea.Model

	frames       map[frameId]*vframe.Model
	models       map[string]tea.Model
	tabPosLookup map[string]int
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
		table.WithStyles(table.Styles{
			Cell:     styles.tableCell,
			Header:   styles.tableHeader,
			Selected: styles.tableSelected,
		}),
		table.WithKeyMap(table.KeyMap{
			LineUp: key.NewBinding(
				key.WithKeys(conf.Keys["global"]["up"]),
				key.WithHelp(fmt.Sprintf("↑/%s", conf.Keys["global"]["up"]), "up"),
			),
			LineDown: key.NewBinding(
				key.WithKeys(conf.Keys["global"]["down"]),
				key.WithHelp(fmt.Sprintf("↓/%s", conf.Keys["global"]["down"]), "down"),
			),
			PageUp: key.NewBinding(
				key.WithKeys(conf.Keys["global"]["pageUp"]),
				key.WithHelp(conf.Keys["global"]["pageUp"], "pageUp"),
			),
			PageDown: key.NewBinding(
				key.WithKeys(conf.Keys["global"]["pageDown"]),
				key.WithHelp(conf.Keys["global"]["pageDown"], "pageDown"),
			),
			HalfPageUp: key.NewBinding(
				key.WithKeys(conf.Keys["global"]["halfPageUp"]),
				key.WithHelp(conf.Keys["global"]["halfPageUp"], "halfPageUp"),
			),
			HalfPageDown: key.NewBinding(
				key.WithKeys(conf.Keys["global"]["halfPageDown"]),
				key.WithHelp(conf.Keys["global"]["halfPageDown"], "halfPageDown"),
			),
			GotoTop: key.NewBinding(
				key.WithKeys("home", "g"),
				key.WithHelp("g/home", "go to start"),
			),
			GotoBottom: key.NewBinding(
				key.WithKeys("end", "G"),
				key.WithHelp("G/end", "go to end"),
			),
		}),
	)

	headings := []string{
		"table",
		"two",
		"three",
		"four",
	}

	models := map[string]tea.Model{
		headings[0]: dataTable,
		headings[1]: s1,
		headings[2]: s2,
		headings[3]: s3,
	}

	tabPosLookup := map[string]int{}
	for i, label := range headings {
		tabPosLookup[label] = i
	}

	tabs := tabs.New(headings)
	tabs = tabs.FocusedStyle(styles.tabFocusedStyle)
	tabs = tabs.BlurredStyle(styles.tabBlurredStyle)
	tabs = tabs.SetFocused(activeTab)

	infoText := `
app template
version: 0.0.1
theme: gruvbox
`

	m := Model{
		headerModel:  text.New().Content(logo),
		footerModel:  text.New().Content("footer"),
		info:         text.New().Content(infoText),
		blank:        text.New().Content(""),
		prompt:       prompt,
		tabs:         tabs,
		models:       models,
		tabPosLookup: tabPosLookup,
	}

	frames := map[frameId]*vframe.Model{
		header: vframe.
			New().
			Kind(vframe.Horizontal).
			Style(styles.headerStyle).
			Content(
				[]tea.Model{
					m.info,
					m.headerModel,
				},
			),
		nav: vframe.
			New().
			Content(
				[]tea.Model{
					m.tabs,
				},
			),
		body: vframe.
			New().
			Style(styles.bodyStyle).
			Content(
				[]tea.Model{m.models["three"]},
			),
		footer: vframe.
			New().
			Style(styles.footerStyle).
			Content(
				[]tea.Model{
					m.footerModel,
					m.blank,
				},
			),
	}

	m.frames = frames

	return m
}

func (m Model) Init() tea.Cmd {
	m.Update(commandprompt.PromptInput(command))
	command = ""

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

		// update nav
		m.frames[nav], cmd = m.frames[nav].Update(tea.WindowSizeMsg{
			Width:  msg.Width,
			Height: frameHeights[nav],
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
		log.Printf("prompt input: %#v", msg)
		m.suggestions = nil
		m.inputSuggestionCounter = 0
		m.inputHint = ""

		if msg == "quit" || msg == "q" {
			m.quitting = true
			return m, tea.Quit

		} else if msg == "dnd" {

			dataTable.SetRows(charRows)
			dataTable.SetColumns(charColumns)
			m.models["table"] = dataTable

			if activeTab == 0 {
				m.SetContent(m.models["table"])
			}

		} else if msg == "city" {
			dataTable.SetColumns(cityColumns)
			dataTable.SetRows(cityRows)
			m.models["table"] = dataTable

			if activeTab == 0 {
				m.SetContent(m.models["table"])
			}

		} else {
			model, ok := m.models[string(msg)]
			if ok {
				m.setActiveTab(m.tabPosLookup[string(msg)])
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
			m.tabPrev()
			m.SetContent(m.models[m.tabs.(tabs.Model).GetHeadings()[activeTab]])
			return m, cmd

		} else if key.Matches(msg, key.NewBinding(key.WithKeys(conf.Keys["global"]["right"]))) && !m.showprompt {
			m.tabNext()
			m.SetContent(m.models[m.tabs.(tabs.Model).GetHeadings()[activeTab]])
			return m, cmd

			// table key bindings
		} else if key.Matches(msg, key.NewBinding(key.WithKeys(conf.Keys["global"]["down"]))) && !m.showprompt {
			m.frames[body], cmd = m.frames[body].Update(msg)
			return m, cmd

		} else if key.Matches(msg, key.NewBinding(key.WithKeys(conf.Keys["global"]["up"]))) && !m.showprompt {
			m.frames[body], cmd = m.frames[body].Update(msg)
			return m, cmd

		} else if key.Matches(msg, key.NewBinding(key.WithKeys(conf.Keys["global"]["halfPageUp"]))) && !m.showprompt {
			m.frames[body], cmd = m.frames[body].Update(msg)
			return m, cmd

		} else if key.Matches(msg, key.NewBinding(key.WithKeys(conf.Keys["global"]["halfPageDown"]))) && !m.showprompt {
			m.frames[body], cmd = m.frames[body].Update(msg)
			return m, cmd

		} else if key.Matches(msg, key.NewBinding(key.WithKeys(conf.Keys["global"]["pageDown"]))) && !m.showprompt {
			m.frames[body], cmd = m.frames[body].Update(msg)
			return m, cmd

		} else if key.Matches(msg, key.NewBinding(key.WithKeys(conf.Keys["global"]["pageUp"]))) && !m.showprompt {
			m.frames[body], cmd = m.frames[body].Update(msg)
			return m, cmd

		} else if key.Matches(msg, key.NewBinding(key.WithKeys("g"))) && !m.showprompt {
			m.frames[body], cmd = m.frames[body].Update(msg)
			return m, cmd

		} else if key.Matches(msg, key.NewBinding(key.WithKeys("G"))) && !m.showprompt {
			m.frames[body], cmd = m.frames[body].Update(msg)
			return m, cmd

			// tab suggestions
		} else if key.Matches(msg, key.NewBinding(key.WithKeys("tab"))) && m.showprompt {

			headings := m.tabs.(tabs.Model).GetHeadings()

			if len(m.inputHint) == 0 {
				m.inputHint = m.prompt.(commandprompt.Model).TextInput.Value()
			}

			log.Println("hint: ", m.inputHint)

			if len(m.suggestions) == 0 {
				for _, heading := range headings {
					if strings.Contains(strings.ToLower(heading), m.inputHint) {
						m.suggestions = append(m.suggestions, heading)
					}
				}
			}

			if len(m.suggestions) > 0 {
				m.prompt = m.prompt.(commandprompt.Model).SetValue(
					m.suggestions[m.inputSuggestionCounter],
				)

				m.inputSuggestionCounter = (m.inputSuggestionCounter + 1) % len(m.suggestions)
				log.Println("suggestion: ", m.prompt.(commandprompt.Model).TextInput.Value())
			}

		} else {
			m.suggestions = nil
			m.inputSuggestionCounter = 0
			m.inputHint = ""
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
			[]tea.Model{
				m.prompt,
				m.blank,
			},
		)
	} else {
		m.frames[footer] = m.frames[footer].Content(
			[]tea.Model{
				m.footerModel,
				m.blank,
			},
		)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.err != nil {
		return m.err.Error()
	}

	return styles.fullscreenStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Center,
			m.frames[header].View(),
			m.frames[nav].View(),
			m.frames[body].View(),
			m.frames[footer].View(),
		),
	)
}

func (m Model) SetContent(tm tea.Model) {
	m.frames[header] = m.frames[header].Content(
		[]tea.Model{
			m.info,
			m.headerModel,
		},
	)

	m.frames[nav] = m.frames[nav].Content(
		[]tea.Model{
			m.tabs,
		},
	)

	m.frames[body] = m.frames[body].Content(
		[]tea.Model{tm},
	)

	m.frames[body], _ = m.frames[body].Update(tea.WindowSizeMsg{
		Width:  styles.bodyStyle.GetWidth(),
		Height: styles.bodyStyle.GetHeight(),
	})
}

func (m *Model) setActiveTab(index int) {
	if index >= 0 && index <= len(m.tabs.(tabs.Model).GetHeadings()) {
		activeTab = index
		m.tabs = m.tabs.(tabs.Model).SetFocused(activeTab)
	}
}

func (m *Model) tabNext() {
	numHeadings := len(m.tabs.(tabs.Model).GetHeadings())
	activeTab = int(math.Min(float64(activeTab+1), float64(numHeadings-1)))
	m.tabs = m.tabs.(tabs.Model).SetFocused(activeTab)
}

func (m *Model) tabPrev() {
	activeTab = int(math.Max(float64(activeTab-1), 0))
	m.tabs = m.tabs.(tabs.Model).SetFocused(activeTab)
}

func main() {
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal: ", err)
			os.Exit(1)
		}
		defer f.Close()
	} else {
		log.SetOutput(io.Discard)
	}

	log.Printf("\n\n\n\n")
	log.Printf("--------------------\n")

	flag.StringVar(&command, "c", "", "help")
	flag.Parse()

	p := tea.NewProgram(New(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
