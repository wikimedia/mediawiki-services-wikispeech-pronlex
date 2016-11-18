package decompounder

import (
	"fmt"
	"strings"
	"testing"
)

var ts = "Wanted '%v' got '%v'\n"

func spunk() { fmt.Println() }

func Test_Tree(t *testing.T) {

	tr := NewtNode()

	if want, got := rune(0), tr.r; want != got {
		t.Errorf(ts, want, got)
	}

	tr = tr.add("strut")
	tr = tr.add("strutnos")
	tr = tr.add("strutnosar")

	// for k, v := range tr.sons {
	//     fmt.Printf("HOJSAN: %#v : %s\n", k, string(v.r))
	// }

	if want, got := rune(0), tr.r; want != got {
		t.Errorf(ts, want, got)
	}

	if want, got := 1, len(tr.sons); want != got {
		t.Errorf(ts, want, got)
	}

	s1 := "strutnosarna"
	prfs := tr.prefixes(s1)
	//fmt.Printf("Arcs: %#v\n", prfs)
	if want, got := 3, len(prfs); want != got {
		t.Errorf(ts, want, got)
	}

	// for _, p := range prfs {
	//     fmt.Printf("PREFIX: %v\n", s1[p.start:p.end])
	// }

	pt := NewPrefixTree()
	pt.Add("ap")
	pt.Add("hund")
	pt.Add("aphund")
	pt.Add("nos")

	s := "aphundar"
	arczz := pt.Prefixes(s)
	if want, got := 2, len(arczz); want != got {
		t.Errorf(ts, want, got)
	}

	prefs := map[string]bool{s[arczz[0].start:arczz[0].end]: true, s[arczz[1].start:arczz[1].end]: true}
	if _, ok := prefs["ap"]; !ok {
		t.Errorf(ts, "'ap'", "nothing")
	}
	if _, ok := prefs["aphund"]; !ok {
		t.Errorf(ts, "'aphund'", "nothing")
	}

	st := NewSuffixTree()

	st.Add("rot")
	st.Add("mos")
	st.Add("nos")

	z := "skrotmos"
	arcs := st.Suffixes(z)
	if want, got := 1, len(arcs); want != got {
		t.Errorf(ts, want, got)
	}

	st.Add("rotmos")
	arcs = st.Suffixes(z)
	//fmt.Printf("ARKZ: %#v\n", arcs)
	if want, got := 2, len(arcs); want != got {
		t.Errorf(ts, want, got)
	}

	suffs := map[string]bool{z[arcs[0].start:arcs[0].end]: true, z[arcs[1].start:arcs[1].end]: true}
	if _, ok := suffs["mos"]; !ok {
		t.Errorf(ts, "'mos'", "nothing")
	}
	if _, ok := suffs["rotmos"]; !ok {
		t.Errorf(ts, "'rotmos'", "nothing")
	}

}

func Test_Paths(t *testing.T) {

	a1 := arc{start: 0, end: 3}
	a2 := arc{start: 3, end: 7}

	res := paths([]arc{a1, a2}, 0, 7)

	if want, got := 1, len(res); want != got {
		t.Errorf("NOOOO! %d %d", want, got)
	}
	p := res[0]
	if want, got := 2, len(p); want != got {
		t.Errorf("AAAA! %d %d", want, got)

	}
	a1_ := p[0]
	if want, got := 0, a1_.start; want != got {
		t.Errorf("AAAA! %d %d", want, got)

	}
	if want, got := 3, p[1].start; want != got {
		t.Errorf("AAAA! %d %d", want, got)

	}

	a3 := arc{start: 3, end: 5}
	a4 := arc{start: 5, end: 7}
	a5 := arc{start: 3, end: 6}

	res = paths([]arc{a1, a2, a3, a4, a5}, 0, 7)
	if want, got := 2, len(res); want != got {
		t.Errorf("Suck %d %d", want, got)
	}
	//fmt.Printf("\n%#v\n", res)
}

func Test_Decompounder(t *testing.T) {

	d := NewDecompounder()

	d.AddPrefix("sylt")
	d.AddPrefix("syl")

	d.AddSuffix("järn")
	d.AddSuffix("tjärn")

	decomps := d.Decomp("syltjärn")
	if w, g := 2, len(decomps); w != g {
		t.Errorf(ts, w, g)
	}
	if w, g := 2, len(decomps[0]); w != g {
		t.Errorf(ts, w, g)
	}
	if w, g := 2, len(decomps[1]); w != g {
		t.Errorf(ts, w, g)
	}

	p1 := decomps[0][0]
	p2 := decomps[1][0]

	if p1 == p2 {
		t.Error("Aouch")
	}

	if p1 != "syl" && p2 != "syl" {
		t.Error("Aouch")
	}
	if p1 != "sylt" && p2 != "sylt" {
		t.Error("Aouch")
	}
}

func Test_Decomp_RecursivePrefixes(t *testing.T) {

	decomp := NewDecompounder()
	decomp.AddPrefix("svavel")
	decomp.AddPrefix("kanin")

	decomp.AddSuffix("förening")

	ds1 := decomp.Decomp("svavelkaninförening")
	//ds1 := decomp.Decomp("svavelförening")
	if w, g := 1, len(ds1); w != g {

		t.Errorf(ts, w, g)
	}

	decomp.AddSuffix("kanin")

	ds2 := decomp.Decomp("kaninkanin")
	if w, g := 1, len(ds2); w != g {

		t.Errorf(ts, w, g)
	}

	ds3 := decomp.Decomp("kaninkaninkaninkaninkanin")
	if w, g := 1, len(ds3); w != g {
		t.Errorf(ts, w, g)

	}
	if w, g := 5, len(ds3[0]); w != g {
		t.Errorf(ts, w, g)
	}

	ds4 := decomp.Decomp("kaninkaninsvavelkaninkanin")
	if w, g := 1, len(ds4); w != g {
		t.Errorf(ts, w, g)

	}
	if w, g := 5, len(ds4[0]); w != g {
		t.Errorf(ts, w, g)
	}

	// Oh my... the following test was made to cath an
	// over-generation error, due to the fact that a prefix
	// initially was allowed to end at the end of the input
	// string. This was changed, so that a prefix must end before
	// the end of the input string.

	decomp.AddPrefix("k")
	decomp.AddPrefix("a")
	decomp.AddPrefix("ka")
	decomp.AddPrefix("kan")
	decomp.AddPrefix("nin")
	decomp.AddPrefix("in")
	decomp.AddPrefix("i")
	decomp.AddPrefix("n")

	ds5 := decomp.Decomp("kaninkanin")
	unique := make(map[string]bool)
	for _, d0 := range ds5 {
		d := strings.Join(d0, "+")
		if unique[d] {
			fmt.Printf("DARN! %v\n", d)
		} else {
			unique[d] = true
		}
	}
	if w, g := len(unique), len(ds5); w != g {
		t.Errorf(ts, w, g)
	}

	n3 := "xnikolaj3000"

	decomp.AddPrefix(n3)
	ds6 := decomp.Decomp(n3)
	if w, g := 0, len(ds6); w != g {
		t.Errorf(ts, w, g)
	}
	p6 := decomp.prefixes.Prefixes(n3)
	if w, g := 0, len(p6); w != g {
		t.Errorf(ts, w, g)
	}
	p6b := decomp.prefixes.RecursivePrefixes(n3)
	if w, g := 0, len(p6b); w != g {
		t.Errorf(ts, w, g)
	}
	p6b2 := decomp.prefixes.RecursivePrefixes(n3 + n3)
	if w, g := 1, len(p6b2); w != g {
		t.Errorf(ts, w, g)
	}

	decomp.AddSuffix(n3)
	ds7 := decomp.Decomp(n3)
	if w, g := 0, len(ds7); w != g {
		t.Errorf(ts, w, g)
	}

	s7 := decomp.suffixes.Suffixes(n3)
	if w, g := 0, len(s7); w != g {
		t.Errorf(ts, w, g)
	}

	ds8 := decomp.Decomp(n3 + n3)
	if w, g := 1, len(ds8); w != g {
		t.Errorf(ts, w, g)
	}

	fmt.Println("OBS: bugg: alen -> ale+n (stockholmsfinalen)")
}

func Test_Alen(t *testing.T) {
	decomp := NewDecompounder()
	decomp.AddPrefix("ale")
	decomp.AddPrefix("n")
	decomp.AddPrefix("fin")
	decomp.AddPrefix("stockholms")

	decomp.AddSuffix("finalen")

	ds1 := decomp.Decomp("alen")
	if w, g := 0, len(ds1); w != g {
		t.Errorf(ts, w, g)
	}

	p1 := decomp.prefixes.Prefixes("alen")
	if w, g := 1, len(p1); w != g {
		t.Errorf(ts, w, g)
	}
	w := arc{start: 0, end: 3}
	g := p1[0]
	if w != g {
		t.Errorf(ts, w, g)
	}

	ds2 := decomp.Decomp("stockholmsfinalen")
	if w, g := 1, len(ds2); w != g {
		t.Errorf(ts, w, g)
	}
	if w, g := 2, len(ds2[0]); w != g {
		t.Errorf(ts, w, g)
	}

}

func Test_LenSort(t *testing.T) {

	// return the versions with fewest compound parts first

	decomp := NewDecompounder()
	decomp.AddPrefix("odalbonde")
	decomp.AddPrefix("bonde")
	decomp.AddPrefix("odal")

	decomp.AddSuffix("husbil")
	decomp.AddSuffix("bil")

	decomp.AddPrefix("hus")

	ds1 := decomp.Decomp("odalbondehusbil")
	soFar := 0
	for _, d := range ds1 {
		if len(d) < soFar {
			t.Errorf("this thingy is not sorted: %#v", ds1)
			return
		}
		soFar = len(d)
	}

}

func Test_InfixS(t *testing.T) {

	decomp := NewDecompounder()
	decomp.AddPrefix("finland")
	decomp.AddSuffix("båt")

	ds1 := decomp.Decomp("finlandbåt")
	if w, g := 1, len(ds1); w != g {
		t.Errorf(ts, w, g)
	}
	if w, g := 2, len(ds1[0]); w != g {
		t.Errorf(ts, w, g)
	}

	decomp.AddInfix("s")

	ds2 := decomp.Decomp("finlandsbåt")
	if w, g := 1, len(ds2); w != g {
		t.Errorf(ts, w, g)
	}
	if w, g := 3, len(ds2[0]); w != g {
		t.Errorf(ts, w, g)
	}

	if w, g := "s", ds2[0][1]; w != g {
		t.Errorf(ts, w, g)
	}

}
