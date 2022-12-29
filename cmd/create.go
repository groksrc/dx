package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "Creates your CLI in the outdir",
	Aliases: []string{"c"},
	Long: `Creates your CLI in the output directory.

Required parameters: cli, outdir, company

cli: The name of your company CLI. Must not contain whitespace.
outdir: The output directory where your CLI will be created. Must be a valid directory.
company: The name of your company.

These parameters can be specified in your dx config file or as arguments to the command.

Arguments to the command take precedence over values in the config file.

See the dx README for more information about configuration.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		create(args)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}

func create(args []string) {
	config := loadDxConfig()

	if len(args) > 0 {
		config.Cli = args[0]
	}

	if len(args) > 1 {
		config.OutDir = args[1]
	}

	if len(args) > 2 {
		config.Company = args[2]
	}

	validate(config)

	outDir := outDirPath(config)

	fmt.Printf("Creating %s CLI in %s for %s\n", config.Cli, outDir, config.Company)

	render(config)
}

func render(config DxConfig) {
	createOutDir(config)
	createCmdOutDir(config)
	if err := initOutConfig(config); err != nil {
		log.Fatal(err)
	}
}

func createOutDir(config DxConfig) {
	path := outDirPath(config)

	if err := os.MkdirAll(path, 0755); err != nil {
		log.Fatal(err)
	}
}

func createCmdOutDir(config DxConfig) {
	path := outDirPath(config) + "/cmd"

	if err := os.MkdirAll(path, 0755); err != nil {
		log.Fatal(err)
	}
}

func initOutConfig(config DxConfig) error {

	outDir := outDirPath(config)
	filePath := filepath.Join(outDir, fmt.Sprintf(".%s.yaml", config.Cli))

	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	stringsToWrite := []string{
		fmt.Sprintf("cli: %s\n", config.Cli),
		fmt.Sprintf("company: %s\n", config.Company),
		fmt.Sprintf("commands:\n")}

	for _, str := range stringsToWrite {
		if _, err = f.WriteString(str); err != nil {
			return err
		}
	}

	return nil
}

func validate(config DxConfig) {
	validateCli(config.Cli)
	validateRequired(config.OutDir, "outdir")
	validateRequired(config.Company, "company")
}

func validateCli(cli string) {
	validateRequired(cli, "cli")
	if regexp.MustCompile(`\s`).MatchString(cli) {
		log.Fatal("cli name cannot contain whitespace")
	}
}

func validateRequired(val string, name string) {
	if val == "" {
		log.Fatalf("%s is a required parameter", name)
	}
}
