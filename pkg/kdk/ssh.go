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
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/codeskyblue/go-sh"
	"github.com/docker/docker/api/types"
)

func Ssh(cfg KdkEnvConfig) {

	log.Info("Connecting to KDK container")

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
		log.Info("KDK is not currently running.  Starting...")
		Pull(&cfg, false)
		Up(cfg)
		Provision(cfg)
	}

	// Build socksString
	var socksString string
	if cfg.ConfigFile.AppConfig.SocksPort != "" {
		socksString = "-D " + cfg.ConfigFile.AppConfig.SocksPort
	}

	// connect to KDK container via ssh
	connectionString := cfg.User() + "@localhost"
	commandString := fmt.Sprintf("ssh %s -A -p %s -i %s -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null %s", connectionString, cfg.ConfigFile.AppConfig.Port, cfg.PrivateKeyPath(), socksString)
	log.Infof("executing ssh command: %s", commandString)
	commandMap := strings.Split(commandString, " ")
	if err := sh.Command(commandMap[0], commandMap[1:]).SetStdin(os.Stdin).Run(); err != nil {
		log.WithField("error", err).Fatal("Failed to ssh to KDK container.")
	}

	log.Info("KDK session exited")
}
