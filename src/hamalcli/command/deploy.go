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
	backend = "http://127.0.0.1:5099"
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
		var hamalJSON types.Hamal
		if err = json.Unmarshal(content, &hamalJSON); err != nil {
			return cli.NewExitError(fmt.Sprintf("%s", err.Error()), 1)
		}
		project, err := getProject(hamalJSON.ProjectName)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("%s", err.Error()), 1)
		}
		if project == nil {
			if err = createProject(content); err != nil {
				return cli.NewExitError(fmt.Sprintf("%s", err.Error()), 1)
			}
		}
	}
	return nil
}

func getProject(projectName string) (interface{}, error) {
	type responseBodyType struct {
		Code int         `json:"code"`
		Data interface{} `json:"data"`
	}

	resp, err := http.Get(backend + "/projects?name=" + projectName)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusOK {
		var respBody responseBodyType
		if err = json.Unmarshal(body, &respBody); err != nil {
			return nil, err
		} else {
			if respBody.Code != 0 {
				return respBody.Data, nil
			}
		}
	} else {
		return nil, errors.New(string(body))
	}
	return nil, nil
}

func createProject(hamalByte []byte) error {
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

func updateProject() {
}
