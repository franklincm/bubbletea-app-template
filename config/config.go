package config

import (
	"errors"
	"io/fs"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/imdario/mergo"
)

type (
	Config struct {
		Keys   map[string]Keymap `toml:"keys"`
		Styles map[string]string `toml:"styles"`
	}

	Keymap map[string]string
)

var defaultConfig = Config{
	Keys: map[string]Keymap{
		"global": {
			"left":         "h",
			"down":         "j",
			"up":           "k",
			"right":        "l",
			"halfPageDown": "ctrl+d",
			"halfPageUp":   "ctrl+u",
			"pageDown":     " ",
			"pageUp":       "b",
		},
	},
	Styles: map[string]string{
		"borderForeground":        "#458588",
		"tabBlurredBackground":    "#282828",
		"tabBlurredForeground":    "#928374",
		"tabFocusedBackground":    "#98971a",
		"tabFocusedForeground":    "#ebdbb2",
		"tableHeaderForeground":   "#3c3836",
		"tableSelectedBackground": "#282828",
		"tableSelectedForeground": "#d3869b",
	},
}

func New() Config {
	var config Config

	_, err := toml.DecodeFile("config.toml", &config)
	if err == nil {
		if err := mergo.Merge(defaultConfig, &config); err != nil {
			return config
		}
	}

	_, err = toml.DecodeFile("/etc/cb/config.toml", &config)
	if err == nil {
		if err := mergo.Merge(defaultConfig, &config); err != nil {
			return config
		}
	}

	_, err = toml.DecodeFile(
		os.Getenv("HOME")+"/.cb.toml",
		&config,
	)
	if err == nil {
		if err := mergo.Merge(defaultConfig, &config); err != nil {
			return config
		}
	}

	_, err = toml.DecodeFile(
		os.Getenv("HOME")+"/.config/cb/config.toml",
		&config,
	)
	if err == nil {
		if err := mergo.Merge(defaultConfig, &config); err != nil {
			return config
		}
	}

	_, err = toml.DecodeFile(os.Getenv("CB_CONFIG_PATH"), &config)
	if err == nil {
		if err := mergo.Merge(defaultConfig, &config); err != nil {
			return config
		}
	}

	var pathError *fs.PathError

	if err != nil {
		log.Println("config not found")

		if errors.As(err, &pathError) {
			log.Println("using defaults")
			return defaultConfig
		}
	}

	return config
}
