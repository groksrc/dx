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

func outConfigFile(config Config) string {
	outDir := outDirPath(config)
	return filepath.Join(outDir, fmt.Sprintf(".%s.yaml", config.cli))
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

func outDirPath(config Config) string {
	path := os.ExpandEnv(config.outdir)

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
