package svnst

import (
	"fmt"
	"regexp"
	"strings"
)

func canMatchChar(c byte) string {
	var res string = "XXXXZZZXXX"

	switch c {
	case 'a':
		res = "(?:a|A:)"
	case 'b':
		res = "b"
	case 'f':
		res = "f"
	case 'p':
		res = "p"
	case 's':
		res = "s"
	case 't':
		res = "t"
	case 'v':
		res = "v"

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

	return trans[0:indxs[2]], trans[indxs[3]:]

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
