package main

import (
	"testing"
)

var m = "wanted '%v' got '%v'"

func Test_t0(t *testing.T) {
	s1 := "% rs rt u: rt"
	w1 := "% s t u: rt"

	if g1 := deRetroflex(s1); w1 != g1 {
		t.Errorf(m, w1, g1)
	}
}

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
		{"standard", "skåpen", `"" s t a n . d a rd . % rs k o: . p e n`, `"" s t a n . d a rd`, `% s k o: . p e n`},
		{"skridsko", "seglare", `"" s k r I . s k U . % s e: . g l a . r e`, `"" s k r I . s k U`, `% s e: . g l a . r e`},
		{"klo", "beväpnad", `"" k l u: . b e . % v E: p . n a d`, `"" k l u:`, `b e . % v E: p . n a d`},
		{"cicero", "stilar", `"" s i: . s e . r O . % s t i: . l a r`, `"" s i: . s e . r O`, `% s t i: . l a r`},
		{"rokoko", "möbels", `r O . k O . "" k o: . % m 2: . b e l s`, `r O . k O . "" k o:`, `% m 2: . b e l s`},
		{"index", "klausal", `"" I n . d e k s . k l au . % s A: l`, `"" I n . d e k s`, `k l au . % s A: l`},
		{"års", "dag", `"" o: rs . % rd A: g`, `"" o: rs`, `% d A: g`},
		// NB: missing retroflexation --- due to '-rr'?:
		{"abborr", "stigen", `"" a . b O r . % s t i: . g e n`, `"" a . b O r`, `% s t i: . g e n`},
		{"alster", "lind", `"" a l . s t e . % rl I n d`, `"" a l . s t e r`, `% l I n d`},
		{"ankar", "lindningar", `"" a N . k a r . % rl I n d . n I N . a r`, `"" a N . k a r`, `% l I n d . n I N . a r`},
		{"marie", "helene", `m a . "" r i: . h e . % l e: n`, `m a . "" r i:`, `h e . % l e: n`},

		// TODO: How to handle silent double chars...?:
		//{"artikel", "löshet", `a . "" rt I . k e . l 2: s . % h e: t`, `a . "" rt I . k e . l 2: s . % h e: t`, `l 2: s . % h e: t`},
	}

	for _, ts := range tsts {
		r1, r2, err := splitTrans(ts.lhs, ts.rhs, ts.tra)
		if err != nil {
			t.Errorf("splitTrans error : %v", err)
		}
		if r1 != ts.res1 {
			t.Errorf(m, ts.res1, r1)
		}
		if r2 != ts.res2 {
			t.Errorf(m, ts.res2, r2)
		}
	}

}
