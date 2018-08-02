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
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/sirupsen/logrus"
	"github.com/codeskyblue/go-sh"
	"github.com/spf13/viper"
)

var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "Connect to running KDK container via ssh",
	Long: `Connect to running KDK container via ssh`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := logrus.New().WithField("command", "ssh")

		connectionString := strings.Join([]string{viper.Get("docker.environment.KDK_USERNAME").(string), "localhost"}, "@")

		logger.Info("Connecting to KDK container")
		sh.Command("ssh", connectionString, "-A", "-p", "2022", "-i", "~/.kdk/ssh/id_rsa", "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null").SetStdin(os.Stdin).Run()
		logger.Info("KDK session exited")
	},
}

func init() {
	rootCmd.AddCommand(sshCmd)
}
