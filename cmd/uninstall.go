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
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/hadley31/cs2pm/util"
	"github.com/spf13/cobra"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstalls plugins from a registry file",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		dest := cmd.Flag("dir").Value.String()
		config, err := util.ReadConfigFile()

		if err != nil {
			panic(err)
		}

		wg := &sync.WaitGroup{}

		for _, config := range config.Plugins {
			wg.Add(1)
			go uninstallPlugin(&config, dest, wg)
		}

		wg.Wait()
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)

	uninstallCmd.Flags().StringP("dir", "d", ".", "Directory to install the plugin to")
}

func uninstallPlugin(config *util.PluginConfig, dest string, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("Uninstalling plugin %s\n", config.Name)

	for _, file := range config.Uninstall.Files {
		filePath := filepath.Join(dest, file)
		fmt.Printf("Removing file %s\n", filePath)
		err := os.Remove(filePath)

		if err != nil {
			fmt.Printf("Error removing file %s: %s\n", filePath, err)
		}
	}

	for _, dir := range config.Uninstall.Directories {
		dirPath := filepath.Join(dest, dir)
		fmt.Printf("Removing directory %s\n", dirPath)
		err := os.RemoveAll(dirPath)

		if err != nil {
			fmt.Printf("Error removing directory %s: %s\n", dirPath, err)
		}
	}
}
