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
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dx",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	Version: "0.1.1",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dx.yaml)")

	log.SetFlags(0)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".dx" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".dx")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func commandExists(data []map[string]interface{}, full string) bool {
	for _, v := range data {
		if v["full"] == full {
			return true
		}
	}
	return false
}

func loadConfig() Config {
	settings := viper.AllSettings()

	cli := ""
	if settings["cli"] != nil {
		cli = settings["cli"].(string)
	}

	company := ""
	if settings["company"] != nil {
		company = settings["company"].(string)
	}

	outdir := ""
	if settings["outdir"] != nil {
		outdir = settings["outdir"].(string)
	}

	return Config{cli, company, outdir}
}

func loadOutConfig(config Config) OutConfig {
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

func prompt(message string, noAnswer bool) bool {
	fmt.Print(message, " ")

	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)

	if text == "" {
		return noAnswer
	}

	return strings.ToLower(text) == "y"
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
