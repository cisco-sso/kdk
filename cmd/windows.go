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
	"strings"
//	"os"

	"github.com/Sirupsen/logrus"
//	"github.com/cisco-sso/kdk/internal/app/kdk"
	"github.com/spf13/cobra"
)

/*
  kdk windows service run-keybase-mirror
  kdk windows service install
  kdk windows service remove
*/

var windowsCmd = &cobra.Command{
	Use:   "windows",
	Short: "Windows subcommands for KDK container",
	Long:  `Windows subcommands KDK container`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := logrus.New().WithField("command", "windows")
		logger.Info("Print: " + strings.Join(args, " "))
	},
}

var windowsServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "Windows Service subcommands for KDK container",
	Long:  `Windows Service subcommands KDK container`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := logrus.New().WithField("command", "windows service")

		logger.Info("Print: " + strings.Join(args, " "))
	},
}

var windowsServiceInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Windows Service Install subcommands for KDK container",
	Long:  `Windows Service Install subcommands KDK container`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := logrus.New().WithField("command", "windows service install")

		logger.Info("Print: " + strings.Join(args, " "))
	},
}

var windowsServiceRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Windows Service Remove subcommands for KDK container",
	Long:  `Windows Service Remove subcommands KDK container`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := logrus.New().WithField("command", "windows service remove")

		logger.Info("Print: " + strings.Join(args, " "))
	},
}

var windowsServiceRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Windows Service Run subcommands for KDK container",
	Long:  `Windows Service Run subcommands KDK container`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := logrus.New().WithField("command", "windows service run")

		logger.Info("Print: " + strings.Join(args, " "))
	},
}

func init() {
	rootCmd.AddCommand(windowsCmd)
	windowsCmd.AddCommand(windowsServiceCmd)
	windowsServiceCmd.AddCommand(windowsServiceInstallCmd)
	windowsServiceCmd.AddCommand(windowsServiceRemoveCmd)
	windowsServiceCmd.AddCommand(windowsServiceRunCmd)

}
