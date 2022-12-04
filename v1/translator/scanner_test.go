package translator_test

import (
	"strings"
	"testing"

	scan "github.com/menabrealabs/marlowe/v1/translator"
)

func testScanner(str string) []scan.Token {
	bufReader := strings.NewReader(str)
	scanner := scan.NewScanner(bufReader)
	var out []scan.Token

	for {
		token := scanner.Scan()
		out = append(out, token)
		if token.Type == scan.EOF {
			break
		}
	}

	return out
}

func TestValidKeywords(t *testing.T) {
	var keywords = []string{
		// Contracts
		"Let", "When", "If", "Pay", "Assert", "Close",
		//Actions
		"Deposit", "Notify", "Choice", "ChoiceId", "Bound",
		//Values
		"AvailableMoney", "Constant", "NegValue", "AddValue", "SubValue", "MulValue", "DivValue",
		"ChoiceValue", "TimeIntervalValue", "UseValue", "Cond",
		// Observations
		"AndObs", "OrObs", "NotObs", "ChoseSomething", "ValueGE", "ValueGT", "ValueLE", "ValueLT", "ValueEQ", "TrueObs", "FalseObs",
	}

	tokens := testScanner(strings.Join(keywords, " "))

	for i, val := range tokens {
		if tokens[i].Type == scan.EOF {
			continue
		}

		if tokens[i].Type == scan.INVALID {
			t.Errorf("FAILED: %v is an invalid keyword", tokens[i].Value)
		}

		if tokens[i].Type != scan.KEYWORD {
			t.Errorf("Failed to tokenize keyword.\nExpected: %v\nGot: %v", val, tokens[i].Value)
		}
	}
}

func TestInvalidKeywords(t *testing.T) {
	tokens := testScanner("InvalidKeyword")

	if tokens[0].Type != scan.INVALID {
		t.Errorf("Failed to identify invalid keyword: InvalidKeyword")
	}
}

func TestValidIntegers(t *testing.T) {
	ints := []string{"1454", "4848844032", "2223454"}
	input := strings.Join(ints, " ")
	tokens := testScanner(input)

	for i, num := range ints {
		if tokens[i].Type != scan.INT || tokens[i].Value != num {
			t.Errorf("Failed to tokenize integer.\nExpected: INT\nGot: %v", tokens[i])
		}
	}
}

func TestInvalidIntegers(t *testing.T) {
	ints := []string{"123.", "1_000", "0b01", "0xff"}
	input := strings.Join(ints, " ")
	tokens := testScanner(input)

	for i := range ints {
		if tokens[i].Type != scan.INVALID {
			t.Errorf("Failed to identify invalid integer\nExpected: INVALID\nGot: %v", tokens[i])
		}
	}
}

func TestValidStrings(t *testing.T) {
	strs := []string{"\"name\"", "\"Buyer\"", "\"L337\"", "\"LeFt & Right3\""}
	input := strings.Join(strs, " ")
	tokens := testScanner(input)

	for i, str := range strs {
		if tokens[i].Type != scan.STRING && tokens[i].Value != str {
			t.Errorf("Failed to tokenize string.\nExpected: STRING %v\nGot: %v", str, tokens[i])
		}
	}
}

func TestValidParentheses(t *testing.T) {
	tokens := testScanner("( )")

	if tokens[0].Type != scan.PARENS_L {
		t.Error("Failed to tokenize left/open parenthesis")
	}

	if tokens[1].Type != scan.PARENS_R {
		t.Error("Failed to tokenize right/close parenthesis")
	}
}

func TestValidSquareBrackets(t *testing.T) {
	tokens := testScanner("[ ]")

	if tokens[0].Type != scan.SQUARE_L {
		t.Error("Failed to tokenize left/open square bracket")
	}

	if tokens[1].Type != scan.SQUARE_R {
		t.Error("Failed to tokenize right/close square bracket")
	}
}

func TestValidComma(t *testing.T) {
	tokens := testScanner(",")

	if tokens[0].Type != scan.COMMA {
		t.Error("Failed to tokenize comma")
	}
}

func TestValidNewlineReset(t *testing.T) {
	tokens := testScanner("( )\n[ ]")

	if tokens[1].Position.Column != 3 {
		t.Errorf("Failed to reset newline.\nColumn expected: 3\nColumn got: %v", tokens[1].Position.Column)
	}

	if tokens[1].Position.Line != 1 {
		t.Errorf("Failed to reset newline.\nLine expected: 1\nLine got: %v", tokens[1].Position.Line)
	}

	if tokens[2].Position.Column != 1 {
		t.Errorf("Failed to reset newline.\nColumn expected: 1\nColumn got: %v", tokens[2].Position.Column)
	}

	if tokens[2].Position.Line != 2 {
		t.Errorf("Failed to reset newline.\nLine expected: 2\nLine got: %v", tokens[2].Position.Line)
	}
}
