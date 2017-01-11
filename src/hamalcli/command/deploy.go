package command

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Dataman-Cloud/hamal/src/models"
	ui "github.com/gizak/termui"
	"github.com/urfave/cli"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const (
	// TODO move me to configfile
	BACKEND        = "http://192.168.1.51:5099/v1/hamal"
	PROJECTEXISTED = 10002
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
		if err = createProject(content); err != nil {
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
		if confirmRollingUpdate(project) {
			rollingUpdateProject(project)
		}
	}
	return nil
}

func getProject(projectName string) (*models.Project, error) {
	resp, err := http.Get(BACKEND + "/projects/" + projectName)
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
			if respBody.Code == 0 {
				return &respBody.Data, nil
			}
		}
	} else {
		return nil, errors.New(string(body))
	}
	return nil, nil
}

func createProject(hamalByte []byte) error {
	req, err := http.NewRequest("POST", BACKEND+"/projects", bytes.NewBuffer(hamalByte))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusCreated {
		fmt.Print(string(body))
	} else {
		var respCode responseCodeType
		if err = json.Unmarshal(body, &respCode); err != nil {
			return err
		}
		if respCode.Code == PROJECTEXISTED {
			return nil
		}
		return errors.New(string(body))
	}

	return nil
}

func rollingUpdateProject(project *models.Project) error {
	client := &http.Client{}
	// TODO (wtzhou) we can support PER-app-PER-project only now
	for _, app := range project.Applications {
		req, err := http.NewRequest("PUT", BACKEND+"/projects/"+project.Name+"/rollingupdate", strings.NewReader(`{"app_id":"`+app.AppId+`"}`))
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

	header := ui.NewPar("Press q to quit, Press h/<- , l/-> to switch options, Press enter to excute the option")
	header.Height = 2
	header.Width = 50
	header.Border = false
	header.TextBgColor = ui.ColorBlue

	// TODO (wtzhou) only support per-app-per-project
	stagesSum := len(project.Applications[0].RollingUpdatePolicy)
	currentStage := int(project.Applications[0].CurrentStage)
	stagesArray := make([]string, stagesSum)

	for i := 0; i < int(currentStage); i++ {
		stagesArray[i] = " [" + strconv.Itoa(i) + "] " + "[Updated " + strconv.Itoa(int(project.Applications[0].RollingUpdatePolicy[i].InstancesToUpdate)) + " instances](fg-blue)"
	}
	stagesArray[currentStage] = "*[" + strconv.Itoa(currentStage) + "] " + "[Pending update " + strconv.Itoa(int(project.Applications[0].RollingUpdatePolicy[currentStage].InstancesToUpdate)) + " instances](fg-white,bg-green)"
	for i := int(currentStage) + 1; i < stagesSum; i++ {
		stagesArray[i] = " [" + strconv.Itoa(i) + "] " + "[Pending update " + strconv.Itoa(int(project.Applications[0].RollingUpdatePolicy[i].InstancesToUpdate)) + " instances](fg-gray)"
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
