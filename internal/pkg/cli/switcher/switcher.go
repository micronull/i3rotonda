package switcher

//go:generate mockgen -destination ./mocks/$GOFILE -package mocks io WriteCloser

import (
	"errors"
	"flag"
	"io"
	"log/slog"

	"github.com/micronull/i3rotonda/internal/pkg/types"
)

const commandName = "switch"

type Command struct {
	fs     *flag.FlagSet
	writer io.WriteCloser
	action string
}

func NewCommand(wr io.WriteCloser) *Command {
	c := &Command{
		writer: wr,
		fs:     flag.NewFlagSet(commandName, flag.ContinueOnError),
	}

	c.fs.StringVar(&c.action, "a", "next", "switch direction, next or prev")

	return c
}

func (c *Command) Init(args []string) error {
	return c.fs.Parse(args)
}

var errWrongAction = errors.New("wrong action")

func (c *Command) Run() error {
	defer func() {
		if err := c.writer.Close(); err != nil {
			slog.Warn("couldn't close socket connected", "error", err.Error())
		}
	}()

	var err error

	switch c.action {
	case "next":
		_, err = c.writer.Write([]byte{types.ActionNext})
	case "prev":
		_, err = c.writer.Write([]byte{types.ActionPrev})
	default:
		err = errWrongAction
	}

	return err
}

func (c *Command) Name() string {
	return commandName
}
