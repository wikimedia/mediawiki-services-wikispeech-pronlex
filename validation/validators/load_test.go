package validators

import (
	"fmt"
	"testing"

	"github.com/stts-se/symbolset"
)

func ss_for_test(name string) (symbolset.SymbolSet, error) {
	syms := []symbolset.Symbol{
		{Desc: "sil", String: "i:", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		{Desc: "sill", String: "I", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		{Desc: "full", String: "u0", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		{Desc: "ful", String: "}:", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		{Desc: "matt", String: "a", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		{Desc: "mat", String: "A:", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		{Desc: "bot", String: "u:", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		{Desc: "bott", String: "U", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		{Desc: "häl", String: "E:", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		{Desc: "häll", String: "E", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		{Desc: "aula", String: "au", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		{Desc: "syl", String: "y:", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		{Desc: "syll", String: "Y", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		{Desc: "hel", String: "e:", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		{Desc: "herr,hett", String: "e", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		{Desc: "nöt", String: "2:", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		{Desc: "mött,förra", String: "9", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		{Desc: "mål", String: "o:", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		{Desc: "moll,håll", String: "O", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		{Desc: "bättre", String: "@", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		{Desc: "europa", String: "eu", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		{Desc: "pol", String: "p", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		{Desc: "bok", String: "b", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		{Desc: "tok", String: "t", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		{Desc: "bort", String: "rt", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		{Desc: "mod", String: "m", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		{Desc: "nod", String: "n", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		{Desc: "dop", String: "d", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		{Desc: "bord", String: "rd", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		{Desc: "fot", String: "k", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		{Desc: "våt", String: "g", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		{Desc: "lång", String: "N", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		{Desc: "forna", String: "rn", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		{Desc: "fot", String: "f", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		{Desc: "våt", String: "v", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		{Desc: "kjol", String: "C", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		{Desc: "fors", String: "rs", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		{Desc: "rov", String: "r", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		{Desc: "lov", String: "l", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		{Desc: "sot", String: "s", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		{Desc: "sjok", String: "x", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		{Desc: "hot", String: "h", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		{Desc: "porla", String: "rl", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		{Desc: "jord", String: "j", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		{Desc: "syllable delimiter", String: ".", Cat: symbolset.SyllableDelimiter, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		{Desc: "accent I", String: `"`, Cat: symbolset.Stress, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		{Desc: "accent II", String: `""`, Cat: symbolset.Stress, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		{Desc: "secondary stress", String: "%", Cat: symbolset.Stress, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		{Desc: "phoneme delimiter", String: " ", Cat: symbolset.PhonemeDelimiter, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		{String: "+", Cat: symbolset.CompoundDelimiter, Desc: "", IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
	}

	return symbolset.NewSymbolSet(name, syms)
}

func TestValidatoFromFile1(t *testing.T) {
	name := "sv_se_nst_test"
	fName := fmt.Sprintf("%s.vd", name)

	ss, err := ss_for_test(name)
	if err != nil {
		t.Errorf("couldn't initialise symbol set : %s", err)
		return
	}

	v, err := LoadValidatorFromFile(ss, fName)
	if err != nil {
		t.Errorf("couldn't load validator from file %s : %s", fName, err)
		return
	}

	nRules := len(v.Rules)
	if nRules != 10 {
		t.Errorf(fsExp, 10, nRules)
	}

	nTests := v.NumberOfTests()
	if nTests != 2 {
		t.Errorf(fsExp, 2, nTests)
	}

}

func TestValidatoFromFile2(t *testing.T) {
	name := "sv_se_nst_test_fail"
	fName := fmt.Sprintf("%s.vd", name)

	ss, err := ss_for_test(name)
	if err != nil {
		t.Errorf("couldn't initialise symbol set : %s", err)
		return
	}

	_, err = LoadValidatorFromFile(ss, fName)
	if err == nil {
		t.Errorf("expected error here")
		return
	}

}
