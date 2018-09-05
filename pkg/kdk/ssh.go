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

	"github.com/Sirupsen/logrus"
	"github.com/codeskyblue/go-sh"
)

func Ssh(cfg KdkEnvConfig, logger logrus.Entry) error {
	logger.Info("Connecting to KDK container")

	connectionString := cfg.User() + "@localhost"
	commandString := fmt.Sprintf("ssh %s -A -p %s -i %s -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null", connectionString, cfg.ConfigFile.AppConfig.Port, cfg.PrivateKeyPath())
	if cfg.ConfigFile.AppConfig.Debug {
		logger.Infof("executing ssh command: %s", commandString)
	}
	commandMap := strings.Split(commandString, " ")
	if err := sh.Command(commandMap[0], commandMap[1:]).SetStdin(os.Stdin).Run(); err != nil {
		logger.WithField("error", err).Fatal("Failed to ssh to KDK container.")
	}

	logger.Info("KDK session exited")
	return nil
}
