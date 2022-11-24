package cmd

import (
	// "bufio"
	// "log"
	"fmt"
	// "os"

	"github.com/spf13/cobra"
	// "github.com/spf13/viper"

)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Builds your CLI in the outdir",
	Long: `dx build creates your CLI in the output directory

See the dx README for more information about configuration.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("building cli")
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}