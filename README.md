# i3wm workspace switcher

## how to use

1. `make install`
2. add in `~/.i3/config`:

```
# workspace back and forth
bindsym $mod+Tab exec "i3rotonda switch -a=prev"
bindsym $mod+Shift+Tab exec "i3rotonda switch -a=next"

exec --no-startup-id i3rotonda serve
```

## Flags

### for `serve` subcommand

* `-e` exclude workspaces from observation, names or numbers separated by commas
* `-d` **TODO** time after which a switch can be considered to have completed