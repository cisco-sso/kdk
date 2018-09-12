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
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/cisco-sso/kdk/pkg/prompt"
	"github.com/cisco-sso/kdk/pkg/utils"
	"github.com/docker/docker/api/types"
	"github.com/ghodss/yaml"
	"github.com/mholt/archiver"
	"github.com/savaki/jq"
)

var (
	latestReleaseVersion = getLatestReleaseVersion()
)

func Update(cfg KdkEnvConfig, debug bool, skipUpdate bool, logger logrus.Entry) error {
	if !skipUpdate {
		doUpdateConfig := updateConfigCheck(cfg)
		if doUpdateConfig {
			logger.Info("A newer version of the kdk binary or image is available")
			logger.Infof("Update will move from version %s -> %s", cfg.ConfigFile.AppConfig.ImageTag, latestReleaseVersion)
			logger.Info("If you would like to skip update, please use --skip-update flag")
			updateConfig(&cfg, debug, logger)
		} else {
			logger.Info("Config has not changed")
		}
		doImageUpdate := updateImageCheck(cfg, logger)
		if doImageUpdate {
			updateImage(cfg, debug, logger)
		} else {
			logger.Infof("Most recent KDK image has already been pulled [%s]", latestReleaseVersion)
		}
		doBinUpdate := updateBinCheck()
		if doBinUpdate {
			updateBin(logger)
		} else {
			logger.Infof("Using most recent version of KDK bin [%s]", latestReleaseVersion)
		}
		return nil
	} else {
		logger.Info("Skipped updated")
	}
	return nil
}

// get latest release version
func getLatestReleaseVersion() (out string) {
	resp, err := http.Get("https://api.github.com/repos/cisco-sso/kdk/releases/latest")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	// get latest tag name
	op, _ := jq.Parse(".tag_name")
	version, _ := op.Apply([]byte(string(body)))
	versionStr := strings.Replace(string(version), "\"", "", -1)
	return versionStr
}

// get kdk docker image on host
func getKdkImages(cfg KdkEnvConfig, logger logrus.Entry) (out []types.ImageSummary) {
	var kdkImages []types.ImageSummary
	images, err := cfg.DockerClient.ImageList(cfg.Ctx, types.ImageListOptions{All: true})
	if err != nil {
		logger.WithField("error", err).Fatal("Failed to list docker images")
	}
	for _, image := range images {
		for key := range image.Labels {
			if key == "kdk" {
				kdkImages = append(kdkImages, image)
				break
			}
		}
	}
	return kdkImages
}

// check if kdk image needs to be updated
func updateImageCheck(cfg KdkEnvConfig, logger logrus.Entry) (out bool) {
	kdkImages := getKdkImages(cfg, logger)

	for _, image := range kdkImages {
		var tags []string
		for _, tag := range image.RepoTags {
			imageTag := strings.Split(tag, ":")[1]
			tags = append(tags, imageTag)
		}
		if utils.Contains(tags, latestReleaseVersion) {
			return false
		}
	}
	return true
}

// update kdk image
func updateImage(cfg KdkEnvConfig, debug bool, logger logrus.Entry) {
	logger.Infof("Update KDK image?")
	p := prompt.Prompt{
		Text:     "Continue? [y/n] ",
		Loop:     true,
		Validate: prompt.ValidateYorN,
	}
	if result, err := p.Run(); err != nil || result == "n" {
		logger.Error("KDK image update canceled or invalid input.")
	} else {
		Pull(cfg, debug)
	}
}

// check if kdk bin needs to be updated
func updateBinCheck() (out bool) {
	if Version == latestReleaseVersion {
		return false
	}
	return true
}

// update kdk bin
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

		downloadBaseName := "kdk-" + latestReleaseVersion + "-" + runtime.GOOS + "-" + runtime.GOARCH
		baseUrl := "https://github.com/cisco-sso/kdk/releases/download/"
		downloadLink := baseUrl + latestReleaseVersion + "/" + downloadBaseName + ".tar.gz"
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
		resp, err := http.Get(downloadLink)
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

// check if kdk config needs to be updated
func updateConfigCheck(cfg KdkEnvConfig) (out bool) {
	if cfg.ConfigFile.AppConfig.ImageTag != latestReleaseVersion ||
		cfg.ConfigFile.ContainerConfig.Labels["kdk"] != latestReleaseVersion ||
		cfg.ConfigFile.ContainerConfig.Image != cfg.ImageCoordinates() {
		return true
	}
	return false
}

// update kdk config
func updateConfig(cfg *KdkEnvConfig, debug bool, logger logrus.Entry) (err error) {
	cfg.ConfigFile.AppConfig.ImageTag = latestReleaseVersion
	cfg.ConfigFile.ContainerConfig.Labels["kdk"] = latestReleaseVersion
	cfg.ConfigFile.ContainerConfig.Image = cfg.ImageCoordinates()

	y, err := yaml.Marshal(cfg.ConfigFile)
	if err != nil {
		logger.Fatal("Failed to create YAML string of configuration", err)
	}
	p := prompt.Prompt{
		Text:     fmt.Sprintf("Update config file [%s]? [y/n]", cfg.ConfigPath()),
		Loop:     true,
		Validate: prompt.ValidateYorN,
	}
	if result, err := p.Run(); err == nil && result == "y" {
		logger.Info("Updating KDK config")
		ioutil.WriteFile(cfg.ConfigPath(), y, 0600)
	} else {
		logger.Fatal("Existing KDK config not overwritten")
	}
	return nil
}
