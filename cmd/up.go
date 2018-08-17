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

		client, err := client.NewEnvClient()
		if err != nil {
			logger.WithField("error", err).Fatal("Failed to create docker client")
		}

		ctx := context.Background()

		logger.Info("Pulling KDK image. This may take a moment...")
		if err := Pull(client, KdkImageCoordinates); err != nil {
			logger.WithField("error", err).Fatal("Failed to pull KDK image")
		}

		var binds []string

		for _, bind := range viper.Get("docker.binds").([]interface{}) {
			configBind := make(map[string]string)
			for key, value := range bind.(map[interface{}]interface{}) {
				configBind[key.(string)] = value.(string)
			}
			binds = append(binds, configBind["source"]+":"+configBind["target"])
		}

		containerCreateResp, err := client.ContainerCreate(
			ctx,
			&container.Config{
				Hostname: viper.Get("docker.hostname").(string),
				Image:    KdkImageCoordinates,
				Tty:      true,
				Env: []string{
					"KDK_USERNAME=" + viper.Get("docker.environment.KDK_USERNAME").(string),
					"KDK_SHELL=" + viper.Get("docker.environment.KDK_SHELL").(string),
					"KDK_DOTFILES_REPO=" + viper.Get("docker.environment.KDK_DOTFILES_REPO").(string),
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

		logger.Info("Provisioning KDK user. This may take a moment...")

		if out, err := Provision(); err != nil {
			logger.WithField("error", err).Fatal("Failed to provision KDK user", err, out)
		}
		logger.Info("Successfully started KDK container")
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
}
