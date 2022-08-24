package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/micronull/i3rotonda/internal/pkg/cli/serve"
	"github.com/micronull/i3rotonda/internal/pkg/cli/switcher"
)

type runner interface {
	Init(args []string) error
	Run() error
	Name() string
}

var errInit = errors.New("run init error")

func main() {
	if err := run(os.Args[1:]); err != nil {
		if !errors.Is(err, errInit) {
			fmt.Println(err)
		}

		os.Exit(1)
	}

	os.Exit(0)
}

func run(args []string) error {
	if len(args) < 1 || args[0] == "-h" || args[0] == "--help" {
		return errors.New("usage: <command> [<args>]\n" +
			"   serve - run observer for switcher by history")
	}

	cmds := []runner{
		serve.NewCommand(),
		switcher.NewCommand(),
	}

	subcommand := args[0]

	for _, cmd := range cmds {
		if cmd.Name() == subcommand {
			if err := cmd.Init(args[1:]); err != nil {
				return fmt.Errorf("%w: %s", errInit, err.Error())
			}

			return cmd.Run()
		}
	}

	return fmt.Errorf("unknown subcommand: %s", subcommand)
}
