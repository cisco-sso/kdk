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
	log "github.com/Sirupsen/logrus"
	"github.com/cisco-sso/kdk/pkg/kdk"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize KDK",
	Long:  `Initialize KDK: Create/recreate KDK configuration and pull latest image`,
	Run: func(cmd *cobra.Command, args []string) {
		CurrentKdkEnvConfig.CreateKdkConfig()
		CurrentKdkEnvConfig.CreateKdkSshKeyPair()
		log.Infof("KDK config written to %s. Modify this file to suit your needs.", CurrentKdkEnvConfig.ConfigPath())
	},
}

func init() {
	initCmd.Flags().StringVarP(&CurrentKdkEnvConfig.ConfigFile.AppConfig.Name, "name", "n", "kdk", "KDK Name")
	initCmd.Flags().StringVarP(&CurrentKdkEnvConfig.ConfigFile.AppConfig.Port, "port", "p", kdk.Port, "KDK Port")
	initCmd.Flags().StringVarP(&CurrentKdkEnvConfig.ConfigFile.AppConfig.ImageRepository, "image-repository", "r", "ciscosso/kdk", "KDK Image Repository")
	initCmd.Flags().StringVarP(&CurrentKdkEnvConfig.ConfigFile.AppConfig.ImageTag, "image-tag", "t", kdk.Version, "KDK Image Tag")
	initCmd.Flags().StringVarP(&CurrentKdkEnvConfig.ConfigFile.AppConfig.DotfilesRepo, "dotfiles-repo", "", "https://github.com/cisco-sso/yadm-dotfiles.git", "KDK Dotfiles Repo")
	initCmd.Flags().StringVarP(&CurrentKdkEnvConfig.ConfigFile.AppConfig.Shell, "shell", "s", "/bin/bash", "KDK shell")

	rootCmd.AddCommand(initCmd)
}
