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

package keybase

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/cisco-sso/kdk/pkg/prompt"
	"github.com/codeskyblue/go-sh"
)

// Write keybase mirror script [windows only]
const mirrorScript = `
@echo off

if "%1"=="" (
  echo "You must pass either start or stop"
  break
)

if "%1"=="start" (
  echo "Starting"
  start "KDK Keybase Mirror" /B "C:\Program Files\Dokan\Dokan Library-1.1.0\sample\mirror\mirror.exe" /r K:\ /l C:\Users\%USERNAME%\.kdk\keybase
  break
)

if "%1"=="stop" (
  echo "stopping"
  tskill.exe mirror
  break
) else (
  echo "Unrecognized parameter %1.  You must pass either start or stop"
)
`

func writeMirrorScript(configDir string) (out string, err error) {
	script := []byte(mirrorScript)

	scriptPath := filepath.Join(configDir, "keybase-mirror.cmd")
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		err := ioutil.WriteFile(scriptPath, script, 0700)
		if err != nil {
			return "", err
		}
	}
	return scriptPath, nil
}

// Start keybase mirror [windows only]
func StartMirror(configDir string, debug bool, logger logrus.Entry) error {

	keybaseTestDir := filepath.Join(configDir, "keybase", "private")

	// keybase mirror already started.  Nothing to do
	if _, err := os.Stat(keybaseTestDir); err == nil {
		logger.Info("Keybase mirror already started")
		return nil
	}
	logger.Info("Writing keybase mirror script")
	scriptPath, err := writeMirrorScript(configDir)
	if err != nil {
		return err
	}
	commandString := fmt.Sprintf("powershell %s %s", scriptPath, "start")
	if debug {
		logrus.Infof("Starting Keybase mirror with command; %s", commandString)
	}
	commandMap := strings.Split(commandString, " ")
	if err := sh.Command(commandMap[0], commandMap[1:]).SetStdin(os.Stdin).Run(); err != nil {
		return err
	}
	return nil
}

// Stop keybase mirror [windows only]
func StopMirror(configDir string, debug bool, logger logrus.Entry) error {
	// TODO(rluckie) Fix StopMirror to work with multiple KDK containers
	// Use docker client to iterate though all running containers and ensure that no containers have mirror dir mounted
	// If so, stop mirror
	logger.Info("Writing keybase mirror script")
	scriptPath, err := writeMirrorScript(configDir)
	if err != nil {
		return err
	}
	commandString := fmt.Sprintf("powershell %s %s", scriptPath, "stop")
	if debug {
		logrus.Infof("Stopping Keybase mirror with command; %s", commandString)
	}
	commandMap := strings.Split(commandString, " ")
	if err := sh.Command(commandMap[0], commandMap[1:]).SetStdin(os.Stdin).Run(); err != nil {
		return err
	}
	return nil
}

// Get keybase mounts
// Linux & OSX: Detect /keybase
// Windows10: Detect k: and /k
func GetMounts(configRootDir string, logger logrus.Entry) (source string, target string, err error) {

	keybaseRoots := []string{"/keybase", "k:", "/k"}
	keybaseTestSubdir := "/private"
	for _, keybaseRoot := range keybaseRoots {
		if absPath, err := filepath.Abs(filepath.Join(keybaseRoot, keybaseTestSubdir)); err == nil {
			if path, err := filepath.EvalSymlinks(absPath); err == nil {
				source := filepath.Dir(path)
				target := "/keybase"

				logger.Infof("Detected keybase filesystem at: %v", source)

				prmpt := prompt.Prompt{
					Text:     "Mount your keybase directory within KDK? [y/n] ",
					Loop:     true,
					Validate: prompt.ValidateYorN,
				}
				if result, err := prmpt.Run(); err == nil && result == "y" {
					logger.Info("Adding /keybase mount to configuration")
					if runtime.GOOS == "windows" {
						source = filepath.Join(configRootDir, "keybase")
						if _, err := os.Stat(source); os.IsNotExist(err) {
							if err := os.Mkdir(source, 0700); err != nil {
								logger.WithField("error", err).Fatalf("Failed to create KDK keybase mirror directory [%s]", source)
								return "", "", err
							}
						}
					}
					return source, target, nil
				}
			}
		}
	}
	return "", "", errors.New("Failed to detect potential keybase filesystem mounts")
}
