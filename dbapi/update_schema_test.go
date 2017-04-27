package dbapi

import (
	"testing"
)

func Test_UpdateSchema(t *testing.T) {

	if err := UpdateSchema("DSDSDS"); err != nil {
		t.Errorf("NOOOOOOOOOOOOOOOOOOOO: %v", err)
	}

}
