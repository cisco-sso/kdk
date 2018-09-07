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
	"testing"
)

func TestSliceContainsString(t *testing.T) {

	strSlice := []string{
		"Hello",
	}
	result := Contains(strSlice, "World")
	if result {
		t.Log("Contains reports having non existing value.")
		t.FailNow()
	}
	strSlice = append(strSlice, "World")
	result = Contains(strSlice, "World")
	if !result {
		t.Log("Contains reports not containing existing value.")
		t.FailNow()
	}
}

func TestSliceContainsInt(t *testing.T) {

	strSlice := []int{
		1,
	}
	result := Contains(strSlice, 2)
	if result {
		t.Log("Contains reports having non existing value.")
		t.FailNow()
	}
	strSlice = append(strSlice, 2)
	result = Contains(strSlice, 2)
	if !result {
		t.Log("Contains reports not containing existing value.")
		t.FailNow()
	}
}

func TestSliceContainsCustom(t *testing.T) {

	type Custom struct {
		StrVal string
		IntVal int
	}

	expected := Custom{
		StrVal: "World",
		IntVal: 2,
	}
	inputArr := []Custom{
		{
			StrVal: "Hello",
			IntVal: 1,
		},
	}
	result := Contains(inputArr, expected)
	if result {
		t.Log("Contains reports having non existing value.")
		t.FailNow()
	}

	inputArr = append(inputArr, expected)
	result = Contains(inputArr, expected)
	if !result {
		t.Log("Contains reports not containing existing value.")
		t.FailNow()
	}

}
