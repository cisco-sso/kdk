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
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/cisco-sso/kdk/pkg/prompt"
	"github.com/cisco-sso/kdk/pkg/utils"
	"github.com/docker/docker/api/types"
)

func Prune(cfg KdkEnvConfig, debug bool) error {
	log.Info("Starting Prune...")

	var (
		imageIds                 []string
		runningContainerImageIds []string
		staleImageIds            []string
	)

	// Get containers
	containers, err := cfg.DockerClient.ContainerList(cfg.Ctx, types.ContainerListOptions{})
	if err != nil {
		log.WithField("error", err).Fatal("Failed to list docker containers")
	}
	// Get images
	images, err := cfg.DockerClient.ImageList(cfg.Ctx, types.ImageListOptions{})
	if err != nil {
		log.WithField("error", err).Fatal("Failed to list docker images")
	}

	// Iterate through containers and track running container imageIds
	for _, container := range containers {
		if strings.Contains(container.Status, "Up") {
			runningContainerImageIds = append(runningContainerImageIds, container.ImageID)
		}
	}

	// Iterate through images and track images that have a `kdk` label key
	for _, image := range images {
		for key := range image.Labels {
			if key == "kdk" {
				imageIds = append(imageIds, image.ID)
				break
			}
		}
	}

	// iterate through imageIds and add imageIds that are NOT associated with currently running containers
	for imageId := range imageIds {
		if utils.Contains(runningContainerImageIds, imageIds[imageId]) {
		} else {
			staleImageIds = append(staleImageIds, imageIds[imageId])
		}
	}

	if len(staleImageIds) > 0 {
		// iterate through staleImageIds, prmpt user to confirm deletion
		for staleImage := range staleImageIds {
			targetImage := staleImageIds[staleImage]
			log.Infof("Delete stale KDK image [%s]?", targetImage)
			prmpt := prompt.Prompt{
				Text:     "Continue? [y/n] ",
				Loop:     true,
				Validate: prompt.ValidateYorN,
			}
			if result, err := prmpt.Run(); err != nil || result == "n" {
				log.Error("KDK stale image deletion canceled or invalid input.")
				return err
			}
			if _, err := cfg.DockerClient.ImageRemove(cfg.Ctx, targetImage, types.ImageRemoveOptions{Force: true, PruneChildren: true}); err != nil {
				log.WithField("error", err).Fatalf("Failed to prune KDK image [%s]", targetImage)
				return err
			} else {
				log.Infof("Deleted stale KDK image [%s]", targetImage)
			}
		}
	} else {
		log.Infof("No stale KDK images to delete")
	}
	return nil
}
