package line

import (
	//"github.com/stts-se/pronlex/lex"

	"testing"
)

func Test_NewWS(t *testing.T) {

	nws, err := NewWS()
	if err != nil {
		t.Errorf("My heart bleeds: %v", err)
	}
	_ = nws // Hooray!
}
