package symbolset

import (
	"testing"
)

var fsExp = "Expected: '%v' got: '%v'"
var fsDidntExp = "Didn't expect: '%v'"

func Test_TODO_IMPLEMENT_TESTS_HERE(t *testing.T) {
	expect := "a"
	result := expect
	if result != expect {
		t.Errorf(fsExp, expect, result)
	}

}
