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
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/cisco-sso/kdk/pkg/prompt"
	"github.com/savaki/jq"
)

func Update(cfg KdkEnvConfig, logger logrus.Entry) error {
	logger.Info("Updating KDK")

	updateImage(cfg, logger)
	updateBin(logger)
	return nil
}

func updateImage(cfg KdkEnvConfig, logger logrus.Entry) {
	logger.Infof("Update KDK image?")
	p := prompt.Prompt{
		Text:     "Continue? [y/n] ",
		Loop:     true,
		Validate: prompt.ValidateYorN,
	}
	if result, err := p.Run(); err != nil || result == "n" {
		logger.Error("KDK image update canceled or invalid input.")
	} else {
		Pull(cfg)
	}
}
func updateBin(logger logrus.Entry) {
	kdkBinPath := "/usr/local/bin/kdk"
	logger.Infof("Update KDK binary?")
	p := prompt.Prompt{
		Text:     "Continue? [y/n] ",
		Loop:     true,
		Validate: prompt.ValidateYorN,
	}
	if result, err := p.Run(); err != nil || result == "n" {
		logger.Error("KDK binary update canceled or invalid input.")
	} else {
		resp, err := http.Get("https://api.github.com/repos/cisco-sso/kdk/releases/latest")
		if err != nil {
			logger.WithField("error", err).Fatal("Failed to update KDK binary")
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		op, _ := jq.Parse(".tag_name")
		version, _ := op.Apply([]byte(string(body)))
		versionStr := strings.Replace(string(version), "\"", "", -1)

		baseUrl := "https://github.com/cisco-sso/kdk/releases/download/"
		downloadLink := baseUrl + versionStr + "/kdk-" + versionStr + "-" + runtime.GOOS + "-" + runtime.GOARCH + ".tar.gz"

		if runtime.GOOS == "windows" {
			kdkBinPath = kdkBinPath + ".exe"
		}

		output, err := os.Create(kdkBinPath)
		if err != nil {
			logger.WithField("error", err).Fatalf("Failed create KDK bin @ %s", kdkBinPath)
		}
		defer output.Close()

		resp, err = http.Get(downloadLink)
		if err != nil {
			logger.WithField("error", err).Fatalf("Failed to download KDK binary from %s", downloadLink)
		}
		defer resp.Body.Close()

		_, err = io.Copy(output, resp.Body)

		if err != nil {
			logger.WithField("error", err).Fatalf("Failed to write KDK binary to %s", kdkBinPath)
		}

	}
}
