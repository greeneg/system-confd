package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/sys/unix"
	"gopkg.in/ini.v1"

	"github.com/greeneg/system-confd/controllers"
	"github.com/greeneg/system-confd/globals"
	"github.com/greeneg/system-confd/routes"
)

func getConfigDir() (string, error) {
	// get our working directory
	appdir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}

	var configDir string
	// get our config dir
	_, err = os.Stat("/etc/system-confd")
	if os.IsNotExist(err) {
		// if doesn't exist, check if appdir config exists
		_, err = os.Stat(filepath.Join(appdir, "config"))
		if os.IsNotExist(err) {
			return "", err
		} else {
			configDir = filepath.Join(appdir, "config")
		}
	} else if err != nil {
		return "", err
	} else {
		// if exists, use it
		configDir = "/etc/system-confd"
	}

	log.Println("INFO Using configuration directory:", configDir)
	return configDir, nil
}

func getConfig(configDir string) (globals.Config, error) {
	configFile := filepath.Join(configDir, "config.ini")
	config := globals.Config{}

	// Load the config file
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return globals.Config{}, err
	}

	if err := ini.MapTo(&config, configFile); err != nil {
		return globals.Config{}, err
	}

	log.Println("INFO Loaded configuration from:", configFile)
	return config, nil
}

func getPluginRegistry(pluginRegistryFile string, configDir string, logger Logger) (globals.PluginRegstry, error) {
	registry := globals.PluginRegstry{}

	// Load the plugin registry file
	var registryFile string

	// Validate and clean the plugin registry file path to prevent directory traversal
	cleanPath := filepath.Clean(pluginRegistryFile)
	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return globals.PluginRegstry{}, err
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		logger.Warn("Plugin registry file not found. File=" + absPath)
		registryFile = filepath.Join(configDir, "plugins.json")
		logger.Info("Attempting to use plugin registry from source tree path: " + registryFile)

		// Validate fallback path as well
		cleanRegistryFile := filepath.Clean(registryFile)
		absRegistryFile, err := filepath.Abs(cleanRegistryFile)
		if err != nil {
			return globals.PluginRegstry{}, err
		}

		if _, err := os.Stat(absRegistryFile); os.IsNotExist(err) {
			return globals.PluginRegstry{}, err
		}
		pluginRegistryFile = absRegistryFile
	} else {
		pluginRegistryFile = absPath
	}

	logger.Info("Loading plugin registry from: " + pluginRegistryFile)
	jsonContent, err := os.ReadFile(pluginRegistryFile)
	if err != nil {
		return globals.PluginRegstry{}, err
	}
	if err := json.Unmarshal([]byte(jsonContent), &registry); err != nil {
		return globals.PluginRegstry{}, err
	}

	logger.Info("Loaded plugin registry from: " + pluginRegistryFile + " with " + strconv.Itoa(len(registry.Plugins)) + " plugins")
	return registry, nil
}

func checkForPluginDir(pluginDir string) error {
	// if the global plugin directory exists, use it. Otherwise, use qualified path
	// of the individual plugin
	if _, err := os.Stat(pluginDir); os.IsNotExist(err) {
		// test for the plugin path from the registry
		for _, plugin := range pluginRegstry.Plugins {
			if plugin.Enabled {
				if _, err := os.Stat(plugin.Path); os.IsNotExist(err) {
					return err
				}
			}
		}
	}

	return nil
}

func getPluginDir(plugin globals.Plugin) string {
	if plugin.Path != "" {
		return filepath.Dir(plugin.Path)
	}
	if config.General.PluginDir != "" {
		return config.General.PluginDir
	}
	return "/usr/lib/system-confd/plugins"
}

// global variables
var config globals.Config
var configDir string
var pluginRegstry globals.PluginRegstry

func main() {
	// are we running as root?
	if os.Geteuid() != 0 {
		panic("This application must be run as root.")
	}

	r := gin.Default()
	err := r.SetTrustedProxies(nil)
	if err != nil {
		log.Fatal("Error setting trusted proxies: " + err.Error())
	}

	configDir, err := getConfigDir()
	if err != nil {
		log.Fatal("Error getting config directory: " + err.Error())
	}

	config, err := getConfig(configDir)
	if err != nil {
		panic("Error loading config: " + err.Error())
	}

	// set up our logger
	logger, err := setupLogger(config)
	if err != nil {
		panic("Error setting up logger: " + err.Error())
	}

	pluginRegstry, err := getPluginRegistry(config.General.PluginRegistryFile, configDir, logger)
	if err != nil {
		panic("Error loading plugin registry: " + err.Error())
	}
	// for ease, pull out the names of the plugins into a slice
	pluginNames := make([]string, len(pluginRegstry.Plugins))
	for i, plugin := range pluginRegstry.Plugins {
		logger.Info("Plugin found: name=" + plugin.Name + ", enabled=" + strconv.FormatBool(plugin.Enabled))
		logger.Info("Plugin path: " + plugin.Path)
		pluginNames[i] = plugin.Name
	}

	// check if the plugin directory exists
	err = checkForPluginDir(config.General.PluginDir)
	if err != nil {
		logger.Error("Plugin directory(s) not found: " + err.Error())
	}

	// if the plugin is enabled, ensure it exists
	for _, plugin := range pluginRegstry.Plugins {
		if plugin.Enabled {
			pluginDir := getPluginDir(plugin)
			pluginPath := filepath.Join(pluginDir, plugin.Name)
			if _, err := os.Stat(pluginDir); os.IsNotExist(err) {
				logger.Error("Plugin directory not found: " + pluginDir)
			} else {
				logger.Info("Plugin directory exists: " + pluginDir)
			}
			if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
				logger.Error("Plugin not found: " + pluginPath)
			} else {
				logger.Info("Plugin exists: " + pluginPath)
			}
		}
	}

	SystemConfd := new(controllers.SystemConfd)

	// set up the non-plugin routes
	staticGroup := r.Group("/api/v1")
	routes.Static(staticGroup, SystemConfd)

	// set up the plugin routes
	pluginGroup := r.Group("/api/v1/system/plugins")
	routes.SetupRoutes(pluginGroup, SystemConfd, pluginRegstry)

	// set up our socket
	socketPath := config.General.SocketPath
	if socketPath == "" {
		socketPath = "/var/run/system-confd/system-confd.sock" // default socket path
	}
	defer func() {
		if err := os.Remove(socketPath); err != nil && !os.IsNotExist(err) {
			logger.Error("Error removing socket file: " + socketPath + " " + string(err.Error()))
		}
	}()

	if err := os.MkdirAll(filepath.Dir(socketPath), 0750); err != nil {
		logger.Error("Error creating socket directory: " + err.Error())
		return
	}

	unix.Umask(0077) // set umask to 0077 to restrict permissions
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		logger.Error("Error creating socket listener: " + err.Error())
		return
	}
	// Cleanup the sockfile.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		err = os.Remove(socketPath)
		if err != nil && !os.IsNotExist(err) {
			if os.IsNotExist(err) {
				logger.Info("Socket does not exist. Ignoring")
				return
			}
			logger.Error("Error removing socket file: " + socketPath + " " + err.Error())
		}
		os.Exit(1)
	}()

	// Create server with timeouts to prevent slowloris attacks
	server := &http.Server{
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	logger.Info("Server starting on socket: " + socketPath)
	if err := server.Serve(listener); err != nil {
		logger.Error("Error starting HTTP server: " + err.Error())
		return
	}
}
