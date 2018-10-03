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
	"strings"

	"github.com/opsgenie/opsgenie-go-sdk/alertsv2"
	ogcli "github.com/opsgenie/opsgenie-go-sdk/client"
	"github.com/opsgenie/opsgenie-go-sdk/userv2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type OpsGenie struct {
	APIKey string
}

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
	AcknowledgedBy string
	Description    string
	Subject        string
}

type AlertIssues []AlertIssue

func getGitHubUsername(userCli *ogcli.OpsGenieUserV2Client, ogUsername string) (string, error) {
	req := userv2.GetUserRequest{
		Identifier: &userv2.Identifier{
			Username: ogUsername,
		},
	}

	var username string
	resp, err := userCli.Get(req)
	if err != nil {
		return username, err
	}

	for _, t := range resp.User.Tags {
		if strings.HasPrefix(t, "github=") {
			username = fmt.Sprintf("%v", strings.TrimPrefix(t, "github="))
		}
	}

	return username, nil
}

func opsGenie(apiKey string) (AlertIssues, error) {
	var list AlertIssues
	cli := new(ogcli.OpsGenieClient)
	cli.SetAPIKey(apiKey)

	alertCli, _ := cli.AlertV2()
	userCli, _ := cli.UserV2()

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

	for i, alert := range response.Alerts {
		fmt.Printf("%v(%v): %v\n", i, alert.Acknowledged, alert.Message)
		response, err := alertCli.Get(alertsv2.GetAlertRequest{
			Identifier: &alertsv2.Identifier{ID: alert.ID},
		})
		if err != nil {
			return nil, err
		}
		if alert.Acknowledged {
			handle, err := getGitHubUsername(userCli, response.Alert.Report.AcknowledgedBy)
			if err != nil {
				return nil, err
			}
			list = append(list, AlertIssue{
				ID:             alert.ID,
				AcknowledgedBy: handle,
				Description:    response.Alert.Description,
				Subject:        alert.Message,
			})
		}
	}

	return nil, nil
}

func ogMain() {
	var ogc OpsGenie

	ogc.APIKey = viper.Get("opsgenie-api-key").(string)

	fmt.Printf("config: %v\n", viper.AllSettings())
	urls, err := ghc.findIssuesWithString(TestUUID)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", urls)
	opsGenie(ogc.APIKey)
}
