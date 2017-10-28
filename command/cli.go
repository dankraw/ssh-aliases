package command

import (
	"fmt"
	"os"
	"os/user"

	"github.com/urfave/cli"
)

const SSH_ALIASES_DIR = ".ssh-aliases"

type CLI struct {
	version     string
	ListCommand *ListCommand
}

func NewCLI(version string) *CLI {
	return &CLI{
		version:     version,
		ListCommand: NewListCommand(),
	}
}

func (c *CLI) ConfigureCLI() error {
	homeDir, err := c.homeDir()
	if err != nil {
		return err
	}
	var source string
	app := cli.NewApp()
	app.Version = c.version
	app.Name = "ssh-aliases"
	app.Usage = "Template driven ssh config generation"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "dir",
			Usage:       "Input HCL configs directory path",
			Value:       fmt.Sprintf("%s/%s", homeDir, SSH_ALIASES_DIR),
			Destination: &source,
		},
	}
	app.Commands = []cli.Command{{
		Name:    "list",
		Aliases: []string{"l"},
		Usage:   "Prints list of aliases with target hosts",
		Action: func(ctx *cli.Context) error {
			return c.ListCommand.List(source)
		},
	}, {
		Name:    "compile",
		Aliases: []string{"c"},
		Usage:   "Compiles aliases and print ssh config output",
		Action: func(ctx *cli.Context) error {
			fmt.Println("compile!")
			return nil
		},
	}}
	app.Run(os.Args)
	return nil
}

func (c *CLI) homeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.HomeDir, nil
}
