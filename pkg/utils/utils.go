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
	"net"
	"reflect"
)

func Contains(input interface{}, target interface{}) bool {
	rValue := reflect.ValueOf(input)
	rTarget := reflect.ValueOf(target)
	if rValue.IsValid() {
		for i := 0; i < rValue.Len(); i++ {
			rVal := rValue.Index(i)
			if rVal.IsValid() {
				if rVal.Interface() == rTarget.Interface() {
					return true
				}
			}
		}
	}
	return false
}

// Ask kernel for a free port
func GetPort() int {
	if out, err := net.ResolveTCPAddr("tcp", "localhost:0"); err != nil {
		return 0
	} else {
		if listen, err := net.ListenTCP("tcp", out); err != nil {
			defer listen.Close()
			return 0
		} else {
			defer listen.Close()
			return listen.Addr().(*net.TCPAddr).Port
		}
	}
}
