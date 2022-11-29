package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct{ cli, company, outdir string }

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Builds your CLI in the outdir",
	Long: `dx build creates your CLI in the output directory

See the dx README for more information about configuration.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		build()
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}

func build() {
	config := loadConfig()
	fmt.Printf("Creating %s CLI in %s for %s\n", config.cli, config.outdir, config.company)

	generate(config)
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

func installCobra() {
	if _, err := exec.LookPath("cobra-cli"); err != nil {
		fmt.Println("Installing Cobra")
		cmd := exec.Command("go", "install", "github.com/spf13/cobra-cli@latest")
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
	}
}

func loadConfig() Config {
	settings := viper.AllSettings()

	if settings["cli"] == "" {
		log.Fatal("cli config value missing")
	}

	if settings["company"] == "" {
		log.Fatal("company config value missing")
	}

	cli := settings["cli"].(string)
	company := settings["company"].(string)
	outdir := settings["outdir"].(string)

	return Config{cli, company, outdir}
}
