package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct{ cli, company, outdir string }

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates your CLI in the outdir",
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
	cliCmd.AddCommand(createCmd)
}

func create(args []string) {
	config := loadConfig()

	if len(args) > 0 {
		config.cli = args[0]
	}

	if len(args) > 1 {
		config.outdir = args[1]
	}

	if len(args) > 2 {
		config.company = args[2]
	}

	validate(config)

	fmt.Printf("Creating %s CLI in %s for %s\n", config.cli, config.outdir, config.company)

	generate(config)
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

func generate(config Config) {
	if err := os.MkdirAll(config.outdir, os.ModePerm); err != nil {
		log.Println(err)
	}
	render(config)
}

func render(config Config) {
	installCobra()
	cobraInit(config)
	if err := viperInit(config); err != nil {
		log.Fatal((err))
	}
}

func installCobra() {
	if _, err := exec.LookPath("cobra-cli"); err != nil {
		fmt.Println("Installing Cobra")
		cmd := exec.Command("go", "install", "github.com/spf13/cobra-cli@latest")
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
	}
}

func cobraInit(config Config) {
	if err := os.Chdir(config.outdir); err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("cobra-cli", "init", config.cli)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func viperInit(config Config) error {

	// TODO: Fix relative paths, ~, $HOME, etc
	filePath := filepath.Join(config.outdir, config.cli, fmt.Sprintf(".%s.yaml", config.cli))

	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	stringsToWrite := []string{
		fmt.Sprintf("cli: %s\n", config.cli),
		fmt.Sprintf("company: %s\n", config.company),
		fmt.Sprintf("commands:\n")}

	for _, str := range stringsToWrite {
		if _, err = f.WriteString(str); err != nil {
			return err
		}
	}

	return nil
}

func validate(config Config) {
	validateCli(config.cli)
	validateOutdir(config.outdir)
	validateCompany(config.company)
}

func validateCli(cli string) {
	validateRequired(cli, "cli")
	if regexp.MustCompile(`\s`).MatchString(cli) {
		log.Fatal("cli name cannot contain whitespace")
	}
}

func validateOutdir(outdir string) {
	validateRequired(outdir, "outdir")
	// TODO: validate path
}

func validateCompany(company string) {
	validateRequired(company, "company")
}

func validateRequired(val string, name string) {
	if val == "" {
		log.Fatalf("%s is a required parameter", name)
	}
}
