/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type OutCommandData struct {
	CommandName,
	CommandVar,
	Parent,
	Body string
}

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:     "generate",
	Short:   "Generates the CLI in the output directory",
	Aliases: []string{"g"},
	Long: `The CLI is generated and the output is saved in the OutDir. Use the -c flag to
automatically compile the program after generation is complete.`,
	Run: func(cmd *cobra.Command, args []string) {
		generate()
	},
}

var Compile bool

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().BoolVarP(&Compile, "compile", "c", false, "Compile the resulting program")
}

func generate() {
	config := loadDxConfig()
	configMap := loadOutConfig(config)

	generateGoMod(config)

	generateMain(config)

	generateRoot(config)

	for _, c := range configMap.Commands {
		generateCommand(config, c)
	}

	tidy(config)

	if Compile {
		compile(config)
	}
}

func compile(config DxConfig) {
	run(config, "build")
}

func generateGoMod(config DxConfig) {
	run(config, "mod", "init", config.Cli)
}

func generateMain(config DxConfig) {
	generateFile("main.go", config)
}

func generateRoot(config DxConfig) {
	generateFile("cmd/root.go", config)
}

func generateCommand(config DxConfig, command map[string]string) {
	cmd := strings.Split(command["full"], " ")
	parent := cmd[:len(cmd)-1]
	outCmdData := OutCommandData{
		CommandName: command["name"],
		CommandVar:  generateCommandName(cmd),
		Parent:      generateCommandName(parent),
		Body:        command["body"],
	}
	// fmt.Println(outCmdData)
	generateCommandFile(config.OutDir, outCmdData)
}

func generateCommandName(commands []string) string {
	if len(commands) == 0 {
		return "root"
	}

	caser := cases.Title(language.English)
	names := make([]string, len(commands))
	for i, x := range commands {
		if i == 0 {
			names[i] = x
		} else {
			names[i] = caser.String(x)
		}
	}

	return strings.Join(names, "")
}

func run(config DxConfig, args ...string) {
	cmd := exec.Command("go", args...)
	cmd.Dir = config.OutDir
	cmd.Run()
}

func tidy(config DxConfig) {
	run(config, "mod", "tidy")
}

func generateFile(file string, config DxConfig) {
	f, err := os.Create(fmt.Sprintf("%s/%s", config.OutDir, file))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	t := template.Must(template.New(file).Parse(Templates[file]))
	err = t.Execute(f, config)
	if err != nil {
		log.Fatal(err)
	}
}

func generateCommandFile(path string, outCmd OutCommandData) {
	f, err := os.Create(fmt.Sprintf("%s/cmd/%s.go", path, outCmd.CommandVar))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	t := template.Must(template.New("command").Parse(Templates["command"]))
	err = t.Execute(f, outCmd)
	if err != nil {
		log.Fatal(err)
	}
}
