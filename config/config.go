package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/BurntSushi/toml"
)

type (
	Config struct {
		Keys map[string]Keymap `toml:"keys"`
	}

	Keymap map[string]string
)

var defaultConfig = Config{
	Keys: map[string]Keymap{
		"global": {},
	},
}

func New() Config {
	var config Config

	_, err := toml.DecodeFile("config.toml", &config)
	if err == nil {
		return config
	}

	_, err = toml.DecodeFile("/etc/cb/config.toml", &config)
	if err == nil {
		return config
	}

	_, err = toml.DecodeFile(
		os.Getenv("HOME")+"/.cb.toml",
		&config,
	)
	if err == nil {
		return config
	}

	_, err = toml.DecodeFile(
		os.Getenv("HOME")+"/.config/cb/config.toml",
		&config,
	)
	if err == nil {
		return config
	}

	_, err = toml.DecodeFile(os.Getenv("CB_CONFIG_PATH"), &config)
	if err == nil {
		return config
	}

	var pathError *fs.PathError

	if err != nil {
		fmt.Fprintln(os.Stderr, "config not found")

		if errors.As(err, &pathError) {
			fmt.Fprintln(os.Stderr, "using defaults")
			return defaultConfig
		}
	}

	return config
}
