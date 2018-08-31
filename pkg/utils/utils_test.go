package utils

import "testing"

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
