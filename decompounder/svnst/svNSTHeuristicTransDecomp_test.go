package svnst

import (
	"fmt"
	"testing"
)

func f() { fmt.Println() }

var m = "wanted '%v' got '%v'"

func Test_t1(t *testing.T) {

	var tsts = []struct {
		lhs  string
		rhs  string
		tra  string
		res1 string
		res2 string
	}{
		{"upp", "fällde", `"" u0 p . % f E l . d e`, `"" u0 p`, `% f E l . d e`},
		{"luft", "strömningarna", `"" l u0 f t . % s t r 2 m . n I N . a . rn a`, `"" l u0 f t`, `% s t r 2 m . n I N . a . rn a`},
	}

	for _, ts := range tsts {
		r1, r2 := splitTrans(ts.lhs, ts.rhs, ts.tra)
		if r1 != ts.res1 {
			t.Errorf(m, ts.res1, r1)
		}
		if r2 != ts.res2 {
			t.Errorf(m, ts.res2, r2)
		}

	}

}
