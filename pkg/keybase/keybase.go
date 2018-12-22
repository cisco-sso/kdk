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

	"github.com/cisco-sso/kdk/pkg/prompt"
	"github.com/codeskyblue/go-sh"
	log "github.com/sirupsen/logrus"
)

// Write keybase mirror script [windows only]
const mirrorScript = `
@ECHO off

if "%1"=="" (
  ECHO You must pass either 'start' or 'stop' as the first script argument
  BREAK
)

if "%1"=="start" (
  ECHO Starting
  FOR /F "tokens=*" %%F IN ('dir "c:\Program Files\Dokan" /S/B ^| findstr /R "\\sample\\mirror\\mirror.exe"') DO @(
    SET MIRROREXE=%%F
    GOTO :once
  )

  :once
  IF DEFINED MIRROREXE (
    ECHO Found mirror.exe at:
    ECHO   %MIRROREXE%
    ECHO Starting Keybase Mirror
    ECHO   Please IGNORE THE BENIGN ERROR below regarding the failure to add security privilege
    START "KDK Keybase Mirror" /B "%MIRROREXE%" /r K:\ /l C:\Users\%USERNAME%\.kdk\keybase
  ) ELSE (
    ECHO Failed to locate mirror.exe within "c:\Program Files\Dokan"
    ECHO   Will not start keybase mirror for KDK
  )

  BREAK
)

IF "%1"=="stop" (
  ECHO Stopping
  tskill.exe mirror
  BREAK
)
`

func writeMirrorScript(configDir string) (out string, err error) {
	script := []byte(mirrorScript)

	scriptPath := filepath.Join(configDir, "keybase-mirror.cmd")
	err = ioutil.WriteFile(scriptPath, script, 0700)
	if err != nil {
		return "", err
	}
	return scriptPath, nil
}

// Start keybase mirror [windows only]
func StartMirror(configDir string) error {

	keybaseTestDir := filepath.Join(configDir, "keybase", "private")

	// keybase mirror already started.  Nothing to do
	if _, err := os.Stat(keybaseTestDir); err == nil {
		log.Info("Keybase mirror already started")
		return nil
	}

	// Write/Overwrite the mirror script every time
	//   in case it changes up upgrades
	log.Info("Writing keybase mirror script")
	scriptPath, err := writeMirrorScript(configDir)
	if err != nil {
		return err
	}

	commandString := fmt.Sprintf("powershell %s %s", scriptPath, "start")
	log.Debugf("Starting Keybase mirror with command; %s", commandString)
	commandMap := strings.Split(commandString, " ")
	if err := sh.Command(commandMap[0], commandMap[1:]).SetStdin(os.Stdin).Run(); err != nil {
		return err
	}
	return nil
}

// Stop keybase mirror [windows only]
func StopMirror(configDir string) error {
	// TODO(rluckie) Fix StopMirror to work with multiple KDK containers
	// Use docker client to iterate though all running containers and ensure that no containers have mirror dir mounted
	// If so, stop mirror
	log.Info("Writing keybase mirror script")
	scriptPath, err := writeMirrorScript(configDir)
	if err != nil {
		return err
	}
	commandString := fmt.Sprintf("powershell %s %s", scriptPath, "stop")
	log.Debugf("Stopping Keybase mirror with command; %s", commandString)
	commandMap := strings.Split(commandString, " ")
	if err := sh.Command(commandMap[0], commandMap[1:]).SetStdin(os.Stdin).Run(); err != nil {
		return err
	}
	return nil
}

// Get keybase mounts
// Linux & OSX: Detect /keybase
// Windows10: Detect k: and /k
func GetMounts(configRootDir string) (source string, target string, err error) {

	keybaseRoots := []string{"/keybase", "k:", "/k"}
	keybaseTestSubdir := "/private"
	for _, keybaseRoot := range keybaseRoots {
		if absPath, err := filepath.Abs(filepath.Join(keybaseRoot, keybaseTestSubdir)); err == nil {
			if path, err := filepath.EvalSymlinks(absPath); err == nil {
				source := filepath.Dir(path)
				target := "/keybase"

				log.Infof("Detected keybase filesystem at: %v", source)

				prmpt := prompt.Prompt{
					Text:     "Mount your keybase directory within KDK? [y/n] ",
					Loop:     true,
					Validate: prompt.ValidateYorN,
				}
				if result, err := prmpt.Run(); err == nil && result == "y" {
					log.Info("Adding /keybase mount to configuration")
					if runtime.GOOS == "windows" {
						source = filepath.Join(configRootDir, "keybase")
						if _, err := os.Stat(source); os.IsNotExist(err) {
							if err := os.Mkdir(source, 0700); err != nil {
								log.WithField("error", err).Fatalf("Failed to create KDK keybase mirror directory [%s]", source)
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
