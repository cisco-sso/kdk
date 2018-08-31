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

package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/cisco-sso/kdk/internal/pkg/utils/simpleprompt"
	"github.com/codeskyblue/go-sh"
)

// Write keybase mirror script [windows only]
func keybaseWriteMirrorScript(configDir string) (out string, err error) {

	cmdString := `
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
	script := []byte(cmdString)

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
func KeybaseStartMirror(configDir string) error {

	keybaseTestDir := filepath.Join(configDir, "keybase", "private")
	if _, err := os.Stat(keybaseTestDir); err == nil {
		return nil
	}
	scriptPath, err := keybaseWriteMirrorScript(configDir)
	if err != nil {
		return err
	}
	commandString := fmt.Sprintf("powershell %s %s", scriptPath, "start")
	commandMap := strings.Split(commandString, " ")
	if err := sh.Command(commandMap[0], commandMap[1:]).SetStdin(os.Stdin).Run(); err != nil {
		return err
	}
	return nil
}

// Stop keybase mirror [windows only]
func KeybaseStopMirror(configDir string) error {

	keybaseTestDir := filepath.Join(configDir, "keybase", "private")
	if _, err := os.Stat(keybaseTestDir); err != nil {
		return nil
	}
	scriptPath, err := keybaseWriteMirrorScript(configDir)
	if err != nil {
		return err
	}
	commandString := fmt.Sprintf("powershell %s %s", scriptPath, "stop")
	commandMap := strings.Split(commandString, " ")
	if err := sh.Command(commandMap[0], commandMap[1:]).SetStdin(os.Stdin).Run(); err != nil {
		return err
	}
	return nil
}

// Get keybase mounts
// Linux & OSX: Detect /keybase
// Windows10: Detect k: and /k
func KeybaseGetMounts(configDir string, logger logrus.Entry) (source string, target string, err error) {

	keybaseRoots := []string{"/keybase", "k:", "/k"}
	keybaseTestSubdir := "/private"
	for _, keybaseRoot := range keybaseRoots {
		if absPath, err := filepath.Abs(filepath.Join(keybaseRoot, keybaseTestSubdir)); err == nil {
			if path, err := filepath.EvalSymlinks(absPath); err == nil {
				source := filepath.Dir(path)
				target := "/keybase"

				logger.Infof("Detected keybase filesystem at: %v", source)

				prompt := simpleprompt.Prompt{
					Text:     "Mount your keybase directory within KDK? [y/n] ",
					Loop:     true,
					Validate: simpleprompt.ValidateYorN,
				}
				if result, err := prompt.Run(); err == nil && result == "y" {
					logger.Info("Adding /keybase mount to configuration")
					if runtime.GOOS == "windows" {
						source = filepath.Join(configDir, "keybase")
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
