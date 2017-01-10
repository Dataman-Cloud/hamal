package api

import (
	"github.com/Dataman-Cloud/hamal/src/models"
	"github.com/Dataman-Cloud/hamal/src/service"
	"github.com/Dataman-Cloud/hamal/src/utils"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

const (
	ParamError      = "503-10001"
	ProjectExist    = "503-10002"
	ProjectNotExist = "503-10003"
	UpdateError     = "503-10004"
)

type HamalControl struct {
	Service *service.HamalService
}

func InitHamalControl() *HamalControl {
	return &HamalControl{
		Service: service.InitHamalService(),
	}
}

func (hc *HamalControl) Ping(ctx *gin.Context) {
	utils.Ok(ctx, "success")
}

func (hc *HamalControl) CreateProject(ctx *gin.Context) {
	var project models.Project
	if err := ctx.BindJSON(&project); err != nil {
		log.Error("invalid param")
		utils.ErrorResponse(ctx, utils.NewError(ParamError, "invalid param"))
		return
	}

	if err := hc.Service.CreateProject(project); err != nil {
		log.Error(err)
		utils.ErrorResponse(ctx, utils.NewError(ProjectExist, err))
		return
	}
	utils.Create(ctx, "success")
}

func (hc *HamalControl) UpdateProject(ctx *gin.Context) {
	var project models.Project
	if err := ctx.BindJSON(&project); err != nil {
		log.Error("invalid param")
		utils.ErrorResponse(ctx, utils.NewError(ParamError, "invalid param"))
		return
	}

	if err := hc.Service.UpdateProject(project); err != nil {
		log.Error(err)
		utils.ErrorResponse(ctx, utils.NewError(ProjectNotExist, err))
		return
	}
	utils.Update(ctx, "success")
}

func (hc *HamalControl) GetProjects(ctx *gin.Context) {
	projects := hc.Service.GetProjects()
	utils.Ok(ctx, projects)
}

func (hc *HamalControl) DeleteProjects(ctx *gin.Context) {
	if err := hc.Service.DeleteProject(ctx.Param("name")); err != nil {
		log.Error(err)
		utils.ErrorResponse(ctx, utils.NewError(ProjectNotExist, err))
		return
	}
	utils.Delete(ctx, "success")
}

func (hc *HamalControl) GetProject(ctx *gin.Context) {
	project, err := hc.Service.GetProject(ctx.Param("name"))
	if err != nil {
		log.Error(err)
		utils.ErrorResponse(ctx, utils.NewError(ProjectNotExist, err))
		return
	}
	utils.Ok(ctx, project)
}

func (hc *HamalControl) UpdateInAction(ctx *gin.Context) {
	projectName := ctx.Query("project_name")
	if projectName == "" {
		utils.ErrorResponse(ctx, utils.NewError(ParamError, "invalid project_name"))
		return
	}

	appName := ctx.Query("app_name")
	if appName == "" {
		utils.ErrorResponse(ctx, utils.NewError(ParamError, "invalid app_name"))
		return
	}

	stage := ctx.Query("stage")
	if stage == "" {
		utils.ErrorResponse(ctx, utils.NewError(ParamError, "invalid stage"))
		return
	}
	err := hc.Service.UpdateInAction(projectName, appName, stage)
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(ParamError, err))
		return
	}
	utils.Ok(ctx, "success")
}
