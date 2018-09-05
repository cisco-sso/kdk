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
	"github.com/Sirupsen/logrus"
	"github.com/codeskyblue/go-sh"
)

func Provision(cfg KdkEnvConfig, logger logrus.Entry) error {
	// TODO (rluckie): replace sh docker sdk
	logger.Info("Starting KDK user provisioning. This may take a moment.  Hang tight...")
	if _, err := sh.Command("docker", "exec", cfg.ConfigFile.AppConfig.Name, "/usr/local/bin/provision-user").Output(); err != nil {
		logger.WithField("error", err).Fatal("Failed to provision KDK user.")
		return err
	} else {
		logger.Info("Completed KDK user provisioning.")
		return nil
	}
}
