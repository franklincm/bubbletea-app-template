package main

import "github.com/charmbracelet/bubbles/key"

var quitKeys = key.NewBinding(
	key.WithKeys("q", "esc", "ctrl+c"),
	key.WithHelp("", "press q to quit"),
)

var cycleKey = key.NewBinding(
	key.WithKeys("c"),
)
