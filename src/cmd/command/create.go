package command

import (
	"github.com/urfave/cli"
)

func NewCreateCommand() cli.Command {
	return cli.Command{
		Name:    "create",
		Aliases: []string{"c"},
		Usage:   "create project",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "config",
				Usage: "config file",
			},
		},
		Action: func(c *cli.Context) error {
			config := c.String("config")
			if config == "" {
				cli.ShowCommandHelp(c, "create")
				return nil
			}
			return nil
		},
	}
}
