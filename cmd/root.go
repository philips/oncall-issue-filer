// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type GitHub struct {
	APIKey string
	Repo   string
}

var ghc GitHub

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "oncall-issue-filer",
	Short: "File issues on GitHub based on acknowledged OpsGenie alerts",
	Long: `On-call rotations where an incident requires multiple stakeholders
coming together over a period of days or weeks need a coordination point. For many
teams that is GitHub Issues.

This tool will assist an on-call person coordinate by automatically filing a
GitHub issue once an alert has been acknowledged.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	opsgenieCmd.PersistentFlags().StringVar(&ghc.APIKey, "github-api-key", "", "API key for GitHub")
	viper.BindPFlag("github-api-key", rootCmd.PersistentFlags().Lookup("github-api-key"))

	opsgenieCmd.PersistentFlags().StringVar(&ghc.Repo, "github-repo", "", "target GitHub Repo for filing")
	viper.BindPFlag("github-repo", rootCmd.PersistentFlags().Lookup("github-api-key"))

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetEnvPrefix("on_call_filer")
	viper.AutomaticEnv()
}
