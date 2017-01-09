package command

import (
	"fmt"
	"github.com/urfave/cli"
	"io/ioutil"
)

func NewDeployCommand() cli.Command {
	return cli.Command{
		Name:    "deploy",
		Aliases: []string{"d"},
		Usage:   "deploy a project",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "file, f",
				Usage: "Load deploy file from `FILE`",
			},
		},
		Action: func(c *cli.Context) error {
			file := c.String("file")
			if file == "" {
				cli.ShowCommandHelp(c, "deploy")
			} else {
				content, err := ioutil.ReadFile(file)
				if err != nil {
					return cli.NewExitError(fmt.Sprintf("%s", err.Error()), 1)
				}
				DeployProject(content)
			}
			return nil
		},
	}
}

func DeployProject(deployfile []byte) error {
	return nil
}
