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
	"time"
	"strings"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Create a snapshot of a running KDK container",
	Long:  `Create a snapshot of a running KDK container`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := logrus.New().WithField("command", "pull")

		ctx := context.Background()
		client, err := client.NewEnvClient()

		if err != nil {
			logger.WithField("error", err).Fatal("Failed to create docker client")
		}

		snapshotName := strings.Join([]string{"kdk", strconv.Itoa(int(time.Now().UnixNano()))}, "-")

		_, err = client.ContainerCommit(ctx, "kdk", types.ContainerCommitOptions{Reference: snapshotName})
		if err != nil {
			logger.WithField("error", err).Fatal("Failed to create snapshot of KDK container")
		}
		logger.Info("Successfully created snapshot of KDK container.", snapshotName)

	},
}

func init() {
	rootCmd.AddCommand(snapshotCmd)
}
