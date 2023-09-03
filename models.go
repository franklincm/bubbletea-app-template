package main

import (
	tea "github.com/charmbracelet/bubbletea"
	spinner "github.com/franklincm/bubbletea-template/components/spinner"
	table "github.com/franklincm/bubbletea-template/components/table"
)

var dataTable table.Model

type Layout struct {
	models         map[string]tea.Model
	headings       []string
	tabNameToIndex map[string]int
}

func NewLayout() Layout {
	headings := []string{
		"table",
		"two",
		"three",
		"four",
	}

	return Layout{
		headings: headings,

		models: map[string]tea.Model{
			headings[0]: NewTable(),
			headings[1]: spinner.New(spinner.WithSpinner(spinner.Dot)),
			headings[2]: spinner.New(spinner.WithSpinner(spinner.Points)),
			headings[3]: spinner.New(spinner.WithSpinner(spinner.Pulse)),
		},

		tabNameToIndex: map[string]int{
			headings[0]: 0,
			headings[1]: 1,
			headings[2]: 2,
			headings[3]: 3,
		},
	}
}
