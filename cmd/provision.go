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

package cmd

import (
	"github.com/Sirupsen/logrus"
	"github.com/codeskyblue/go-sh"
	"github.com/spf13/cobra"
)

func Provision() (out []byte, err error) {
	out, err = sh.Command("docker", "exec", "kdk", "/usr/local/bin/provision-user").Output()
	return out, err
}

var provisionCmd = &cobra.Command{
	Use:   "provision",
	Short: "Provision KDK user",
	Long:  `Provision KDK user`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := logrus.New().WithField("command", "provision")

		logger.Info("Provisioning KDK user. This may take a moment...")

		// TODO (rluckie): replace sh docker sdk
		if out, err := Provision(); err != nil {
			logger.WithField("error", err).Fatal("Failed to provision KDK user", err, out)
		}
		logger.Info("Successfully provisioned KDK user")
	},
}

func init() {
	rootCmd.AddCommand(provisionCmd)
}
