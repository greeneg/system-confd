package globals

import "errors"

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
	Enabled     bool   `json:"enabled"`
	Description string `json:"description"`
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
