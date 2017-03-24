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
		{"grad", "spår", `"" g r A: d . % s p o: r`, `"" g r A: d`, `% s p o: r`},
		{"sprit", "tillståndet", `"" s p r i: t . t I l . % s t O n . d e t`, `"" s p r i: t`, `t I l . % s t O n . d e t`},
		{"kvicksilver", "utsläpps", `"" k v I k . s I l . v e r . }: t . % s l E p s`, `"" k v I k . s I l . v e r`, `}: t . % s l E p s`},
		{"kräl", "djur", `"" k r E: l . % j }: r`, `"" k r E: l`, `% j }: r`},
		{"militär", "hatet", `m I . l I . "" t {: r . % h A: . t e t`, `m I . l I . "" t {: r`, `% h A: . t e t`},
		//	{"beaufort", "skalans", `b O . "" f o: . % rs k A: . l a n s`, `b O . "" f o: . % rs k A: . l a n s`, `b O . "" f o: . % rs k A: . l a n s`},
		{"standard", "skåpen", `"" s t a n . d a rd . % rs k o: . p e n`, `"" s t a n . d a rd`, `% rs k o: . p e n`},
		{"skridsko", "seglare", `"" s k r I . s k U . % s e: . g l a . r e`, `"" s k r I . s k U`, `% s e: . g l a . r e`},
		{"klo", "beväpnad", `"" k l u: . b e . % v E: p . n a d`, `"" k l u:`, `b e . % v E: p . n a d`},
		{"cicero", "stilar", `"" s i: . s e . r O . % s t i: . l a r`, `"" s i: . s e . r O`, `% s t i: . l a r`},
		{"rokoko", "möbels", `r O . k O . "" k o: . % m 2: . b e l s`, `r O . k O . "" k o:`, `% m 2: . b e l s`},
		{"index", "klausal", `"" I n . d e k s . k l au . % s A: l`, `"" I n . d e k s`, `k l au . % s A: l`},
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
