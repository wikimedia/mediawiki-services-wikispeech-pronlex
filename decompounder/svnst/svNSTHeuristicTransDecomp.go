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
	'c': "[skxC]", // charlott c -> x, check c -> C
	'd': "[dj]",   // 'djur' d -> j
	'e': "(?:e:?|\\}:|E)",
	'é': "e:",
	'f': "f",
	'g': "[jgNx]", // -ing g -> N, gelé g -> x
	'h': "[hj]",   // hjort h -> j
	'i': "(?:I|i:)",
	'j': "[jx]", // Jeanette j -> x
	'k': "[kC]", // kärr k -> C
	'l': "[lj]", // ljud l -> j
	'm': "m",
	'n': "[nN]", // nypon+gränd
	'o': "(?:U|u:|O|o:)",
	'p': "p",
	'q': "k",
	'r': "r",
	's': "[sx]", // måls+skytt s -> x
	't': "[tC]", // tjänst t -> C
	'u': "(?:\\}:|u0)",
	'v': "v",
	'w': "v",
	'x': "s", // k
	'y': "(?:Y|y:)",
	'z': "s",
	'å': "(?:O|o:)",
	'ä': "(?:E:?|\\{:?|e)",
	'ö': "(?:2:?|9:?)",

	'7': "n", // 7-eleven
	'4': "a", // a4

	//TODO: Is this OK?
	'-': "",
}

type sm struct {
	suff    string
	matcher string
}

var suffixMatchers = []sm{
	sm{"ampere", "p {: r"},
	sm{"baby", "b I"},
	sm{"baisse", "b E: s"},
	sm{"bearnaise", "n E: s"},
	sm{"bayonne", "j O n"},
	sm{"bordeaux", "rd o:"},
	sm{"bourgogne", "g O n j"},
	sm{"braille", "b r a j"},
	sm{"boy", "b O j"},
	sm{"champagne", "p a n j"},
	sm{"college", "l I t C"},
	sm{"collage", "l A: rs"},
	sm{"colli", "k O [.] l I"},
	sm{"cockney", "k n I"},
	sm{"dessert", "s {: r"},
	sm{"nnie", "I"}, // Annie
	sm{"brie", "b r i:"},
	sm{"vue", "v (?:y:|Y)"}, //Bellevue
	sm{"vägen", "v E: [.] g e n"},
	sm{"bridge", "I (?:r[dt] rs|t C)"},
	sm{"frey", "E j"},
	//sm{"arbitrage", "A: rs"},
	//sm{"garage", "A: rs"},
	//sm{"plantage", "A: rs"},
	sm{"age", "A: rs"},
	sm{"allonge", "N rs"},
	sm{"geneve", "n E: v"},
	sm{"gustaf", "s t A: v"},
	sm{"hav", "h a f"},
	sm{"horney", "rn I"},
	sm{"ville", "v I l"},
	sm{"jacob", "k O p"},
	sm{"jakob", "k O p"},
	sm{"suu", "s u:"},
	sm{"kai", "k a j"},
	sm{"mai", "m a j"},
	sm{"may", "m a j"},
	sm{"thai", "t a j"},
	sm{"rose", "r o: s"}, // 'rose-marie'
	sm{"konsert", "s {: r"},
	sm{"träd", "t r E:?"},
	sm{"marie", "r (?:i:|I)"},
	sm{"ry", "r I"}, // 'mary'
	sm{"sch", "rs"}, // 'marsch'
	sm{"oecd", "d e:"},
	sm{"ph", "f"},
	sm{"renault", "n o:"},
	sm{"service", "I s"},
	sm{"skrid", "s k r I"},
	sm{"svend", "s v e n"},
	sm{"blind", "b l I n"},
	sm{"bond", "b U n"},
	sm{"stad", "s t a"}, // 'Vrigstad'
	sm{"zenith", "n I t"},
	sm{"ou", "u:"},
	sm{"ai", "a j"}, //Nicolai
}

func canMatchPrefix(p string) string {
	res := ""

	//fmt.Println(p)

	if p == "hundra" {
		return "h u0 n . d r a"
	}
	if p == "åttio" {
		return "O . t I . U"
	}

	if p == "g" { // 'g+dur'
		return "g e:"
	}
	if p == "hp" { // 'hp+färg'
		return "h o: . p e:"
	}
	if p == "pk" { // 'pk+banken'
		return "p e: . k o:"
	}
	if p == "tt" { // 'tt+afp'
		return "t e: . t e:"
	}
	if p == "tv" { // 'tv+shop'
		return "t e: . v e:"
	}

	if p == "och" { // 'berg+och+dal+bana'
		return "O(?: k)?"
	}

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

	//fmt.Println(res)
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
	if strings.HasPrefix(s, "c") { //motor+Cykel
		return true
	}
	if strings.HasPrefix(s, "z") { //buffert+zon
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
	"rn": true,
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

// doubleChar returns true iff lhs ends with the same rune as rhs
// starts with
func doubleChar(lhs, rhs string) (rune, bool) {
	fmt.Printf("LHS: %s\tRHS: %s\n", lhs, rhs)

	var res rune
	var isSame bool

	if lhs == "" || rhs == "" {
		return res, isSame
	}

	lhsRunes := []rune(lhs)
	lastLhs := lhsRunes[len(lhsRunes)-1]

	rhsRunes := []rune(rhs)
	firstRhs := rhsRunes[0]

	if lastLhs == firstRhs {
		return lastLhs, true
	}

	return res, isSame
}

func splitTrans(lhs, rhs, trans string) (string, string, error) {

	lhs = strings.ToLower(lhs)
	rhs = strings.ToLower(rhs)
	lhs0 := lhs

	//TODO: Does this work? Should replace strings with []rune, to be shure that things work?

	//fmt.Println("lhs: " + lhs)
	//fmt.Println("rhs: " + rhs)

	// If retroflexation spanns the compound boundary, matching is
	// a bit more tricky
	retro := retroflexation(lhs0, rhs)
	//fmt.Printf("retro: %v\n", retro)

	// Final -r "jumps" to rhs transcription, thus the character
	// before final 'r' should be used for matching lhs
	//TODO: Does this work? Should replace strings with []rune, to be shure that things work?
	if retro && strings.HasSuffix(lhs0, "r") && !strings.HasSuffix(lhs0, "rr") && len(lhs0) > 1 {
		lhs0 = lhs0[0 : len(lhs0)-1]
	}

	//TODO: Handling double chars over compound boundaries more complicated than I first though.
	// One reason is that the last char of lhs is not always silent.

	// Chech whether lhs ends with the same rune that rhs starts
	// with. If so, the last char of lhs may be silent in the
	// transcription of lhs
	// _, isDouble := doubleChar(lhs, rhs)
	// //fmt.Printf(">>>> %v %v\n", char, isDouble)
	// if !retro && isDouble {
	// 	// Knock of the last char, and hope that is enough
	// 	lhs0 = lhs0[0 : len(lhs0)-1]
	// }

	lhsM, lErr := canMatchLhs(lhs0)
	//fmt.Printf("lhsM: %s\n", lhsM)
	//fmt.Printf("lErr: %v\n", lErr)

	// if final -r has been removed from lhs, we still need
	// optional matching of ' r'
	if retro && len(lhs) > len(lhs0) {
		lhsM = lhsM + "(?: r)?"
	}
	//	if !retro && isDouble && len(lhs) > len(lhs0) {
	//		lhsM = lhsM + "(?: " + +")?"
	//	}

	rhsM, rErr := canMatchRhs(rhs)

	// TODO This is not how you handle disjunctive errors
	var errr error
	if lErr != nil {
		if errr != nil {
			errr = fmt.Errorf("%v : %v", errr, lErr)
		} else {
			errr = fmt.Errorf("%v", lErr)
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

	//fmt.Println(reStrn)

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
	//t2 = deRetroflex(t2)
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
