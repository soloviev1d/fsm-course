package lexer

import (
	"bufio"
	"io"
	"unicode"
)

type Token int

const (
	EOF = iota
	TYPE
	OPERATOR
	PUNCTUATION
	KEYWORD
	IDENTIFIER
	COMMENT
	LITERAL
	NUM
)

var tokens = []string{
	EOF:         "EOF",
	TYPE:        "DEFAULT TYPE",
	OPERATOR:    "OPERATOR",
	PUNCTUATION: "PUNCTUATION",
	KEYWORD:     "KEYWORD",
	IDENTIFIER:  "IDENTIFIER",
	COMMENT:     "COMMENT",
	LITERAL:     "LITERAL",
	NUM:         "NUM",
}

func (t Token) String() string {
	return tokens[t]
}

type Position struct {
	Line int
	Col  int
}

type Lexer struct {
	pos *Position
	r   *bufio.Reader
}

func NewLexer(reader io.Reader) *Lexer {
	i := 2.5
	i++
	return &Lexer{
		pos: &Position{Line: 1, Col: 1},
		r:   bufio.NewReader(reader),
	}
}

func (l *Lexer) Lex() (*Position, Token, string) {
	for {
		r, _, err := l.r.ReadRune()
		if err != nil {
			if err == io.EOF {
				return l.pos, EOF, ""
			}
		}

		l.pos.Col++
		// reset if \n
		if r == '\n' {
			l.pos.Col = 0
			l.pos.Line++
			continue
		}

		// ignore comments
		if r == '/' {
			r, _, err := l.r.ReadRune()
			if err != nil {
				if err == io.EOF {
					return l.pos, EOF, ""
				}
			}
			l.pos.Col++
			switch r {
			case '*':
				for {
					r, _, err := l.r.ReadRune()

					if err != nil {
						if err == io.EOF {
							return l.pos, EOF, ""
						}
					}

					l.pos.Col++

					if r != '*' {
						continue
					}
					if r == '\n' {
						l.pos.Col = 0
						l.pos.Line++
					}

					r, _, err = l.r.ReadRune()
					if err != nil {
						if err == io.EOF {
							return l.pos, EOF, ""
						}
					}
					l.pos.Col++
					if r == '/' {
						break
					}
				}
				break
			case '/':
				for {
					r, _, err := l.r.ReadRune()
					if err != nil {
						if err == io.EOF {
							return l.pos, EOF, ""
						}
					}
					l.pos.Col++
					if r == '\n' {
						break

					}
				}

			}
		}
		if unicode.IsSpace(r) {
			continue
		} else if isIdentifier(r) {
			sp := l.pos
			l.Unread()
			lit := l.Identifier()
			if _, ok := KeywordMap[lit]; ok {
				return sp, KEYWORD, lit
			}
			if _, ok := DefaultTypeMap[lit]; ok {
				return sp, TYPE, lit
			}
			return sp, IDENTIFIER, lit
		} else if isOperator(r) {
			sp := l.pos

			l.Unread()
			lit := l.Operator()
			return sp, OPERATOR, lit
		} else if isPunct(r) {

			sp := l.pos

			l.Unread()
			lit := l.Punctuation()
			return sp, PUNCTUATION, lit

		} else if unicode.IsDigit(r) {
			sp := l.pos
			l.Unread()
			lit := l.Number()
			return sp, NUM, lit
		} else if r == '"' || r == '\'' || r == '`' {
			sp := l.pos
			l.Unread()
			lit := l.Literal()
			return sp, LITERAL, lit
		}
		// panic(err)

	}
}

func (l *Lexer) Unread() {
	if err := l.r.UnreadRune(); err != nil {
		panic(err)
	}

	l.pos.Col--
}
func (l *Lexer) Identifier() string {
	var ident string
	for {
		r, _, err := l.r.ReadRune()
		if err != nil {
			if err == io.EOF {
				return ident
			}
		}

		l.pos.Col++
		if unicode.IsLetter(r) || r == '_' {
			ident += string(r)
		} else {
			l.Unread()
			return ident
		}
	}
}

func (l *Lexer) Literal() string {
	var lit string

	for {
		r, _, err := l.r.ReadRune()
		if err != nil {
			if err == io.EOF {
				return lit
			}
		}
		if r == '\\' {
			lit += string(r)
			r, _, err = l.r.ReadRune()
			if err != nil {
				if err == io.EOF {
					return lit
				}
			}
			l.pos.Col++
			lit += string(r)
			continue
		}

		l.pos.Col++
		lit += string(r)
		if rune(lit[0]) == r && len(lit) > 1 {
			return lit
		}

	}
}
func (l *Lexer) Number() string {
	var num string
	for {
		r, _, err := l.r.ReadRune()
		if err != nil {
			if err == io.EOF {
				return num
			}
		}

		l.pos.Col++
		if unicode.IsDigit(r) || r == '.' {
			num += string(r)
		} else {
			l.Unread()
			return num
		}
	}
}

func (l *Lexer) Punctuation() string {
	var p string
	for {
		r, _, err := l.r.ReadRune()
		if err != nil {
			if err == io.EOF {
				return p
			}
		}

		l.pos.Col++
		if r == '\n' {
			l.Unread()
			return p
		}
		if _, ok := PunctMap[r]; ok {
			p += string(r)
		} else {
			l.Unread()
			return p
		}
	}

}

func (l *Lexer) Operator() string {
	var op string
	for {
		r, _, err := l.r.ReadRune()
		if err != nil {
			if err == io.EOF {
				return op
			}
		}

		l.pos.Col++
		if r == '\n' {
			l.Unread()
			return op
		}
		if _, ok := OpMap[r]; ok {
			op += string(r)
		} else {
			l.Unread()
			return op
		}
	}

}
