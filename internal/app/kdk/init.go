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
	"errors"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"

	"github.com/docker/go-connections/nat"
	"github.com/ghodss/yaml"
	"github.com/Sirupsen/logrus"
	"github.com/cisco-sso/kdk/internal/pkg/utils"
	"github.com/docker/docker/client"
	"github.com/docker/docker/api/types/container"
	"github.com/manifoldco/promptui"
)

var (
	ConfigDir        string
	ConfigName       string
	ConfigPath       string
	DockerClient     *client.Client
	Ctx              context.Context
	KdkConfig        *kdkConfig
)

// Struct of configs from the docker lib, to be saved directly as ~/.kdk/config.yaml
type kdkConfig struct {
	AppConfig appConfig
	ContainerConfig container.Config `json:",omitempty"`
	HostConfig container.HostConfig  `json:",omitempty"`
}

type appConfig struct {
	Name string
	Port string
	Username string
}

func InitKdkConfig(logger logrus.Entry) error {

	currentUser, _ := user.Current()

	// Define volume bindings for mounting the ssh pub key into authorized keys
	binds := []string{fmt.Sprintf("/Users/%v/.kdk/ssh/id_rsa.pub:/home/%v/.ssh/authorized_keys",
		currentUser.Username, currentUser.Username)}

	// Define volume bindings for the keybase directory
	//   Linux & OSX: Detec /keybase
	//   Windows10: Detect k: and /k
	keybaseRoots := []string{ "/keybase", "k:", "/k" }
	keybaseTestSubdir := "/private"
	for _, keybaseRoot := range keybaseRoots {
		if absPath, err := filepath.Abs(filepath.Join(keybaseRoot, keybaseTestSubdir)); err == nil {
			if path, err := filepath.EvalSymlinks(absPath); err == nil {
				target := filepath.Dir(path)+":/keybase"

				logger.Info("Detected /keybase filesystem")
				prompt := promptui.Prompt {
					Label: "Mount your /keybase directory within KDK? [y/n]",
				Default: "y",
					IsVimMode: true,
					Validate: promptuiValidateYorN,
				}
				if result, err := prompt.Run(); err == nil && result == "y" {
					logger.Info(fmt.Sprintf("Adding Bind target `%v` to configuration", target))
					binds = append(binds, path+":/keybase")
				} else {
					logger.Info(fmt.Sprintf("Skip Adding of Bind target `%v` to configuration", target))
				}
			}
		}
	}

	// Create the Default configuration struct that will be written as the config file
	KdkConfig = &kdkConfig{
		AppConfig: appConfig{
			Name: "kdk",
			Port: "2022",
			Username: currentUser.Username,
		},
		ContainerConfig: container.Config{
			Hostname: "kdk",
			Image:    "ciscosso/kdk:debian-latest",
			Tty:      true,
			Env: []string{
				"KDK_USERNAME=" + currentUser.Username,
				"KDK_SHELL=/bin/bash",
				"KDK_DOTFILES_REPO=https://github.com/cisco-sso/yadm-dotfiles.git",
			},
			ExposedPorts: nat.PortSet{
				"2022/tcp": struct{}{},
			},
		},
		HostConfig: container.HostConfig{
			// TODO (rluckie): shouldn't default to privileged -- issue with ssh cmd
			Privileged: true,
			PortBindings: nat.PortMap{
				"2022/tcp": []nat.PortBinding{
					{
						HostPort: "2022",
					},
				},
			},
			Binds: binds,
		},
	}

	// Ensure that the ~/.kdk directory exists
	if _, err := os.Stat(ConfigDir); os.IsNotExist(err) {
		if err := os.Mkdir(ConfigDir, 0700); err != nil {
			logger.WithField("error", err).Fatal("Failed to create KDK config directory")
			return err
		}
	}

	// Create the ~/.kdk/config.yaml file if it doesn't exist
	y, err := yaml.Marshal(KdkConfig)
	if err != nil {
		logger.Fatal("Failed to create YAML string of configuration", err)
	}
	if _, err := os.Stat(ConfigPath); os.IsNotExist(err) {
		logger.Warn("KDK config does not exist")
		logger.Info("Creating KDK config")

		ioutil.WriteFile(ConfigPath, y, 0600)
	} else {
		logger.Warn("KDK config exists")
		prompt := promptui.Prompt {
			Label: "Overwrite existing KDK config? [y/n]",
		Default: "n",
			IsVimMode: true,
			Validate: promptuiValidateYorN,
		}
		if result, err := prompt.Run(); err == nil && result == "y" {
			logger.Info("Creating KDK config")
			ioutil.WriteFile(ConfigPath, y, 0600)
		} else {
			logger.Info("Existing KDK config not overwritten")
			return err
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

func promptuiValidateYorN(input string) error {
	if input == "y" || input == "n" {
		return nil
	}
	return errors.New("Input must be 'y' or 'n'")
}
