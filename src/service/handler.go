package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
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
	DeploySuccess = "success"
	DeployCreated = "created"
	DeployIng     = "updateing"
	Undefined     = "undefined"
)

type HamalService struct {
	SwanHost     string
	Projects     map[string]*models.Project
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
		Projects:     make(map[string]*models.Project),
		CurrentStage: make(map[string]int64),
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
		PMutex: new(sync.Mutex),
	}
}

func (hs *HamalService) CreateOrUpdateProject(project *models.Project) error {
	hs.PMutex.Lock()
	defer hs.PMutex.Unlock()
	if _, ok := hs.Projects[project.Name]; ok {
		return errors.New("project is exist")
	}
	for _, app := range project.Applications {
		as, err := hs.GetApp(app.AppId)
		if err != nil {
			return err
		}
		if as.State != "normal" {
			return errors.New("app state is't normal can't update")
		}
	}

	project.CreateTime = time.Now().Format(time.RFC3339Nano)
	hs.Projects[project.Name] = project
	return nil
}

func (hs *HamalService) UpdateProject(project *models.Project) error {
	hs.PMutex.Lock()
	defer hs.PMutex.Unlock()
	if _, ok := hs.Projects[project.Name]; !ok {
		return errors.New("project " + project.Name + " is not exist")
	}

	project.CreateTime = time.Now().Format(time.RFC3339Nano)
	hs.Projects[project.Name] = project
	return nil
}

func (hs *HamalService) GetProjects() []*models.Project {
	hs.PMutex.Lock()
	defer hs.PMutex.Unlock()
	var projects []*models.Project
	for _, v := range hs.Projects {
		hs.GetProjectDeployStatus(v)
		projects = append(projects, v)
	}
	return projects
}

func (hs *HamalService) DeleteProject(name string) error {
	hs.PMutex.Lock()
	defer hs.PMutex.Unlock()
	if _, ok := hs.Projects[name]; !ok {
		return errors.New("project " + name + " is not exist")
	}

	delete(hs.Projects, name)
	return nil
}

func (hs *HamalService) GetProject(name string) (*models.Project, error) {
	hs.PMutex.Lock()
	defer hs.PMutex.Unlock()
	project, ok := hs.Projects[name]
	if !ok {
		return project, errors.New("project " + name + " is not exist")
	}

	hs.GetProjectDeployStatus(project)
	return project, nil
}

func (hs *HamalService) GetProjectDeployStatus(project *models.Project) {
	for n, application := range project.Applications {
		status, stage := hs.GetAppDeployStatus(project, application)
		project.Applications[n].NextStage = stage
		project.Applications[n].Status = status
	}

}

func (hs *HamalService) GetAppDeployStatus(project *models.Project, application models.AppUpdateStage) (string, int64) {
	app, err := hs.GetApp(application.AppId)
	if err != nil {
		return Undefined, 0
	}

	if app.ProposedVersion == nil {
		if project.Status == 0 {
			return DeployCreated, int64(0)
		}
		return DeploySuccess, int64(0)
	}

	var appCurrentVersion int64
	for _, task := range app.Tasks {
		if app.ProposedVersion != nil && task.VersionID == app.ProposedVersion.ID {
			appCurrentVersion += 1
		}
	}

	var stageCount int64
	for stageNum, rp := range application.RollingUpdatePolicy {
		stageCount += rp.InstancesToUpdate
		/*if appCurrentVersion == 1 {
			return app.State, int64(0)
		} else if appCurrentVersion-1 == stageCount {
			return app.State, int64(stageNum + 1)
		}*/

		if appCurrentVersion == stageCount {
			return app.State, int64(stageNum + 1)
		}
	}

	return Undefined, 0
}

func (hs *HamalService) RollingUpdate(projectName, appName string) error {
	hs.PMutex.Lock()
	defer hs.PMutex.Unlock()

	project, ok := hs.Projects[projectName]
	if !ok {
		return errors.New("project " + projectName + " not exist")
	}

	var application models.AppUpdateStage
	instance := int64(0)
	for _, app := range project.Applications {
		state, stage := hs.GetAppDeployStatus(project, app)
		if app.AppId == appName && int(stage) < len(app.RollingUpdatePolicy) && state != DeploySuccess {
			instance = app.RollingUpdatePolicy[stage].InstancesToUpdate
			application = app
			break
		}
	}

	if instance == 0 {
		return errors.New("invalid stage")
	}

	app, err := hs.GetApp(appName)
	if err != nil {
		return err
	}
	project.Status = 1
	log.Info(hs.SwanHost + Apps + "/" + application.AppId)
	if app.State == "normal" && app.ProposedVersion == nil {
		body, _ := json.Marshal(application.App)
		req, err := http.NewRequest("PUT",
			hs.SwanHost+Apps+"/"+application.AppId,
			bytes.NewReader(body))
		req.Header.Add("Content-Type", "application/json")
		if err != nil {
			log.Error(err)
			return err
		}

		resp, err := hs.Client.Do(req)
		if err != nil {
			log.Error(err)
			return err
		}

		if resp.StatusCode != http.StatusOK {
			data, _ := utils.ReadResponseBody(resp)
			log.Errorf("%s", data)
			return errors.New(string(data))
		}
		return nil
	}

	req, err := http.NewRequest("PATCH",
		fmt.Sprintf("%s%s/%s%s", hs.SwanHost, Apps, appName, ProceedUpdate),
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
	data, _ := utils.ReadResponseBody(resp)
	err = json.Unmarshal(data, &app)
	return app, err
}

func (hs *HamalService) GetAppVersions(appId string) (map[string]types.Version, error) {
	m := make(map[string]types.Version)

	app, err := hs.GetApp(appId)
	if err != nil {
		return m, err
	}

	var newVersionId string
	var oldVersionId string
	if app.ProposedVersion != nil {
		newVersionId = app.ProposedVersion.PreviousVersionID
		oldVersionId = app.ProposedVersion.ID
	} else {
		if app.CurrentVersion.PreviousVersionID != "" {
			oldVersionId = app.CurrentVersion.ID
			newVersionId = app.CurrentVersion.PreviousVersionID
		} else {
			newVersionId = app.CurrentVersion.ID
		}
	}

	resp, err := hs.Client.Get(fmt.Sprintf("%s%s/%s/versions/%s", hs.SwanHost, Apps, appId, newVersionId))
	if err == nil {
		var newVersion types.Version
		data, _ := utils.ReadResponseBody(resp)
		json.Unmarshal(data, &newVersion)
		m["new_version"] = newVersion
	}

	if oldVersionId != "" {
		m["old_version"] = *app.ProposedVersion
	}

	return m, err
}
