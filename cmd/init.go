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
	"github.com/Sirupsen/logrus"
	"github.com/cisco-sso/kdk/internal/app/kdk"
	"github.com/spf13/cobra"
)

var (
	port            string
	imageRepository string
	imageTag        string
	dotfilesRepo    string
	shell           string
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize KDK",
	Long:  `Initialize KDK: Create/recreate KDK configuration and pull latest image`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := logrus.New().WithField("command", "init")

		kdk.InitKdkConfig(KdkName, port, imageRepository, imageTag, dotfilesRepo, shell, *logger)
		kdk.InitKdkSshKeyPair(*logger)
		logger.Infof("KDK config written to %s. Modify this file to suit your needs.", kdk.ConfigPath)
	},
}

func init() {
	initCmd.Flags().StringVarP(&KdkName, "name", "n", "kdk", "KDK Name")
	initCmd.Flags().StringVarP(&port, "port", "p", "2022", "KDK Port")
	initCmd.Flags().StringVarP(&imageRepository, "image-repository", "r", "ciscosso/kdk", "KDK Image Repository")
	initCmd.Flags().StringVarP(&imageTag, "image-tag", "t", "debian-latest", "KDK Image Tag")
	initCmd.Flags().StringVarP(&dotfilesRepo, "dotfiles-repo", "d", "https://github.com/cisco-sso/yadm-dotfiles.git", "KDK Dotfiles Repo")
	initCmd.Flags().StringVarP(&shell, "shell", "s", "/bin/bash", "KDK shell")

	rootCmd.AddCommand(initCmd)
}
