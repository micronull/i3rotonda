package serve

import (
	"flag"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/micronull/i3rotonda/internal/pkg/socket"
	"github.com/micronull/i3rotonda/internal/pkg/types"
	"github.com/micronull/i3rotonda/internal/pkg/wm"
)

func NewCommand(wm wm.WorkspaceManager) *Command {
	c := &Command{
		wm: wm,
		fs: flag.NewFlagSet("serve", flag.ContinueOnError),
		m:  &sync.RWMutex{},
	}

	c.fs.StringVar(&c.exclude, "e", "", "exclude workspaces from observation, names or numbers separated by commas")
	c.fs.DurationVar(&c.d, "d", time.Second, "time after which a switch can be considered to have completed")

	return c
}

type Command struct {
	wm wm.WorkspaceManager
	fs *flag.FlagSet

	exclude string
	d       time.Duration

	packet [10]wm.Workspace

	m *sync.RWMutex
}

func (c *Command) Name() string {
	return c.fs.Name()
}

func (c *Command) Init(args []string) error {
	return c.fs.Parse(args)
}

func (c *Command) Run() error {
	log.Println("INFO: observer is running")

	go c.runServer()
	go c.runListenWorkspace()

	ch := make(chan struct{})
	defer close(ch)

	<-ch

	return nil
}

func (c *Command) runServer() {
	addr := socket.Run(func(conn net.Conn) {
		d := make([]byte, 1)

		for {
			_, err := conn.Read(d)

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatal("ERROR: read socket: " + err.Error())
			}
		}

		_ = conn.Close()

		c.action(types.Action(d[0]))
	})

	log.Printf("INFO: listing: %s\n", addr)
}

func (c *Command) action(a types.Action) {
	if a == types.ACTION_NONE {
		return
	}

	ws := c.wm.GetCurrentWorkspace()

	if ws == nil || ws.GetName() == "" {
		return
	}

	if isExclude(ws.GetName(), strings.Split(c.exclude, ",")) {
		c.wm.Switch(c.packet[0].GetName())

		return
	}

	if a == types.ACTION_NEXT {
		if c.packet[0] == nil {
			return
		}

		var first wm.Workspace
		for i := 0; i < len(c.packet); i++ {
			if i == 0 {
				first = c.packet[i]
			}

			if i+1 < len(c.packet) {
				c.packet[i] = c.packet[i+1]
			}

			if len(c.packet) == i+1 || c.packet[i] == nil {
				c.packet[i] = first

				break
			}
		}

		c.wm.Switch(c.packet[0].GetName())
	}
}

func (c *Command) runListenWorkspace() {
	e := strings.Split(c.exclude, ",")

	for ws := range c.wm.OnChangeWorkspace() {
		if isExclude(ws.GetName(), e) || c.packet[0] != nil && c.packet[0].GetName() == ws.GetName() {
			continue
		}

		c.m.Lock()

		for i := cap(c.packet) - 1; i > 0; i-- {
			c.packet[i] = c.packet[i-1]
		}

		c.packet[0] = ws

		c.m.Unlock()

		log.Printf("DEBUG: added packet: %s", ws.GetName())
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
