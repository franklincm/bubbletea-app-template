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
	nav:    1,
	footer: 2,
}

type Model struct {
	err      error
	quitting bool

	width  int
	height int

	inputSuggestionCounter int
	inputHint              string
	suggestions            []string
	showprompt             bool

	footerText tea.Model
	prompt     tea.Model
	tabs       tea.Model

	layout Layout
}

func New() Model {
	prompt := commandprompt.New(":")
	prompt.InputShow = key.NewBinding(key.WithKeys(":"))

	layout := NewLayout()

	footerFrame := NewFooter()

	m := Model{
		footerText: footerFrame.GetContent()[0],
		prompt:     prompt,
		tabs:       layout.tabs,
		layout:     layout,
	}

	m.setActiveTab(m.layout.tabNameToIndex["table"])

	return m
}

func (m Model) Init() tea.Cmd {
	m.Update(commandprompt.PromptInput(command))
	command = ""

	return tea.Batch(
		tea.EnterAltScreen,
		m.layout.models["two"].Init(),
		m.layout.models["three"].Init(),
		m.layout.models["four"].Init(),
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
		m.layout.frames[header], cmd = m.layout.frames[header].Update(tea.WindowSizeMsg{
			Width:  msg.Width,
			Height: frameHeights[header],
		})
		cmds = append(cmds, cmd)

		// update nav
		m.layout.frames[nav], cmd = m.layout.frames[nav].Update(tea.WindowSizeMsg{
			Width:  msg.Width,
			Height: frameHeights[nav],
		})
		cmds = append(cmds, cmd)

		// update footer
		m.layout.frames[footer], cmd = m.layout.frames[footer].Update(tea.WindowSizeMsg{
			Width:  msg.Width,
			Height: frameHeights[footer],
		})
		cmds = append(cmds, cmd)

		// update body
		m.layout.frames[body], cmd = m.layout.frames[body].Update(tea.WindowSizeMsg{
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
			m.layout.models["table"] = dataTable

			if activeTab == 0 {
				m.SetContent(m.layout.models["table"])
			}
		} else if msg == "city" {
			dataTable.SetColumns(cityColumns)
			dataTable.SetRows(cityRows)
			m.layout.models["table"] = dataTable

			if activeTab == 0 {
				m.SetContent(m.layout.models["table"])
			}
		} else {
			model, ok := m.layout.models[string(msg)]
			if ok {
				m.setActiveTab(m.layout.tabNameToIndex[string(msg)])
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
			m.SetContent(m.layout.models[m.layout.tabs.(tabs.Model).GetHeadings()[activeTab]])
			return m, cmd
		} else if key.Matches(msg, key.NewBinding(key.WithKeys(conf.Keys["global"]["right"]))) && !m.showprompt {
			m.tabNext()
			m.SetContent(m.layout.models[m.layout.tabs.(tabs.Model).GetHeadings()[activeTab]])
			return m, cmd

			// table key bindings
		} else if key.Matches(msg, key.NewBinding(key.WithKeys(conf.Keys["global"]["down"]))) && !m.showprompt {
			m.layout.frames[body], cmd = m.layout.frames[body].Update(msg)
			return m, cmd
		} else if key.Matches(msg, key.NewBinding(key.WithKeys(conf.Keys["global"]["up"]))) && !m.showprompt {
			m.layout.frames[body], cmd = m.layout.frames[body].Update(msg)
			return m, cmd
		} else if key.Matches(msg, key.NewBinding(key.WithKeys(conf.Keys["global"]["halfPageUp"]))) && !m.showprompt {
			m.layout.frames[body], cmd = m.layout.frames[body].Update(msg)
			return m, cmd
		} else if key.Matches(msg, key.NewBinding(key.WithKeys(conf.Keys["global"]["halfPageDown"]))) && !m.showprompt {
			m.layout.frames[body], cmd = m.layout.frames[body].Update(msg)
			return m, cmd
		} else if key.Matches(msg, key.NewBinding(key.WithKeys(conf.Keys["global"]["pageDown"]))) && !m.showprompt {
			m.layout.frames[body], cmd = m.layout.frames[body].Update(msg)
			return m, cmd
		} else if key.Matches(msg, key.NewBinding(key.WithKeys(conf.Keys["global"]["pageUp"]))) && !m.showprompt {
			m.layout.frames[body], cmd = m.layout.frames[body].Update(msg)
			return m, cmd
		} else if key.Matches(msg, key.NewBinding(key.WithKeys("g"))) && !m.showprompt {
			m.layout.frames[body], cmd = m.layout.frames[body].Update(msg)
			return m, cmd
		} else if key.Matches(msg, key.NewBinding(key.WithKeys("G"))) && !m.showprompt {
			m.layout.frames[body], cmd = m.layout.frames[body].Update(msg)
			return m, cmd

			// tab suggestions
		} else if key.Matches(msg, key.NewBinding(key.WithKeys("tab"))) && m.showprompt {
			headings := m.layout.tabs.(tabs.Model).GetHeadings()

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
		for f := range m.layout.frames {
			m.layout.frames[f], cmd = m.layout.frames[f].Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	m.layout.tabs, cmd = m.layout.tabs.Update(msg)
	cmds = append(cmds, cmd)

	// update main models
	for model := range m.layout.models {
		m.layout.models[model], cmd = m.layout.models[model].Update(msg)
		cmds = append(cmds, cmd)
	}

	if m.showprompt {
		if err := m.layout.frames[footer].SetContentModel(0, m.prompt); err != nil {
			log.Println(err)
		}
	} else {
		if err := m.layout.frames[footer].SetContentModel(0, m.footerText); err != nil {
			log.Println(err)
		}
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
			m.layout.frames[header].View(),
			m.layout.frames[nav].View(),
			m.layout.frames[body].View(),
			m.layout.frames[footer].View(),
		),
	)
}

func (m Model) SetContent(tm tea.Model) {
	m.layout.frames[nav] = m.layout.frames[nav].Content(
		[]tea.Model{
			m.layout.tabs,
		},
	)

	m.layout.frames[body] = m.layout.frames[body].Content(
		[]tea.Model{tm},
	)

	m.layout.frames[body], _ = m.layout.frames[body].Update(tea.WindowSizeMsg{
		Width:  styles.bodyStyle.GetWidth(),
		Height: styles.bodyStyle.GetHeight(),
	})
}

func (m *Model) setActiveTab(index int) {
	if index >= 0 && index <= len(m.layout.tabs.(tabs.Model).GetHeadings()) {
		activeTab = index
		m.layout.tabs = m.layout.tabs.(tabs.Model).SetFocused(activeTab)
	}
}

func (m *Model) tabNext() {
	numHeadings := len(m.layout.tabs.(tabs.Model).GetHeadings())
	tabNext := int(math.Min(float64(activeTab+1), float64(numHeadings-1)))
	activeTab = tabNext
	m.layout.tabs = m.layout.tabs.(tabs.Model).SetFocused(activeTab)
}

func (m *Model) tabPrev() {
	tabPrev := int(math.Max(float64(activeTab-1), 0))
	activeTab = tabPrev
	m.layout.tabs = m.layout.tabs.(tabs.Model).SetFocused(activeTab)
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
