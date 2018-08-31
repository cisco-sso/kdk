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
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/cisco-sso/kdk/pkg/keybase"
	"github.com/cisco-sso/kdk/pkg/prompt"
	"github.com/cisco-sso/kdk/pkg/ssh"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
	"github.com/ghodss/yaml"
)

// TODO (rluckie) refactor to be a method of kdkConfig type
func InitKdkConfig(cfg KdkEnvConfig, logger logrus.Entry) error {

	// Initialize storage mounts/volumes
	var mounts []mount.Mount         // hostConfig
	volumes := map[string]struct{}{} // containerConfig
	labels := map[string]string{"kdk": Version}

	// Define mount configurations for mounting the ssh pub key into a tmp location where the bootstrap script may
	//   copy into <userdir>/.ssh/authorized keys.  This is required because Windows mounts squash permissions to
	//   777 which makes ssh fail a strict check on pubkey permissions.
	source := cfg.PublicKeyPath()
	target := "/tmp/id_rsa.pub"
	mounts = append(mounts, mount.Mount{Type: mount.TypeBind, Source: source, Target: target, ReadOnly: true})
	volumes[target] = struct{}{}

	// Keybase mounts
	source, target, err := keybase.GetMounts(cfg.ConfigRootDir(), logger)
	if err != nil {
		logger.Warn("Failed to add keybase mount:", err)
	} else {
		mounts = append(mounts, mount.Mount{Type: mount.TypeBind, Source: source, Target: target, ReadOnly: false})
		volumes[target] = struct{}{}
	}

	// Define Additional volume bindings
	for {
		prmpt := prompt.Prompt{
			Text:     "Would you like to mount additional docker host directories into the KDK? [y/n] ",
			Loop:     true,
			Validate: prompt.ValidateYorN,
		}
		if result, err := prmpt.Run(); err == nil && result == "y" {
			prmpt = prompt.Prompt{
				Text:     "Please enter the docker host source directory (e.g. /Users/<username>/Projects) ",
				Loop:     true,
				Validate: prompt.ValidateDirExists,
			}
			source, err := prmpt.Run()
			if err == nil {
				logger.Infof("Entered host source directory mount %v", source)
			}

			prmpt = prompt.Prompt{
				Text:     "Please enter the docker container target directory (e.g. /home/<username>/Projects) ",
				Loop:     false,
				Validate: nil,
			}
			target, err := prmpt.Run()
			if err == nil {
				logger.Infof("Entered container target directory mount %v", target)
			}

			mounts = append(mounts, mount.Mount{Type: mount.TypeBind, Source: source, Target: target, ReadOnly: false})
			volumes[target] = struct{}{}
		} else {
			break
		}
	}

	// Create the Default configuration struct that will be written as the config file
	cfg.KdkCfg = &kdkConfig{
		ContainerConfig: container.Config{
			Hostname: cfg.Name,
			Image:    cfg.ImageCoordinates(),
			Tty:      true,
			Env: []string{
				"KDK_USERNAME=" + cfg.User(),
				"KDK_SHELL=" + cfg.Shell,
				"KDK_DOTFILES_REPO=" + cfg.DotfilesRepo,
			},
			ExposedPorts: nat.PortSet{
				"2022/tcp": struct{}{},
			},
			Volumes: volumes,
			Labels:  labels,
		},
		HostConfig: container.HostConfig{
			// TODO (rluckie): shouldn't default to privileged -- issue with ssh cmd
			Privileged: true,
			PortBindings: nat.PortMap{
				"2022/tcp": []nat.PortBinding{
					{
						HostPort: cfg.Port,
					},
				},
			},
			Mounts: mounts,
		},
	}

	// Ensure that the ~/.kdk directory exists
	if _, err := os.Stat(cfg.ConfigRootDir()); os.IsNotExist(err) {
		if err := os.Mkdir(cfg.ConfigRootDir(), 0700); err != nil {
			logger.WithField("error", err).Fatalf("Failed to create KDK config directory [%s]", cfg.ConfigRootDir())
			return err
		}
	}

	// Ensure that the ~/.kdk/<kdkName> directory exists
	if _, err := os.Stat(cfg.ConfigDir()); os.IsNotExist(err) {
		if err := os.Mkdir(cfg.ConfigDir(), 0700); err != nil {
			logger.WithField("error", err).Fatalf("Failed to create KDK config directory", filepath.Dir(cfg.ConfigDir()))
			return err
		}
	}

	// Create the ~/.kdk/<kdkName>/config.yaml file if it doesn't exist
	y, err := yaml.Marshal(&cfg)
	if err != nil {
		logger.Fatal("Failed to create YAML string of configuration", err)
	}
	if _, err := os.Stat(cfg.ConfigPath()); os.IsNotExist(err) {
		logger.Warn("KDK config does not exist")
		logger.Info("Creating KDK config")

		ioutil.WriteFile(cfg.ConfigPath(), y, 0600)
	} else {
		logger.Warn("KDK config exists")
		prmpt := prompt.Prompt{
			Text:     "Overwrite existing KDK config? [y/n] ",
			Loop:     true,
			Validate: prompt.ValidateYorN,
		}
		if result, err := prmpt.Run(); err == nil && result == "y" {
			logger.Info("Creating KDK config")
			ioutil.WriteFile(cfg.ConfigPath(), y, 0600)
		} else {
			logger.Info("Existing KDK config not overwritten")
			return err
		}
	}
	return nil
}

// TODO (rluckie) refactor to be a method of kdkConfig type
func InitKdkSshKeyPair(cfg KdkEnvConfig, logger logrus.Entry) error {

	if _, err := os.Stat(cfg.ConfigRootDir()); os.IsNotExist(err) {
		if err := os.Mkdir(cfg.ConfigRootDir(), 0700); err != nil {
			logger.WithField("error", err).Fatal("Failed to create KDK config directory")
		}
	}
	if _, err := os.Stat(cfg.KeypairDir()); os.IsNotExist(err) {
		if err := os.Mkdir(cfg.KeypairDir(), 0700); err != nil {
			logger.WithField("error", err).Fatal("Failed to create ssh key directory")
		}
	}
	if _, err := os.Stat(cfg.PrivateKeyPath()); os.IsNotExist(err) {
		logger.Warn("KDK ssh key pair not found.")
		logger.Info("Generating ssh key pair...")
		privateKey, err := ssh.GeneratePrivateKey(4096)
		if err != nil {
			logger.WithField("error", err).Fatal("Failed to generate ssh private key")
			return err
		}
		publicKeyBytes, err := ssh.GeneratePublicKey(&privateKey.PublicKey)
		if err != nil {
			logger.WithField("error", err).Fatal("Failed to generate ssh public key")
			return err
		}
		err = ssh.WriteKeyToFile(ssh.EncodePrivateKey(privateKey), cfg.PrivateKeyPath())
		if err != nil {
			logger.WithField("error", err).Fatal("Failed to write ssh private key")
			return err
		}
		err = ssh.WriteKeyToFile([]byte(publicKeyBytes), cfg.PublicKeyPath())
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
