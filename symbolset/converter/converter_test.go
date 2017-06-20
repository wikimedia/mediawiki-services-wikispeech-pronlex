package converter

import (
	"fmt"
	"testing"

	"github.com/stts-se/pronlex/symbolset"
)

func TestLoadFromDir(t *testing.T) {
	sSets, err := symbolset.LoadSymbolSetsFromDir("../test_data")
	if err != nil {
		t.Errorf("LoadSymbolSetsFromDir() didn't expect error here : %v", err)
		return
	}
	_, testRes, err := LoadFromDir(sSets, "./test_data")
	if err != nil {
		t.Errorf("LoadSymbolSetsFromDir() didn't expect error here : %v", err)
		return
	}
	if !testRes.OK {
		for _, err := range testRes.Errors {
			fmt.Println(err)
		}
		t.Errorf("FAIL")
	}
}
