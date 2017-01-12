package command

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	cfg "github.com/Dataman-Cloud/hamal/src/hamalcli/config"
	"github.com/Dataman-Cloud/hamal/src/models"
	ui "github.com/gizak/termui"
	"github.com/urfave/cli"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const (
	CODE_SUCCESS           = 0
	PROJECT_NOT_EXIST      = 10003
	PROJECT_STATUS_SUCCESS = "success"
	PROJECT_STATUS_CREATED = "created"
)

type responseCodeType struct {
	Code int `json:"code"`
}

type responseBodyType struct {
	Code int            `json:"code"`
	Data models.Project `json:"data"`
}

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
		var hamalJSON models.Hamal
		if err = json.Unmarshal(content, &hamalJSON); err != nil {
			return cli.NewExitError(fmt.Sprintf("%s", err.Error()), 1)
		}
		project, err := getProject(hamalJSON.Name)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("%s", err.Error()), 1)
		}
		if project == nil {
			if err = createProject(content); err != nil {
				return cli.NewExitError(fmt.Sprintf("%s", err.Error()), 1)
			}
			// TODO (wtzhou) we can bypass the duplicated getProject call if createProject return the object
			if project, err = getProject(hamalJSON.Name); err != nil {
				return cli.NewExitError(fmt.Sprintf("%s", err.Error()), 1)
			}
		}
		if project.Applications[0].Status == PROJECT_STATUS_SUCCESS {
			fmt.Print("Have been updated to current version")
		} else if confirmRollingUpdate(project) {
			rollingUpdateProject(project)
		}
	}
	return nil
}

func getProject(projectName string) (*models.Project, error) {
	resp, err := http.Get(cfg.GetServerFullUrl() + "/projects/" + projectName)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var respCode responseCodeType
	if err = json.Unmarshal(body, &respCode); err != nil {
		return nil, err
	}
	switch respCode.Code {
	case PROJECT_NOT_EXIST:
		return nil, nil
	case CODE_SUCCESS:
		var respBody responseBodyType
		if err = json.Unmarshal(body, &respBody); err != nil {
			return nil, err
		}
		return &respBody.Data, nil
	default:
		return nil, errors.New("Unknown Error")
	}
}

func createProject(hamalByte []byte) error {
	req, err := http.NewRequest("POST", cfg.GetServerFullUrl()+"/projects", bytes.NewBuffer(hamalByte))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusCreated {
		return nil
	} else {
		return errors.New(string(body))
	}
	return nil
}

func rollingUpdateProject(project *models.Project) error {
	client := &http.Client{}
	// TODO (wtzhou) we can support PER-app-PER-project only now
	for _, app := range project.Applications {
		req, err := http.NewRequest("PUT", cfg.GetServerFullUrl()+"/projects/"+project.Name+"/rollingupdate", strings.NewReader(`{"app_id":"`+app.AppId+`"}`))
		req.Header.Set("Content-Type", "application/json")
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
	}
	return nil
}

func confirmRollingUpdate(project *models.Project) bool {
	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	header := ui.NewPar("Press [q](fg-magenta) to quit, Press [h/<-](fg-magenta), [l/->](fg-magenta) to switch options, Press [enter](fg-magenta) to excute the option")
	header.Height = 2
	header.Width = 50
	header.Border = false
	header.TextBgColor = ui.ColorBlue

	// TODO (wtzhou) only support per-app-per-project
	stagesSum := len(project.Applications[0].RollingUpdatePolicy)
	currentStage := int(project.Applications[0].NextStage)
	stagesArray := make([]string, stagesSum)

	for i := 0; i < int(currentStage); i++ {
		stagesArray[i] = " [" + strconv.Itoa(i) + "] " + "[Updated " + strconv.Itoa(int(project.Applications[0].RollingUpdatePolicy[i].InstancesToUpdate)) + " instances](fg-blue)"
	}
	if stagesSum > currentStage {
		stagesArray[currentStage] = "*[" + strconv.Itoa(currentStage) + "] " + "[Pending update " + strconv.Itoa(int(project.Applications[0].RollingUpdatePolicy[currentStage].InstancesToUpdate)) + " instances](fg-white,bg-green)"
		for i := int(currentStage) + 1; i < stagesSum; i++ {
			stagesArray[i] = " [" + strconv.Itoa(i) + "] " + "[Pending update " + strconv.Itoa(int(project.Applications[0].RollingUpdatePolicy[i].InstancesToUpdate)) + " instances](fg-white)"
		}
	}

	stagesUI := ui.NewList()
	stagesUI.Items = stagesArray
	stagesUI.ItemFgColor = ui.ColorYellow
	stagesUI.BorderLabel = project.Name + " Progress..."
	stagesUI.Height = 10
	stagesUI.Width = 20

	confirmBar := ui.NewPar("Confirm")
	confirmBar.TextFgColor = ui.ColorWhite
	confirmBar.TextBgColor = ui.ColorGreen
	confirmBar.Height = 2
	confirmBar.Width = 5
	confirmBar.Border = false

	cancelBar := ui.NewPar("Cancel")
	cancelBar.Height = 2
	cancelBar.TextFgColor = ui.ColorWhite
	cancelBar.TextBgColor = ui.ColorDefault
	cancelBar.Width = 5
	cancelBar.Border = false

	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(3, 2, header),
		),
		ui.NewRow(
			ui.NewCol(4, 2, stagesUI),
		),
		ui.NewRow(
			ui.NewCol(2, 2, cancelBar),
			ui.NewCol(1, 1, confirmBar),
		),
	)

	ui.Body.Y = 3
	ui.Body.Align()
	ui.Render(ui.Body)

	confirm := true

	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
		confirm = false
	})
	ui.Handle("/sys/kbd/h", func(ui.Event) {
		highlightToggle(cancelBar, confirmBar)
		ui.Render(ui.Body)
		confirm = false
	})
	ui.Handle("/sys/kbd/l", func(ui.Event) {
		highlightToggle(confirmBar, cancelBar)
		ui.Render(ui.Body)
		confirm = true
	})
	ui.Handle("/sys/kbd/<left>", func(ui.Event) {
		highlightToggle(cancelBar, confirmBar)
		ui.Clear()
		ui.Render(ui.Body)
		confirm = false
	})
	ui.Handle("/sys/kbd/<right>", func(ui.Event) {
		highlightToggle(confirmBar, cancelBar)
		ui.Clear()
		ui.Render(ui.Body)
		confirm = true
	})
	ui.Handle("/sys/kbd/<enter>", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Loop()
	return confirm
}

func highlightToggle(parA *ui.Par, parB *ui.Par) {
	parA.TextFgColor = ui.ColorWhite
	parA.TextBgColor = ui.ColorGreen
	parB.TextFgColor = ui.ColorWhite
	parB.TextBgColor = ui.ColorDefault
}
