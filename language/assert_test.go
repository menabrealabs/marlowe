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
		t.Error("Incorrect JSON format: ", string(jbytes))
	}

	t.Logf("Marshalled JSON: %v", string(jbytes))
}
