package switcher

import (
	"flag"
	"net"
	"time"

	"github.com/micronull/i3rotonda/internal/pkg/socket"
	"github.com/micronull/i3rotonda/internal/pkg/types"
)

type Command struct {
	fs *flag.FlagSet
	c  net.Conn
	a  string
	d  time.Duration
}

func NewCommand() *Command {
	c := &Command{
		fs: flag.NewFlagSet("switch", flag.ContinueOnError),
	}

	c.fs.StringVar(&c.a, "a", "next", "switch direction, next or prev")
	c.fs.DurationVar(&c.d, "d", time.Millisecond * 500, "time after which a switch can be considered to have completed")

	return c
}

func (c *Command) Init(args []string) (err error) {
	if err = c.fs.Parse(args); err != nil {
		return err
	}

	c.c, err = socket.Connect()

	return
}

func (c *Command) Run() (err error) {
	_, err = c.c.Write([]byte{byte(types.ACTION_NEXT)})

	_ = c.c.Close()

	return
}

func (c *Command) Name() string {
	return c.fs.Name()
}
