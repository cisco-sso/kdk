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
	"github.com/docker/docker/api/types"
	log "github.com/sirupsen/logrus"
)

func Kubesync(cfg KdkEnvConfig) {

	// Check if KDK container is running
	kdkRunning := false

	containers, err := cfg.DockerClient.ContainerList(cfg.Ctx, types.ContainerListOptions{All: true})

	if err != nil {
		log.WithField("error", err).Fatal("Failed to list docker containers")
	}

	for _, container := range containers {
		for _, name := range container.Names {
			if name == "/"+cfg.ConfigFile.AppConfig.Name {
				if container.State == "running" {
					kdkRunning = true
					break
				}
			}
		}
	}

	// if KDK container is not running, start it and provision KDK user
	if !kdkRunning {
		log.Error("KDK is not currently running.")
	}

	kubeconfigHostPath := cfg.Home() + "/.kube/config"
	kubeconfigKDKPath := ".kube/docker-for-desktop.example.org"

	// Create ~/.kube directory inside KDK if it doesn't already exist.
	remoteCommand := "mkdir -p ~/.kubetest"
	if err = cfg.Exec(remoteCommand); err != nil {
		log.WithField("error", err).Fatal("Failed to mkdir in KDK container.")
	}

	// Sync default KUBECONFIG to KDK
	if err = cfg.SCPTo(kubeconfigHostPath, kubeconfigKDKPath); err != nil {
		log.WithField("error", err).Fatal("Failed to scp to KDK container.")
	}

	// Tune Docker for Desktop's Kubernetes API hostname in KUBECONFIG
	remoteCommand = "sed -i -e 's@localhost@host.docker.internal@g' -e 's@docker-for-desktop.*@docker-for-desktop.example.org@g' " + kubeconfigKDKPath
	if err = cfg.Exec(remoteCommand); err != nil {
		log.WithField("error", err).Fatal("Failed to transform KUBECONFIG in KDK container.")
	}
	log.Info("Docker for Desktop KUBECONFIG synchronized to KDK.")
}
