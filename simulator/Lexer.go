package simulator

import "fmt"

const (
	WHITESPACE = iota
	TEXT
	NUMBER
	SYMBOL
	KEYWORD
)

var Categories = map[int]string{
	0: "WS",
	1: "TXT",
	2: "NUM",
	3: "SYM",
	4: "KW",
}

type Token struct {
	Category int
	ID       string
}

type Lexer struct {
	Raw    []byte
	Tokens []Token
}

func NewLexer() Lexer {
	return Lexer{
		Tokens: make([]Token, 0),
	}
}

func (l *Lexer) Lex() []Token {
	for i := 0; i < len(l.Raw); {
		c := l.Raw[i]

		if isSymbol(c) {
			i = l.LexSymbol(i)
		} else if isWS(c) {
			i = l.LexWS(i)
		} else if isNumber(c) {
			i = l.LexNumber(i)
		} else if c == '#' {
			i = l.SkipComment(i)
		} else {
			i = l.LexWord(i)
		}
	}

	return l.Tokens
}

func isSymbol(c byte) bool {
	return (33 <= c && c <= 47 && c != 35) || (58 <= c && c <= 64) || ((91 <=
		c && c <= 96) && c != 95) || (123 <= c && c <= 126)
}

func isWS(c byte) bool {
	return c == 0x20 || c == 0x0A
}

func isNumber(c byte) bool {
	return c >= '0' && c <= '9'
}

func isKeyword(s string) bool {
	for _, k := range Keywords {
		if k == s {
			return true
		}
	}
	return false
}

func (l *Lexer) SkipComment(index int) int {

	newIndex := index

	for i := index; i < len(l.Raw); i++ {
		c := l.Raw[i]
		if c == '\n' {
			break
		}
		newIndex = i
	}

	return newIndex + 1
}

func (l *Lexer) LexSymbol(index int) int {
	var collected string
	newIndex := index

	for i := index; i < len(l.Raw); i++ {
		c := l.Raw[i]

		if isSymbol(c) {
			collected += string(c)
		} else {
			break
		}

		newIndex = i
	}

	token := Token{SYMBOL, collected}
	l.Tokens = append(l.Tokens, token)

	return newIndex + 1
}

func (l *Lexer) LexWS(index int) int {
	var collected string
	newIndex := index

	for i := index; i < len(l.Raw); i++ {
		c := l.Raw[i]

		if isWS(c) {
			collected += string(c)
		} else {
			break
		}

		newIndex = i
	}

	return newIndex + 1
}

func (l *Lexer) LexNumber(index int) int {
	var collected string
	newIndex := index

	for i := index; i < len(l.Raw); i++ {
		c := l.Raw[i]

		if c >= '0' && c <= '9' {
			collected += string(c)
		} else if c == 'E' || c == 'e' {
			i++
			if l.Raw[i] == '+' || l.Raw[i] == '-' || isNumber(l.Raw[i]) {
				collected += string(l.Raw[i-1])
				collected += string(l.Raw[i])
			} else {
				break
			}
		} else if c == '.' {
			if !isNumber(l.Raw[i+1]) {
				break
			} else {
				collected += string(l.Raw[i])
			}
		} else if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') {
			collected += string(l.Raw[i])
		} else {
			break
		}

		newIndex = i
	}

	token := Token{NUMBER, collected}
	l.Tokens = append(l.Tokens, token)
	return newIndex + 1
}

func (l *Lexer) LexWord(index int) int {
	var collected string
	newIndex := index

	for i := index; i < len(l.Raw); i++ {
		c := l.Raw[i]

		if !isSymbol(c) && !isWS(c) {
			collected += string(c)
		} else {
			break
		}

		newIndex = i
	}

	category := TEXT
	if isKeyword(collected) {
		category = KEYWORD
	}

	token := Token{category, collected}
	l.Tokens = append(l.Tokens, token)

	return newIndex + 1
}

func (l *Lexer) PrintTokens() {
	fmt.Println("Tokens: ")
	for _, tk := range l.Tokens {
		fmt.Println("\t", Categories[tk.Category], " ", tk.ID)
	}
}
