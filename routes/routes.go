package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/greeneg/system-confd/controllers"
	"github.com/greeneg/system-confd/globals"
	"github.com/greeneg/system-confd/simplelogger"
)

func Static(g *gin.RouterGroup, s *controllers.SystemConfd) {
	g.GET("/", s.GetRoot)

	// Static routes
	g.GET("/health", s.HealthCheck)
	g.GET("/version", s.Version)
	g.GET("/config", s.GetConfig)
}

func SetupRoutes(g *gin.RouterGroup, s *controllers.SystemConfd, pRegistry globals.PluginRegistry, logger simplelogger.Logger) {
	// Plugin routes
	for _, plugin := range pRegistry.Plugins {
		if plugin.IsEnabled() {
			logger.Info("Setting up routes for plugin: " + plugin.Name)
			// read in the plugin meta data
			pluginMeta, err := globals.LoadPluginMeta(plugin.Path, plugin.MetaData)
			if err != nil {
				continue
			}
			pluginMeta.SetupRoutes(g, &plugin, logger)
		}
	}
}
