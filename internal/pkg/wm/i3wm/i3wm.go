package i3wm

import (
	"fmt"
	"log"

	"go.i3wm.org/i3/v4"
	"golang.org/x/exp/slog"

	"github.com/micronull/i3rotonda/internal/pkg/wm"
)

type I3wm struct {
	lg *log.Logger
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

type ws struct {
	name string
}

func (w *ws) GetName() string {
	return w.name
}

func (i *I3wm) GetCurrentWorkspace() wm.Workspace {
	wss, _ := i3.GetWorkspaces()

	for _, w := range wss {
		if w.Focused {
			return &ws{w.Name}
		}
	}

	return nil
}

func (i *I3wm) OnChangeWorkspace() <-chan wm.Workspace {
	ch := make(chan wm.Workspace)

	go func() {
		recv := i3.Subscribe(i3.WorkspaceEventType)

		for recv.Next() {
			ev := recv.Event().(*i3.WorkspaceEvent)

			if ev.Change != "focus" {
				continue
			}

			ch <- &ws{ev.Current.Name}
		}
	}()

	return ch
}
