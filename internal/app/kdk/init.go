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
	"path/filepath"
	"strings"

	"github.com/docker/go-connections/nat"
	"github.com/ghodss/yaml"
	"github.com/Sirupsen/logrus"
	"github.com/cisco-sso/kdk/internal/pkg/utils"
	"github.com/docker/docker/client"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/manifoldco/promptui"
)

var (
	ConfigDir        string
	ConfigName       string
	ConfigPath       string
	Verbose          bool
	KeypairDir       string
	PrivateKeyPath   string
	PublicKeyPath    string
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
	username := currentUser.Username

	//user.Current().Username can contain '\' on windows.
	if strings.Contains(username, "\\") {
		username = strings.Split(username, "\\")[1]
	}

	// Initialize storage mounts/volumes
	mounts  := []mount.Mount{}       // hostConfig
	volumes := map[string]struct{}{} // containerConfig

	// Define mount configurations for mounting the ssh pub key into a tmp location where the bootstrap script may
	// copy into <userdir>/.ssh/authorized keys.  This is required because Windows mounts squash permissions to 777
	// which makes ssh fail a strict check on pubkey permissions.
	source := PublicKeyPath
	target := "/tmp/id_rsa.pub"
	mounts = append(mounts, mount.Mount{Type: mount.TypeBind, Source: source, Target: target, ReadOnly: true})
	volumes[target] = struct{}{}

	// Define volume binding for public key

	// Define volume bindings for the keybase directory
	//   Linux & OSX: Detec /keybase
	//   Windows10: Detect k: and /k
	keybaseRoots := []string{ "/keybase", "k:", "/k" }
	keybaseTestSubdir := "/private"
	for _, keybaseRoot := range keybaseRoots {
		if absPath, err := filepath.Abs(filepath.Join(keybaseRoot, keybaseTestSubdir)); err == nil {
			if path, err := filepath.EvalSymlinks(absPath); err == nil {
				source := filepath.Dir(path)
				target := "/keybase"

				logger.Infof("Detected /keybase filesystem at: %v", source)

				prompt := promptui.Prompt {
					Label: "Mount your /keybase directory within KDK? [y/n]",
				Default: "y",
					IsVimMode: true,
					Validate: promptuiValidateYorN,
				}
				if result, err := prompt.Run(); err == nil && result == "y" {
					logger.Info("Adding /keybase mount to configuration")
					mounts = append(mounts, mount.Mount{Type: mount.TypeBind, Source: source, Target: target})
					volumes[target] = struct{}{}

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
			Username: username,
		},
		ContainerConfig: container.Config{
			Hostname: "kdk",
			Image:    "ciscosso/kdk:debian-latest",
			Tty:      true,
			Env: []string{
				"KDK_USERNAME=" + username,
				"KDK_SHELL=/bin/bash",
				"KDK_DOTFILES_REPO=https://github.com/cisco-sso/yadm-dotfiles.git",
			},
			ExposedPorts: nat.PortSet{
				"2022/tcp": struct{}{},
			},
			Volumes: volumes,
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
			Mounts: mounts,
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

	if _, err := os.Stat(ConfigDir); os.IsNotExist(err) {
		if err := os.Mkdir(ConfigDir, 0700); err != nil {
			logger.WithField("error", err).Fatal("Failed to create KDK config directory")
		}
	}
	if _, err := os.Stat(KeypairDir); os.IsNotExist(err) {
		if err := os.Mkdir(KeypairDir, 0700); err != nil {
			logger.WithField("error", err).Fatal("Failed to create ssh key directory")
		}
	}
	if _, err := os.Stat(PrivateKeyPath); os.IsNotExist(err) {
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
		err = utils.WriteKeyToFile(utils.EncodePrivateKey(privateKey), PrivateKeyPath)
		if err != nil {
			logger.WithField("error", err).Fatal("Failed to write ssh private key")
			return err
		}
		err = utils.WriteKeyToFile([]byte(publicKeyBytes), PublicKeyPath)
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
