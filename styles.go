package main

import (
	"github.com/charmbracelet/lipgloss"
	config "github.com/franklincm/bubbletea-template/config"
)

type Styles struct {
	bodyStyle       lipgloss.Style
	footerStyle     lipgloss.Style
	fullscreenStyle lipgloss.Style
	headerStyle     lipgloss.Style
	navStyle        lipgloss.Style
	tabBlurredStyle lipgloss.Style
	tabFocusedStyle lipgloss.Style
	tableCell       lipgloss.Style
	tableHeader     lipgloss.Style
	tableSelected   lipgloss.Style
}

func NewStyle(conf config.Config) Styles {
	return Styles{
		bodyStyle: lipgloss.NewStyle().
			Margin(0).
			Padding(0).
			BorderBottom(true).
			BorderTop(true).
			BorderRight(true).
			BorderLeft(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color(conf.Styles["borderForeground"])).
			Align(lipgloss.Center).
			AlignVertical(lipgloss.Center),

		footerStyle: lipgloss.NewStyle().
			Margin(0).
			Padding(0).
			Align(lipgloss.Left),

		fullscreenStyle: lipgloss.NewStyle().
			Padding(0).
			Margin(0).
			Align(lipgloss.Center),

		headerStyle: lipgloss.NewStyle().
			Margin(0).
			Padding(0).
			Align(lipgloss.Right),

		navStyle: lipgloss.NewStyle().
			BorderForeground(lipgloss.Color(conf.Styles["borderForeground"])).
			BorderTop(true).
			BorderStyle(lipgloss.NormalBorder()),

		tabBlurredStyle: lipgloss.NewStyle().
			MarginRight(1).
			Background(lipgloss.Color(conf.Styles["tabBlurredBackground"])).
			Foreground(lipgloss.Color(conf.Styles["tabBlurredForeground"])).
			Align(lipgloss.Center),

		tabFocusedStyle: lipgloss.NewStyle().
			MarginRight(1).
			Background(lipgloss.Color(conf.Styles["tabFocusedBackground"])).
			Foreground(lipgloss.Color(conf.Styles["tabFocusedForeground"])).
			Align(lipgloss.Center),

		tableCell: lipgloss.
			NewStyle().
			Padding(0, 0),

		tableHeader: lipgloss.
			NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color(conf.Styles["tableHeaderForeground"])).
			BorderBottom(true).
			Bold(false),

		tableSelected: lipgloss.
			NewStyle().
			Bold(true).
			Background(lipgloss.Color(conf.Styles["tableSelectedBackground"])).
			Foreground(lipgloss.Color(conf.Styles["tableSelectedForeground"])),
	}
}
