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
	"github.com/cisco-sso/kdk/pkg/kdk"
	"github.com/spf13/cobra"
)

var kubesyncCmd = &cobra.Command{
	Use:   `kubesync`,
	Short: "Sync default KUBECONFIG to KDK",
	Long:  "Sync default KUBECONFIG to KDK and tune Docker Kubernetes API hostname",
	Run: func(cmd *cobra.Command, args []string) {
		kdk.Kubesync(CurrentKdkEnvConfig)
	},
}

func init() {
	rootCmd.AddCommand(kubesyncCmd)
}
