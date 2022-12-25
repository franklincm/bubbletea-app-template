package main

import "github.com/charmbracelet/lipgloss"

var headerStyle = lipgloss.NewStyle().
	Margin(0).
	Padding(0)

var bodyStyle = lipgloss.NewStyle().
	Margin(0).
	Padding(0).
	BorderBottom(true).
	BorderTop(true).
	BorderStyle(lipgloss.NormalBorder()).
	Align(lipgloss.Center).
	AlignVertical(lipgloss.Center)

var footerStyle = lipgloss.NewStyle().
	Margin(0).
	Padding(0)

var fullscreenStyle = lipgloss.NewStyle().
	Padding(0).
	Margin(0).
	Align(lipgloss.Center)
