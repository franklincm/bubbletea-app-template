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
```

## Debug

```
make && DEBUG=1 ./dist/template
```

In a separate terminal:

```
tail -f debug.log
```
