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
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Start KDK container",
	Long:  `Start KDK container`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := logrus.New().WithField("command", "up")

		imageCoordinates := []string{viper.Get("image.repository").(string), viper.Get("image.tag").(string)}

		client, err := client.NewEnvClient()
		if err != nil {
			logger.WithField("error", err).Fatal("Failed to create docker client")
		}

		ctx := context.Background()

		var binds []string

		for _, bind := range viper.Get("docker.binds").([]interface{}) {
			configBind := make(map[string]string)
			for key, value := range bind.(map[interface{}]interface{}) {
				configBind[key.(string)] = value.(string)
			}
			binds = append(binds, fmt.Sprintf("%s:%s", configBind["source"], configBind["target"]))
		}

		containerCreateResp, err := client.ContainerCreate(
			ctx,
			&container.Config{
				Hostname: viper.Get("docker.hostname").(string),
				Image:    strings.Join(imageCoordinates, ":"),
				Tty:      true,
				Env: []string{
					fmt.Sprintf("KDK_USERNAME=%s", viper.Get("docker.environment.KDK_USERNAME").(string)),
					fmt.Sprintf("KDK_SHELL=%s", viper.Get("docker.environment.KDK_SHELL").(string)),
					fmt.Sprintf("KDK_DOTFILES_REPO=%s", viper.Get("docker.environment.KDK_DOTFILES_REPO").(string)),
				},
				ExposedPorts: nat.PortSet{
					"2022": struct{}{},
				},
			},
			&container.HostConfig{
				// TODO (rluckie): shouldn't default to privileged -- issue with ssh cmd
				Privileged: true,
				PortBindings: nat.PortMap{
					"2022": []nat.PortBinding{
						{
							HostPort: "2022",
						},
					},
				},
				Binds: binds,
			},
			nil,
			viper.Get("docker.hostname").(string))
		if err != nil {
			logger.WithField("error", err).Fatal("Failed to create KDK container")
		}
		if err := client.ContainerStart(ctx, containerCreateResp.ID, types.ContainerStartOptions{}); err != nil {
			logger.WithField("error", err).Fatal("Failed to start KDK container")
		}
		logger.Info("Successfully started KDK container")
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
}
