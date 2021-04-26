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
	"github.com/tidwall/pretty"
	"github.com/spf13/viper"
	"net/http"
	"io/ioutil"
)

// catalogItemCmd represents the catalogItem command
var catalogItemCmd = &cobra.Command{
	Use:   "catalogItem",
	Short: "Use this command with the necessary inputs to deploy a catalog item from vRA.",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.

For example: vra-cli deploy catalogItem -d deployment1 -n catalogitem1 -p project1 -i input1:foo -i input2:bar
`,
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
		proj_id := getProjectIdByName(project)
		catalogItem_id := getCatalogItemIdByName(name)
	  url := "https://"+server+"/catalog/api/items/"+catalogItem_id+"/request"
	  method := "POST"
		var token = getToken()
	  payload := strings.NewReader(`{
								  "deploymentName": "`+deployment+`",
								  "inputs": `+input_str+`,
								  "projectId": "`+proj_id+`"
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
	deployCmd.AddCommand(catalogItemCmd)
	catalogItemCmd.Flags().StringVarP(&name, "name", "n", "", "Deploy catalog item with this name")
	catalogItemCmd.Flags().StringVarP(&deployment, "deployment", "d", "", "Name for the deployment")
	catalogItemCmd.Flags().StringVarP(&project, "project", "p", "", "vRA Project name")
	catalogItemCmd.Flags().StringArrayVarP(&inputs, "inputs", "i", []string{}, "Inputs for the catalog deployment request. Can be used mulitple times")
	catalogItemCmd.MarkFlagRequired("name")
	catalogItemCmd.MarkFlagRequired("deployment")
	catalogItemCmd.MarkFlagRequired("project")
}
