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
	"github.com/spf13/viper"
	"strings"
	"crypto/tls"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

func getToken() string {
	if viper.GetString("target."+currentTargetName+".apitoken") == "" {
		user := viper.GetString("target."+currentTargetName+".username")
		pass := viper.GetString("target."+currentTargetName+".password")
		server := viper.GetString("target."+currentTargetName+".server")
		domain := viper.GetString("target."+currentTargetName+".domain")

	  url := "https://"+server+"/csp/gateway/am/api/login"
	  method := "POST"

	  payload := strings.NewReader(`{
	  "username": "`+user+`",
	  "password": "`+pass+`",
		"domain": "`+domain+`"
	  }`)
		tr := &http.Transport{
	        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	    }
	  client := &http.Client {Transport: tr}
	  req, err := http.NewRequest(method, url, payload)

	  if err != nil {
	    fmt.Println(err)
	    return ("failed")
	  }
	  req.Header.Add("Content-Type", "application/json")

	  res, err := client.Do(req)
	  if err != nil {
	    fmt.Println(err)
	    return ("failed")
	  }
	  defer res.Body.Close()

	  body, err := ioutil.ReadAll(res.Body)
	  if err != nil {
	    fmt.Println(err)
	    return ("failed")
	  }
	  var token Token
	  json.Unmarshal([]byte(body), &token)
	  var access_token = token.Token
	  return access_token

	} else if viper.GetString("target."+currentTargetName+".apitoken") != "" {
		vrac_token := viper.GetString("target."+currentTargetName+".apitoken")
    server := viper.GetString("target."+currentTargetName+".server")
		client := resty.New()
		queryResponse, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: ignoreCert}).R().
			SetBody(AuthenticationRequestCloud{vrac_token}).
			SetResult(&AuthenticationResponseCloud{}).
			SetError(&AuthenticationError{}).
			Post("https://" + server + "/iaas/api/login")
		if err != nil {
		  fmt.Println(err)
		  return ("failed")
		}
		if queryResponse.IsError() {
			log.Debugln("Authentication failed!", queryResponse.RawResponse)
			return "Authentication failed!"
		}
		return queryResponse.Result().(*AuthenticationResponseCloud).Token
		}
	return ""
}

func getProjectIdByName(name string) string {
	server := viper.GetString("target."+currentTargetName+".server")

  url := "https://"+server+"/iaas/api/projects"
  method := "GET"
	var token = getToken()
	tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
  client := &http.Client {Transport: tr}
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
  var project Projects
  json.Unmarshal([]byte(body), &project)
  for i := 0; i < len(project.Project); i++ {
		var proj_name = project.Project[i].Name
		if proj_name == name {
			var cat_id = project.Project[i].ID
			return cat_id
		} else if i == (len(project.Project)-1) {
			fmt.Println("Did not find project: " + name)
		}
  }
	return ""
}

func getDeploymentIdByName(name string) string {
	server := viper.GetString("target."+currentTargetName+".server")

  url := "https://"+server+"/deployment/api/deployments"
  method := "GET"
	var token = getToken()
	tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
  client := &http.Client {Transport: tr}
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
  var deployment Deployments
  json.Unmarshal([]byte(body), &deployment)
  for i := 0; i < len(deployment.Deployment); i++ {
		var dep_name = deployment.Deployment[i].Name
		if dep_name == name {
			var dep_id = deployment.Deployment[i].ID
			return dep_id
		} else if i == (len(deployment.Deployment)-1) {
			fmt.Println("Did not find deployment: " + name)
		}
  }
	return ""
}

func getDeploymentResourceIdByName(name, resource string) string {
	server := viper.GetString("target."+currentTargetName+".server")
  deployment_id := getDeploymentIdByName(name)
  url := "https://"+server+"/deployment/api/deployments/"+deployment_id+"/resources"
  method := "GET"
	var token = getToken()
	tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
  client := &http.Client {Transport: tr}
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
  var deployment Deployments
  json.Unmarshal([]byte(body), &deployment)
  for i := 0; i < len(deployment.Deployment); i++ {
		var dep_res_name = deployment.Deployment[i].Name
		if dep_res_name == resource {
			var dep_res_id = deployment.Deployment[i].ID
			return dep_res_id
		} else if i == (len(deployment.Deployment)-1) {
			fmt.Println("Did not find deployment resource: " + name)
		}
  }
	return ""
}

func getDeploymentResourceActionIdByName(name, resource, action string) string {
	server := viper.GetString("target."+currentTargetName+".server")
  deployment_id := getDeploymentIdByName(name)
	dep_res_id := getDeploymentResourceIdByName(name, resource)
  url := "https://"+server+"/deployment/api/deployments/"+deployment_id+"/resources/"+dep_res_id+"/actions"
  method := "GET"
	var token = getToken()
	tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
  client := &http.Client {Transport: tr}
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
	var actions []Action
	json.Unmarshal([]byte(body), &actions)
  for i := 0; i < len(actions); i++ {
		var action_name = actions[i].Name
		if action_name == action {
			var action_id = actions[i].ID
			return action_id
		} else if i == (len(actions)-1) {
			fmt.Println("Did not find deployment resource action: " + action)
		}
  }
	return ""
}

type arrayFlags []string

func (i *arrayFlags) String() string {
    return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
    *i = append(*i, value)
    return nil
}
