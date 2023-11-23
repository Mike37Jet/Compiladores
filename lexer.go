package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode"
)

type Token int

var contador int
var b int

const (
	EOF = iota
	ILLEGAL
	IDENT
	INT
	SEMI   // ;
	STRING // "string"

	// Infix ops
	ADD // +
	SUB // -
	MUL // *
	DIV // /

	ASSIGN  // =
	LPAREN  // (
	RPAREN  // )
	COMILLA // "
)

var tokens = map[Token]string{
	EOF:     "EOF",
	ILLEGAL: "ILLEGAL",
	IDENT:   "IDENT",
	INT:     "INT",
	SEMI:    ";",
	STRING:  "STRING",

	ADD:     "+",
	SUB:     "-",
	MUL:     "*",
	DIV:     "/",
	ASSIGN:  "=",
	LPAREN:  "(",
	RPAREN:  ")",
	COMILLA: "COMILLA",
}

func (t Token) String() string {
	if str, ok := tokens[t]; ok {
		return str
	}
	return "DESCONOCIDO"
}

type Position struct {
	line   int
	column int
}

type Lexer struct {
	pos    Position
	reader *bufio.Reader
}

func NewLexer(reader io.Reader) *Lexer {
	return &Lexer{
		pos:    Position{line: 1, column: 0},
		reader: bufio.NewReader(reader),
	}
}

func (l *Lexer) Lex() (Position, Token, string) {
	for {
		r, _, err := l.reader.ReadRune()

		if err != nil {
			if err == io.EOF {
				return l.pos, EOF, ""
			}

			panic(err)
		}

		l.pos.column++

		switch r {
		case '\n':
			l.resetPosition()
		case ';':
			return l.pos, SEMI, ";"
		case '+':
			return l.pos, ADD, "+"
		case '-':
			return l.pos, SUB, "-"
		case '*':
			return l.pos, MUL, "*"
		case '/':
			return l.pos, DIV, "/"
		case '=':
			return l.pos, ASSIGN, "="
		case '(':
			return l.pos, LPAREN, "("
		case ')':
			return l.pos, RPAREN, ")"
		case '"':
			if contador == 1 {
				startPos := l.pos
				lit := l.lexString()
				contador = 0
				return startPos, STRING, lit
			}

			if b == 0 {
				contador++
				b++
				l.reader.UnreadRune()
				return l.pos, COMILLA, "\""
			}
			b--
			return l.pos, COMILLA, "\""

		default:
			if unicode.IsSpace(r) {
				continue
			} else if unicode.IsDigit(r) {
				startPos := l.pos
				l.backup()
				lit := l.lexInt()
				return startPos, INT, lit
			} else if unicode.IsLetter(r) {
				startPos := l.pos
				l.backup()
				lit := l.lexIdent()
				return startPos, IDENT, lit
			} else {
				return l.pos, ILLEGAL, string(r)
			}

		}
	}
}

func (l *Lexer) resetPosition() {
	l.pos.line++
	l.pos.column = 0
}

func (l *Lexer) backup() {
	err := l.reader.UnreadRune()
	if err != nil {
		return
	}
	l.pos.column--
}

func (l *Lexer) lexInt() string {
	var lit string
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return lit
			}
		}

		l.pos.column++
		if unicode.IsDigit(r) {
			lit = lit + string(r)
		} else {
			l.backup()
			return lit
		}
	}
}

func (l *Lexer) lexIdent() string {

	var lit string
	r, _, err := l.reader.ReadRune()
	if err != nil {
		if err == io.EOF {
			return lit
		}
		panic(err)
	}

	l.pos.column++
	if unicode.IsLetter(r) {
		lit = string(r)
	}
	return lit

}
func (l *Lexer) lexString() string {

	var lit string
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return lit
			}
			panic(err)
		}

		l.pos.column++

		if r == '"' && len(lit) >= 1 {
			l.backup()
			return lit
		} else {
			if r != '"' {
				lit += string(r)
			}
		}
	}
	return ""
}

func main() {
	file, err := os.Open("input.test")
	if err != nil {
		panic(err)
	}

	lexer := NewLexer(file)
	for {
		pos, tok, lit := lexer.Lex()
		if tok == EOF {
			break
		}
		fmt.Printf("%d:%d\t%s\t%s\n", pos.line, pos.column, tok, lit)
	}
}
