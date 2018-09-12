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
	"io"
	"io/ioutil"
	"os"

	"github.com/docker/docker/api/types"
)

func Pull(cfg KdkEnvConfig, debug bool) error {
	out, err := cfg.DockerClient.ImagePull(cfg.Ctx, cfg.ImageCoordinates(), types.ImagePullOptions{})
	defer out.Close()
	if debug {
		io.Copy(os.Stdout, out)
	} else {
		io.Copy(ioutil.Discard, out)
	}
	return err
}
