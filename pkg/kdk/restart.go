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
	"github.com/ghodss/yaml"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"strings"
)

func Restart(cfg KdkEnvConfig) {

	log.Info("Restarting KDK container")

	// Create snapshot of running KDK container
	snapshotName, err := Snapshot(cfg)
	if err != nil {
		log.WithField("err", err).Fatal("Failed to create KDK image snapshot")
	}

	// Destroy running KDK container
	Destroy(cfg)

	// Save config with snapshot image tag
	cfg.ConfigFile.AppConfig.ImageTag = strings.Split(snapshotName, ":")[1]
	cfg.ConfigFile.ContainerConfig.Image = snapshotName
	y, err := yaml.Marshal(cfg.ConfigFile)
	if err != nil {
		log.WithField("error", err).Error("Failed to create YAML string of configuration")
	}

	err = ioutil.WriteFile(cfg.ConfigPath(), y, 0600)
	if err != nil {
		log.WithField("error", err).Error("Failed to write new config file")
	}
	// Start KDK container with snapshot image
	cfg.Start()
	log.Info("KDK container restarted")
}
