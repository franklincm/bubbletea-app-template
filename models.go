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

func NewNav() Layout {
	headings := []string{
		"table",
		"two",
		"three",
		"four",
	}

	s1 := spinner.New(spinner.WithSpinner(spinner.Dot))
	s2 := spinner.New(spinner.WithSpinner(spinner.Points))
	s3 := spinner.New(spinner.WithSpinner(spinner.Pulse))

	return Layout{
		headings: headings,

		models: map[string]tea.Model{
			headings[0]: NewTable(),
			headings[1]: s1,
			headings[2]: s2,
			headings[3]: s3,
		},

		tabNameToIndex: map[string]int{
			"table": 0,
			"two":   1,
			"three": 2,
			"four":  3,
		},
	}
}
