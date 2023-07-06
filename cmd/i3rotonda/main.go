package main

import (
	"errors"
	"fmt"
	"os"
	"runtime/debug"

	logic "github.com/micronull/i3rotonda/internal/app/switcher"
	"github.com/micronull/i3rotonda/internal/pkg/socket"
	"github.com/micronull/i3rotonda/internal/pkg/wm/i3wm"

	"github.com/micronull/i3rotonda/internal/pkg/cli/serve"
	"github.com/micronull/i3rotonda/internal/pkg/cli/switcher"
)

type command interface {
	Init(args []string) error
	Run() error
	Name() string
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Println(err)

		os.Exit(1)
	}

	os.Exit(0)
}

func run(args []string) error {
	if len(args) < 1 || args[0] == "-h" || args[0] == "--help" {
		return errors.New("usage: <command> [<args>]\n" +
			"   serve  - run observer for switcher by history\n" +
			"   switch - switch the current workspace")
	}

	var cmd command

	switch sc := args[0]; sc {
	case "switch":
		debug.SetGCPercent(-1) // disable GC.

		conn, err := socket.Connect()
		if err != nil {
			return fmt.Errorf("couldn't connect to socket: %w", err)
		}

		cmd = switcher.NewCommand(conn)
	case "serve":
		cmd = serve.NewCommand(i3wm.NewI3wm(), logic.NewSwitcher(32))
	default:
		return fmt.Errorf("unknown subcommand: %s", sc)
	}

	if err := cmd.Init(args[1:]); err != nil {
		return fmt.Errorf("couldn't initialize: %w", err)
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("couldn't run: %w", err)
	}

	return nil
}
