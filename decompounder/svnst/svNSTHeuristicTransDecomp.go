package svnst

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var char2phon = map[byte]string{

	'a': "(?:a|A:)",
	'b': "b",
	'c': "k",
	'd': "[dj]", // 'djur' d -> j
	'e': "e:?",
	'f': "f",
	'g': "[jg]",
	'h': "h",
	'i': "(?I|i:)",
	'j': "j",
	'k': "k",
	'l': "l",
	'm': "m",
	'n': "n",
	'o': "(?:U|u:|O|o:)",
	'p': "p",
	'q': "k",
	'r': "r",
	's': "[sr]", // standard+skåp rd+s -> rd + rs
	't': "t",
	'u': "\\}:",
	'v': "v",
	'x': "s", // k
}

func canMatchChar(c byte) string {
	var res string = "XXXXZZZXXX"

	if phon, ok := char2phon[c]; ok {
		res = phon
	} else {
		fmt.Fprintf(os.Stderr, "svNSTHeuristicTransDecomp.CanMatchChar: unknown char: '%c'\n", c)
	}

	return res
}

func canMatchLhs(lhs string) string {
	var res string
	if lhs == "" {
		return res
	}

	canMatch := canMatchChar(lhs[len(lhs)-1])
	if canMatch != "" {
		res = canMatch
	}

	return res
}

func canMatchRhs(rhs string) string {
	var res string
	if rhs == "" {
		return res
	}

	canMatch := canMatchChar(rhs[0])
	if canMatch != "" {
		res = canMatch
	}

	return res
}

var nonPhonemes = `[ .%"]+`

func splitTrans(lhs, rhs, trans string) (string, string) {

	lhsM := canMatchLhs(lhs)
	rhsM := canMatchRhs(rhs)

	reStrn := lhsM + "([ .]+)" + nonPhonemes + rhsM
	//fmt.Printf("re: %s\n", reStrn)

	// TODO Panics on incorrect RE
	re := regexp.MustCompile(reStrn)

	indxs := re.FindStringSubmatchIndex(trans)
	if len(indxs) != 4 {
		return "", trans
	}

	t1 := strings.TrimSpace(trans[0:indxs[2]])
	t2 := strings.TrimSpace(trans[indxs[3]:])

	return t1, t2
}

func HerusticSvNSTTransDecomp(orthDecomp, trans string) []string {
	var res []string

	decomps := strings.Split(orthDecomp, "+")
	// Loop over pairs of decomps

	transSoFar := trans
	for i := 0; i < len(decomps)-1; i++ {
		lhs := strings.TrimSpace(decomps[i])
		rhs := strings.TrimSpace(decomps[i+1])
		if lhs == "-" || lhs == "" {
			continue
		}

		tLhs, tRhs := splitTrans(lhs, rhs, transSoFar)
		res = append(res, tLhs)
		transSoFar = tRhs

	}
	fmt.Println()

	return res
}
