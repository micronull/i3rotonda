package serve

import (
	"errors"
	"flag"
	"io"
	"log/slog"
	"net"
	"strings"
	"time"

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
	sw switcher
	wm wm.WorkspaceManager
	fs *flag.FlagSet

	exclude string
	d       time.Duration
}

func NewCommand(wm wm.WorkspaceManager, sw switcher) *Command {
	c := &Command{
		wm: wm,
		sw: sw,
		fs: flag.NewFlagSet(commandName, flag.ContinueOnError),
	}

	c.fs.StringVar(&c.exclude, "e", "", "exclude workspaces from observation, names or numbers separated by commas")
	c.fs.DurationVar(&c.d, "d", time.Second, "time after which a switch can be considered to have completed")

	return c
}

func (c *Command) Name() string {
	return commandName
}

func (c *Command) Init(args []string) error {
	return c.fs.Parse(args)
}

func (c *Command) Run() error {
	slog.Info("observer is running")

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

	if cws == nil || cws.GetName() == "" {
		return
	}

	if isExclude(cws.GetName(), strings.Split(c.exclude, ",")) {
		if cr := c.sw.Current(); cr != nil {
			c.wm.Switch(cr.GetName())
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
		c.wm.Switch(ws.GetName())
	}
}

func (c *Command) runListenWorkspace() {
	e := strings.Split(c.exclude, ",")

	for ws := range c.wm.OnChangeWorkspace() {
		if isExclude(ws.GetName(), e) {
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
