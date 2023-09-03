package main

import (
	tea "github.com/charmbracelet/bubbletea"
	frame "github.com/franklincm/bubbletea-template/components/frame"
	text "github.com/franklincm/bubbletea-template/components/text"
)

func NewFooterText() tea.Model {
	return text.New().Content("footer")
}

func NewBlank() tea.Model {
	return text.New().Content("")
}

func NewFooter() *frame.Model {
	return frame.
		New().
		Style(styles.footerStyle).
		Content(
			[]tea.Model{
				NewFooterText(),
				NewBlank(),
			},
		)
}
