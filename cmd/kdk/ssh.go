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
	"github.com/cisco-sso/kdk/pkg/kdk"
	"github.com/spf13/cobra"
)

var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "Connect to running KDK container via ssh",
	Long:  `Connect to running KDK container via ssh`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := logrus.New().WithField("command", "ssh")

		kdk.Ssh(CurrentKdkEnvConfig, *logger)

	},
}

func init() {
	rootCmd.AddCommand(sshCmd)
}
