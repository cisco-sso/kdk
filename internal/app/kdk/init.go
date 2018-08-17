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
	"context"
	"io/ioutil"
	"os"
	"os/user"
	"path"

	"github.com/Sirupsen/logrus"
	"github.com/cisco-sso/kdk/internal/pkg/utils"
	"github.com/docker/docker/client"
	"github.com/manifoldco/promptui"
)

var (
	ConfigDir        string
	ConfigName       string
	ConfigPath       string
	ImageCoordinates string
	Name             string
	Port             string
	DockerClient     *client.Client
	Ctx              context.Context
)

func InitKdkConfig(logger logrus.Entry) error {

	currentUser, _ := user.Current()

	configTemplate := `image:
  repository: {{ .imageRepository }}
  tag: {{ .imageTag }}
docker:
  name: {{ .name }}
  hostname: {{ .name }}
  environment:
    KDK_DOTFILES_REPO: {{ .dotfilesRepo }}
    KDK_SHELL: {{ .shell }}
    KDK_USERNAME: {{ .username }}
    KDK_PORT: "{{ .port }}"
  binds:
    - source: /Users/{{ .username }}/.kdk/ssh/id_rsa.pub
      target: /home/{{ .username }}/.ssh/authorized_keys
`

	configDefaults := map[string]interface{}{
		"imageRepository": "ciscosso/kdk",
		"imageTag":        "debian-latest",
		"name":            "kdk",
		"port":            "2022",
		"username":        currentUser.Username,
		"dotfilesRepo":    "https://github.com/cisco-sso/yadm-dotfiles.git",
		"shell":           "/bin/bash",
	}
	if _, err := os.Stat(ConfigDir); os.IsNotExist(err) {
		if err := os.Mkdir(ConfigDir, 0700); err != nil {
			logger.WithField("error", err).Fatal("Failed to create KDK config directory")
			return err
		}
	}
	if _, err := os.Stat(ConfigPath); os.IsNotExist(err) {
		logger.Warn("KDK config does not exist")
		logger.Info("Creating KDK config")
		ioutil.WriteFile(ConfigPath, []byte(utils.Tprintf(configTemplate, configDefaults)), 0600)
	} else {
		logger.Warn("KDK config exists")
		logger.Info("Overwrite existing KDK config?")
		prompt := promptui.Prompt{Label: "Continue", IsConfirm: true}

		if _, err := prompt.Run(); err != nil {
			logger.Info("Existing KDK config not overwritten")
			return err
		} else {
			logger.Info("Creating KDK config")
			ioutil.WriteFile(ConfigPath, []byte(utils.Tprintf(configTemplate, configDefaults)), 0600)
		}
	}
	return nil
}

func InitKdkSshKeyPair(logger logrus.Entry) error {
	keypairName := "id_rsa"
	keypairDir := path.Join(ConfigDir, "ssh")

	privateKeyPath := path.Join(keypairDir, keypairName)

	if _, err := os.Stat(ConfigDir); os.IsNotExist(err) {
		if err := os.Mkdir(ConfigDir, 0700); err != nil {
			logger.WithField("error", err).Fatal("Failed to create KDK config directory")
		}
	}
	if _, err := os.Stat(keypairDir); os.IsNotExist(err) {
		if err := os.Mkdir(keypairDir, 0700); err != nil {
			logger.WithField("error", err).Fatal("Failed to create ssh key directory")
		}
	}
	if _, err := os.Stat(privateKeyPath); os.IsNotExist(err) {
		logger.Warn("KDK ssh key pair not found.")
		logger.Info("Generating ssh key pair...")
		privateKey, err := utils.GeneratePrivateKey(4096)
		if err != nil {
			logger.WithField("error", err).Fatal("Failed to generate ssh private key")
			return err
		}
		publicKeyBytes, err := utils.GeneratePublicKey(&privateKey.PublicKey)
		if err != nil {
			logger.WithField("error", err).Fatal("Failed to generate ssh public key")
			return err
		}
		err = utils.WriteKeyToFile(utils.EncodePrivateKey(privateKey), privateKeyPath)
		if err != nil {
			logger.WithField("error", err).Fatal("Failed to write ssh private key")
			return err
		}
		err = utils.WriteKeyToFile([]byte(publicKeyBytes), privateKeyPath+".pub")
		if err != nil {
			logger.WithField("error", err).Fatal("Failed to write ssh public key")
			return err
		}
		logger.Info("Successfully generated ssh key pair.")

	} else {
		logger.Info("KDK ssh key pair exists.")
	}
	return nil
}
