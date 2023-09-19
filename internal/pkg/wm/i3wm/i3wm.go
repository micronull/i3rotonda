package i3wm

import (
	"fmt"
	"log/slog"

	"go.i3wm.org/i3/v4"

	"github.com/micronull/i3rotonda/internal/pkg/wm"
)

type I3wm struct {
	wsts map[string]bool
}

func NewI3wm() *I3wm {
	return &I3wm{
		wsts: make(map[string]bool, 5),
	}
}

var _ wm.WorkspaceManager = &I3wm{}

func (i *I3wm) Switch(target string) {
	if target == "" {
		return
	}

	cmd := fmt.Sprintf("workspace %s", target)

	if _, err := i3.RunCommand(cmd); err != nil {
		slog.Error("i3wm run command failed", "error", err.Error())
	}
}

func (i *I3wm) GetCurrentWorkspace() wm.Workspace {
	wss, _ := i3.GetWorkspaces()

	for _, w := range wss {
		if w.Focused {
			return &ws{w.Name, i}
		}
	}

	return nil
}

func (i *I3wm) OnChangeWorkspace() <-chan wm.Workspace {
	ch := make(chan wm.Workspace)

	go func() {
		recv := i3.Subscribe(i3.WorkspaceEventType)

		for recv.Next() {
			ev, ok := recv.Event().(*i3.WorkspaceEvent)
			if !ok {
				continue
			}

			name := ev.Current.Name

			if ev.Change == "empty" {
				i.wsts[name] = false
			}

			if ev.Change != "focus" {
				continue
			}

			i.wsts[name] = true

			ch <- &ws{name, i}
		}
	}()

	return ch
}

func (i *I3wm) isEmptyWorkspace(ws wm.Workspace) bool {
	return i.wsts[ws.GetName()]
}

type ws struct {
	name string
	wm   *I3wm
}

var _ wm.Workspace = &ws{}

func (w *ws) GetName() string {
	return w.name
}

func (w *ws) IsEmpty() bool {
	return !w.wm.isEmptyWorkspace(w)
}
