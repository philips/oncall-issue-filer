// Copyright © 2018 Brandon Philips <brandon@ifup.org>
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
	"net/http"
	"os"
	"strings"

	"github.com/google/go-github/v18/github"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

var ghc GitHub

// TODO: put into new package
type GitHub struct {
	Client *github.Client
	Repo   string
}

// findIssues finds issues that contain the given OpsGenie Alert ID and returns
// a list of issue URLs. NOTE: we actually don't care about the URLs and only
// request the first result.
func (g GitHub) findIssuesWithString(id string) ([]string, error) {
	query := fmt.Sprintf("repo:%s %s", g.Repo, id)
	opts := &github.SearchOptions{
		Sort:        "date",
		Order:       "desc",
		ListOptions: github.ListOptions{Page: 1, PerPage: 1},
	}
	result, _, err := g.Client.Search.Issues(context.Background(), query, opts)
	if err != nil {
		return nil, err
	}
	for _, issue := range result.Issues {
		return []string{*issue.URL}, nil
	}

	return nil, nil
}

func (g GitHub) fileIssueFromAlert(alert AlertIssue) (*github.Issue, *github.Response, error) {
	ctx := context.Background()
	body := fmt.Sprintf("%v\n\nAlert ID: %v\nAcknowledgedBy: %v\n",
		alert.Description,
		alert.ID,
		alert.AcknowledgedBy)

	repo := strings.Split(g.Repo, "/")

	return g.Client.Issues.Create(ctx, repo[0], repo[1], &github.IssueRequest{
		Title:     &alert.Subject,
		Body:      &body,
		Assignees: &alert.Assignees,
	})
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

func initConfig() {
	viper.SetEnvPrefix("oncall_issue_filer")
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()

	ghc.Repo = viper.Get("github-repo").(string)
	apiKey := viper.Get("github-api-key").(string)
	var tc *http.Client
	if apiKey != "" {
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: viper.Get("github-api-key").(string)},
		)
		tc = oauth2.NewClient(ctx, ts)
	} else {
		fmt.Printf("INFO: github-api-key unset, using anonymous API auth\n")
	}

	ghc.Client = github.NewClient(tc)
}
