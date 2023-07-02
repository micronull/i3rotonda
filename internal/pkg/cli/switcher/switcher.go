package switcher

import (
	"errors"
	"flag"
	"log"
	"net"
	"time"

	"github.com/micronull/i3rotonda/internal/pkg/types"
)

type Command struct {
	connector connector

	fs *flag.FlagSet
	c  net.Conn
	a  string
	d  time.Duration
}

type connector func() (net.Conn, error)

func NewCommand(connector connector) *Command {
	c := &Command{
		connector: connector,
		fs:        flag.NewFlagSet("switch", flag.ContinueOnError),
	}

	c.fs.StringVar(&c.a, "a", "next", "switch direction, next or prev")
	c.fs.DurationVar(&c.d, "d", time.Millisecond*500, "time after which a switch can be considered to have completed")

	return c
}

func (c *Command) Init(args []string) (err error) {
	if err = c.fs.Parse(args); err != nil {
		return err
	}

	c.c, err = c.connector()

	return err
}

var errWrongAction = errors.New("wrong action")

func (c *Command) Run() error {
	defer func() {
		if err := c.c.Close(); err != nil {
			log.Println("WARNING: couldn't close socket connected")
		}
	}()

	var err error

	switch c.a {
	case "next":
		_, err = c.c.Write([]byte{byte(types.ACTION_NEXT)})
	case "prev":
		_, err = c.c.Write([]byte{byte(types.ACTION_PREV)})
	default:
		err = errWrongAction
	}

	return err
}

func (c *Command) Name() string {
	return c.fs.Name()
}
