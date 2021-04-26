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
var projectCmd = &cobra.Command{
	Use:   "projects",
	Short: "Lists the projects you are a member in vRA",
	Long: `Lists the projects you are a member in vRA.

	For example: vra-cli list prjects`,
	Run: func(cmd *cobra.Command, args []string) {
		listProjects()
	},
}

func listProjects() {
	server := viper.GetString("target."+currentTargetName+".server")

  url := "https://"+server+"/iaas/api/projects"
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
  var project Projects
  json.Unmarshal([]byte(body), &project)
  for i := 0; i < len(project.Project); i++ {
    var content = project.Project[i].Name
    //fmt.Println(content)
		color.Green(content)
  }
}

func init() {
	listCmd.AddCommand(projectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// catalogItemsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// catalogItemsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
