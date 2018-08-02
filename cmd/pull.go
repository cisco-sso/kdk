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
	"strings"

	"github.com/spf13/cobra"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull KDK docker image",
	Long:  `Pull the latest/configured KDK docker image`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := logrus.New().WithField("command", "pull")

		imageCoordinates := []string{viper.Get("image.repository").(string), viper.Get("image.tag").(string)}

		ctx := context.Background()
		cli, err := client.NewEnvClient()

		if err != nil {
			logger.WithField("error", err).Fatal("Failed to create docker client")
		}

		logger.Info("Pulling KDK image. This may take a few minutes...")

		out, err := cli.ImagePull(ctx, strings.Join(imageCoordinates, ":"), types.ImagePullOptions{})
		if err != nil {
			logger.WithField("error", err).Fatal("Failed to pull KDK image")
			panic(err)
		}
		defer out.Close()

		logger.Info("Successfully pulled KDK image.")
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)
}
