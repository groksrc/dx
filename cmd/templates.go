package cmd

var Templates = map[string]string{
	// main.go
	"main.go": `/*
*/
package main

import "{{ .Cli }}/cmd"

func main() {
	cmd.Execute()
}
`,

	// cmd/root.go
	"cmd/root.go": `/*
*/
package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "{{ .Cli }}",
	Short: "Short Description of your CLI",
	Long: ` + "`" + `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.` + "`" + `,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
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
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	// Find home directory.
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	// Search config in home directory with name ".{{ .Cli }}" (without extension).
	viper.AddConfigPath(home)
	viper.SetConfigType("yaml")
	viper.SetConfigName(".{{ .Cli }}")

	viper.AutomaticEnv() // read in environment variables that match

	viper.ReadInConfig()
}
`,
	// command
	"command": `
	package cmd

	import (
		"fmt"

		"github.com/spf13/cobra"
	)

	// {{ .CommandVar }}Cmd represents the {{ .CommandName }} command
	var {{ .CommandVar }}Cmd = &cobra.Command{
		Use:   "{{ .CommandName }}",
		Short: "A brief description of your command",
		Long: ` + "`" + `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:

	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.` + "`" + `,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("{{ .CommandName }} called")
		},
	}

	func init() {
		{{ .Parent }}Cmd.AddCommand({{ .CommandVar }}Cmd)
	}
`,
}
