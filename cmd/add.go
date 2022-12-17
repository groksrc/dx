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
	"log"
	"path/filepath"
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
	cliCmd.AddCommand(addCmd)

	addCmd.Flags().StringVarP(&Body, "body", "b", "", "The body of the command")
	addCmd.MarkFlagRequired("body")

	addCmd.Flags().StringVarP(&Parent, "parent", "p", "", "The parent command, if any")
}

type OutConfig struct {
	Cli      string                   `yaml:"cli"`
	Company  string                   `yaml:"company"`
	Commands []map[string]interface{} `yaml:"commands"`
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

	config := loadConfig()
	outConfig := loadOutConfig(config)

	validateAdd(outConfig.Commands, data)

	outCmdMap := outCommandToMap(data)
	outConfig.Commands = append(outConfig.Commands, outCmdMap)

	save(outConfig, config)
}

func outCommandToMap(data OutCommand) map[string]interface{} {
	outCmd := make(map[string]interface{})
	outCmd["name"] = data.Name
	outCmd["body"] = data.Body
	outCmd["full"] = data.Full
	return outCmd
}

func outConfigFile(config Config) string {
	return filepath.Join(config.outdir, config.cli, fmt.Sprintf(".%s.yaml", config.cli))
}

func validateAdd(commands []map[string]interface{}, data OutCommand) {
	if commandExists(commands, data.Full) {
		// TODO: ask if they want to overwrite
		log.Fatal(fmt.Sprintf("A command named '%s' already exists", data.Full))
	}

	if Parent != "" && !commandExists(commands, Parent) {
		log.Fatal(fmt.Sprintf("A parent command named '%s' was not found", Parent))
	}
}
