package util

import (
	"os"

	"gopkg.in/yaml.v2"
)

const manifestFileName = "cs2pm.yaml"

type PluginManifestConfig struct {
	Plugins []PluginConfig
}

type PluginConfig struct {
	Name          string
	Description   string
	DownloadUrl   string `yaml:"downloadUrl"`
	ExtractPrefix string `yaml:"extractPrefix"`
	Uninstall     struct {
		Files       []string
		Directories []string
	}
}

func ReadManifestFile() (*PluginManifestConfig, error) {
	yamlFile, err := os.ReadFile(manifestFileName)
	if err != nil {
		return nil, err
	}

	v := PluginManifestConfig{}
	if err := yaml.Unmarshal(yamlFile, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

func WriteManifestFile(config *PluginManifestConfig) error {
	yamlBytes, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	err = os.WriteFile(manifestFileName, yamlBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}
