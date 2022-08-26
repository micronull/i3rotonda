package i3wm

import (
	"fmt"
	"log"

	"github.com/micronull/i3rotonda/internal/pkg/wm"
	"go.i3wm.org/i3/v4"
)

type i3wm struct {
	lg *log.Logger
}

var _ wm.WorkspaceManager = &i3wm{}

func New(lgr *log.Logger) *i3wm {
	return &i3wm{lgr}
}

func (i *i3wm) log(msg string, v ...any) {
	if i.lg == nil {
		return
	}

	i.lg.Printf(msg, v...)
}

func (i *i3wm) Switch(target string) {
	cmd := fmt.Sprintf("workspace %s", target)

	i.log("INFO: cmd next to: %s", cmd)

	if _, err := i3.RunCommand(cmd); err != nil {
		i.log("ERROR: running command: %s", err.Error())
	}
}

type ws struct {
	name string
}

func (w *ws) GetName() string {
	return w.name
}

func (i *i3wm) GetCurrentWorkspace() wm.Workspace {
	wss, _ := i3.GetWorkspaces()

	for _, w := range wss {
		if w.Focused {
			return &ws{w.Name}
		}
	}

	return nil
}

func (i *i3wm) OnChangeWorkspace() <-chan wm.Workspace {
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
