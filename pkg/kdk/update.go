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
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
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

	// Add a sudo hint for unix-based OS's
	sudo := ""
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		sudo = "sudo " // trailing space is intentional
	}
	if needsUpdateBin() || needsUpdateImage(cfg) || needsUpdateConfig(cfg) {
		log.Warn("Upgrade Available\n" + strings.Join([]string{
			"***************************************",
			"Some KDK components are out of date",
			"  Latest Version:                      " + latestReleaseVersion,
			"  Binary Version:                      " + Version,
			"  Config Version:                      " + cfg.ConfigFile.AppConfig.ImageTag,
			"  Container Present at Config Version: " + strconv.FormatBool(!needsUpdateImage(cfg)),
			"",
			"Please upgrade the KDK with the commands:",
			"  " + sudo + "kdk update",
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
	return !hasKdkImageWithTag(cfg, latestReleaseVersion)
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
		return
	}

	if !(needsUpdateBin() || needsUpdateImage(cfg) || needsUpdateConfig(cfg)) {
		log.Warn("Upgrade Unavailable.  Already at latest versions")
		return
	}

	if needsUpdateBin() {
		if (runtime.GOOS == "linux" || runtime.GOOS == "darwin") && os.Geteuid() != 0 {
			log.Fatal("Please execute the update command with `sudo` or as the `root` user")
		}

		log.Info("Updating KDK binary")
		err := updateBin()
		if err != nil {
			log.WithField("error", err).Fatal("Failed to update KDK bin")
		}
	} else {
		log.Info("Updating KDK binary skipped: Already at latest version")
	}

	if needsUpdateImage(cfg) {
		log.Info("Updating KDK container image")
		err := pullImage(cfg, cfg.ConfigFile.AppConfig.ImageRepository+":"+latestReleaseVersion)
		if err != nil {
			log.WithField("error", err).Fatal("Failed to update KDK image")
		}
	} else {
		log.Info("Updating KDK container image skipped: Already at latest version")
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
	//// Construct all of the paths upfront

	// Calculate the download url
	baseUrl := "https://github.com/cisco-sso/kdk/releases/download"
	downloadBaseName := "kdk-" + latestReleaseVersion + "-" + runtime.GOOS + "-" + runtime.GOARCH
	downloadLink := baseUrl + "/" + latestReleaseVersion + "/" + downloadBaseName + ".tar.gz"

	// Calculate the temporary download and unpacking location
	tmpDir := filepath.Join(os.TempDir(), "kdk-install")
	tgzFile := filepath.Join(tmpDir, downloadBaseName+".tar.gz")

	kdkBinFile, _ := os.Executable() // this currently running binary will be overwritten
	kdkBinFileUnpacked := filepath.Join(tmpDir, filepath.Base(kdkBinFile))
	kdkBinFileTrash := filepath.Join(os.TempDir(), filepath.Base(kdkBinFile)+".old")
	// ^ Some filesystems do not allow replacing or deleting a currently
	//   running binary.  We'll move it out of the way instead of deletion

	log.WithField("file", kdkBinFile).Info("Bin File Location")

	//// download tgz file to the tmp dir
	err := downloadFile(downloadLink, tmpDir, tgzFile)
	if err != nil {
		log.WithField("error", err).WithField("file", tgzFile).WithField("url", downloadLink).Fatal("Failed to download file")
	}
	log.WithField("file", tgzFile).WithField("url", downloadLink).Info("Successfully downloaded file")

	// extract tgz
	err = archiver.TarGz.Open(tgzFile, tmpDir)
	if err != nil {
		log.WithField("error", err).WithField("file", tgzFile).Fatal("Failed to extract tgz")
		return err
	}
	log.WithField("file", tgzFile).Info("Successfully extracted tgz file")

	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		// copy the new file next to the org binary, so it is on the same partition/filesystem so that moves work
		err = copyFile(kdkBinFileUnpacked, kdkBinFile+".new")
		if err != nil {
			log.WithField("error", err).WithField("fileSrc", kdkBinFileUnpacked).WithField("fileDst", kdkBinFile+".new").Fatal("Failed to copy file")
			return err
		}
		log.WithField("fileSrc", kdkBinFileUnpacked).WithField("fileDst", kdkBinFile+".new").Info("Successfully copied file")

		// set the copy to be executable
		err = os.Chmod(kdkBinFile+".new", 0755)
		if err != nil {
			log.WithField("error", err).WithField("file", kdkBinFile+".new").Fatal("Failed to chmod file")
			return err
		}
		log.WithField("file", kdkBinFile+".new").Info("Successfully chmod'd file")

		// remove the original bin file
		err = os.Remove(kdkBinFile)
		if err != nil {
			log.WithField("error", err).WithField("file", kdkBinFile).Fatal("Failed to delete file")
		}
		log.WithField("file", kdkBinFile).Info("Successfully deleted file")

		// rename the new file to be the the executable file
		err = os.Rename(kdkBinFile+".new", kdkBinFile)
		if err != nil {
			log.WithField("error", err).WithField("fileSrc", kdkBinFile+".new").WithField("fileDst", kdkBinFile).Fatal("Failed to rename file")
		}
		log.WithField("fileSrc", kdkBinFile+".new").WithField("fileDst", kdkBinFile).Info("Successfully renamed file")
	} else if runtime.GOOS == "windows" {
		// rename the bin file to a trash location out of the way
		err = os.Rename(kdkBinFile, kdkBinFileTrash)
		if err != nil {
			log.WithField("error", err).WithField("fileSrc", kdkBinFile).WithField("fileDst", kdkBinFileTrash).Fatal("Failed to rename file")
		}
		log.WithField("fileSrc", kdkBinFile).WithField("fileDst", kdkBinFileTrash).Info("Successfully renamed file")

		// copy the new file next to the org binary, so it is on the same partition/filesystem
		err = copyFile(kdkBinFileUnpacked, kdkBinFile)
		if err != nil {
			log.WithField("error", err).WithField("fileSrc", kdkBinFileUnpacked).WithField("fileDst", kdkBinFile).Fatal("Failed to copy file")
			return err
		}
		log.WithField("fileSrc", kdkBinFileUnpacked).WithField("fileDst", kdkBinFile).Info("Successfully copied file")
	} else {
		log.Fatal("Unhandled code path")
	}

	// remove temp dir
	err = os.RemoveAll(tmpDir)
	if err != nil {
		log.WithField("error", err).WithField("dir", tmpDir).Error("Failed to remove directory")
		return err
	}
	log.WithField("dir", tmpDir).Info("Successfully removed directory")

	return nil
}

// update kdk image
func updateImage(cfg *KdkEnvConfig) error {
	Pull(cfg, true)
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
		Timeout: time.Duration(3 * time.Second),
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

func hasKdkImageWithTag(cfg *KdkEnvConfig, tagSearch string) bool {
	kdkImages := getKdkImages(cfg)

	for _, image := range kdkImages {
		var tags []string
		for _, tag := range image.RepoTags {
			imageTag := strings.Split(tag, ":")[1]
			tags = append(tags, imageTag)
		}
		if utils.Contains(tags, tagSearch) {
			return true
		}
	}
	return false
}

func copyFile(src, dst string) error {
	// Copy the src file to dst. Any existing file will be overwritten and will not
	// copy file attributes.
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func downloadFile(url string, dir string, file string) error {
	// Open the TCP stream for downloading download latest tgz release
	resp, err := http.Get(url)
	if err != nil {
		log.WithField("error", err).WithField("url", url).Error("Failed to reach Url")
		return err
	}
	defer resp.Body.Close()

	// create dir if it doesn't exist
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			log.WithField("error", err).WithField("dir", dir).Error("Failed to create directory")
			return err
		}
	}

	// open file file for eventual writing
	fd, err := os.Create(file)
	if err != nil {
		log.WithField("error", err).WithField("file", file).Error("Failed to open file for writing")
		return err
	}
	defer fd.Close()

	// write the TCP stream to the file
	_, err = io.Copy(fd, resp.Body)
	if err != nil {
		log.WithField("error", err).WithField("file", file).Errorf("Failed to write to file")
		return err
	}

	return nil
}
