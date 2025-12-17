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

	/*
		// Plugin routes
		g.GET("/plugins", s.GetPlugins)
		g.POST("/plugins/enable/:name", s.EnablePlugin)
		g.POST("/plugins/disable/:name", s.DisablePlugin)

		// System information routes
		g.GET("/system/info", s.GetSystemInfo)
		g.GET("/system/services", s.GetServices)
	*/
}

func SetupRoutes(g *gin.RouterGroup, s *controllers.SystemConfd, pRegistry globals.PluginRegstry) {
}
