package main

import (
	"math"

	tea "github.com/charmbracelet/bubbletea"
	tabs "github.com/franklincm/bubbles/tabs"
	frame "github.com/franklincm/bubbletea-template/components/frame"
	spinner "github.com/franklincm/bubbletea-template/components/spinner"
	table "github.com/franklincm/bubbletea-template/components/table"
)

var dataTable table.Model

type Layout struct {
	models         map[string]tea.Model
	headings       []string
	tabNameToIndex map[string]int
	tabs           tea.Model
	frames         map[frameId]*frame.Model
	activeTab      int
}

func (l *Layout) tabNext() {
	numHeadings := len(l.tabs.(tabs.Model).GetHeadings())
	tabNext := int(math.Min(float64(l.activeTab+1), float64(numHeadings-1)))
	l.activeTab = tabNext
	l.tabs = l.tabs.(tabs.Model).SetFocused(l.activeTab)
}

func (l *Layout) tabPrev() {
	tabPrev := int(math.Max(float64(l.activeTab-1), 0))
	l.activeTab = tabPrev
	l.tabs = l.tabs.(tabs.Model).SetFocused(l.activeTab)
}

func initTabs(headings []string, activeTab int) tea.Model {
	tabs := tabs.New(headings)
	tabs = tabs.FocusedStyle(styles.tabFocusedStyle)
	tabs = tabs.BlurredStyle(styles.tabBlurredStyle)
	tabs = tabs.SetFocused(activeTab)

	return tabs
}

func NewLayout() Layout {
	activeTab := 0
	headings := []string{
		"table",
		"two",
		"three",
		"four",
	}
	models := map[string]tea.Model{
		headings[0]: NewTable(),
		headings[1]: spinner.New(spinner.WithSpinner(spinner.Dot)),
		headings[2]: spinner.New(spinner.WithSpinner(spinner.Points)),
		headings[3]: spinner.New(spinner.WithSpinner(spinner.Pulse)),
	}

	tabs := initTabs(headings, activeTab)

	return Layout{
		activeTab: activeTab,
		headings:  headings,

		models: models,

		tabNameToIndex: map[string]int{
			headings[0]: 0,
			headings[1]: 1,
			headings[2]: 2,
			headings[3]: 3,
		},
		tabs: tabs,
		frames: map[frameId]*frame.Model{
			header: NewHeader(),
			nav: frame.
				New().
				Content(
					[]tea.Model{
						tabs,
					},
				).Style(styles.navStyle),
			body: frame.
				New().
				Style(styles.bodyStyle).
				Content(
					[]tea.Model{models["table"]},
				),
			footer: NewFooter(),
		},
	}
}
