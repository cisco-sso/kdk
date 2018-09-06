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
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/cisco-sso/kdk/pkg/prompt"
	"github.com/mholt/archiver"
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
	kdkBinDir := "/usr/local/bin"
	kdkBinName := "kdk"
	if runtime.GOOS == "windows" {
		kdkBinName = kdkBinName + ".exe"
	}
	kdkBinPath := filepath.Join(kdkBinDir, kdkBinName)
	logger.Infof("Update KDK binary?")
	p := prompt.Prompt{
		Text:     "Continue? [y/n] ",
		Loop:     true,
		Validate: prompt.ValidateYorN,
	}
	if result, err := p.Run(); err != nil || result == "n" {
		logger.Error("KDK binary update canceled or invalid input.")
	} else {
		// get latest release
		resp, err := http.Get("https://api.github.com/repos/cisco-sso/kdk/releases/latest")
		if err != nil {
			logger.WithField("error", err).Fatal("Failed to update KDK binary")
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		// get latest tag name
		op, _ := jq.Parse(".tag_name")
		version, _ := op.Apply([]byte(string(body)))
		versionStr := strings.Replace(string(version), "\"", "", -1)

		downloadBaseName := "kdk-" + versionStr + "-" + runtime.GOOS + "-" + runtime.GOARCH
		baseUrl := "https://github.com/cisco-sso/kdk/releases/download/"
		downloadLink := baseUrl + versionStr + "/" + downloadBaseName + ".tar.gz"
		downloadDir := filepath.Join("/tmp", downloadBaseName)
		downloadPath := filepath.Join(downloadDir, downloadBaseName+".tar.gz")

		// create downloadDir if it doesn't exist
		if _, err := os.Stat(downloadDir); os.IsNotExist(err) {
			os.Mkdir(downloadDir, 0700)
		}

		// create downloadPath file
		output, err := os.Create(downloadPath)
		if err != nil {
			logger.WithField("error", err).Fatal("Failed to download KDK tgz")
		}
		defer output.Close()

		// download latest release for arch/os to temp dir
		resp, err = http.Get(downloadLink)
		if err != nil {
			logger.WithField("error", err).Fatalf("Failed to download KDK tgz from %s", downloadLink)
		}
		defer resp.Body.Close()

		_, err = io.Copy(output, resp.Body)

		if err != nil {
			logger.WithField("error", err).Fatalf("Failed to write KDK binary to %s", kdkBinPath)
		}

		// extract tgz
		err = archiver.TarGz.Open(downloadPath, downloadDir)
		if err != nil {
			logger.WithField("error", err).Fatal("Failed to extract KDK tgz")
		}

		binSrcPath := filepath.Join(downloadDir, kdkBinName)

		binDestPath := filepath.Join("/usr/local/bin", kdkBinName)

		src, err := os.Open(binSrcPath)
		if err != nil {
			logger.WithField("error", err).Fatalf("Failed to read KDK binary @", binSrcPath)
		}
		defer src.Close()

		// copy bin to appropriate location
		dest, err := os.OpenFile(binDestPath, os.O_RDWR|os.O_CREATE, 0700)
		if err != nil {
			logger.WithField("error", err).Fatalf("Failed to write KDK binary @", binDestPath)
		}
		defer dest.Close()

		_, err = io.Copy(dest, src)
		if err != nil {
			logger.WithField("error", err).Fatalf("Failed to write KDK binary @", binDestPath)
		}

		// remove temp dir
		os.RemoveAll(downloadDir)
	}

}
