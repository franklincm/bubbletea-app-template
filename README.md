# bubbletea-app-template

## Options

```
-c <command>
```

## Config (TOML)

Config paths:

- `/etc/cb/config.toml`
- `~/.cb.toml`
- `~/.config/cb/config.toml`
- Environment: `CB_CONFIG_PATH`


```toml
[keys.global]
left = "h"
down = "j"
up = "k"
right = "l"

halfPageDown = "ctrl+d"
halfPageUp = "ctrl+u"

pageDown = " "
pageUp = "b"

describe = "d"
edit = "e"
refresh = "r"

[styles]
borderForeground = "#458588"
tabBlurredBackground = "#282828"
tabBlurredForeground = "#928374"
tabFocusedBackground = "#98971a"
tabFocusedForeground = "#ebdbb2"
tableHeaderForeground = "#3c3836"
tableSelectedBackground = "#3c3836"
tableSelectedForeground = "#d3869b"
```

## Debug

```
make && DEBUG=1 ./dist/template
```

In a separate terminal:

```
tail -f debug.log
```
