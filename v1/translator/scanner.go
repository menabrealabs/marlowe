// Copyright 2022 Menabrea Labs Inc.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// The Marlowe translator translates Marlowe script code into the Go internal
// representation (IR) as defined in the marlowe/v1/language package.
package translator

import (
	"bufio"
	"errors"
	"io"
	"unicode"
)

type TokenType uint8

const (
	EOF = iota
	INVALID
	KEYWORD  // Pay, When, If, etc.
	STRING   // "name", etc. (comma delimited string)
	INT      // Marlowe only supports integer numeric values
	PARENS_L // (
	PARENS_R // )
	SQUARE_L // [
	SQUARE_R // ]
	COMMA    // ,
)

var tokens = [...]string{
	EOF:      "EOF",
	INVALID:  "INVALID",
	KEYWORD:  "KEYWORD",
	STRING:   "STRING",
	INT:      "INT",
	PARENS_L: "(",
	PARENS_R: ")",
	SQUARE_L: "[",
	SQUARE_R: "]",
	COMMA:    ",",
}

var validKeywords = [...]string{
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

func (t TokenType) String() string {
	return tokens[t]
}

type Token struct {
	Type     TokenType
	Value    string
	Position Position
}

type Position struct {
	Line   int
	Column int
}

type Scanner struct {
	position Position
	reader   *bufio.Reader
}

func NewScanner(reader io.Reader) *Scanner {
	return &Scanner{
		position: Position{Line: 1, Column: 0},
		reader:   bufio.NewReader(reader),
	}
}

func (scan *Scanner) Scan() Token {
	for {
		rune, _, err := scan.reader.ReadRune()

		// Return EOF when we get an io.EOF from the reader
		if err == io.EOF {
			return Token{Type: EOF}
		}

		// Panic on any other unhandled error
		if err != nil {
			panic(err)
		}

		scan.position.Column++

		switch rune {
		case '\n':
			scan.resetPosition()
		case '(':
			return Token{Type: PARENS_L, Value: "(", Position: scan.position}
		case ')':
			return Token{Type: PARENS_R, Value: ")", Position: scan.position}
		case '[':
			return Token{Type: SQUARE_L, Value: "[", Position: scan.position}
		case ']':
			return Token{Type: SQUARE_R, Value: "]", Position: scan.position}
		case ',':
			return Token{Type: COMMA, Value: ",", Position: scan.position}
		case '"':
			// Scan the string including quotes
			scan.backup()
			str := scan.str()
			return Token{Type: STRING, Value: str, Position: scan.position}
		default:
			// Ignore spaces
			if unicode.IsSpace(rune) {
				continue
			}

			// Tokenize INT
			if unicode.IsDigit(rune) {
				scan.backup()
				num, err := scan.integer()

				if err != nil {
					return Token{Type: INVALID, Value: num, Position: scan.position}
				}

				return Token{Type: INT, Value: num, Position: scan.position}
			}

			// Tokenize keywords
			if unicode.IsLetter(rune) {
				scan.backup()
				kw := scan.keyword()
				if scan.isValidKeyword(kw) {
					return Token{Type: KEYWORD, Value: kw, Position: scan.position}
				} else {
					return Token{Type: INVALID, Value: kw}
				}
			}

			// If the token is not caught, it is an invalid token
			return Token{Type: INVALID, Value: string(rune), Position: scan.position}
		}
	}
}

func (scan *Scanner) isValidKeyword(word string) bool {
	for _, kw := range validKeywords {
		if word == kw {
			return true
		}
	}
	return false
}

func (scan *Scanner) backup() {
	if err := scan.reader.UnreadRune(); err != nil {
		panic(err)
	}
	scan.position.Column--
}

func (scan *Scanner) integer() (string, error) {
	var number string

	for {
		rune, _, err := scan.reader.ReadRune()
		if err == io.EOF {
			return number, nil
		}

		scan.position.Column++

		if unicode.IsLetter(rune) || unicode.IsPunct(rune) {
			scan.backup()
			return number, errors.New("invalid character in an integer")
		}

		if unicode.IsDigit(rune) {
			number += string(rune)
			continue
		}

		scan.backup()
		return number, nil

	}
}

func (scan *Scanner) str() string {
	var str string
	var quote uint8

	for {
		rune, _, err := scan.reader.ReadRune()
		if err == io.EOF {
			return str
		}

		scan.position.Column++

		if rune == '"' {
			quote++
		}

		if quote < 2 {
			str += string(rune)
			continue
		}

		scan.backup()
		return str
	}
}

func (scan *Scanner) keyword() string {
	var str string

	for {
		rune, _, err := scan.reader.ReadRune()
		if err == io.EOF {
			return str
		}

		scan.position.Column++

		if unicode.IsLetter(rune) {
			str += string(rune)
			continue
		}

		scan.backup()
		return str
	}
}

func (scan *Scanner) resetPosition() {
	scan.position.Line++
	scan.position.Column = 0
}
