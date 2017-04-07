package svnst

import (
	"fmt"
	//"os"
	"regexp"
	"strings"
)

var char2phon = map[rune]string{

	'a': "(?:a|A:)",
	'b': "b",
	'c': "[skx]", // charlott c -> x
	//'d': "[djr]", // 'djur' d -> j, års+dag d -> rd
	'd': "[dj]", // 'djur' d -> j
	'e': "(?:e:?|\\}:|E)",
	'é': "e:",
	'f': "f",
	'g': "[jgNx]", // -ing g -> N, gelé g -> x
	'h': "h",
	'i': "(?:I|i:)",
	'j': "[jx]", // Jeanette j -> x
	'k': "[kC]", // kärr k -> C
	//'l': "[lr]", // års+lopp, l -> rl
	'l': "[lj]", // ljud l -> j
	'm': "m",
	'n': "[nN]", // nypon+gränd
	'o': "(?:U|u:|O|o:)",
	'p': "p",
	'q': "k",
	'r': "r",
	//'s': "[srx]", // standard+skåp rd+s -> rd + rs, måls+skytt s -> x
	's': "[sx]", // måls+skytt s -> x
	//'t': "[tr]", // aborr+träsk t -> rt
	't': "[tC]", // tjänst t -> C
	'u': "(?:\\}:|u0)",
	'v': "v",
	'w': "v",
	'x': "s", // k
	'y': "(?:Y|y:)",
	'z': "s",
	'å': "(?:O|o:)",
	'ä': "(?:E|\\{:)",
	'ö': "(?:2:?|9:?)",

	'7': "n", // 7-eleven
	'4': "a", // a4

}

type sm struct {
	suff    string
	matcher string
}

var suffixMatchers = []sm{
	sm{"ie", "I"},           // Annie
	sm{"vue", "v (?:y:|Y)"}, //Bellevue
}

func canMatchPrefix(p string) string {
	res := ""

	fmt.Println(p)

	if p == "anne" {
		return "a [nN]"
	}
	if p == "bernadotte" {
		return "d O t"
	}

	for _, sm := range suffixMatchers {
		if strings.HasSuffix(p, sm.suff) {
			return sm.matcher
		}
	}

	return res
}

func canMatchChar(c rune) (string, error) {
	var res string = "XXXXZZZXXX"
	var err error

	if phon, ok := char2phon[c]; ok {
		res = phon
	} else {
		msg := fmt.Sprintf("svNSTHeuristicTransDecomp.CanMatchChar: unknown char: '%s'\n", string(c))
		//fmt.Fprintf(os.Stderr, msg)
		err = fmt.Errorf("%s", msg)
	}

	return res, err
}

func makeRes(matchStrings ...string) string {
	res := ""
	res0 := []string{}
	for _, m := range matchStrings {
		if m != "" {
			res0 = append(res0, m)
		}
	}

	if len(res0) == 1 {
		return res0[0]
	}
	if len(res0) > 0 {
		res1 := strings.Join(res0, "|")
		res = "(?:" + res1 + ")"
	}

	fmt.Println(res)
	return res
}

func canMatchLhs(lhs string) (string, error) {
	var res string
	var err error
	if lhs == "" {
		return res, err
	}

	pMatcher := canMatchPrefix(lhs)

	runes := []rune(lhs)
	charMatcher, err := canMatchChar(runes[len(runes)-1])
	//if canMatch != "" {
	//	res = canMatch
	//}

	res = makeRes(pMatcher, charMatcher)

	return res, err
}

func canMatchRhs(rhs string) (string, error) {
	var res string
	var err error
	if rhs == "" {
		return res, err
	}

	runes := []rune(rhs)
	canMatch, err := canMatchChar(runes[0])
	if canMatch != "" {
		res = canMatch
	}

	return res, err
}

func endsInPotentialRetroflex(s string) bool {
	if strings.HasSuffix(s, "rs") {
		return true
	}
	if strings.HasSuffix(s, "rt") {
		return true
	}
	if strings.HasSuffix(s, "rts") {
		return true
	}
	if strings.HasSuffix(s, "rd") {
		return true
	}
	if strings.HasSuffix(s, "rds") {
		return true
	}

	if strings.HasSuffix(s, "rn") {
		return true
	}
	if strings.HasSuffix(s, "rns") {
		return true
	}

	if strings.HasSuffix(s, "rl") {
		return true
	}
	if strings.HasSuffix(s, "rls") {
		return true
	}

	if strings.HasSuffix(s, "r") {
		return true
	}

	return false
}

func startsWithPotentialRetroflex(s string) bool {
	if strings.HasPrefix(s, "s") {
		return true
	}
	if strings.HasPrefix(s, "t") {
		return true
	}
	if strings.HasPrefix(s, "d") {
		return true
	}
	if strings.HasPrefix(s, "n") {
		return true
	}
	if strings.HasPrefix(s, "l") {
		return true
	}

	return false
}
func retroflexation(lhs, rhs string) bool {
	lhsRetro := endsInPotentialRetroflex(lhs)
	rhsRetro := startsWithPotentialRetroflex(rhs)

	return lhsRetro && rhsRetro
}

var nonPhoneme = map[string]bool{

	".":  true,
	"%":  true,
	`"`:  true,
	`""`: true,
}

var isRetroflex = map[string]bool{
	"rs": true,
	"rt": true,
	"rd": true,
	"rl": true,
}

// remove initial retroflexation
func deRetroflex(s string) string {
	var res []string
	syms := strings.Split(strings.TrimSpace(s), " ")
	done := false
	for _, sym := range syms {
		if nonPhoneme[sym] {
			res = append(res, sym)
			continue
		}
		if !done && isRetroflex[sym] {
			sym = strings.Replace(sym, "r", "", 1)
			res = append(res, sym)
		} else {
			done = true
			res = append(res, sym)
		}
	}

	return strings.Join(res, " ")
}

var nonPhonemes = `[ .%"]*`

func splitTrans(lhs, rhs, trans string) (string, string, error) {

	lhs0 := lhs

	// If retroflexation spanns the compound boundary, matching is
	// a bit more tricky
	retro := retroflexation(lhs0, rhs)

	// Final -r "jumps" to rhs transcription, thus the character
	// before final 'r' should be used for matching lhs
	if retro && strings.HasSuffix(lhs0, "r") && !strings.HasSuffix(lhs0, "rr") && len(lhs0) > 1 {
		lhs0 = lhs0[0 : len(lhs0)-1]
	}

	lhsM, lErr := canMatchLhs(lhs0)
	// if final -r has been removed from lhs, we still need
	// optional matching of ' r'
	if retro && len(lhs) > len(lhs0) {
		lhsM = lhsM + "(?: r)?"
	}

	rhsM, rErr := canMatchRhs(rhs)

	// TODO This is not how you handle disjunctive errors
	var errr error
	if lErr != nil {
		if errr != nil {
			errr = fmt.Errorf("%s : %s", errr, lErr)
		} else {
			errr = fmt.Errorf("%s", lErr)
		}
	}
	if rErr != nil {
		if errr != nil {
			errr = fmt.Errorf("%v : %v", errr, rErr)
		} else {
			errr = fmt.Errorf("%v", rErr)
		}
	}

	if retro {
		rhsM = "r?" + rhsM
	}

	reStrn := lhsM + "([ .]+)" + nonPhonemes + rhsM

	fmt.Println(reStrn)

	// TODO Panics on incorrect RE
	re := regexp.MustCompile(reStrn)

	indxs := re.FindStringSubmatchIndex(trans)
	if len(indxs) != 4 {
		return "", trans, errr
	}

	t1 := strings.TrimSpace(trans[0:indxs[2]])
	t2 := strings.TrimSpace(trans[indxs[3]:])

	if retro {
		// We might have to add final ' r' to lhs if it has
		// hopped compound boundary:
		if strings.HasSuffix(lhs, "r") && !strings.HasSuffix(t1, " r") {
			t1 = t1 + " r"
		}

		t2 = deRetroflex(t2)
	}
	return t1, t2, errr
}

func HerusticSvNSTTransDecomp(orthDecomp, trans string) ([]string, error) {
	var res []string
	var err error
	decomps := strings.Split(orthDecomp, "+")
	// Loop over pairs of decomps

	transSoFar := trans
	for i := 0; i < len(decomps)-1; i++ {
		lhs := strings.TrimSpace(decomps[i])
		rhs := strings.TrimSpace(decomps[i+1])
		if lhs == "-" || lhs == "" {
			continue
		}

		tLhs, tRhs, err0 := splitTrans(lhs, rhs, transSoFar)
		if err0 != nil {
			if err == nil {
				err = err0
			} else {
				err = fmt.Errorf("%v : %v", err, err0)
			}
		}
		if tLhs != "" {
			res = append(res, tLhs)
		}
		transSoFar = tRhs
	}
	res = append(res, transSoFar)
	//fmt.Println()

	if len(res) != len(decomps) {
		msg := "failed to match decomps and transparts"
		if err == nil {
			err = fmt.Errorf(msg)
		} else {
			err = fmt.Errorf("%v : %s", err, msg)
		}
	}

	return res, err
}
