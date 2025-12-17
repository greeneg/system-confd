package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/greeneg/system-confd/controllers"
	"github.com/greeneg/system-confd/globals"
)

func Static(g *gin.RouterGroup, s *controllers.SystemConfd) {
	g.GET("/", s.GetRoot)

	// Static routes
	g.GET("/health", s.HealthCheck)
	g.GET("/version", s.Version)
	g.GET("/config", s.GetConfig)
}

func SetupRoutes(g *gin.RouterGroup, s *controllers.SystemConfd, pRegistry globals.PluginRegistry) {
}
