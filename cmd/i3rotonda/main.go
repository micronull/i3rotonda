package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/micronull/i3rotonda/internal/pkg/cli/serve"
)

type runner interface {
	Init([]string) error
	Run() error
	Name() string
}

var errInited = errors.New("run init error")

func root(args []string) error {
	if len(args) < 1 || args[0] == "-h" || args[0] == "--help" {
		return errors.New("usage: <command> [<args>]\n" +
			"   serve - run observer for switch by history")
	}

	cmds := []runner{
		serve.NewCommand(),
	}

	subcommand := args[0]

	for _, cmd := range cmds {
		if cmd.Name() == subcommand {
			if err := cmd.Init(args[1:]); err != nil {
				return fmt.Errorf("%w: %s", errInited, err.Error())
			}

			return cmd.Run()
		}
	}

	return fmt.Errorf("unknown subcommand: %s", subcommand)
}

func main() {
	if err := root(os.Args[1:]); err != nil {
		if !errors.Is(err, errInited) {
			fmt.Println(err)
		}

		os.Exit(1)
	}

	os.Exit(0)
}
