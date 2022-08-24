package serve

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"

	"github.com/micronull/i3rotonda/internal/pkg/socket"
	"github.com/micronull/i3rotonda/internal/pkg/types"

	"go.i3wm.org/i3/v4"
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

	ws := getCurrentWorkspace()
	if ws.Focused && isExclude(ws.Name, strings.Split(c.exclude, ",")) {
		i3CmdSwitch(c.packet[0].Name)

		return
	}

	if a == types.ACTION_NEXT {
		if c.packet[0] == nil {
			return
		}

		var first *i3.Node
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

		debugPacket(c.packet[:])
		i3CmdSwitch(c.packet[0].Name)
	}
}

func i3CmdSwitch(wsName string) {
	cmd := fmt.Sprintf("workspace %s", wsName)

	log.Printf("INFO: cmd next to: %s", cmd)

	if _, err := i3.RunCommand(cmd); err != nil {
		log.Printf("ERROR: running command: %s", err.Error())
	}
}

func getCurrentWorkspace() i3.Workspace {
	ws, _ := i3.GetWorkspaces()

	for _, w := range ws {
		if w.Focused {
			return w
		}
	}

	return i3.Workspace{}
}

func debugPacket(packet []*i3.Node) {
	var v []string

	for _, node := range packet {
		if node != nil {
			v = append(v, node.Name)
		}
	}

	log.Printf("DEBUG: %v", v)
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

		log.Printf("DEBUG: added packet: %s", ws.Name)
	}
}

func check(ev *i3.WorkspaceEvent, e []string) bool {
	if ev.Change != "focus" {
		return false
	}

	if isExclude(ev.Current.Name, e) {
		return false
	}

	return true
}

func isExclude(wsName string, e []string) bool {
	for _, s := range e {
		if s == wsName {
			return true
		}
	}

	return false
}
