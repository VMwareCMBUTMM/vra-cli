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
	"io/ioutil"
	"net/http"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/fatih/color"
	"encoding/json"
)

// catalogItemsCmd represents the catalogItems command
var deploymentResourceActionsCmd = &cobra.Command{
	Use:   "deploymentResourceActions",
	Short: "Lists the Day 2 Actions available for the specified deployment resource",
	Long: `Lists the Day 2 Actions available for the specified deployment resource

	For example: vra-cli list deploymentResourceActions -n deployment1 -r "Deploy HELM"`,
	Run: func(cmd *cobra.Command, args []string) {
		server := viper.GetString("target."+currentTargetName+".server")
		deployment_id := getDeploymentIdByName(name)
		dep_res_id := getDeploymentResourceIdByName(name, resource)
		url := "https://"+server+"/deployment/api/deployments/"+deployment_id+"/resources/"+dep_res_id+"/actions"
		method := "GET"
		var token = getToken()

		client := &http.Client {
		}

		req, err := http.NewRequest(method, url, nil)

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
		var action []Action
	  json.Unmarshal([]byte(body), &action)
	  for i := 0; i < len(action); i++ {
	    var content = action[i].Name
			color.Green(content)
	  }
	},
}

func init() {
	listCmd.AddCommand(deploymentResourceActionsCmd)

	deploymentResourceActionsCmd.Flags().StringVarP(&name, "name", "n", "","Name of the deployment.")
	deploymentResourceActionsCmd.MarkFlagRequired("name")
	deploymentResourceActionsCmd.Flags().StringVarP(&resource, "resource", "r", "","Name of the deployment resource.")
	deploymentResourceActionsCmd.MarkFlagRequired("resource")
}
