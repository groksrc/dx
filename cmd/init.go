/*
Copyright © 2022 Drew Cain -- @groksrc

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"golang.org/x/exp/slices"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize dx for your environment",
	Long: `dx init will create a config file at $HOME/.dx.yaml

See the dx README for more information about configuration.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				initializeConfig(cmd, args)
			} else {
				overwrite("A config was found but could not be read. Overwrite? [Y/n]", cmd, args)
			}
		} else {
			overwrite("A config was found. Overwrite? [Y/n]", cmd, args)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func setDefault(params ...string) {
	fmt.Println(params[1])
	reader := bufio.NewReader(os.Stdin)
	val, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	if val == "\n" && len(params) == 3 {
		val = params[2]
	}
	viper.Set(params[0], strings.TrimSpace(val))
}

// TODO: Set cobra config
func initializeConfig(cmd *cobra.Command, args []string) {
	viper.SafeWriteConfig()
	setDefault("company", "Enter your company name:")
	setDefault("cli", "Enter cli name:")
	setDefault("description", "Enter a description for your CLI. Displayed by running the root command:")
	setDefault("outdir", "Enter the output directory: [./out]", "./out")
	if err := viper.WriteConfig(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("dx init success!")
	}
}

func overwrite(message string, cmd *cobra.Command, args []string) {
	fmt.Println(message)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	yes := scanner.Text()
	if slices.Contains([]string{"Y", "y", ""}, yes) {
		initializeConfig(cmd, args)
	}
}
