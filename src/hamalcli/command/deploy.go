package command

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Dataman-Cloud/hamal/src/types"
	"github.com/urfave/cli"
	"io/ioutil"
	"net/http"
)

const (
	// TODO move me to configfile
	backend = "127.0.0.1"
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
		Action: DeployAction,
	}
}

func DeployAction(c *cli.Context) error {
	file := c.String("file")
	if file == "" {
		cli.ShowCommandHelp(c, "deploy")
	} else {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("%s", err.Error()), 1)
		}
		if err = deployProject(content); err != nil {
			return cli.NewExitError(fmt.Sprintf("%s", err.Error()), 1)
		}
	}
	return nil
}

func getProject(hamalByte []byte) error {
	var hamalJSON types.Hamal
	if err := json.Unmarshal(hamalByte, &hamalJSON); err != nil {
		return cli.NewExitError(fmt.Sprintf("%s", err.Error()), 1)
	}
	return nil
}

func deployProject(hamalByte []byte) error {
	req, err := http.NewRequest("POST", backend+"/projects", bytes.NewBuffer(hamalByte))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusOK {
		fmt.Print(string(body))
	} else {
		return errors.New(string(body))
	}

	return nil
}
