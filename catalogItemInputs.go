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
	"github.com/Jeffail/gabs"
	"github.com/tidwall/pretty"
	"github.com/spf13/viper"
)

// catalogItemInputsCmd represents the catalogItemInputs command
var catalogItemInputsCmd = &cobra.Command{
	Use:   "catalogItemInputs",
	Short: "Lists the inputs available for the catalog item in Service Broker.",
	Long: `Lists the inputs available for the catalog item in Service Broker.

	For example: vra-cli list catalogItemInputs -n "Deploy Apache Server"`,
	Run: func(cmd *cobra.Command, args []string) {
		name, _:= cmd.Flags().GetString("name")
		getCatalogItemInputs(name)
	},
}

func getCatalogItemIdByName(name string) string {
	server := viper.GetString("target."+currentTargetName+".server")

  url := "https://"+server+"/catalog/api/items"
  method := "GET"
	var token = getToken()
  client := &http.Client {
  }
  req, err := http.NewRequest(method, url, nil)

  if err != nil {
    fmt.Println(err)
    return ("Failed")
  }
  req.Header.Add("Authorization", "Bearer " + token)
  req.Header.Add("Content-Type", "application/json")

  res, err := client.Do(req)
  if err != nil {
    fmt.Println(err)
    return ("Failed")
  }
  defer res.Body.Close()

  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    fmt.Println(err)
    return ("Failed")
  }
  var catalog Catalog
  json.Unmarshal([]byte(body), &catalog)
  for i := 0; i < len(catalog.Catalog); i++ {
		var cat_name = catalog.Catalog[i].Name
		if cat_name == name {
			var cat_id = catalog.Catalog[i].ID
			return cat_id
		} else if i == (len(catalog.Catalog)-1) {
			fmt.Println("Did not find catalog item: " + name)
		}
  }
	return ""
}

func getCatalogItemInputs(name string) {
	var cat_id = getCatalogItemIdByName(name)
	server := viper.GetString("target."+currentTargetName+".server")

  url := "https://"+server+"/catalog/api/items/" + cat_id
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

	jsonParsed, err := gabs.ParseJSON(body)
	if err != nil {
		fmt.Println(err)
    return
	}

	data := jsonParsed.Path("schema.properties").String()
	test := []byte(data)
	tryit := pretty.Pretty(test)
	fmt.Println(string(pretty.Color(tryit, nil)))
}

func init() {
	listCmd.AddCommand(catalogItemInputsCmd)

	catalogItemInputsCmd.Flags().StringP("name", "n","","Name of the catalog item.")
}
