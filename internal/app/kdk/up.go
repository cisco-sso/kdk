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

package kdk

import (
	"context"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/spf13/viper"
)

func Up(ctx context.Context, dockerClient *client.Client, imageCoordinates string, logger logrus.Entry) error {
	var binds []string

	for _, bind := range viper.Get("docker.binds").([]interface{}) {
		configBind := make(map[string]string)
		for key, value := range bind.(map[interface{}]interface{}) {
			configBind[key.(string)] = value.(string)
		}
		binds = append(binds, configBind["source"]+":"+configBind["target"])
	}

	containerCreateResp, err := dockerClient.ContainerCreate(
		ctx,
		&container.Config{
			Hostname: viper.Get("docker.hostname").(string),
			Image:    imageCoordinates,
			Tty:      true,
			Env: []string{
				"KDK_USERNAME=" + viper.Get("docker.environment.KDK_USERNAME").(string),
				"KDK_SHELL=" + viper.Get("docker.environment.KDK_SHELL").(string),
				"KDK_DOTFILES_REPO=" + viper.Get("docker.environment.KDK_DOTFILES_REPO").(string),
			},
			ExposedPorts: nat.PortSet{
				"2022/tcp": struct{}{},
			},
		},
		&container.HostConfig{
			// TODO (rluckie): shouldn't default to privileged -- issue with ssh cmd
			Privileged: true,
			PortBindings: nat.PortMap{
				"2022/tcp": []nat.PortBinding{
					{
						HostPort: Port,
					},
				},
			},
			Binds: binds,
		},
		nil,
		viper.Get("docker.hostname").(string))
	if err != nil {
		logger.WithField("error", err).Fatal("Failed to create KDK container")
		return err
	}
	if err := DockerClient.ContainerStart(ctx, containerCreateResp.ID, types.ContainerStartOptions{}); err != nil {
		logger.WithField("error", err).Fatal("Failed to start KDK container")
		return err
	}
	logger.Info("Successfully started KDK container")
	return nil
}
