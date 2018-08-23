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
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/cisco-sso/kdk/internal/app/kdk"
	"github.com/docker/docker/client"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	versionNumber string
	cfgFile       string
	verbose       bool
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

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.WithFields(logrus.Fields{"err": err}).Fatal("Failed to execute RootCmd.")
	}
}

func init() {
	versionNumber = "0.5.3"
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kdk.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}

func initConfig() {
	kdk.Verbose = verbose

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		kdk.ConfigPath = cfgFile
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		kdk.ConfigDir = filepath.Join(home, ".kdk")
		kdk.ConfigName = "config"
		kdk.ConfigPath = filepath.Join(kdk.ConfigDir, kdk.ConfigName+".yaml")
		kdk.KeypairDir = filepath.Join(kdk.ConfigDir, "ssh")
		kdk.PrivateKeyPath = filepath.Join(kdk.KeypairDir, "id_rsa")
		kdk.PublicKeyPath = filepath.Join(kdk.KeypairDir, "id_rsa.pub")

		if _, err := os.Stat(kdk.ConfigDir); os.IsNotExist(err) {
			err = os.Mkdir(kdk.ConfigDir, 0700)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		viper.AddConfigPath(kdk.ConfigDir)
		viper.SetConfigName("config")
	}

	viper.SetEnvPrefix("kdk")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if viper.GetBool("json") {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
	if err != nil {
		logrus.WithFields(logrus.Fields{"configFileUsed": viper.ConfigFileUsed(), "err": err}).Warnln("Failed to load KDK config.")
	}
	if _, err := os.Stat(kdk.ConfigPath); err == nil {
		kdk.ImageCoordinates = viper.Get("image.repository").(string) + ":" + viper.Get("image.tag").(string)
		kdk.Name = viper.Get("docker.name").(string)
		kdk.Port = viper.Get("docker.environment.KDK_PORT").(string)
	}
	kdk.Ctx = context.Background()

	kdk.DockerClient, err = client.NewEnvClient()
	if err != nil {
		logrus.WithField("error", err).Fatal("Failed to create docker client")
	}
}
