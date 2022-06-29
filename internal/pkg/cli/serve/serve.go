package serve

import (
	"flag"
	"strings"

	"github.com/go-pkgz/lgr"
	"go.i3wm.org/i3"
)

func NewCommand(logger *lgr.Logger) *Command {
	c := &Command{
		fs:  flag.NewFlagSet("serve", flag.ContinueOnError),
		log: logger,
	}

	c.fs.StringVar(&c.exclude, "e", "", "exclude workspaces from observation, names or numbers separated by commas")

	return c
}

type Command struct {
	fs  *flag.FlagSet
	log *lgr.Logger

	exclude string
	packet  [10]*i3.Node
}

func (c *Command) Name() string {
	return c.fs.Name()
}

func (c *Command) Init(args []string) error {
	return c.fs.Parse(args)
}

func (c *Command) Run() error {
	c.log.Logf("INFO observer is running")

	go c.runListenWorkspace()

	ch := make(chan struct{})
	defer close(ch)

	<-ch

	return nil
}

func (c *Command) runListenWorkspace() {
	e := strings.Split(c.exclude, ",")

	recv := i3.Subscribe(i3.WorkspaceEventType)

	for recv.Next() {
		ev := recv.Event().(*i3.WorkspaceEvent)

		ws := ev.Current

		if !check(ev, e) || c.packet[0] != nil && c.packet[0].Name == ws.Name {
			continue
		}

		for i := cap(c.packet) - 1; i > 0; i-- {
			c.packet[i] = c.packet[i-1]
		}

		c.packet[0] = &ws
	}
}

func check(ev *i3.WorkspaceEvent, e []string) bool {
	if ev.Change != "focus" {
		return false
	}

	var isExclude bool

	for _, s := range e {
		if s == ev.Current.Name {
			isExclude = true

			break
		}
	}

	if isExclude {
		return false
	}

	return true
}
