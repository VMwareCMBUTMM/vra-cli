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
	"encoding/json"
	"github.com/spf13/cobra"
	"github.com/fatih/color"
	"github.com/spf13/viper"
)

// catalogItemsCmd represents the catalogItems command
var deploymentsCmd = &cobra.Command{
	Use:   "deployments",
	Short: "Lists al deployments you are allowed to see.",
	Long: `Lists al deployments you are allowed to see.

	For example: vra-cli list deploymentDetails -n deployment1`,
	Run: func(cmd *cobra.Command, args []string) {
		server := viper.GetString("target."+currentTargetName+".server")
		url := "https://"+server+"/deployment/api/deployments"
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
		var deployment Deployments
	  json.Unmarshal([]byte(body), &deployment)
	  for i := 0; i < len(deployment.Deployment); i++ {
	    var content = deployment.Deployment[i].Name
			color.Green(content)
	  }
	},
}

func init() {
	listCmd.AddCommand(deploymentsCmd)
}
