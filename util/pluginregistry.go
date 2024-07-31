package util

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

func GetPluginRegistryFilePath(pluginName string) string {
	return filepath.Join(GetLocalRegistryRepoPath(), "registry", fmt.Sprintf("%s.yaml", pluginName))
}

func ReadPluginRegistryFile(pluginName string) (*PluginConfig, error) {
	yamlFile, err := os.ReadFile(GetPluginRegistryFilePath(pluginName))
	if err != nil {
		return nil, err
	}

	v := PluginConfig{}
	if err := yaml.Unmarshal(yamlFile, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

func AddPluginToRegistry(plugin *PluginConfig) error {
	config, err := ReadManifestFile()
	if err != nil {
		return err
	}

	for _, p := range config.Plugins {
		if p.Name == plugin.Name {
			return fmt.Errorf("plugin %s already exists in the registry", plugin.Name)
		}
	}

	config.Plugins = append(config.Plugins, *plugin)

	err = WriteManifestFile(config)
	if err != nil {
		return err
	}

	return nil
}

func RemovePluginFromRegistry(pluginName string) error {
	config, err := ReadManifestFile()
	if err != nil {
		return err
	}

	filtered := []PluginConfig{}

	for _, p := range config.Plugins {
		if p.Name != pluginName {
			filtered = append(filtered, p)
		}
	}

	config.Plugins = filtered

	err = WriteManifestFile(config)
	if err != nil {
		return err
	}

	return nil
}
