# i3wm workspace switcher

![coverage](https://raw.githubusercontent.com/micronull/i3rotonda/badges/.badges/main/coverage.svg)

## how to use

1. `make install`
2. add in `~/.i3/config`:

```
# workspace back and forth
bindsym $mod+Tab exec "i3rotonda switch -a=prev"
bindsym $mod+Shift+Tab exec "i3rotonda switch -a=next"

exec --no-startup-id i3rotonda serve
```

## Config

Create config file into `/home/user/.config/i3rotonda/config.yml`.

### Example

```yaml
debug: true
workspaces:
  exclude:
    - 1
    - 2
    - 3
```