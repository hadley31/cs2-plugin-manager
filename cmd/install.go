/*
Copyright © 2024 Nicholas Hadley <contact@nicholashadley.dev>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs a command by name",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please provide a plugin to install")
			os.Exit(1)
		}

		plugin := args[0]
		dest := cmd.Flag("dir").Value.String()

		pluginConfig, err := readYamlFile(fmt.Sprintf("%s.yaml", plugin))

		if err != nil {
			fmt.Println("Error reading plugin.yaml file:", err)
			os.Exit(1)
		}

		extractDir := filepath.Join(dest, pluginConfig.ExtractPrefix)

		fmt.Printf("Installing plugin %s\n", pluginConfig.Name)

		tempFile, err := downloadPlugin(pluginConfig.DownloadUrl)

		if err != nil {
			fmt.Println("Error downloading plugin:", err)
			os.Exit(1)
		}

		unzipPlugin(tempFile, extractDir)
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.Flags().StringP("dir", "d", "", "Directory to install the plugin to")
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

func readYamlFile(yamlFilePath string) (*PluginConfig, error) {
	yamlFile, err := os.ReadFile(yamlFilePath)
	if err != nil {
		return nil, err
	}

	v := PluginConfig{}
	if err := yaml.Unmarshal(yamlFile, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

func downloadPlugin(url string) (string, error) {
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

func unzipPlugin(tempFilePath string, dest string) {
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
