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
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v18/github"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ghc GitHub

// TODO: put into new package
type GitHub struct {
	APIKey string
	Repo   string
}

// findIssues finds issues that contain the given OpsGenie Alert ID and returns
// a list of issue URLs. NOTE: we actually don't care about the URLs and only
// request the first result.
func (g GitHub) findIssuesWithString(id string) ([]string, error) {
	client := github.NewClient(nil)

	query := fmt.Sprintf("is:open repo:%s %s", g.Repo, id)
	println(query)
	opts := &github.SearchOptions{
		Sort:        "date",
		Order:       "desc",
		ListOptions: github.ListOptions{Page: 1, PerPage: 1},
	}
	result, _, err := client.Search.Issues(context.Background(), query, opts)
	if err != nil {
		return nil, err
	}
	for _, issue := range result.Issues {
		return []string{*issue.URL}, nil
	}

	return nil, nil
}

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
	rootCmd.PersistentFlags().String("github-api-key", "", "API key for GitHub")
	viper.BindPFlag("github-api-key", rootCmd.PersistentFlags().Lookup("github-api-key"))

	rootCmd.PersistentFlags().String("github-repo", "", "target GitHub Repo for filing")
	viper.BindPFlag("github-repo", rootCmd.PersistentFlags().Lookup("github-repo"))

	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetEnvPrefix("oncall_issue_filer")
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()

	ghc.Repo = viper.Get("github-repo").(string)
	ghc.APIKey = viper.Get("github-api-key").(string)
}
