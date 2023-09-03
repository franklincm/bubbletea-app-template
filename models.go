package main

import (
	"fmt"

	key "github.com/charmbracelet/bubbles/key"
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

	dataTable = table.New(
		table.WithColumns(cityColumns),
		table.WithRows(cityRows),
		table.WithFocused(true),
		table.WithStyles(table.Styles{
			Cell:     styles.tableCell,
			Header:   styles.tableHeader,
			Selected: styles.tableSelected,
		}),
		table.WithKeyMap(table.KeyMap{
			LineUp: key.NewBinding(
				key.WithKeys(conf.Keys["global"]["up"]),
				key.WithHelp(fmt.Sprintf("↑/%s", conf.Keys["global"]["up"]), "up"),
			),
			LineDown: key.NewBinding(
				key.WithKeys(conf.Keys["global"]["down"]),
				key.WithHelp(fmt.Sprintf("↓/%s", conf.Keys["global"]["down"]), "down"),
			),
			PageUp: key.NewBinding(
				key.WithKeys(conf.Keys["global"]["pageUp"]),
				key.WithHelp(conf.Keys["global"]["pageUp"], "pageUp"),
			),
			PageDown: key.NewBinding(
				key.WithKeys(conf.Keys["global"]["pageDown"]),
				key.WithHelp(conf.Keys["global"]["pageDown"], "pageDown"),
			),
			HalfPageUp: key.NewBinding(
				key.WithKeys(conf.Keys["global"]["halfPageUp"]),
				key.WithHelp(conf.Keys["global"]["halfPageUp"], "halfPageUp"),
			),
			HalfPageDown: key.NewBinding(
				key.WithKeys(conf.Keys["global"]["halfPageDown"]),
				key.WithHelp(conf.Keys["global"]["halfPageDown"], "halfPageDown"),
			),
			GotoTop: key.NewBinding(
				key.WithKeys("home", "g"),
				key.WithHelp("g/home", "go to start"),
			),
			GotoBottom: key.NewBinding(
				key.WithKeys("end", "G"),
				key.WithHelp("G/end", "go to end"),
			),
		}),
	)

	s1 := spinner.New()
	s1.Spinner = spinner.Dot

	s2 := spinner.New()
	s2.Spinner = spinner.Points

	s3 := spinner.New()
	s3.Spinner = spinner.Pulse

	return Layout{
		headings: headings,

		models: map[string]tea.Model{
			headings[0]: dataTable,
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
