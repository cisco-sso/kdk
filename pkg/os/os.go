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

package os

import (
	log "github.com/Sirupsen/logrus"

	"os"
	"runtime"
)

var (
	CurrentUser string
	LineSep     string
	TmpDir      string
)

func init() {
	CurrentUser = getCurrentUser()
	LineSep = getLineSeparator()
	TmpDir = getTmpDir()
}

func getCurrentUser() string {
	user, userOk := os.LookupEnv("USER")
	if userOk {
		return user
	}
	username, usernameOk := os.LookupEnv("USERNAME")
	if usernameOk {
		return username
	}

	log.Fatal("Unable to identify current USER or USERNAME")
	return ""
}

func getLineSeparator() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	} else {
		return "\n"
	}
}

func getTmpDir() string {
	if runtime.GOOS == "windows" {
		tmp, tmpOk := os.LookupEnv("TMP")
		if tmpOk {
			return tmp
		}
		temp, tempOk := os.LookupEnv("TEMP")
		if tempOk {
			return temp
		}
		log.Fatal("Unhandled code path")
	} else {
		return "/tmp"
	}

	log.Fatal("Unhandled code path")
	return ""
}
