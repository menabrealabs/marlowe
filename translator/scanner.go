package translator

import (
	"bufio"
	"io"
	"unicode"
)

type Token uint8

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
	QUOTE    // "
	COMMA    // ,
)

var tokens = []string{
	EOF:      "EOF",
	INVALID:  "INVALID",
	KEYWORD:  "KEYWORD",
	STRING:   "STRING",
	INT:      "INT",
	PARENS_L: "(",
	PARENS_R: ")",
	SQUARE_L: "[",
	SQUARE_R: "]",
	QUOTE:    string('"'),
	COMMA:    ",",
}

func (t Token) String() string {
	return tokens[t]
}

type Position struct {
	line   int
	column int
}

type Scanner struct {
	position Position
	reader   *bufio.Reader
}

func NewScanner(reader io.Reader) *Scanner {
	return &Scanner{
		position: Position{line: 1, column: 0},
		reader:   bufio.NewReader(reader),
	}
}

func (scan *Scanner) Scan() (Position, Token, string) {
	for {
		rune, _, err := scan.reader.ReadRune()

		// Return EOF when we get an io.EOF from the reader
		if err == io.EOF {
			return scan.position, EOF, ""
		}

		// Panic on any other unhandled error
		if err != nil {
			panic(err)
		}

		scan.position.column++

		switch rune {
		case '\n':
			scan.resetPosition()
		case '(':
			return scan.position, PARENS_L, "("
		case ')':
			return scan.position, PARENS_R, ")"
		case '[':
			return scan.position, SQUARE_L, "["
		case ']':
			return scan.position, SQUARE_R, "]"
		case ',':
			return scan.position, COMMA, ","
		case '"':
			return scan.position, STRING, string(rune) // TO-DO: Strings need to be analysed fully
		default:
			// Ignore spaces
			if unicode.IsSpace(rune) {
				continue
			}

			// Tokenize INT
			if unicode.IsDigit(rune) {
				scan.backup()
				scan.integer()
				continue
			}

			// Tokenize keywords
			if unicode.IsTitle(rune) {
				scan.backup()
				scan.keyword()
				continue
			}
		}
	}
}

func (scan *Scanner) backup() {
	if err := scan.reader.UnreadRune(); err != nil {
		panic(err)
	}
	scan.position.column--
}

func (scan *Scanner) integer() string {
	var number string

	for {
		rune, _, err := scan.reader.ReadRune()
		if err == io.EOF {
			return number
		}

		scan.position.column++

		if unicode.IsDigit(rune) {
			number += string(rune)
		} else {
			scan.backup()
			return number
		}
	}
}

func (scan *Scanner) str() string {
	var str string

	for {
		rune, _, err := scan.reader.ReadRune()
		if err == io.EOF {
			return str
		}

		scan.position.column++

		if unicode.IsLetter(rune) {
			str += string(rune)
		} else {
			scan.backup()
			return str
		}
	}
}

func (scan *Scanner) keyword() string {
	return "fix-me"
}

func (scan *Scanner) resetPosition() {
	scan.position.line++
	scan.position.column = 0
}
