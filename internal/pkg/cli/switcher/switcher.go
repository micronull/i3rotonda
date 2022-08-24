package switcher

import (
	"flag"
	"net"

	"github.com/micronull/i3rotonda/internal/pkg/socket"
	"github.com/micronull/i3rotonda/internal/pkg/types"
)

type Command struct {
	fs *flag.FlagSet
	c  net.Conn
}

func NewCommand() *Command {
	return &Command{
		fs: flag.NewFlagSet("switcher", flag.ContinueOnError),
	}
}

func (c *Command) Init(_ []string) (err error) {
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
