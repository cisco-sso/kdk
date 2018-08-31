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
	"runtime"

	"github.com/Sirupsen/logrus"
	"github.com/cisco-sso/kdk/internal/pkg/utils"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func Up(
	ctx context.Context,
	dockerClient *client.Client,
	cfg *kdkConfig,
	logger logrus.Entry) error {
	if runtime.GOOS == "windows" {
		if err := utils.KeybaseStartMirror(ConfigDir); err != nil {
			logger.WithField("error", err).Fatal("Failed to start keybase mirror")
			return err
		}
	}
	containerCreateResp, err := dockerClient.ContainerCreate(
		ctx,
		&cfg.ContainerConfig,
		&cfg.HostConfig,
		nil,
		cfg.AppConfig.Name,
	)
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
