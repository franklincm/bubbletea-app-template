package text

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	style   lipgloss.Style
	width   int
	height  int
	content string
}

func New() Model {
	return Model{
		// style:   style,
		content: "",
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m Model) View() string {
	return m.content
}

func (m Model) Content(content string) Model {
	m.content = content
	return m
}
