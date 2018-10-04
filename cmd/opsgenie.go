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
	"fmt"
	"strings"

	"github.com/opsgenie/opsgenie-go-sdk/alertsv2"
	ogcli "github.com/opsgenie/opsgenie-go-sdk/client"
	"github.com/opsgenie/opsgenie-go-sdk/userv2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// opsgenieCmd represents the opsgenie command
var opsgenieCmd = &cobra.Command{
	Use:   "opsgenie",
	Short: "Issue filer for OpsGenie",
	Long: `Use the OpsGenie API integration to read alerts and file
issues against all acknowledged alerts. Uses the AlertID as the deduplication
key.`,
	Run: func(cmd *cobra.Command, args []string) {
		ogMain()
	},
}

func init() {
	rootCmd.AddCommand(opsgenieCmd)

	opsgenieCmd.PersistentFlags().String("opsgenie-api-key", "", "API key for OpsGenie")
	viper.BindPFlag("opsgenie-api-key", opsgenieCmd.PersistentFlags().Lookup("opsgenie-api-key"))
}

const (
	TestUUID             = "f879bc7a-3ee7-4af8-bc98-29bcc3dc3b12"
	MaxOutstandingIssues = 20 // TODO(philips): handle pagination, etc
)

type AlertIssue struct {
	ID             string
	Assignees      []string
	AcknowledgedBy string
	Description    string
	Subject        string
}

type AlertIssues []AlertIssue

type OpsGenie struct {
	cli *ogcli.OpsGenieClient
}

func (o OpsGenie) gitHubUsername(ogUsername string) (string, error) {
	userCli, _ := o.cli.UserV2()

	req := userv2.GetUserRequest{
		Identifier: &userv2.Identifier{
			Username: ogUsername,
		},
	}

	var username string
	resp, err := userCli.Get(req)
	if err != nil {
		return "", err
	}

	for _, t := range resp.User.Tags {
		if strings.HasPrefix(t, "github=") {
			username = fmt.Sprintf("%v", strings.TrimPrefix(t, "github="))
		}
	}

	return username, nil
}

func (o OpsGenie) acknowledgedAlerts() (AlertIssues, error) {
	var list AlertIssues

	alertCli, _ := o.cli.AlertV2()

	response, err := alertCli.List(alertsv2.ListAlertRequest{
		Limit:                MaxOutstandingIssues,
		Offset:               0,
		SearchIdentifierType: alertsv2.Name,
		Query:                "acknowledged=true",
		// TODO(philips): sort by date asc
	})
	if err != nil {
		return nil, err
	}

	for _, alert := range response.Alerts {
		response, err := alertCli.Get(alertsv2.GetAlertRequest{
			Identifier: &alertsv2.Identifier{ID: alert.ID},
		})
		if err != nil {
			return nil, err
		}
		handle, err := o.gitHubUsername(response.Alert.Report.AcknowledgedBy)
		if err != nil {
			return nil, err
		}
		list = append(list, AlertIssue{
			ID:             alert.ID,
			Assignees:      []string{handle},
			AcknowledgedBy: response.Alert.Report.AcknowledgedBy,
			Description:    response.Alert.Description,
			Subject:        alert.Message,
		})
	}

	return list, nil
}

func ogMain() {
	apiKey := viper.Get("opsgenie-api-key").(string)
	if apiKey == "" {
		fmt.Printf("ERROR: opsgenie-api-key is unset\n")
	}

	var ogc OpsGenie
	cli := new(ogcli.OpsGenieClient)
	cli.SetAPIKey(apiKey)
	ogc.cli = cli

	list, err := ogc.acknowledgedAlerts()
	if err != nil {
		panic(err)
	}

	for {
		for _, alert := range list {
			urls, err := ghc.findIssuesWithString(alert.ID)
			if err != nil {
				panic(err)
			}
			if len(urls) > 0 {
				fmt.Printf("INFO: alert %s has existing issue %s\n", alert.ID, urls[0])
				continue
			}

			fmt.Printf("INFO: filing issue for %v\n", alert.ID)
			issue, _, err := ghc.fileIssueFromAlert(alert)
			if err != nil {
				panic(err)
			}

			fmt.Printf("INFO: filed issue for %v at %s\n", alert.ID, *issue.URL)
		}
	}
}
