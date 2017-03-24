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
	'c': "k",
	'd': "[dj]", // 'djur' d -> j
	'e': "e:?",
	'é': "e:",
	'f': "f",
	'g': "[jg]",
	'h': "h",
	'i': "(?:I|i:)",
	'j': "j",
	'k': "k",
	'l': "l",
	'm': "m",
	'n': "[nN]", // nypon+gränd
	'o': "(?:U|u:|O|o:)",
	'p': "p",
	'q': "k",
	'r': "r",
	's': "[sr]", // standard+skåp rd+s -> rd + rs
	't': "t",
	'u': "\\}:",
	'v': "v",
	'w': "v",
	'x': "s", // k
	'y': "y:",
	'z': "s",
	'å': "O",
	'ä': "E",
	'ö': "2:",

	'7': "n", // 7-eleven
	'4': "a", // a4

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

func canMatchLhs(lhs string) (string, error) {
	var res string
	var err error
	if lhs == "" {
		return res, err
	}

	runes := []rune(lhs)
	canMatch, err := canMatchChar(runes[len(runes)-1])
	if canMatch != "" {
		res = canMatch
	}

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

var nonPhonemes = `[ .%"]+`

func splitTrans(lhs, rhs, trans string) (string, string, error) {

	lhsM, lErr := canMatchLhs(lhs)
	rhsM, rErr := canMatchRhs(rhs)

	reStrn := lhsM + "([ .]+)" + nonPhonemes + rhsM
	//fmt.Printf("re: %s\n", reStrn)

	// TODO Panics on incorrect RE
	re := regexp.MustCompile(reStrn)

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

	indxs := re.FindStringSubmatchIndex(trans)
	if len(indxs) != 4 {
		return "", trans, errr
	}

	t1 := strings.TrimSpace(trans[0:indxs[2]])
	t2 := strings.TrimSpace(trans[indxs[3]:])

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
		res = append(res, tLhs)
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
