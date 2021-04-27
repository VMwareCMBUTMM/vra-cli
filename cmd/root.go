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
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/mrz1836/go-sanitize"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Configuration
	cfgFile           string
	currentTargetName string
	targetConfig      config
	version           = "dev"
	commit            = "none"
	date              = "unknown"
	builtBy           = "unknown"
	// Global Flags
	debug      bool
	ignoreCert bool
	// Command Flags
	id          string
	name        string
	project     string
	deployment  string
	resource    string
	action      string
	actions     string
	inputs      []string
	typename    string
	value       string
	description string
	status      string
	// exportFile  string
	// importFile  string
	printJson  bool
	exportPath string
	importPath string
)

type config struct {
	password    string
	server      string
	username    string
	accesstoken string
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "vra-cli",
	Short: "This is the vRealize Automation CLI to interface with vRA.",
	Long: `This is the vRealize Automation CLI to interface with vRA.`,
	// Run: func(cmd *cobra.Command, args []string) {},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.vra-cli.yaml)")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug logging")
	rootCmd.PersistentFlags().BoolVar(&ignoreCert, "ignoreCertificateWarnings", false, "Disable HTTPS Certificate Validation")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // If the user has specified a config file
		if file, err := os.Stat(cfgFile); err == nil { // Check if it exists
			viper.SetConfigFile(file.Name())
		} else {
			log.Fatalln("File specified with --config does not exist")
		}
	}

	// Home directory
	home, err := homedir.Dir()
	if err != nil {
		log.Fatalln(err)
	}

	viper.SetConfigName(".vra-cli")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(home)

	// Attempt to read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			viper.SetConfigType("yaml")
			viper.WriteConfigAs(filepath.Join(home, ".vra-cli"))
			viper.ReadInConfig()
		} else {
			log.Fatalln(err)
		}
	}
	currentTargetName = viper.GetString("currentTargetName")
	if currentTargetName != "" {
		log.Println("Using config:", viper.ConfigFileUsed(), "Target:", currentTargetName)
		configuration := viper.Sub("target." + currentTargetName)
		if configuration == nil { // Sub returns nil if the key cannot be found
			log.Fatalln("Target configuration not found")
		}
		targetConfig = config{
			server:      sanitize.URL(configuration.GetString("server")),
			username:    configuration.GetString("username"),
			password:    configuration.GetString("password"),
			accesstoken: configuration.GetString("accesstoken"),
		}
	}
}
