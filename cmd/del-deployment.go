/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
	"io/ioutil"
	"github.com/tidwall/pretty"
	"strings"
)

// deploymentCmd represents the deployment command
var deldeploymentCmd = &cobra.Command{
	Use:   "deployment",
	Short: "Deletes an existing deployment.",
	Long: `Deletes an existing deployment.
	For example:
	vra-cli delete deployment -n deployment1`,
	Run: func(cmd *cobra.Command, args []string) {
		server := viper.GetString("target."+currentTargetName+".server")
		deployment_id := getDeploymentIdByName(name)
		url := "https://"+server+"/deployment/api/deployments/"+deployment_id+"/requests"
		method := "POST"
		var token = getToken()
		payload := strings.NewReader(`{
		              "actionId": "Deployment.Delete",
									"inputs": {
								    "ignoreDeleteFailures": false
								  },
		              "targetId": "`+deployment_id+`"
		            }`)

		client := &http.Client {
		}

		req, err := http.NewRequest(method, url, payload)

		if err != nil {
		  fmt.Println(err)
		  return
		}
		req.Header.Add("Authorization", "Bearer " + token)
		req.Header.Add("Content-Type", "application/json")

		res, err := client.Do(req)
		if err != nil {
		  fmt.Println(err)
		  return
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
		  fmt.Println(err)
		  return
		}
		fmt.Println(string(pretty.Color(pretty.Pretty(body), nil)))
	},
}

func init() {
	deleteCmd.AddCommand(deldeploymentCmd)

	deldeploymentCmd.Flags().StringVarP(&name, "name", "n", "","Name of the deployment to delete.")
	deldeploymentCmd.MarkFlagRequired("name")
}
