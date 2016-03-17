package main

import "testing"

var fsExpField = "For field %v, expected: '%v' got: '%v'"
var fsExp = "Expected: '%v' got: '%v'"

func Test_LoadLineFmt(t *testing.T) {
	_, err := loadLineFmt()
	if err != nil {
		t.Errorf("didn't expect error here : %s", err)
	}
}
