package command

import (
	"os"
	"os/user"

	"path/filepath"

	"github.com/urfave/cli"
)

const sshAliasesDir = ".ssh_aliases"

type CLI struct {
	version string
}

func NewCLI(version string) *CLI {
	return &CLI{
		version: version,
	}
}

func (c *CLI) ConfigureCLI() error {
	homeDir, err := c.homeDir()
	if err != nil {
		return err
	}
	var scanDir string
	var save bool
	var force bool
	var file string

	app := cli.NewApp()
	app.Version = c.version
	app.Name = "ssh-aliases"
	app.Usage = "Template driven SSH config generation"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "scan",
			Usage:       "Path to scan for input HCL config files",
			Value:       filepath.Join(homeDir, sshAliasesDir),
			Destination: &scanDir,
		},
	}
	app.Commands = []cli.Command{{
		Name:    "list",
		Aliases: []string{"l"},
		Usage:   "Prints to stdout the list of aliases and target hosts",
		Action: func(ctx *cli.Context) error {
			err := NewListCommand(os.Stdout).Execute(scanDir)
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}
			return nil
		},
	}, {
		Name:    "compile",
		Aliases: []string{"c"},
		Usage:   "Prints compiled SSH config file to stdout (default) or writes it to a file",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:        "save",
				Usage:       "Write compilation output to file instead of printing to stdout",
				Destination: &save,
			},
			cli.StringFlag{
				Name:        "file",
				Usage:       "Destination file path",
				Destination: &file,
				Value:       filepath.Join(homeDir, ".ssh", "config"),
			},
			cli.BoolFlag{
				Name:        "force",
				Usage:       "Overwrite existing file without confirmation",
				Destination: &force,
			},
		},
		Action: func(ctx *cli.Context) error {
			var err error
			if save {
				err = NewCompileSaveCommand(file).Execute(scanDir, force)
			} else {
				err = NewCompileCommand(os.Stdout).Execute(scanDir)
			}
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}
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
