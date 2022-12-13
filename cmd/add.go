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
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
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
	Name   string `yaml:"name"`
	Body   string `yaml:"body"`
	Parent string `yaml:"parent"`
}

func addOutCommand(cmd *cobra.Command, args []string) {
	data := OutCommand{
		Name:   args[0],
		Body:   Body,
		Parent: Parent,
	}

	config := loadConfig()
	configMap := loadConfigMap(config)

	validateAdd(configMap.Commands, data)

	outCmdMap := outCommandToMap(data)
	configMap.Commands = append(configMap.Commands, outCmdMap)

	save(configMap, config)
}

func outCommandToMap(data OutCommand) map[string]interface{} {
	outCmd := make(map[string]interface{})
	outCmd["name"] = data.Name
	outCmd["body"] = data.Body
	outCmd["parent"] = data.Parent
	return outCmd
}

func outConfigFile(config Config) string {
	return filepath.Join(config.outdir, config.cli, fmt.Sprintf(".%s.yaml", config.cli))
}

func loadConfigMap(config Config) OutConfig {
	// Open the YAML file
	yamlFile, err := ioutil.ReadFile(outConfigFile(config))
	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal the YAML file into a map
	var outConfig OutConfig
	err = yaml.Unmarshal(yamlFile, &outConfig)
	if err != nil {
		log.Fatal(err)
	}

	return outConfig
}

// Recursive function to search a yaml file for a specific key/value pair
func findOutCommand(data []map[string]interface{}, key string) bool {
	for _, v := range data {
		if v["name"] == key {
			return true
		}
	}
	return false
}

func save(configMap OutConfig, config Config) {
	// Marshal the map back into YAML
	yamlData, err := yaml.Marshal(configMap)
	if err != nil {
		log.Fatal(err)
	}

	// Write the YAML data to the file
	err = ioutil.WriteFile(outConfigFile(config), yamlData, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func validateAdd(commands []map[string]interface{}, data OutCommand) {
	if findOutCommand(commands, data.Name) {
		// TODO: ask if they want to overwrite
		log.Fatal(fmt.Sprintf("A command with the name '%s' already exists", data.Name))
	}

	if Parent != "" && !findOutCommand(commands, data.Parent) {
		log.Fatal(fmt.Sprintf("A parent command with the name '%s' was not found", data.Parent))
	}
}
