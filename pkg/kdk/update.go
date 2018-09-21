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
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/cisco-sso/kdk/pkg/utils"
	"github.com/docker/docker/api/types"
	"github.com/ghodss/yaml"
	"github.com/mholt/archiver"
	"github.com/savaki/jq"
)

var (
	latestReleaseVersion = getLatestReleaseVersion()
)

func WarnIfUpdateAvailable(cfg *KdkEnvConfig) {
	if latestReleaseVersion == "" {
		return
	}

	if needsUpdateBin() || needsUpdateImage(cfg) || needsUpdateConfig(cfg) {
		log.Warn("Upgrade Available\n" + strings.Join([]string{
			"***************************************",
			"The installed KDK version is out of date",
			"  Current: " + Version,
			"  Latest : " + latestReleaseVersion,
			"",
			"Please upgrade the KDK with the commands:",
			"  kdk update",
			"  kdk destroy",
			"  kdk ssh",
			"***************************************"}, "\n"))
	}
	return
}

// check if kdk bin needs to be updated
func needsUpdateBin() bool {
	return Version != latestReleaseVersion
}

// check if kdk image needs to be updated
func needsUpdateImage(cfg *KdkEnvConfig) bool {
	kdkImages := getKdkImages(cfg)

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

// check if kdk config needs to be updated
func needsUpdateConfig(cfg *KdkEnvConfig) bool {
	if cfg.ConfigFile.AppConfig.ImageTag != latestReleaseVersion ||
		cfg.ConfigFile.ContainerConfig.Image != cfg.ImageCoordinates() ||
		cfg.ConfigFile.ContainerConfig.Labels["kdk"] != latestReleaseVersion {
		return true
	}
	return false
}

func Update(cfg *KdkEnvConfig) {
	if latestReleaseVersion == "" {
		log.Warn("Upgrade Unavailable.  Unable to fetch latest version")
	}

	if !(needsUpdateBin() || needsUpdateImage(cfg) || needsUpdateConfig(cfg)) {
		log.Warn("Upgrade Unavailable.  Already at latest versions")
	}
	log.Info("Upgrade Available\n" + strings.Join([]string{
		"***************************************",
		"The installed KDK version is out of date",
		"  Current: " + Version,
		"  Latest : " + latestReleaseVersion,
		"",
		"Upgrading the KDK binary, image, and config to latest version",
		"",
		"After upgrade, restart the the kdk with the commands:",
		"  kdk update",
		"  kdk destroy",
		"  kdk ssh",
		"***************************************"}, "\n"))

	if needsUpdateBin() {
		log.Info("Updating KDK binary")
		err := updateBin()
		if err != nil {
			log.WithField("error", err).Fatal("Failed to update KDK bin")
		}
	} else {
		log.Info("Updating KDK binary skipped: Already at latest version")
	}

	if needsUpdateImage(cfg) {
		log.Info("Updating KDK image")
		err := updateImage(cfg)
		if err != nil {
			log.WithField("error", err).Fatal("Failed to update KDK image")
		}
	} else {
		log.Info("Updating KDK image skipped: Already at latest version")
	}

	if needsUpdateConfig(cfg) {
		log.Info("Updating KDK config")
		err := updateConfig(cfg)
		if err != nil {
			log.WithField("error", err).Fatal("Failed to update KDK config")
		}
	} else {
		log.Info("Updating KDK config skipped: Already at latest version")
	}
	return
}

// update kdk bin
func updateBin() error {
	// TODO: This must be fixed if the binary is ever to run in powershell instead of bash

	// Figure out the binary path
	//   TODO: Don't make assumptions on binary path, like /usr/local/bin
	kdkBinName := "kdk"
	if runtime.GOOS == "windows" {
		kdkBinName = kdkBinName + ".exe"
	}
	kdkBinPath := filepath.Join("usr", "local", "bin", kdkBinName)

	// Calculate the download urls and tmp locations
	baseUrl := "https://github.com/cisco-sso/kdk/releases/download/"
	downloadBaseName := "kdk-" + latestReleaseVersion + "-" + runtime.GOOS + "-" + runtime.GOARCH
	downloadLink := baseUrl + latestReleaseVersion + "/" + downloadBaseName + ".tar.gz"
	downloadDir := filepath.Join("/tmp", downloadBaseName)
	downloadPath := filepath.Join(downloadDir, downloadBaseName+".tar.gz")

	// create downloadDir if it doesn't exist
	if _, err := os.Stat(downloadDir); os.IsNotExist(err) {
		err = os.MkdirAll(downloadDir, 0700)
		if err != nil {
			log.WithField("error", err).Error("Failed to download dir")
			return err
		}
	}

	// create downloadPath file
	output, err := os.Create(downloadPath)
	if err != nil {
		log.WithField("error", err).Error("Failed to download KDK tgz")
		return err
	}
	defer output.Close()

	// download latest release for arch/os to temp dir
	resp, err := http.Get(downloadLink)
	if err != nil {
		log.WithField("error", err).Errorf("Failed to download KDK tgz from %s", downloadLink)
		return err
	}
	defer resp.Body.Close()

	// write the kdk binary
	_, err = io.Copy(output, resp.Body)
	if err != nil {
		log.WithField("error", err).Errorf("Failed to write KDK binary to %s", kdkBinPath)
		return err
	}

	// extract tgz
	err = archiver.TarGz.Open(downloadPath, downloadDir)
	if err != nil {
		log.WithField("error", err).Fatal("Failed to extract KDK tgz")
		return err
	}

	// copy bin to appropriate location
	binSrcPath := filepath.Join(downloadDir, kdkBinName)
	binDestPath := filepath.Join("/usr/local/bin", kdkBinName)

	//   open the bin source
	src, err := os.Open(binSrcPath)
	if err != nil {
		log.WithField("error", err).Errorf("Failed to read KDK binary @", binSrcPath)
		return err
	}
	defer src.Close()

	//   open the bin destination
	dest, err := os.OpenFile(binDestPath, os.O_RDWR|os.O_CREATE, 0700)
	if err != nil {
		log.WithField("error", err).Errorf("Failed to write KDK binary @", binDestPath)
		return err
	}
	defer dest.Close()

	//   Do the copy
	_, err = io.Copy(dest, src)
	if err != nil {
		log.WithField("error", err).Errorf("Failed to write KDK binary @", binDestPath)
		return err
	}

	// remove temp dir
	err = os.RemoveAll(downloadDir)
	if err != nil {
		log.WithField("error", err).Errorf("Failed to remove download directory @", downloadDir)
		return err
	}

	return nil
}

// update kdk image
func updateImage(cfg *KdkEnvConfig) error {
	Pull(cfg)
	return nil
}

// update kdk config
func updateConfig(cfg *KdkEnvConfig) error {
	cfg.ConfigFile.AppConfig.ImageTag = latestReleaseVersion
	cfg.ConfigFile.ContainerConfig.Labels["kdk"] = latestReleaseVersion
	cfg.ConfigFile.ContainerConfig.Image = cfg.ImageCoordinates()

	y, err := yaml.Marshal(cfg.ConfigFile)
	if err != nil {
		log.WithField("error", err).Error("Failed to create YAML string of configuration")
		return err
	}

	err = ioutil.WriteFile(cfg.ConfigPath(), y, 0600)
	if err != nil {
		log.WithField("error", err).Error("Failed to write new config file")
		return err
	}

	return nil
}

func getLatestReleaseVersion() string {
	client := http.Client{
		Timeout: time.Duration(1 * time.Second),
	}

	// Fetch the informational json blob
	resp, err := client.Get("https://api.github.com/repos/cisco-sso/kdk/releases/latest")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	// Parse the latest tag name
	op, _ := jq.Parse(".tag_name")
	version, err := op.Apply([]byte(string(body)))
	if err != nil {
		log.WithField("error", err).Error("Failed to check latest release version", err)
		return ""
	}

	// Remove the bookend quotes
	versionStr := strings.Replace(string(version), "\"", "", -1)
	return versionStr
}

// get kdk docker image on host
func getKdkImages(cfg *KdkEnvConfig) (out []types.ImageSummary) {
	var kdkImages []types.ImageSummary
	images, err := cfg.DockerClient.ImageList(cfg.Ctx, types.ImageListOptions{All: true})
	if err != nil {
		log.WithField("error", err).Fatal("Failed to list docker images")
	}
	for _, image := range images {
		for key := range image.Labels {
			// The "kdk" label is set at the top of the Dockerfile
			if key == "kdk" {
				kdkImages = append(kdkImages, image)
				break
			}
		}
	}
	return kdkImages
}
