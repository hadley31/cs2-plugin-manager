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

func ReadYamlFile(yamlFilePath string) (*PluginRegistry, error) {
	yamlFile, err := os.ReadFile(yamlFilePath)
	if err != nil {
		return nil, err
	}

	v := PluginRegistry{}
	if err := yaml.Unmarshal(yamlFile, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

func DownloadPlugin(url string) (string, error) {
	// Download the plugin from the URL
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error downloading plugin:", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error downloading plugin from %s. Status code: %s", url, resp.Status)
	}

	// Create the file
	out, err := os.CreateTemp("", "cs2pm-plugin-")
	if err != nil {
		fmt.Printf("err: %s", err)
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Printf("err: %s", err)
	}

	return out.Name(), nil
}

func UnzipPlugin(tempFilePath string, dest string) {
	// Unzip the plugin
	archive, err := zip.OpenReader(tempFilePath)
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
