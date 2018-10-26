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
	"runtime"

	"github.com/cisco-sso/kdk/pkg/keybase"
	"github.com/cisco-sso/kdk/pkg/prompt"
	"github.com/docker/docker/api/types"
	log "github.com/sirupsen/logrus"
)

func Up(cfg KdkEnvConfig) (err error) {

	if runtime.GOOS == "windows" {
		if err := keybase.StartMirror(cfg.ConfigRootDir()); err != nil {
			log.WithField("error", err).Fatal("Failed to start keybase mirror")
			return err
		}
	}

	containers, err := cfg.DockerClient.ContainerList(cfg.Ctx, types.ContainerListOptions{All: true})
	if err != nil {
		log.WithField("error", err).Fatal("Failed to list docker containers")
	}
	for _, container := range containers {
		for _, name := range container.Names {
			if name == "/"+cfg.ConfigFile.AppConfig.Name {
				if container.State == "exited" {
					log.Infof("An exited KDK container exists")
					p := prompt.Prompt{
						Text:     "Restart exited KDK container? [y/n] ",
						Loop:     true,
						Validate: prompt.ValidateYorN,
					}
					if result, err := p.Run(); err == nil && result == "y" {
						log.Info("Restarting exited KDK container")
						containerStart(cfg, container.ID)
						return nil
					} else {
						p := prompt.Prompt{
							Text:     "Delete exited KDK container? [y/n] ",
							Loop:     true,
							Validate: prompt.ValidateYorN,
						}
						if result, err := p.Run(); err != nil || result == "n" {
							log.Fatal("KDK exited image deletion canceled or invalid input.")
						}
						log.Info("Removing exited KDK container")
						if err := cfg.DockerClient.ContainerRemove(cfg.Ctx, container.ID, types.ContainerRemoveOptions{Force: true}); err != nil {
							log.WithField("error", err).Fatalf("Failed to remove exited KDK container [%s]", container.ImageID)
						}
					}
				}
			}
		}
	}
	containerID, err := containerCreate(cfg)
	if err != nil {
		log.WithField("error", err).Fatal("Failed to create KDK container")
	} else {
		if err := containerStart(cfg, containerID); err != nil {
			log.WithField("error", err).Fatal("Failed to start KDK container")
		}
	}
	return nil
}
func containerCreate(cfg KdkEnvConfig) (string, error) {
	containerCreateResp, err := cfg.DockerClient.ContainerCreate(
		cfg.Ctx,
		cfg.ConfigFile.ContainerConfig,
		cfg.ConfigFile.HostConfig,
		nil,
		cfg.ConfigFile.AppConfig.Name,
	)
	if err != nil {
		return "", err
	}
	return containerCreateResp.ID, nil
}

func containerStart(cfg KdkEnvConfig, containerID string) (err error) {
	if err := cfg.DockerClient.ContainerStart(cfg.Ctx, containerID, types.ContainerStartOptions{}); err != nil {
		return err
	}
	log.Info("Successfully started KDK container")
	return nil
}
