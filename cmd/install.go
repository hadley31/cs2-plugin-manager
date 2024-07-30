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

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs plugins from a registry file",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			plugin, err := util.ReadPluginRegistryFile(args[0])
			if err != nil {
				panic(err)
			}

			err = util.AddPluginToRegistry(plugin)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Added %s to plugin manifest\n", plugin.Name)
			return
		}

		dest := cmd.Flag("dir").Value.String()

		config, err := util.ReadConfigFile()

		if err != nil {
			panic(err)
		}

		wg := &sync.WaitGroup{}

		for _, plugin := range config.Plugins {
			wg.Add(1)
			installPlugin(&plugin, dest, wg)
		}

		wg.Wait()
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.Flags().StringP("dir", "d", "", "Directory to install the plugin to")
}

func installPlugin(plugin *util.PluginConfig, dest string, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("Installing plugin %s\n", plugin.Name)

	extractDir := filepath.Join(dest, plugin.ExtractPrefix)

	tempFile, err := os.CreateTemp("", "cs2pm-")
	util.DownloadFile(plugin.DownloadUrl, tempFile)
	defer tempFile.Close()

	if err != nil {
		panic(err)
	}

	util.UnzipFile(tempFile.Name(), extractDir)
}
