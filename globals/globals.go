package globals

import (
	"encoding/json"
	"os"
)

func LoadJSONFromFile(filePath string, v any) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, v)
	if err != nil {
		return err
	}
	return nil
}

func LoadPluginMeta(pluginPath, metaData string) (PluginMeta, error) {
	pluginMeta := PluginMeta{}
	err := LoadJSONFromFile(metaData, &pluginMeta)
	if err != nil {
		return PluginMeta{}, err
	}
	return pluginMeta, nil
}
