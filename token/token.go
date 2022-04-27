package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
	Line    int
}

const (
	//ILLEGAL is unknown
	ILLEGAL = "ILLEGAL"
	//EOF is end of file
	EOF = "EOF"

	//IDENT is identifier
	IDENT = "IDENT"
	//INT is literal
	INT    = "INT"
	FLOAT  = "FLOAT"
	TRUE   = "TRUE"
	FALSE  = "FALSE"
	STRING = "STRING"

	//Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"
	PERCENT  = "%"
	EXCLA    = "!"

	//Delimiters
	COMMA     = ","
	SEMICOLON = ";"

	//PARENs
	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	//Comparison operator
	Equal  = "=="
	Nequal = "!="
	LT     = "<"
	RT     = ">"
	ERT    = ">="
	ELT    = "<="
	AND    = "AND"
	OR     = "OR"
	//Keywords
	FUNCTION = "FUNCTION"
	MAKE     = "MAKE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	STOP     = "STOP"
	LOOP     = "LOOP"
)

//Keywords keywords(fn,let etc...)
var Keywords = map[string]TokenType{
	"func":   FUNCTION,
	"make":   MAKE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"stop":   STOP,
	"true":   TRUE,
	"false":  FALSE,
	"and":    AND,
	"or":     OR,
	"loop":   LOOP,
}

//IsKeywords If keywords,return the tokentype,else,return IDENT
func IsKeywords(word string) TokenType {
	if tok, ok := Keywords[word]; ok {
		return tok
	}
	return IDENT
}
