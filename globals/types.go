package globals

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

type PluginRegstry struct {
	Version int      `json:"version"`
	Plugins []Plugin `json:"plugins"`
}
