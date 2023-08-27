package vframe

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type kind int

const (
	Horizontal kind = iota
	Vertical
)

var style = lipgloss.NewStyle().
	Align(lipgloss.Center).
	Margin(0).
	Padding(0)

type Model struct {
	style   lipgloss.Style
	width   int
	height  int
	content []tea.Model
	kind    kind
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
		for i := range m.content {
			m.content[i], cmd = m.content[i].Update(msg)
			cmds = append(cmds, cmd)
		}

	default:
		for i := range m.content {
			m.content[i], cmd = m.content[i].Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return &m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	var strs []string

	if m.kind == Horizontal {
		for i := range m.content {
			if i < len(m.content)-1 {
				strs = append(strs, lipgloss.PlaceHorizontal(m.width/len(m.content), lipgloss.Left, m.content[i].View()))
			} else if i == len(m.content)-1 {
				strs = append(strs, lipgloss.PlaceHorizontal(m.width/len(m.content), lipgloss.Right, m.content[i].View()))
			}
		}
		return m.style.
			Width(m.width).
			Height(m.height).
			Render(
				lipgloss.JoinHorizontal(m.style.GetAlign(), strs...),
			)
	} else {
		for i := range m.content {
			strs = append(strs, m.content[i].View())
		}
		return m.style.
			Width(m.width).
			Height(m.height).
			Render(
				lipgloss.JoinVertical(m.style.GetAlign(), strs...),
			)
	}
}

func (m Model) Append(content tea.Model) *Model {
	m.content = append(m.content, content)
	return &m
}

func (m Model) Content(content []tea.Model) *Model {
	m.content = content
	return &m
}

func (m Model) GetContent() []tea.Model {
	return m.content
}

func (m Model) Style(style lipgloss.Style) *Model {
	m.style = style
	return &m
}

func (m Model) GetStyle() lipgloss.Style {
	return m.style
}

func (m Model) Kind(t kind) *Model {
	m.kind = t
	return &m
}
