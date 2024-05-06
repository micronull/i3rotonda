package serve

import (
	"errors"
	"io"
	"log/slog"
	"net"

	"github.com/micronull/i3rotonda/internal/pkg/config"
	"github.com/micronull/i3rotonda/internal/pkg/socket"
	"github.com/micronull/i3rotonda/internal/pkg/types"
	"github.com/micronull/i3rotonda/internal/pkg/wm"
)

const commandName = "serve"

type switcher interface {
	Add(ws wm.Workspace)
	Current() wm.Workspace
	Prev() wm.Workspace
	Next() wm.Workspace
}

type Command struct {
	sw  switcher
	wm  wm.WorkspaceManager
	cfg config.Config
}

func NewCommand(wm wm.WorkspaceManager, sw switcher) *Command {
	c := &Command{
		wm: wm,
		sw: sw,
	}

	return c
}

func (c *Command) Name() string {
	return commandName
}

func (c *Command) Init([]string) (err error) {
	c.cfg, err = config.Load()

	return err
}

func (c *Command) Run() error {
	if c.cfg.Debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	slog.Debug("observer is running")

	go c.runSocketServer()
	go c.runListenWorkspace()

	ch := make(chan struct{})
	defer close(ch)

	<-ch

	return nil
}

func (c *Command) runSocketServer() {
	addr := socket.Run(func(read net.Conn) {
		d := make([]byte, 1)

		for {
			_, err := read.Read(d)

			if errors.Is(err, io.EOF) {
				break
			}

			if err != nil {
				slog.Error("read error", "error", err.Error())

				break
			}
		}

		_ = read.Close()

		c.action(d[0])
	})

	slog.Info("listing", "addr", addr)
}

func (c *Command) action(a types.Action) {
	if a == types.ActionNone {
		return
	}

	cws := c.wm.GetCurrentWorkspace()

	if cws == nil || cws.Name() == "" {
		return
	}

	if isExclude(cws.Name(), c.cfg.Workspaces.Exclude) {
		if cr := c.sw.Current(); cr != nil {
			c.wm.Switch(cr.Name())
		}

		return
	}

	var ws wm.Workspace

	switch a {
	case types.ActionNext:
		ws = c.sw.Next()
	case types.ActionPrev:
		ws = c.sw.Prev()
	}

	if ws != nil {
		c.wm.Switch(ws.Name())
	}
}

func (c *Command) runListenWorkspace() {
	for ws := range c.wm.OnChangeWorkspace() {
		if isExclude(ws.Name(), c.cfg.Workspaces.Exclude) {
			continue
		}

		c.sw.Add(ws)
	}
}

func isExclude(wsName string, e []string) bool {
	for _, s := range e {
		if s == wsName {
			return true
		}
	}

	return false
}
