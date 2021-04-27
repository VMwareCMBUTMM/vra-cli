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
	"strings"
	"github.com/spf13/viper"
	"net/http"
	"io/ioutil"
	"github.com/tidwall/pretty"
)

// deploymentCmd represents the deployment command
var deploymentResourceCmd = &cobra.Command{
	Use:   "deploymentResource",
	Short: "Updates an existing deployments resource withan available day 2 action",
	Long: `Updates an existing deployments resource withan available day 2 action.
	For example:
	vra-cli update deploymentResource -n deployment1 -r k8s-clst1 -a deploy-helm -i helm-name:apache -i helm-release:my-apache`,
	Run: func(cmd *cobra.Command, args []string) {
		var str strings.Builder
		for i := 0; i < len(inputs); i++ {
		  s := strings.Split(inputs[i], ":")
		  if i == 0 {
		    str.WriteString(`"`+s[0]+`":"`+s[1]+`"`)
		  }
		  if i > 0 {
		    str.WriteString(`,"`+s[0]+`":"`+s[1]+`"`)
		  }
		}
		input_str := "{"+str.String()+"}"
		server := viper.GetString("target."+currentTargetName+".server")
		deployment_id := getDeploymentIdByName(name)
		dep_res_id := getDeploymentResourceIdByName(name, resource)
    action_id := getDeploymentResourceActionIdByName(name, resource, action)
		url := "https://"+server+"/deployment/api/deployments/"+deployment_id+"/resources/"+dep_res_id+"/requests"
		method := "POST"
		var token = getToken()
		payload := strings.NewReader(`{
		              "actionId": "`+action_id+`",
		              "inputs": `+input_str+`,
		              "targetId": "`+dep_res_id+`"
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
	updateCmd.AddCommand(deploymentResourceCmd)

  deploymentResourceCmd.Flags().StringVarP(&name, "name", "n", "","Name of the deployment.")
	deploymentResourceCmd.MarkFlagRequired("name")
	deploymentResourceCmd.Flags().StringVarP(&resource, "resource", "r", "","Name of the deployment resource.")
	deploymentResourceCmd.MarkFlagRequired("resource")
	deploymentResourceCmd.Flags().StringVarP(&action, "action", "a", "","Name of the deployment resource action.")
	deploymentResourceCmd.MarkFlagRequired("action")
	deploymentResourceCmd.Flags().StringArrayVarP(&inputs, "inputs", "i", []string{}, "Inputs for the deployment resource day 2 action request. Can be used mulitple times")
	deploymentResourceCmd.MarkFlagRequired("inputs")
}
