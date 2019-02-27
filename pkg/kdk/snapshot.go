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
	"time"

	"github.com/docker/docker/api/types"
	log "github.com/sirupsen/logrus"
)

func Snapshot(cfg KdkEnvConfig) (string, error) {
	snapshotName := "ciscosso/kdk" + ":" + cfg.User() + "-" + cfg.ConfigFile.AppConfig.Name + "-" + time.Now().Format("20060102150405")
	_, err := cfg.DockerClient.ContainerCommit(cfg.Ctx, cfg.ConfigFile.AppConfig.Name, types.ContainerCommitOptions{Reference: snapshotName})
	if err != nil {
		log.WithField("error", err).Fatal("Failed to create snapshot of KDK container")
		return "", err
	}
	log.Info("Successfully created snapshot of KDK container.", snapshotName)
	return snapshotName, nil
}
