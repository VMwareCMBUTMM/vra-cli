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
	"github.com/Jeffail/gabs"
	"github.com/tidwall/pretty"
	"crypto/tls"
)

// catalogItemsCmd represents the catalogItems command
var deploymentResourceActionInputsCmd = &cobra.Command{
	Use:   "deploymentResourceActionInputs",
	Short: "Lists the inputs for the Day 2 Actions available for the specified deployment resource",
	Long: `Lists the inputs for the Day 2 Actions available for the specified deployment resource

	For example: vra-cli list deploymentResourceActionInputs -n deployment1 -r "dev-clst01" -a deploy-helm`,
	Run: func(cmd *cobra.Command, args []string) {
		server := viper.GetString("target."+currentTargetName+".server")
		deployment_id := getDeploymentIdByName(name)
		dep_res_id := getDeploymentResourceIdByName(name, resource)
    action_id := getDeploymentResourceActionIdByName(name, resource, action)
		url := "https://"+server+"/deployment/api/deployments/"+deployment_id+"/resources/"+dep_res_id+"/actions/"+action_id
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

		jsonParsed, err := gabs.ParseJSON(body)
		if err != nil {
			fmt.Println(err)
	    return
		}

		data := jsonParsed.Path("schema.properties").String()
		test := []byte(data)
		tryit := pretty.Pretty(test)
		fmt.Println(string(pretty.Color(tryit, nil)))
	},
}

func init() {
	listCmd.AddCommand(deploymentResourceActionInputsCmd)

	deploymentResourceActionInputsCmd.Flags().StringVarP(&name, "name", "n", "","Name of the deployment.")
	deploymentResourceActionInputsCmd.MarkFlagRequired("name")
	deploymentResourceActionInputsCmd.Flags().StringVarP(&resource, "resource", "r", "","Name of the deployment resource.")
	deploymentResourceActionInputsCmd.MarkFlagRequired("resource")
	deploymentResourceActionInputsCmd.Flags().StringVarP(&action, "action", "a", "","Name of the deployment resource action.")
	deploymentResourceActionInputsCmd.MarkFlagRequired("action")
}
