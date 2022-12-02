package header

import (
	"github.com/charmbracelet/lipgloss"
)

var style = lipgloss.NewStyle().
	Padding(0).
	Margin(0).
	Align(lipgloss.Center)

type Header struct {
	text string
}
