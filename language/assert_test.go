package language_test

import (
	"encoding/json"
	"testing"
)

func assertJson[T any](t *testing.T, contract T, target string) {
	jbytes, err := json.Marshal(contract)
	if err != nil {
		t.Error(err)
	}

	if string(jbytes) != target {
		t.Errorf("%v [Expected]", target)
		t.Errorf("%v [Got]", string(jbytes))
	} else {
		t.Logf("Marshalled JSON: %v", string(jbytes))
	}
}
