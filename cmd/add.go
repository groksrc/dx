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
	"log"
	"strings"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a new command to your CLI",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		addOutCommand(cmd, args)
	},
}

var Parent string
var Body string

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringVarP(&Body, "body", "b", "", "The body of the command")
	addCmd.MarkFlagRequired("body")

	addCmd.Flags().StringVarP(&Parent, "parent", "p", "", "The parent command, if any")
}

type OutConfig struct {
	Cli      string              `yaml:"cli"`
	Company  string              `yaml:"company"`
	Commands []map[string]string `yaml:"commands"`
}

type OutCommand struct {
	Name string `yaml:"name"`
	Body string `yaml:"body"`
	Full string `yaml:"full"`
}

func addOutCommand(cmd *cobra.Command, args []string) {
	full := []string{Parent, args[0]}
	data := OutCommand{
		Name: args[0],
		Body: Body,
		Full: strings.TrimLeft(strings.Join(full, " "), " "),
	}

	config := loadDxConfig()
	outConfig := loadOutConfig(config)

	validateAdd(outConfig.Commands, data)

	outCmdMap := outCommandToMap(data)
	outConfig.Commands = append(outConfig.Commands, outCmdMap)

	save(outConfig, config)
}

func outCommandToMap(data OutCommand) map[string]string {
	return map[string]string{
		"name": data.Name,
		"body": data.Body,
		"full": data.Full,
	}
}

func validateAdd(commands []map[string]string, data OutCommand) {
	if commandExists(commands, data.Full) {
		// TODO: ask if they want to overwrite
		log.Fatalf("A command named '%s' already exists", data.Full)
	}

	if Parent != "" && !commandExists(commands, Parent) {
		log.Fatalf("A parent command named '%s' was not found", Parent)
	}
}
