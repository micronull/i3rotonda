package serve

import (
	"flag"
	"io"
	"log"
	"net"
	"strings"
	"sync"

	"go.i3wm.org/i3"

	"github.com/micronull/i3rotonda/internal/pkg/socket"
)

func NewCommand() *Command {
	c := &Command{
		fs: flag.NewFlagSet("serve", flag.ContinueOnError),
		m:  &sync.RWMutex{},
	}

	c.fs.StringVar(&c.exclude, "e", "", "exclude workspaces from observation, names or numbers separated by commas")

	return c
}

type Command struct {
	fs *flag.FlagSet

	exclude string
	packet  [10]*i3.Node

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
	addr := socket.Run(func(c net.Conn) {
		var d []byte

		for {
			_, err := c.Read(d)

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatal("ERROR: read socket: " + err.Error())
			}
		}

		log.Printf("DEBUG: read: %v", d)
	})

	log.Printf("INFO: listing: %s\n", addr)
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

		c.m.Lock()

		for i := cap(c.packet) - 1; i > 0; i-- {
			c.packet[i] = c.packet[i-1]
		}

		c.packet[0] = &ws

		c.m.Unlock()
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
