/*
Package cmd Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cmd

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// currentTargetCmd represents the current-target command
var currentTargetCmd = &cobra.Command{
	Use:   "current-target",
	Short: "Display the current-target",
	Long: `Displays the current-target

Examples:
	# Display the current-target
	vra-cli config current-target
`,
	Run: func(cmd *cobra.Command, args []string) {
		var currentTargetName = viper.GetString("currentTargetName")
		if currentTargetName != "" {
			fmt.Println(currentTargetName)
		}
	},
}

// useTargetCmd represents the use-target command
var useTargetCmd = &cobra.Command{
	Use:   "use-target",
	Short: "Set the current target",
	Long: `Set the current target

Examples:
	# Display the current-target
	vra-cli config use-target --name vra8-test-ga
`,
	Run: func(cmd *cobra.Command, args []string) {
		var target = viper.Get("target." + name)
		if target == nil {
			log.Warningln("Target not found! Current target is", viper.GetString("currentTargetName"))
			return
		}
		viper.Set("currentTargetName", name)
		viper.WriteConfig()
		fmt.Println("Current target: ", name)
	},
}

// getConfigTargetCmd represents the get-target command
var getConfigTargetCmd = &cobra.Command{
	Use:   "get-target",
	Short: "Display available target configs",
	Long: `Displays a list of the available target configs

Examples:
	vra-cli config get-target
`,
	Run: func(cmd *cobra.Command, args []string) {
		if name != "" {
			var target = viper.Get("target." + name)
			if target == nil {
				log.Warningln("Target not found.")
			} else {
				PrettyPrint(target)
			}
		} else {
			var targets = viper.GetStringMapString("target")
			for key := range targets {
				fmt.Println(key)
			}
		}
	},
}

var (
	newTargetName string
	newServer     string
	newUsername   string
	newPassword   string
	newDomain     string
	newAPIToken   string
)

// setTargetCmd represents the set-target command
var setTargetCmd = &cobra.Command{
	Use:   "set-target",
	Short: "Creates or updates a target config",
	Long: `Creates or updates a target configuration.

Examples:
  (On-prem Example)
	vra-cli config set-target --name vra-test-ga --server vra8-test-ga.local.com --username test-user --password VMware1! --domain domain.local
	(vRA Cloud Example)
	vra-cli config set-target --name vrac-org --server api.mgmt.cloud.vmware.com --apitoken JhbGciOiJSUzI1NiIsImtpZCI6IjEzNjY3NDcwMTA2Mzk2MTUxNDk0In0
`, Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if viper.IsSet("target." + newTargetName) {
			fmt.Println("Updating", newTargetName)
		} else {
			fmt.Println("Creating new target", newTargetName)
		}
		fmt.Println("Use `vra-cli config use-target --name " + newTargetName + "` to use this target")
		if newServer != "" {
			viper.Set("target."+newTargetName+".server", newServer)
		}
		if newUsername != "" {
			viper.Set("target."+newTargetName+".username", newUsername)
		}
		if newPassword != "" {
			viper.Set("target."+newTargetName+".password", newPassword)
		}
		if newDomain != "" {
			viper.Set("target."+newTargetName+".domain", newDomain)
		}
		if newAPIToken != "" {
			viper.Set("target."+newTargetName+".apitoken", newAPIToken)
		}
		viper.SetConfigType("yaml")
		err := viper.SafeWriteConfig()
		if err != nil {
			log.Info(err)
			viper.WriteConfig()
		}
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Create configuration for vRealize Automation Access",
	Long: `Create a config for vRealize Automation which allows making calls to the vRealize Automation API`,
}

func init() {
	rootCmd.AddCommand(configCmd)
	// current-target
	configCmd.AddCommand(currentTargetCmd)
	// use-target
	configCmd.AddCommand(useTargetCmd)
	useTargetCmd.Flags().StringVarP(&name, "name", "n", "", "Use the target with this name")
	useTargetCmd.MarkFlagRequired("name")
	// get-target
	configCmd.AddCommand(getConfigTargetCmd)
	getConfigTargetCmd.Flags().StringVarP(&name, "name", "n", "", "Display the target with this name")
	// set-target
	configCmd.AddCommand(setTargetCmd)
	setTargetCmd.Flags().StringVarP(&newTargetName, "name", "n", "", "Name of the target configuration")
	setTargetCmd.Flags().StringVarP(&newServer, "server", "s", "", "Server FQDN of the vRealize Automation instance")
	setTargetCmd.Flags().StringVarP(&newUsername, "username", "u", "", "Username to authenticate.")
	setTargetCmd.Flags().StringVarP(&newPassword, "password", "p", "", "Password to authenticate")
	setTargetCmd.Flags().StringVarP(&newDomain, "domain", "d", "", "User domain (ex domain.local). <Optional>")
	setTargetCmd.Flags().StringVarP(&newAPIToken, "apitoken", "a", "", "API token for vRealize Automation Cloud.")
	setTargetCmd.MarkFlagRequired("name")
}
