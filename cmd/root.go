// Copyright © 2018 Cisco Systems, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	log "github.com/Sirupsen/logrus"
)

var (
	versionNumber string
	cfgFile       string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kdk",
	Short: "Kubernetes Development Kit",
	Long: `

 _  __ ____  _  __
/ |/ //  _ \/ |/ /
|   / | | \||   / 
|   \ | |_/||   \ 
\_|\_\\____/\_|\_\
                  

A full kubernetes development environment in a container`,
}

var KdkConfigDir string

var KdkImageCoordinates string

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatal("Failed to execute RootCmd.")
	}
}

func init() {
	versionNumber = "0.5.1"
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kdk.yaml)")
}

func initConfig() {

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		KdkConfigDir = path.Join(home, ".kdk")

		if _, err := os.Stat(KdkConfigDir); os.IsNotExist(err) {
			err = os.Mkdir(KdkConfigDir, 0700)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		viper.AddConfigPath(KdkConfigDir)
		viper.SetConfigName("config")
	}

	// TODO (rluckie) allow config to be set from env var
	viper.SetEnvPrefix("kdk")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if viper.GetBool("json") {
		log.SetFormatter(&log.JSONFormatter{})
	}
	if err != nil {
		log.WithFields(log.Fields{
			"configFileUsed": viper.ConfigFileUsed(),
			"err":            err,
		}).Warnln("Failed to load KDK config.")
	}
	// TODO (rluckie) move KdkImageCoordinates
	KdkImageCoordinates = viper.Get("image.repository").(string) + ":" + viper.Get("image.tag").(string)
}
