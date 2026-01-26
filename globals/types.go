package globals

import (
	"encoding/json"
	"errors"
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/greeneg/system-confd/simplelogger"
)

type General struct {
	SocketPath         string `ini:"socket_path"`
	LogLevel           string `ini:"log_level"`
	LogFile            string `ini:"log_file"`
	Debug              bool   `ini:"debug"`
	DebugLogFile       string `ini:"debug_log_file"`
	PluginDir          string `ini:"plugin_dir"`
	PidFile            string `ini:"pid_file"`
	PluginRegistryFile string `ini:"plugin_registry_file"`
}

type Config struct {
	General General `ini:"General"`
}

type Plugin struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	MetaData    string `json:"metaData"`
	Enabled     bool   `json:"enabled"`
	Description string `json:"description"`
}

type PluginMeta struct {
	Name          string `json:"name"`
	APIMountPoint string `json:"apiMountPoint"`
	APIName       string `json:"apiName"`
	Description   string `json:"description"`
	Author        string `json:"author"`
	License       string `json:"license"`
	ApiPaths      any    `json:"apiPaths"`
}

type PluginRegistry struct {
	Version int      `json:"version"`
	Plugins []Plugin `json:"plugins"`
}

func (r *PluginRegistry) ValidatePluginRegistryVersion(fileName string) (bool, error) {
	if r.Version != 1 {
		return false, errors.New("invalid registry format version")
	}
	return true, nil
}

func (p *Plugin) ValidatePlugin() error {
	if p.Name == "" {
		return errors.New("plugin name is required")
	}
	if p.Path == "" {
		return errors.New("plugin path is required")
	}
	return nil
}

func (r *PluginRegistry) ValidatePlugins() error {
	for _, plugin := range r.Plugins {
		if err := plugin.ValidatePlugin(); err != nil {
			return err
		}
	}
	return nil
}

func (p *Plugin) IsEnabled() bool {
	return p.Enabled
}

func (p *PluginMeta) ValidateApiMountPoint(mountPoint string) (bool, error) {
	if mountPoint == "" {
		return false, errors.New("API mount point is required")
	}
	if mountPoint[0] != '/' {
		return false, errors.New("API mount point must start with '/'")
	}
	switch mountPoint {
	case "/hardware", "/network", "/security", "/services", "/software":
		return true, nil
	default:
		return false, errors.New("API mount point is not one of the allowed reserved paths")
	}
}

func ExecutePluginCommand(pluginPath string, input any) (*exec.Cmd, error) {
	cmd := exec.Command(pluginPath)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	go func() {
		defer stdin.Close()
		jsonInput, _ := json.Marshal(input)
		stdin.Write(jsonInput)
	}()

	return cmd, nil
}

func (p *PluginMeta) SetupRoutes(g *gin.RouterGroup, pPlugin *Plugin, logger simplelogger.Logger) {
	_, err := p.ValidateApiMountPoint(p.APIMountPoint)
	if err != nil {
		return
	}

	// now that we have validated the mount point, set up the routes
	routePrefix := p.APIMountPoint + "/" + p.APIName
	logger.Debugf("route prefix: %s", routePrefix)
	// get the route names from paths
	for path, details := range p.ApiPaths.(map[string]any) {
		detailMap := details.(map[string]any)
		method := detailMap["method"].(string)
		switch method {
		case "GET":
			g.GET(routePrefix+path, func(c *gin.Context) {
				requestData := gin.H{}
				logger.Debugf("Plugin %s route %s called", p.Name, routePrefix+path)
				if p.APIName == "/discover" {
					// generate discovery request JSON for the plugin
					requestData = gin.H{"version": 1, "action": "discover"}
				} else if p.APIName == "/readConfig" {
					// generate readConfig request JSON for the plugin
					requestData = gin.H{"version": 1, "action": "readConfig"}
				}
				// send request JSON as stdin to the plugin and get response
				cmd, err := ExecutePluginCommand(pPlugin.Path+"/"+p.Name+".plugin", requestData)
				if err != nil {
					logger.Warn("Error executing plugin command: " + err.Error())
					c.IndentedJSON(500, gin.H{"error": "internal server error"})
					return
				}
				response, err := cmd.Output()
				if err != nil {
					logger.Warn("Error getting plugin output: " + err.Error())
					c.IndentedJSON(500, gin.H{"error": "internal server error"})
					return
				}
				// return the plugin response as JSON
				var jsonResponse any
				err = json.Unmarshal(response, &jsonResponse)
				if err != nil {
					logger.Warn("Error unmarshaling plugin response: " + err.Error())
					c.IndentedJSON(500, gin.H{"error": "internal server error"})
					return
				}
				c.IndentedJSON(200, jsonResponse)
			})
		case "POST":
			g.POST(routePrefix+path, func(c *gin.Context) {

				c.IndentedJSON(200, gin.H{
					"message": "POST " + routePrefix + path + " called",
				})
			})
		}
	}
}
