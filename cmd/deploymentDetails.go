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
	"github.com/tidwall/pretty"
	"crypto/tls"
)

// catalogItemsCmd represents the catalogItems command
var deploymentDetailsCmd = &cobra.Command{
	Use:   "deploymentDetails",
	Short: "Lists the details of a spcific deployment by name",
	Long: `Lists the details of a spcific deployment by name

	For example: vra-cli list deploymentDetails -n deployment1`,
	Run: func(cmd *cobra.Command, args []string) {
		server := viper.GetString("target."+currentTargetName+".server")
		deployment_id := getDeploymentIdByName(name)
		url := "https://"+server+"/deployment/api/deployments/"+deployment_id
		method := "GET"
		var token = getToken()

		tr := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client {Transport: tr}

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
		fmt.Println(string(pretty.Color(pretty.Pretty(body), nil)))
	},
}

func init() {
	listCmd.AddCommand(deploymentDetailsCmd)

	deploymentDetailsCmd.Flags().StringVarP(&name, "name", "n", "","Name of the deployment.")
	deploymentDetailsCmd.MarkFlagRequired("name")
}
