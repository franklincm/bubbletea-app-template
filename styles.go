package main

import "github.com/charmbracelet/lipgloss"

var headerStyle = lipgloss.NewStyle().
	Margin(0).
	Padding(0).
	Align(lipgloss.Left)

var bodyStyle = lipgloss.NewStyle().
	Margin(0).
	Padding(0).
	BorderBottom(true).
	BorderTop(true).
	BorderRight(true).
	BorderLeft(true).
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("39")).
	Align(lipgloss.Center).
	AlignVertical(lipgloss.Center)

var footerStyle = lipgloss.NewStyle().
	Margin(0).
	Padding(0).
	Align(lipgloss.Left)

var fullscreenStyle = lipgloss.NewStyle().
	Padding(0).
	Margin(0).
	Align(lipgloss.Center)
