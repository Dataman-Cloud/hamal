package router

import (
	"time"

	"github.com/Dataman-Cloud/hamal/src/api"
	"github.com/Dataman-Cloud/hamal/src/router/middleware"
	"github.com/Dataman-Cloud/hamal/src/utils"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

// Router add router function
func Router(middlewares ...gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(utils.Ginrus(log.StandardLogger(), time.RFC3339Nano, false))
	r.Use(middleware.CORSMiddleware())
	r.Use(middlewares...)

	service := api.InitHamalControl()
	hv1 := r.Group("/v1/hamal")
	{
		hv1.GET("/ping", service.Ping)
		hv1.POST("/projects", service.CreateOrUpdateProject)
		hv1.PUT("/projects", service.UpdateProject)
		hv1.GET("/projects", service.GetProjects)
		//hv1.DELETE("/projects/:name", service.DeleteProjects)
		hv1.GET("/projects/:name", service.GetProject)
		hv1.PUT("/projects/:name/rollingupdate", service.RollingUpdate)
	}

	return r
}
