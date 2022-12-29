package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type DxConfig struct{ Cli, Company, Description, OutDir string }

func commandExists(data []map[string]interface{}, full string) bool {
	for _, v := range data {
		if v["full"] == full {
			return true
		}
	}
	return false
}

func loadDxConfig() DxConfig {
	settings := viper.AllSettings()

	cli := ""
	if settings["cli"] != nil {
		cli = settings["cli"].(string)
	}

	company := ""
	if settings["company"] != nil {
		company = settings["company"].(string)
	}

	description := ""
	if settings["description"] != nil {
		description = settings["description"].(string)
	}

	outdir := ""
	if settings["outdir"] != nil {
		outdir = settings["outdir"].(string)
	}

	return DxConfig{cli, company, description, outdir}
}

func outConfigFile(config DxConfig) string {
	outDir := outDirPath(config)
	return filepath.Join(outDir, fmt.Sprintf(".%s.yaml", config.Cli))
}

func loadOutConfig(config DxConfig) OutConfig {
	// Open the YAML file
	yamlFile, err := ioutil.ReadFile(outConfigFile(config))
	if err != nil {
		log.Println(err)
		log.Fatal("Does the CLI need to be created?")
	}

	// Unmarshal the YAML file into a map
	var outConfig OutConfig
	err = yaml.Unmarshal(yamlFile, &outConfig)
	if err != nil {
		log.Fatal(err)
	}

	return outConfig
}

func outDirPath(config DxConfig) string {
	path := os.ExpandEnv(config.OutDir)

	if !filepath.IsAbs(path) {
		if path[0] == '~' {
			usr, err := user.Current()
			if err != nil {
				log.Fatal(err)
			}
			path = filepath.Join(usr.HomeDir, path[1:])
		} else {
			workDir, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}
			path = filepath.Join(workDir, path)
		}
	}

	return filepath.Clean(path)
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

func save(configMap OutConfig, config DxConfig) {
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
