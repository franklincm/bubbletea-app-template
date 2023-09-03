package main

import (
	tea "github.com/charmbracelet/bubbletea"
	frame "github.com/franklincm/bubbletea-template/components/frame"
	text "github.com/franklincm/bubbletea-template/components/text"
)

func NewInfoModel() tea.Model {
	infoText := `
app template
version: 0.0.1
theme: gruvbox
`

	return text.New().Content(infoText)
}

func NewLogo() tea.Model {
	return text.New().Content(logo)
}

func NewHeader() *frame.Model {
	return frame.
		New().
		Kind(frame.Horizontal).
		Style(styles.headerStyle).
		Content(
			[]tea.Model{
				NewInfoModel(),
				NewLogo(),
			},
		)
}
