package router

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"hecate/internal/app/handler"
	"hecate/internal/app/service"
)

func RegisterRoutes(router *gin.Engine, projectService service.ProjectService, log *logrus.Logger) {
	projectHandler := handler.NewProjectHandler(projectService, log)

	// API V1 Group
	apiV1 := router.Group("/api/v1")
	{
		projects := apiV1.Group("/projects")
		{
			projects.POST("", projectHandler.CreateProject)
			// 未来可以添加其他路由: GET, PUT, DELETE等
			projects.GET("", projectHandler.ListProjects)
			projects.GET("/:id", projectHandler.GetProjectByID)
		}
	}
}
