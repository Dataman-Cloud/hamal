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
	Apps = "/v_beta/apps"
)

type HamalService struct {
	SwanHost           string
	Projects           map[string]models.Project
	ProjectExecHistory map[string][]models.ExecHistory
	Client             *http.Client
	PMutex             *sync.Mutex
}

func InitHamalService() *HamalService {
	u, err := url.Parse(config.GetConfig().SwanAddr)
	if err != nil {
		log.Fatalf("invalid swan url: %s", config.GetConfig().SwanAddr)
		return nil
	}
	return &HamalService{
		SwanHost:           u.String(),
		Projects:           make(map[string]models.Project),
		ProjectExecHistory: make(map[string][]models.ExecHistory),
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
			hs.SwanHost+Apps+"/"+app.App.AppID,
			bytes.NewReader(body))
		if err != nil {
			log.Error(err)
			break
		}

		_, err = hs.Client.Do(req)
		if err != nil {
			log.Error(err)
			break
		}
	}

	project.CreateTime = time.Now().Format(time.RFC3339Nano)
	hs.Projects[project.Name] = project
	hs.ProjectExecHistory[project.Name] = append(
		hs.ProjectExecHistory[project.Name],
		models.ExecHistory{
			Time: time.Now().Format(time.RFC3339Nano),
		},
	)
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
		v.UpdateHistory = hs.ProjectExecHistory[v.Name]
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
		project.UpdateHistory = hs.ProjectExecHistory[project.Name]
		return project, errors.New("project is not exist")
	}

	return project, nil
}

func (hs *HamalService) ExecuteUpdate(project_name, app_name, stage string) error {
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

	instance := int64(-1)
	for _, app := range project.Applications {
		if app.App.AppID == app_name {
			instance = app.RollingUpdatePolicy[stageNum].InstancesToUpdate
			break
		}
	}

	if instance == -1 {
		return errors.New("invalid stage")
	}

	req, err := http.NewRequest("PATCH",
		fmt.Sprintf("%s%s/%s/proceed-update", hs.SwanHost, Apps, project_name),
		bytes.NewReader([]byte(fmt.Sprintf("{\"instances\": %d}", instance))))

	if err != nil {
		return err
	}

	_, err = hs.Client.Do(req)
	if err != nil {
		return err
	}

	hs.ProjectExecHistory[project_name] = append(
		hs.ProjectExecHistory[project_name],
		models.ExecHistory{
			Time: time.Now().Format(time.RFC3339Nano),
		},
	)

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
