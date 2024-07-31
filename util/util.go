package util

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type PluginRegistry struct {
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

func ReadConfigFile() (*PluginRegistry, error) {
	yamlFile, err := os.ReadFile("cs2pm.yaml")
	if err != nil {
		return nil, err
	}

	v := PluginRegistry{}
	if err := yaml.Unmarshal(yamlFile, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

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

func WriteConfigFile(config *PluginRegistry) error {
	yamlFile, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	err = os.WriteFile("cs2pm.yaml", yamlFile, 0644)
	if err != nil {
		return err
	}

	return nil
}

func AddPluginToRegistry(plugin *PluginConfig) error {
	config, err := ReadConfigFile()
	if err != nil {
		return err
	}

	for _, p := range config.Plugins {
		if p.Name == plugin.Name {
			return fmt.Errorf("plugin %s already exists in the registry", plugin.Name)
		}
	}

	config.Plugins = append(config.Plugins, *plugin)

	err = WriteConfigFile(config)
	if err != nil {
		return err
	}

	return nil
}

func RemovePluginFromRegistry(pluginName string) error {
	config, err := ReadConfigFile()
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

	err = WriteConfigFile(config)
	if err != nil {
		return err
	}

	return nil
}

func DownloadFile(url string, out *os.File) (*os.File, error) {
	// Download the plugin from the URL
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error downloading plugin:", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error downloading plugin from %s. Status code: %s", url, resp.Status)
	}

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Printf("err: %s", err)
	}

	return out, nil
}

func UnzipFile(source string, dest string) {
	// Unzip the plugin
	archive, err := zip.OpenReader(source)
	if err != nil {
		panic(err)
	}
	defer archive.Close()

	for _, f := range archive.File {
		filePath := filepath.Join(dest, f.Name)

		if !strings.HasPrefix(filePath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return
		}
		if f.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			panic(err)
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			panic(err)
		}

		fileInArchive, err := f.Open()
		if err != nil {
			panic(err)
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			panic(err)
		}

		dstFile.Close()
		fileInArchive.Close()
	}
}
