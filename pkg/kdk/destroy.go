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
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/cisco-sso/kdk/pkg/prompt"
	"github.com/docker/docker/api/types"
)

func Destroy(cfg KdkEnvConfig) error {

	var containerIds []string

	containers, err := cfg.DockerClient.ContainerList(cfg.Ctx, types.ContainerListOptions{})

	if err != nil {
		log.WithField("error", err).Fatal("Failed to list docker containers")
	}
	for _, container := range containers {
		for _, name := range container.Names {
			if name == "/"+cfg.ConfigFile.AppConfig.Name {
				containerIds = append(containerIds, container.ID)
				break
			}
		}
	}
	if len(containerIds) > 0 {
		log.Info("Destroying KDK container(s)...")
		for _, containerId := range containerIds {
			fmt.Printf("Delete KDK container [%s][%v]\n", cfg.ConfigFile.AppConfig.Name, containerId[:8])
			prmpt := prompt.Prompt{
				Text:     "Continue? [y/n] ",
				Loop:     true,
				Validate: prompt.ValidateYorN,
			}
			if result, err := prmpt.Run(); err != nil || result == "n" {
				log.Error("KDK container deletion canceled or invalid input.")
				return nil
			}
			if err := cfg.DockerClient.ContainerRemove(cfg.Ctx, containerId, types.ContainerRemoveOptions{Force: true}); err != nil {
				log.WithField("error", err).Fatal("Failed to remove KDK container")
			}
		}
		log.Info("KDK destroy complete.")
	} else {
		log.Info("No KDK containers found. Nothing to destroy...")
	}
	return nil
}
