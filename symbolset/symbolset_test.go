package symbolset

import (
	"testing"
)

var fsExp = "Expected: '%v' got: '%v'"
var fsDidntExp = "Didn't expect: '%v'"

func testEq(t *testing.T, expect []Symbol, result []Symbol) {
	if len(expect) != len(result) {
		t.Errorf(fsExp, expect, result)
		return
	}
	for i, ex := range expect {
		re := result[i]
		if ex != re {
			t.Errorf(fsExp, expect, result)
			return
		}
	}
}
