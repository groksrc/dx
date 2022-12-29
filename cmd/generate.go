/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"text/template"

	"github.com/spf13/cobra"
)

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
		generateCommand(c)
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
	// cmd := exec.Command("go", "mod", "init", config.Cli)
	// cmd.Dir = config.OutDir
	// cmd.Run()
	run(config, "mod", "init", config.Cli)
}

func generateMain(config DxConfig) {
	generateFile("main.go", config)
}

func generateRoot(config DxConfig) {
	generateFile("cmd/root.go", config)
}

func generateCommand(command map[string]interface{}) {
	// open file

	//
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
