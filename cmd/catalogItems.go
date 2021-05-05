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
	"crypto/tls"
)

// catalogItemsCmd represents the catalogItems command
var catalogItemsCmd = &cobra.Command{
	Use:   "catalogItems",
	Short: "Lists the catalog items from vRA Service Broker",
	Long: `Lists the catalog items from vRA Service Broker catalog which the user credentials have access to deploy.

	For example: vra-cli list catalogItems`,
	Run: func(cmd *cobra.Command, args []string) {
		listCatalogItems()
	},
}

func listCatalogItems() {
	server := viper.GetString("target."+currentTargetName+".server")

  url := "https://"+server+"/catalog/api/items"
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
  var catalog Catalog
  json.Unmarshal([]byte(body), &catalog)
  for i := 0; i < len(catalog.Catalog); i++ {
    var content = catalog.Catalog[i].Name
		color.Green(content)
  }
}

func init() {
	listCmd.AddCommand(catalogItemsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// catalogItemsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// catalogItemsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
