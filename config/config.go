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
		Keys map[string]Keymap `toml:"keys"`
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
