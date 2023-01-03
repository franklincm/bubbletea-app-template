package frame

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var style = lipgloss.NewStyle().
	Align(lipgloss.Center).
	Margin(0).
	Padding(0)

type Model struct {
	style   lipgloss.Style
	width   int
	height  int
	content tea.Model
}

func New() *Model {
	return &Model{
		style: style,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.content, cmd = m.content.Update(msg)
		cmds = append(cmds, cmd)

	default:
		m.content, cmd = m.content.Update(msg)
		cmds = append(cmds, cmd)
	}

	return &m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	return m.style.
		Width(m.width).
		Height(m.height).
		Render(m.content.View())
}

func (m Model) Content(content tea.Model) *Model {
	m.content = content
	return &m
}

func (m Model) GetContent() tea.Model {
	return m.content
}

func (m Model) Style(style lipgloss.Style) *Model {
	m.style = style
	return &m
}

func (m Model) GetStyle() lipgloss.Style {
	return m.style
}
