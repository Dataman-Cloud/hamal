package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/Dataman-Cloud/hamal/src/config"
	"github.com/Dataman-Cloud/hamal/src/models"
	"github.com/Dataman-Cloud/hamal/src/utils"
	"github.com/Dataman-Cloud/swan/src/types"

	log "github.com/Sirupsen/logrus"
)

const (
	Apps          = "/v_beta/apps"
	ProceedUpdate = "/proceed-update"
)

const (
	DeploySuccess = iota + 1
	DeployIng
)

type HamalService struct {
	SwanHost     string
	Projects     map[string]models.Project
	CurrentStage map[string]int64
	Client       *http.Client
	PMutex       *sync.Mutex
}

func InitHamalService() *HamalService {
	u, err := url.Parse(config.GetConfig().SwanAddr)
	if err != nil {
		log.Fatalf("invalid swan url: %s", config.GetConfig().SwanAddr)
		return nil
	}
	return &HamalService{
		SwanHost:     u.String(),
		Projects:     make(map[string]models.Project),
		CurrentStage: make(map[string]int64),
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
		PMutex: new(sync.Mutex),
	}
}

func (hs *HamalService) CreateProject(project models.Project) error {
	hs.PMutex.Lock()
	defer hs.PMutex.Unlock()
	if _, ok := hs.Projects[project.Name]; ok {
		return errors.New("project is exist")
	}

	for _, app := range project.Applications {
		body, _ := json.Marshal(app.App)
		req, err := http.NewRequest("PUT",
			hs.SwanHost+Apps+"/"+app.AppId,
			bytes.NewReader(body))
		req.Header.Add("Content-Type", "application/json")
		if err != nil {
			log.Error(err)
			continue
		}

		resp, err := hs.Client.Do(req)
		if err != nil {
			log.Error(err)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			data, _ := utils.ReadResponseBody(resp)
			log.Errorf("%s", data)
			continue
		}
	}

	project.CreateTime = time.Now().Format(time.RFC3339Nano)
	hs.Projects[project.Name] = project
	return nil
}

func (hs *HamalService) UpdateProject(project models.Project) error {
	hs.PMutex.Lock()
	defer hs.PMutex.Unlock()
	if _, ok := hs.Projects[project.Name]; !ok {
		return errors.New("project is not exist")
	}

	project.CreateTime = time.Now().Format(time.RFC3339Nano)
	hs.Projects[project.Name] = project
	return nil
}

func (hs *HamalService) GetProjects() []models.Project {
	hs.PMutex.Lock()
	defer hs.PMutex.Unlock()
	var projects []models.Project
	for _, v := range hs.Projects {
		projects = append(projects, v)
	}
	return projects
}

func (hs *HamalService) DeleteProject(name string) error {
	hs.PMutex.Lock()
	defer hs.PMutex.Unlock()
	if _, ok := hs.Projects[name]; !ok {
		return errors.New("project is not exist")
	}

	delete(hs.Projects, name)
	return nil
}

func (hs *HamalService) GetProject(name string) (models.Project, error) {
	hs.PMutex.Lock()
	defer hs.PMutex.Unlock()
	project, ok := hs.Projects[name]
	if !ok {
		return project, errors.New("project is not exist")
	}

	hs.GetProjectDeployStatus(&project)
	return project, nil
}

func (hs *HamalService) GetProjectDeployStatus(project *models.Project) {
	for n, application := range project.Applications {
		status, stage := hs.GetAppDeployStatus(project.Name, application)
		project.Applications[n].CurrentStage = stage
		project.Applications[n].Status = status
	}

}

func (hs *HamalService) GetAppDeployStatus(projectName string, application models.AppUpdateStage) (string, int64) {
	app, err := hs.GetApp(application.AppId)
	if err != nil {
		return "not_found", 0
	}

	var appCurrentVersion int64
	for _, task := range app.Tasks {
		if task.VersionID == app.ProposedVersion.ID {
			appCurrentVersion += 1
		}
	}

	var stageCount int64
	for stageNum, rp := range application.RollingUpdatePolicy {
		stageCount += rp.InstancesToUpdate
		if appCurrentVersion <= stageCount+1 {
			return app.State, int64(stageNum)
		}
	}

	return "unknown", 0
}

func (hs *HamalService) UpdateInAction(project_name, app_name, stage string) error {
	hs.PMutex.Lock()
	defer hs.PMutex.Unlock()

	stageNum, err := strconv.Atoi(stage)
	if err != nil {
		return err
	}

	project, ok := hs.Projects[project_name]
	if !ok {
		return errors.New("project " + project_name + " not exist")
	}

	instance := int64(0)
	for _, app := range project.Applications {
		if stageNum > len(app.RollingUpdatePolicy) {
			continue
		}

		if app.AppId == app_name {
			instance = app.RollingUpdatePolicy[stageNum].InstancesToUpdate
			break
		}
	}

	if instance == 0 {
		return errors.New("invalid stage")
	}

	req, err := http.NewRequest("PATCH",
		fmt.Sprintf("%s%s/%s%s", hs.SwanHost, Apps, app_name, ProceedUpdate),
		bytes.NewReader([]byte(fmt.Sprintf("{\"instances\": %d}", instance))))

	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := hs.Client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		data, _ := utils.ReadResponseBody(resp)
		return errors.New(string(data))
	}

	return nil
}

func (hs *HamalService) GetApp(id string) (types.App, error) {
	var app types.App
	resp, err := hs.Client.Get(hs.SwanHost + Apps + "/" + id)
	if err != nil {
		return app, err
	}
	data, err := utils.ReadResponseBody(resp)
	if err != nil {
		return app, err
	}
	err = json.Unmarshal(data, &app)
	return app, err
}
