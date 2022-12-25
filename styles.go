package main

import "github.com/charmbracelet/lipgloss"

var headerStyle = lipgloss.NewStyle().
	Margin(0).
	Padding(0).
	BorderBottom(true).
	BorderTop(false).
	BorderLeft(false).
	BorderRight(false).
	BorderStyle(lipgloss.NormalBorder())

var bodyStyle = lipgloss.NewStyle().
	Margin(0).
	Padding(0)

var footerStyle = lipgloss.NewStyle().
	Margin(0).
	Padding(0).
	BorderBottom(false).
	BorderTop(true).
	BorderLeft(false).
	BorderRight(false).
	BorderStyle(lipgloss.NormalBorder())

var fullscreenStyle = lipgloss.NewStyle().
	Padding(0).
	Margin(0).
	Align(lipgloss.Center)
