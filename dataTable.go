package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	table "github.com/franklincm/bubbletea-template/components/table"
)

func NewTable() tea.Model {
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
	return dataTable
}
