package api

import (
	"github.com/Dataman-Cloud/hamal/src/models"
	"github.com/Dataman-Cloud/hamal/src/service"
	"github.com/Dataman-Cloud/hamal/src/utils"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

const (
	ParamError         = "503-10001"
	ProjectExist       = "503-10002"
	ProjectNotExist    = "503-10003"
	UpdateError        = "503-10004"
	GetAppError        = "503-10005"
	GetAppVersionError = "503-10006"
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

func (hc *HamalControl) CreateOrUpdateProject(ctx *gin.Context) {
	var project models.Project
	if err := ctx.BindJSON(&project); err != nil {
		log.Error("invalid param")
		utils.ErrorResponse(ctx, utils.NewError(ParamError, "invalid param"))
		return
	}

	if err := hc.Service.CreateOrUpdateProject(&project); err != nil {
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

	if err := hc.Service.UpdateProject(&project); err != nil {
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

func (hc *HamalControl) RollingUpdate(ctx *gin.Context) {
	projectName := ctx.Param("name")
	var data models.RollPolicy
	if err := ctx.BindJSON(&data); err != nil {
		utils.ErrorResponse(ctx, utils.NewError(ParamError, err))
		return
	}

	appId := data.AppId
	if appId == "" {
		utils.ErrorResponse(ctx, utils.NewError(ParamError, "invalid app_id"))
		return
	}

	err := hc.Service.RollingUpdate(projectName, appId)
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(ParamError, err))
		return
	}
	utils.Ok(ctx, "success")
}

func (hc *HamalControl) GetApp(ctx *gin.Context) {
	app, err := hc.Service.GetApp(ctx.Param("app_id"))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GetAppError, err))
		return
	}

	utils.Ok(ctx, app)
}

func (hc *HamalControl) Rollback(ctx *gin.Context) {
	projectName := ctx.Param("name")
	var data models.RollPolicy
	if err := ctx.BindJSON(&data); err != nil {
		utils.ErrorResponse(ctx, utils.NewError(ParamError, err))
		return
	}

	appId := data.AppId
	if appId == "" {
		utils.ErrorResponse(ctx, utils.NewError(ParamError, "invalid app_id"))
		return
	}

	err := hc.Service.Rollback(projectName, appId)
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(ParamError, err))
		return
	}
	utils.Ok(ctx, "success")
}

func (hc *HamalControl) GetAppVersions(ctx *gin.Context) {
	version, err := hc.Service.GetAppVersions(ctx.Param("app_id"))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GetAppVersionError, err))
		return
	}

	utils.Ok(ctx, version)
}
