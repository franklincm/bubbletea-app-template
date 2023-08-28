package tabs

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	labels       []string
	focused      int
	focusedStyle lipgloss.Style
	blurredStyle lipgloss.Style
	width        int
	height       int
	margin       int
}

func New(labels []string) Model {
	margin := 1

	return Model{
		focused: 0,
		labels:  labels,
		margin:  margin,
		blurredStyle: lipgloss.NewStyle().
			MarginRight(margin).
			Background(lipgloss.Color("7")).
			Foreground(lipgloss.Color("0")).
			Align(lipgloss.Center),
		focusedStyle: lipgloss.NewStyle().
			MarginRight(margin).
			Background(lipgloss.Color("4")).
			Foreground(lipgloss.Color("255")).
			Align(lipgloss.Center),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {

	margin := 1
	cellWidth := ((m.width) / len(m.labels)) - margin
	var elems []string

	for i, v := range m.labels {
		if i == m.focused {
			elems = append(
				elems,
				m.focusedStyle.Render(
					lipgloss.PlaceHorizontal(
						cellWidth,
						lipgloss.Center, v,
					),
				),
			)
		} else {
			elems = append(
				elems,
				m.blurredStyle.Render(
					lipgloss.PlaceHorizontal(
						cellWidth,
						lipgloss.Center, v,
					),
				),
			)
		}
	}

	return lipgloss.NewStyle().
		Width(m.width).
		Render(lipgloss.JoinHorizontal(lipgloss.Center, elems...))
}

func (m Model) FocusedStyle(style lipgloss.Style) Model {
	m.focusedStyle = style
	return m
}

func (m Model) BlurredStyle(style lipgloss.Style) Model {
	m.blurredStyle = style
	return m
}

func (m Model) SetFocused(index int) Model {
	if index < len(m.labels) && index >= 0 {
		m.focused = index
	}

	return m
}

func (m Model) GetHeadings() []string {
	return m.labels
}
