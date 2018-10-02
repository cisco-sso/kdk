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
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/cisco-sso/kdk/pkg/keybase"
	"github.com/cisco-sso/kdk/pkg/prompt"
	"github.com/cisco-sso/kdk/pkg/ssh"
	"github.com/cisco-sso/kdk/pkg/utils"
	"github.com/codeskyblue/go-sh"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/ghodss/yaml"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
)

var (
	Version = "undefined"
	Port    = strconv.Itoa(utils.GetPort())
)

type KdkEnvConfig struct {
	DockerClient *client.Client
	Ctx          context.Context
	ConfigFile   configFile
	SocksPort    string
}

// Struct of all configs to be saved directly as ~/.kdk/<NAME>/config.yaml
type configFile struct {
	AppConfig       AppConfig
	ContainerConfig *container.Config     `json:",omitempty"`
	HostConfig      *container.HostConfig `json:",omitempty"`
}

type AppConfig struct {
	Name            string
	Port            string
	ImageRepository string
	ImageTag        string
	DotfilesRepo    string
	Shell           string
	SocksPort       string
}

// create docker client and context for easy reuse
func (c *KdkEnvConfig) Init() {
	c.Ctx = context.Background()
	dockerClient, err := client.NewEnvClient()
	if err != nil {
		log.Warn("Failed to create docker client.")
		log.Fatal("Ensure that docker is running.")
	}

	c.DockerClient = dockerClient
}

// current username
func (c *KdkEnvConfig) User() (out string) {
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}
	username := currentUser.Username
	// Windows usernames are `domain\username`.  Strip the domain in case we are running on Windows.
	if strings.Contains(username, "\\") {
		username = strings.Split(username, "\\")[1]
	}
	return username
}

// users home directory
func (c *KdkEnvConfig) Home() (out string) {
	out, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	return out
}

// kdk root config path (~/.kdk)
func (c *KdkEnvConfig) ConfigRootDir() (out string) {
	return filepath.Join(c.Home(), ".kdk")
}

// kdk keypair path path (~/.kdk/ssh)
func (c *KdkEnvConfig) KeypairDir() (out string) {
	return filepath.Join(c.ConfigRootDir(), "ssh")
}

// kdk private key path (~/.kdk/ssh/id_rsa)
func (c *KdkEnvConfig) PrivateKeyPath() (out string) {
	return filepath.Join(c.KeypairDir(), "id_rsa")
}

// kdk public key path (~/.kdk/ssh/id_rsa.pub)
func (c *KdkEnvConfig) PublicKeyPath() (out string) {
	return filepath.Join(c.KeypairDir(), "id_rsa.pub")
}

// kdk container config dir (~/.kdk/<KDK_NAME>)
func (c *KdkEnvConfig) ConfigDir() (out string) {
	return filepath.Join(c.ConfigRootDir(), c.ConfigFile.AppConfig.Name)
}

// kdk container config path (~/.kdk/<KDK_NAME>/config.yaml)
func (c *KdkEnvConfig) ConfigPath() (out string) {
	return filepath.Join(c.ConfigDir(), "config.yaml")
}

// kdk image coordinates (ciscosso/kdk:debian-latest)
func (c *KdkEnvConfig) ImageCoordinates() (out string) {
	return c.ConfigFile.AppConfig.ImageRepository + ":" + c.ConfigFile.AppConfig.ImageTag
}

func (c *KdkEnvConfig) CreateKdkConfig() (err error) {

	// Initialize storage mounts/volumes
	var mounts []mount.Mount         // hostConfig
	volumes := map[string]struct{}{} // containerConfig
	labels := map[string]string{"kdk": Version}

	// Define mount configurations for mounting the ssh pub key into a tmp location where the bootstrap script may
	//   copy into <userdir>/.ssh/authorized keys.  This is required because Windows mounts squash permissions to
	//   777 which makes ssh fail a strict check on pubkey permissions.
	source := c.PublicKeyPath()
	target := "/tmp/id_rsa.pub"
	mounts = append(mounts, mount.Mount{Type: mount.TypeBind, Source: source, Target: target, ReadOnly: true})
	volumes[target] = struct{}{}

	// Keybase mounts
	source, target, err = keybase.GetMounts(c.ConfigRootDir())
	if err != nil {
		log.Warn("Failed to add keybase mount:", err)
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
				log.Infof("Entered host source directory mount %v", source)
			}

			prmpt = prompt.Prompt{
				Text:     "Please enter the docker container target directory (e.g. /home/<username>/Projects) ",
				Loop:     false,
				Validate: nil,
			}
			target, err := prmpt.Run()
			if err == nil {
				log.Infof("Entered container target directory mount %v", target)
			}

			mounts = append(mounts, mount.Mount{Type: mount.TypeBind, Source: source, Target: target, ReadOnly: false})
			volumes[target] = struct{}{}
		} else {
			break
		}
	}

	// Prompt for SOCKS proxy options.
	if c.SocksPort == "" {
		socksPort := ""
		prmpt := prompt.Prompt{
			Text:     "Would you like to enable SOCKS proxy? [y/n] ",
			Loop:     true,
			Validate: prompt.ValidateYorN,
		}
		if result, err := prmpt.Run(); err == nil && result == "y" {
			prmpt = prompt.Prompt{
				Text:     "Please enter SOCKS port number [8000] ",
				Loop:     true,
				Validate: prompt.ValidateIntOrEmptyString,
			}
			socksPort, _ = prmpt.Run()
		}
		if socksPort == "" {
			socksPort = "8000"
		}
		c.ConfigFile.AppConfig.SocksPort = socksPort
	} else {
		c.ConfigFile.AppConfig.SocksPort = c.SocksPort
	}
	log.Infof("Set SOCKS port %v", c.ConfigFile.AppConfig.SocksPort)

	// Create the Default configuration struct that will be written as the config file
	c.ConfigFile.ContainerConfig = &container.Config{
		Hostname: c.ConfigFile.AppConfig.Name,
		Image:    c.ImageCoordinates(),
		Tty:      true,
		Env: []string{
			"KDK_USERNAME=" + c.User(),
			"KDK_SHELL=" + c.ConfigFile.AppConfig.Shell,
			"KDK_DOTFILES_REPO=" + c.ConfigFile.AppConfig.DotfilesRepo,
		},
		ExposedPorts: nat.PortSet{
			"2022/tcp": struct{}{},
		},
		Volumes: volumes,
		Labels:  labels,
	}
	c.ConfigFile.HostConfig = &container.HostConfig{
		// TODO (rluckie): shouldn't default to privileged -- issue with ssh cmd
		Privileged: true,
		PortBindings: nat.PortMap{
			"2022/tcp": []nat.PortBinding{
				{
					HostPort: c.ConfigFile.AppConfig.Port,
				},
			},
		},
		Mounts: mounts,
	}

	// Ensure that the ~/.kdk directory exists
	if _, err := os.Stat(c.ConfigRootDir()); os.IsNotExist(err) {
		if err := os.Mkdir(c.ConfigRootDir(), 0700); err != nil {
			log.WithField("error", err).Fatalf("Failed to create KDK config directory [%s]", c.ConfigRootDir())
			return err
		}
	}

	// Ensure that the ~/.kdk/<kdkName> directory exists
	if _, err := os.Stat(c.ConfigDir()); os.IsNotExist(err) {
		if err := os.Mkdir(c.ConfigDir(), 0700); err != nil {
			log.WithField("error", err).Fatalf("Failed to create KDK config directory", filepath.Dir(c.ConfigDir()))
			return err
		}
	}

	// Create the ~/.kdk/<kdkName>/config.yaml file if it doesn't exist
	y, err := yaml.Marshal(&c.ConfigFile)
	if err != nil {
		log.Fatal("Failed to create YAML string of configuration", err)
	}
	if _, err := os.Stat(c.ConfigPath()); os.IsNotExist(err) {
		log.Warn("KDK config does not exist")
		log.Info("Creating KDK config")

		ioutil.WriteFile(c.ConfigPath(), y, 0600)
	} else {
		log.Warn("KDK config exists")
		prmpt := prompt.Prompt{
			Text:     "Overwrite existing KDK config? [y/n] ",
			Loop:     true,
			Validate: prompt.ValidateYorN,
		}
		if result, err := prmpt.Run(); err == nil && result == "y" {
			log.Info("Creating KDK config")
			ioutil.WriteFile(c.ConfigPath(), y, 0600)
		} else {
			log.Info("Existing KDK config not overwritten")
			return err
		}
	}
	return nil
}

// Creates KDK ssh keypair
func (c *KdkEnvConfig) CreateKdkSshKeyPair() (err error) {

	if _, err := os.Stat(c.ConfigRootDir()); os.IsNotExist(err) {
		if err := os.Mkdir(c.ConfigRootDir(), 0700); err != nil {
			log.WithField("error", err).Fatal("Failed to create KDK config directory")
		}
	}
	if _, err := os.Stat(c.KeypairDir()); os.IsNotExist(err) {
		if err := os.Mkdir(c.KeypairDir(), 0700); err != nil {
			log.WithField("error", err).Fatal("Failed to create ssh key directory")
		}
	}
	if _, err := os.Stat(c.PrivateKeyPath()); os.IsNotExist(err) {
		log.Warn("KDK ssh key pair not found.")
		log.Info("Generating ssh key pair...")
		privateKey, err := ssh.GeneratePrivateKey(4096)
		if err != nil {
			log.WithField("error", err).Fatal("Failed to generate ssh private key")
			return err
		}
		publicKeyBytes, err := ssh.GeneratePublicKey(&privateKey.PublicKey)
		if err != nil {
			log.WithField("error", err).Fatal("Failed to generate ssh public key")
			return err
		}
		err = ssh.WriteKeyToFile(ssh.EncodePrivateKey(privateKey), c.PrivateKeyPath())
		if err != nil {
			log.WithField("error", err).Fatal("Failed to write ssh private key")
			return err
		}
		err = ssh.WriteKeyToFile([]byte(publicKeyBytes), c.PublicKeyPath())
		if err != nil {
			log.WithField("error", err).Fatal("Failed to write ssh public key")
			return err
		}
		log.Info("Successfully generated ssh key pair.")

	} else {
		log.Info("KDK ssh key pair exists.")
	}
	return nil
}

// Returns SSH connection string
func (c *KdkEnvConfig) SSHConnectionString() string {
	return c.User() + "@localhost"
}

// Returns SSH command string
func (c *KdkEnvConfig) SSHCommandString() string {
	return fmt.Sprintf("ssh %s -A -p %s -i %s -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null",
		c.SSHConnectionString(), c.ConfigFile.AppConfig.Port, c.PrivateKeyPath())
}

// Returns SCP command string
func (c *KdkEnvConfig) SCPCommandString() string {
	return fmt.Sprintf("scp -P %s -i %s -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null",
		c.ConfigFile.AppConfig.Port, c.PrivateKeyPath())
}

// SCP's a file into the KDK container
func (c *KdkEnvConfig) SCPTo(hostPath, kdkPath string) error {
	commandString := fmt.Sprintf("%s %s %s:%s", c.SCPCommandString(), hostPath, c.SSHConnectionString(), kdkPath)
	log.Infof("executing scp command: %s", commandString)
	commandMap := strings.Split(commandString, " ")
	if err := sh.Command(commandMap[0], commandMap[1:]).SetStdin(os.Stdin).Run(); err != nil {
		return err
	}
	return nil
}

// Executes a command on the KDK container
func (c *KdkEnvConfig) Exec(command string) error {
	commandString := fmt.Sprintf("%s %s", c.SSHCommandString(), command)
	log.Infof("executing ssh command: %s", commandString)
	commandMap := strings.Split(commandString, " ")
	return sh.Command(commandMap[0], commandMap[1:]).SetStdin(os.Stdin).Run()
}

// Checks that KDK container is running
func (c *KdkEnvConfig) IsRunning() bool {
	kdkRunning := false

	containers, err := c.DockerClient.ContainerList(c.Ctx, types.ContainerListOptions{All: true})
	if err != nil {
		log.WithField("error", err).Fatal("Failed to list docker containers")
	}

	for _, container := range containers {
		for _, name := range container.Names {
			if name == "/"+c.ConfigFile.AppConfig.Name {
				if container.State == "running" {
					kdkRunning = true
					break
				}
			}
		}
	}
	return kdkRunning
}

// If KDK container is not running, start it and provision KDK user.
func (c *KdkEnvConfig) Start() {
	if !c.IsRunning() {
		log.Info("KDK is not currently running.  Starting...")
		Pull(c, false)
		Up(*c)
		Provision(*c)
	}
}
