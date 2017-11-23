package command

import (
	"os/user"

	"path/filepath"

	"io"

	"github.com/urfave/cli"
)

const sshAliasesDir = ".ssh_aliases"

// CLI stands for Command Line Interface
// CLI interprets user input and executes commands
type CLI struct {
	app *cli.App
}

// NewCLI creates new CLI instance
// provided version will be printed with --version
// CLI will write output to provided writer
func NewCLI(version string, writer io.Writer) (*CLI, error) {
	app, err := configureCLI(version, writer)
	if err != nil {
		return nil, err
	}
	return &CLI{
		app: app,
	}, nil
}

func configureCLI(version string, writer io.Writer) (*cli.App, error) {
	homeDir, err := homeDir()
	if err != nil {
		return nil, err
	}
	var scanDir string
	var save bool
	var force bool
	var file string

	app := cli.NewApp()
	app.Version = version
	app.Name = "ssh-aliases"
	app.Usage = "template driven ssh config generation"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "scan, s",
			Usage:       "input files dir",
			Value:       filepath.Join(homeDir, sshAliasesDir),
			Destination: &scanDir,
		},
	}
	app.Commands = []cli.Command{{
		Name:    "list",
		Aliases: []string{"l"},
		Usage:   "Prints the list of host definitions",
		Action: func(ctx *cli.Context) error {
			err := newListCommand(writer).execute(scanDir)
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}
			return nil
		},
	}, {
		Name:    "compile",
		Aliases: []string{"c"},
		Usage:   "Prints compiled ssh config file or writes it to a file",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:        "save",
				Usage:       "write compilation output to file instead of printing to stdout",
				Destination: &save,
			},
			cli.StringFlag{
				Name:        "file",
				Usage:       "destination file path",
				Destination: &file,
				Value:       filepath.Join(homeDir, ".ssh", "config"),
			},
			cli.BoolFlag{
				Name:        "force",
				Usage:       "overwrite existing file without confirmation",
				Destination: &force,
			},
		},
		Action: func(ctx *cli.Context) error {
			var err error
			if save {
				err = newCompileSaveCommand(file).execute(scanDir, force)
			} else {
				err = newCompileCommand(writer).execute(scanDir)
			}
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}
			return nil
		},
	}}
	return app, nil
}

func homeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.HomeDir, nil
}

// ApplyArgs runs the CLI against provided args
func (c *CLI) ApplyArgs(args []string) error {
	return c.app.Run(args)
}
