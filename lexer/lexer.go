package lexer

import (
	"github.com/hmwri/peridot/token"
)

//Lexer struct
type Lexer struct {
	//input Target source
	input string
	//runeinput Input change to rune(Unicode support)
	runeinput []rune
	//position Current position
	position int
	//readPosition Reading position
	readPosition int
	//line
	line int
	//ch Testing character
	ch rune
}

//New Make Lexer struct and call readChar(Lexer format)
func New(input string) *Lexer {
	input = DeleteComment(input)
	l := &Lexer{input: input, runeinput: []rune(input), line: 1}
	l.runeinput = append(l.runeinput, '\n')
	l.readChar()
	return l
}

//readChar Read next character
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.runeinput) {
		l.ch = 0 //EOF
	} else {
		l.ch = l.runeinput[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

//NextToken Check next Lexer.ch and make token
func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skip()
	switch l.ch {
	case '=':
		if next := l.nextRead(); next == '=' {
			tok = token.Token{Type: token.Equal, Literal: "==", Line: l.line}
			l.readChar()
		} else {
			tok = newToken(token.ASSIGN, l.ch, l.line)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch, l.line)
	case '(':
		tok = newToken(token.LPAREN, l.ch, l.line)
	case ')':
		tok = newToken(token.RPAREN, l.ch, l.line)
		l.autoSemicolon(tok.Type)
	case ',':
		tok = newToken(token.COMMA, l.ch, l.line)
	case '{':
		tok = newToken(token.LBRACE, l.ch, l.line)
	case '}':
		tok = newToken(token.RBRACE, l.ch, l.line)
		l.autoSemicolon(tok.Type)
	case '[':
		tok = newToken(token.LBRACKET, l.ch, l.line)
	case ']':
		tok = newToken(token.RBRACKET, l.ch, l.line)
		l.autoSemicolon(tok.Type)
	case '+':
		tok = newToken(token.PLUS, l.ch, l.line)
	case '-':
		tok = newToken(token.MINUS, l.ch, l.line)
	case '*':
		tok = newToken(token.ASTERISK, l.ch, l.line)
	case '/':

		tok = newToken(token.SLASH, l.ch, l.line)
	case '%':
		tok = newToken(token.PERCENT, l.ch, l.line)
	case '!':
		if next := l.nextRead(); next == '=' {
			tok = token.Token{Type: token.Nequal, Literal: "!=", Line: l.line}
			l.readChar()
		} else {
			tok = newToken(token.EXCLA, l.ch, l.line)
		}

	case '<':
		if next := l.nextRead(); next == '=' {
			tok = token.Token{Type: token.ELT, Literal: "<=", Line: l.line}
			l.readChar()
		} else {
			tok = newToken(token.LT, l.ch, l.line)
		}

	case '>':
		if next := l.nextRead(); next == '=' {
			tok = token.Token{Type: token.ERT, Literal: ">=", Line: l.line}
			l.readChar()
		} else {
			tok = newToken(token.RT, l.ch, l.line)
		}
	case '"':
		tok = l.readString()
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
		tok.Line = l.line

	default:
		if isLetter(l.ch) {
			tok = l.readLetter()
			l.autoSemicolon(tok.Type)
			return tok
		}
		if isNumber(l.ch) {
			tok = l.readNumber()
			l.autoSemicolon(tok.Type)
			return tok
		}
		tok = newToken(token.ILLEGAL, l.ch, l.line)
	}
	l.readChar()
	return tok
}

//readLetter Return Letters and Tokentype
func (l *Lexer) readLetter() token.Token {
	start := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	word := string(l.runeinput[start:l.position])
	tokenType := token.IsKeywords(word)
	return token.Token{Type: tokenType, Literal: word, Line: l.line}
}

//readNumber Return Numbers and Tokentype
func (l *Lexer) readNumber() token.Token {
	isFloat := false
	start := l.position
	for isNumber(l.ch) || (l.ch == '.' && !isFloat) {
		if l.ch == '.' {
			isFloat = true
		}
		l.readChar()
	}
	if isFloat {
		return token.Token{Type: token.FLOAT, Literal: string(l.runeinput[start:l.position]), Line: l.line}
	}
	return token.Token{Type: token.INT, Literal: string(l.runeinput[start:l.position]), Line: l.line}
}

//readString Return Letters(String) and Tokentype
func (l *Lexer) readString() token.Token {
	start := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	str := string(l.runeinput[start:l.position])
	return token.Token{Type: token.STRING, Literal: str, Line: l.line}
}

//nextRead read next character
func (l *Lexer) nextRead() rune {
	if l.readPosition >= len(l.runeinput) {
		return 0 //EOF
	}
	return l.runeinput[l.readPosition]
}

//beforeRead read before character
func (l *Lexer) beforeRead() rune {
	if l.readPosition >= len(l.runeinput) {
		return 0 //EOF
	}
	return l.runeinput[l.position-1]
}

//isLetter If a character is letter return true
func isLetter(ch rune) bool {
	if 'A' <= ch && !(ch == '[' || ch == ']') {
		return true
	}
	return false
}

//isNumber If a character is number return true
func isNumber(ch rune) bool {
	if '0' <= ch && ch <= '9' {
		return true
	}
	return false
}

//newToken make token
func newToken(tokenType token.TokenType, ch rune, line int) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch), Line: line}
}

func (l *Lexer) skip() {
	for l.ch == '\n' || l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
		if l.ch == '\n' {
			l.line++
		}
		l.readChar()
	}
}
func (l *Lexer) spaceskip() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}

}

//autoSemicolon put semicolon on end
func (l *Lexer) autoSemicolon(tok token.TokenType) {
	l.spaceskip()
	switch tok {
	case token.MAKE, token.RETURN, token.IF, token.ELSE, token.FUNCTION, token.LOOP:
	case token.RBRACE, token.RPAREN, token.RBRACKET:
		if l.nextRead() == '\n' {
			l.line++
			l.runeinput[l.readPosition] = ';'
		}
	default:
		if l.ch == '\n' {
			l.line++
			l.ch = ';'
		}
	}
}
