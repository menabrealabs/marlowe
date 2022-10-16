package language_test

import (
	"encoding/json"
	"testing"

	m "github.com/menabrealabs/marlowe/language"
)

func TestTypes_CloseContract(t *testing.T) {
	contract := m.Close

	var testC m.Contract = contract
	_, ok := testC.(m.Contract)
	if !ok {
		t.Error("Close does not implement Contract interface.")
	}

	jbytes, err := json.Marshal(contract)
	if err != nil {
		t.Error(err)
	}

	if string(jbytes) != "\"close\"" {
		t.Error("Incorrect JSON format: ", string(jbytes))
	}

	t.Logf("Marshalled JSON: %v", string(jbytes))
}

func TestTypes_LetContract(t *testing.T) {
	// Should generate JSON: {"then":"close","let":"Number","be":1}
	contract := m.Let{
		"Number",
		m.Constant(1),
		m.Close,
	}

	var testC m.Contract = m.Contract(contract)

	_, ok := testC.(m.Contract)
	if !ok {
		t.Error("Let does not implement Contract interface.")
	}

	jbytes, err := json.Marshal(contract)
	if err != nil {
		t.Error(err)
	}

	if string(jbytes) != `{"let":"Number","be":1,"then":"close"}` {
		t.Error("Incorrect JSON format: ", string(jbytes))
	}

	t.Logf("Marshalled JSON: %v", string(jbytes))
}

func TestTypes_IfContract(t *testing.T) {
	// Should generate JSON: {"then":"close","if":{"value":1,"gt":0},"else":"close"}
	contract := m.If{
		m.ValueGT{m.Constant(1), m.Constant(0)},
		m.Close,
		m.Close,
	}

	var testC m.Contract = m.Contract(contract)

	_, ok := testC.(m.Contract)
	if !ok {
		t.Error("Let does not implement Contract interface.")
	}

	jbytes, err := json.Marshal(contract)
	if err != nil {
		t.Error(err)
	}

	if string(jbytes) != `{"if":{"value":1,"gt":0},"then":"close","else":"close"}` {
		t.Error("Incorrect JSON format: ", string(jbytes))
	}

	t.Logf("Marshalled JSON: %v", string(jbytes))
}
