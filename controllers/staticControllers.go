package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/greeneg/system-confd/model"
)

func (s *SystemConfd) GetRoot(c *gin.Context) {
	rootPath, err := model.GetRootPath()
	if err != nil {
		c.IndentedJSON(http.StatusOK, gin.H{"error": string(err.Error())})
	}

	c.IndentedJSON(http.StatusOK, gin.H{"metadata": rootPath})
}

func (s *SystemConfd) HealthCheck(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "System Confd is running",
	})
}

func (s *SystemConfd) Version(c *gin.Context) {
	version := model.GetVersion()

	c.IndentedJSON(http.StatusOK, gin.H{
		"version": version,
	})
}

func (s *SystemConfd) GetConfig(c *gin.Context) {}
