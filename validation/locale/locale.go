package locale

import (
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

func SelfName(t language.Tag) string {
	en := display.English.Tags()
	return en.Name(t)
}

//var localeMatcher = language.NewMatcher(SupportedLocales)

func LookUp(s string) (language.Tag, error) {
	for _, t := range SupportedLocales {
		if t.String() == s {
			return t, nil
		}
	}
	return language.Tag{}, fmt.Errorf("no such locale: %s", s)
}
