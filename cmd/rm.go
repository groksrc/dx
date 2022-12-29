/*
Copyright © 2022 Drew Cain -- @groksrc

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm",
	Args:  cobra.MinimumNArgs(1),
	Short: "Removes a command from your CLI",
	Long: `Removes the command specified from your CLI

If the command is a subcommand pass the entire command name, excluding the name
of your CLI. For example:

$ dx cli rm mycmd mychild mygrandchild  # removes the mygrandchild command

If the command has children, the children will be removed as well.
`,
	Run: func(cmd *cobra.Command, args []string) {
		remove(args[0])
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}

func remove(command string) {
	config := loadDxConfig()
	outConfig := loadOutConfig(config)

	message := fmt.Sprintf("This will remove the command '%s' and any children of this command. Continue? [Y/n]", command)
	if !prompt(message, true) {
		return
	}

	outConfig.Commands = removeChildren(outConfig.Commands, command)

	save(outConfig, config)
}

func removeChildren(data []map[string]interface{}, command string) []map[string]interface{} {
	var commands []map[string]interface{}

	for _, m := range data {
		if !strings.HasPrefix(m["full"].(string), command) {
			commands = append(commands, m)
		}
	}

	return commands
}
