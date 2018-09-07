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
	"os/user"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/mitchellh/go-homedir"
)

var Version = "undefined"

// Struct of all configs to be saved directly as ~/.kdk/<NAME>/config.yaml
type KdkEnvConfig struct {
	DockerClient *client.Client
	Ctx          context.Context
	ConfigFile   configFile
}

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
	Debug           bool
}

// create docker client and context for easy reuse
func (c *KdkEnvConfig) Init() {
	c.Ctx = context.Background()
	dockerClient, err := client.NewEnvClient()
	if err != nil {
		panic(err)
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
