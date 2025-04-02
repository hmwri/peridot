package parser

import (
	"fmt"
	"github.com/hmwri/peridot/ast"
	"github.com/hmwri/peridot/errorwords"
	"github.com/hmwri/peridot/lexer"
	"github.com/hmwri/peridot/token"
	"strconv"
	"strings"
)

type (
	//Parser struct
	Parser struct {
		//lexer struct
		l *lexer.Lexer
		//Current token
		nowToken token.Token
		//Reading token
		readToken token.Token
		//erros
		errors []Err
		//fixparsefunctions map
		prefixParseFns map[token.TokenType]prefixParseFn
		infixParseFns  map[token.TokenType]infixParseFn
	}
	//Err error struct
	Err struct {
		Message string
		Line    int
	}
	//prefix parse function
	prefixParseFn func() ast.Expression
	//infix parse function
	infixParseFn func(ast.Expression) ast.Expression
)

//level
const (
	_        int = iota //auto increment
	LOWEST              //LOWEST lEVEL
	OR                  //or
	AND                 //and
	EQUAL               //==
	INEQUAL             //<,>
	ADDSUB              //+,-
	MULTIDIV            //*,/
	PREFIX              //-x !x
	FUNC                //function(x)
	INDEX               //array[index]
)

//priority
var priority = map[token.TokenType]int{
	token.Equal:    EQUAL,
	token.Nequal:   EQUAL,
	token.AND:      AND,
	token.OR:       OR,
	token.RT:       INEQUAL,
	token.LT:       INEQUAL,
	token.ERT:      INEQUAL,
	token.ELT:      INEQUAL,
	token.PLUS:     ADDSUB,
	token.MINUS:    ADDSUB,
	token.ASTERISK: MULTIDIV,
	token.SLASH:    MULTIDIV,
	token.PERCENT:  MULTIDIV,
	token.LPAREN:   FUNC,
	token.LBRACKET: INDEX,
}

//New Make Parser struct and call nextToken(Parser format)
func New(l *lexer.Lexer) *Parser {
	errorwords.Jerror()
	p := &Parser{l: l, errors: []Err{}}
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.addPrefix(token.IF, p.parseIf)
	p.addPrefix(token.FUNCTION, p.parseFunction)
	//values
	p.addPrefix(token.IDENT, p.parseIdent)
	p.addPrefix(token.INT, p.parseInt)
	p.addPrefix(token.FLOAT, p.parseFloat)
	p.addPrefix(token.TRUE, p.parseBool)
	p.addPrefix(token.FALSE, p.parseBool)
	p.addPrefix(token.STRING, p.parseString)
	//
	p.addPrefix(token.LPAREN, p.parsegroup)
	p.addPrefix(token.LBRACKET, p.parseArray)
	//prefix
	p.addPrefix(token.MINUS, p.parsePrefix)
	p.addPrefix(token.EXCLA, p.parsePrefix)
	//infix
	p.addInfix(token.Equal, p.parseInfix)
	p.addInfix(token.Nequal, p.parseInfix)
	p.addInfix(token.AND, p.parseInfix)
	p.addInfix(token.OR, p.parseInfix)
	p.addInfix(token.RT, p.parseInfix)
	p.addInfix(token.LT, p.parseInfix)
	p.addInfix(token.ERT, p.parseInfix)
	p.addInfix(token.ELT, p.parseInfix)
	p.addInfix(token.PLUS, p.parseInfix)
	p.addInfix(token.MINUS, p.parseInfix)
	p.addInfix(token.ASTERISK, p.parseInfix)
	p.addInfix(token.SLASH, p.parseInfix)
	p.addInfix(token.PERCENT, p.parseInfix)
	//
	p.addInfix(token.LPAREN, p.parseCall)
	p.addInfix(token.LBRACKET, p.parseIndex)
	//
	p.nextToken()
	p.nextToken()
	return p
}

//nextToken Get next token and input to Parser
func (p *Parser) nextToken() {
	p.nowToken = p.readToken
	p.readToken = p.l.NextToken()
}

//Parse Main parse program
func (p *Parser) Parse() *ast.Root {
	root := &ast.Root{Statements: []ast.Statement{}}

	for p.nowToken.Type != token.EOF {

		if stmt := p.parseStmt(); stmt != nil {
			root.Statements = append(root.Statements, stmt)
		}
		p.nextToken()
	}
	return root
}

//parseStmt Do parsing (choose correct parsefunc)
func (p *Parser) parseStmt() ast.Statement {
	switch p.nowToken.Type {
	case token.MAKE:
		return p.parseMakeStmt()
	case token.RETURN:
		return p.parseReturnStmt()
	case token.LOOP:
		return p.parseLoop()
	case token.STOP:
		return p.parseStopStmt()
	default:
		return p.parseExprStmt()
	}
}

//parseMakeStmt return parsed statement and check statement
func (p *Parser) parseMakeStmt() *ast.Make {
	makestmt := &ast.Make{Token: p.nowToken}
	//expect tokenType is IDENT

	if !p.expect(token.IDENT) {
		return nil
	}

	makestmt.Name = &ast.Identifier{Token: p.nowToken, Value: p.nowToken.Literal}

	if !p.expect(token.ASSIGN) {

		return nil
	}
	p.nextToken()
	makestmt.Value = p.parseExpression(LOWEST)
	if p.readToken.Type == token.SEMICOLON {
		p.nextToken()
	}
	return makestmt
}

//parseRetutnStmt return  parsed statement and check statement
func (p *Parser) parseReturnStmt() *ast.Return {
	returnstmt := &ast.Return{Token: p.nowToken}
	p.nextToken()
	returnstmt.Value = p.parseExpression(LOWEST)
	if p.readToken.Type == token.SEMICOLON {
		p.nextToken()
	}
	return returnstmt
}

//parseStopStmt return  parsed statement and check statement
func (p *Parser) parseStopStmt() *ast.Stop {
	stopstmt := &ast.Stop{Token: p.nowToken}
	if p.readToken.Type == token.SEMICOLON {
		p.nextToken()
	}
	return stopstmt
}

//parse expressionstatement
func (p *Parser) parseExprStmt() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.nowToken}
	stmt.Expression = p.parseExpression(LOWEST)
	if p.readToken.Type == token.SEMICOLON {
		p.nextToken()
	}
	return stmt
}

//parse expression
func (p *Parser) parseExpression(level int) ast.Expression {
	// search prefix parse functions
	prefn := p.prefixParseFns[p.nowToken.Type]
	if prefn == nil {
		if p.nowToken.Type == token.ASSIGN {
			p.setError(122, p.nowToken.Line, p.nowToken.Literal)
			return nil
		}
		p.setError(121, p.nowToken.Line, p.nowToken.Literal)
		return nil
	}
	//call prefix parse function and put it into left
	left := prefn()
	for level < p.nextTokenPriority() {
		infn := p.infixParseFns[p.readToken.Type]
		if infn == nil {
			return left
		}
		p.nextToken()
		left = infn(left)
	}
	return left
}

//parseIfExpression return parsed expression and check exoression
func (p *Parser) parseIf() ast.Expression {
	ifexp := &ast.If{Token: p.nowToken}
	p.nextToken()
	ifexp.Condition = p.parseExpression(LOWEST)
	if !p.expect(token.LBRACE) {
		return nil
	}
	ifexp.Consequence = p.parseBlockstmt()

	if p.nextTokenType(token.ELSE) {
		p.nextToken()
		if p.nextTokenType(token.IF) {
			ifexp.Alternative = p.parseBlockstmt()
			return ifexp
		}
		if !p.expect(token.LBRACE) {
			return nil
		}
		ifexp.Alternative = p.parseBlockstmt()
	}
	return ifexp
}

//parseLoop Statement
func (p *Parser) parseLoop() ast.Statement {
	lpexp := &ast.Loop{Token: p.nowToken}
	if p.nextTokenType(token.LBRACE) {
		truetoken := token.Token{Type: token.TRUE, Literal: "true", Line: p.nowToken.Line}
		lpexp.Condition = &ast.Bool{Token: truetoken, Value: true}
	} else {
		p.nextToken()
		lpexp.Condition = p.parseExpression(LOWEST)
	}
	if !p.expect(token.LBRACE) {
		return nil
	}
	lpexp.Process = p.parseBlockstmt()
	if p.readToken.Type == token.SEMICOLON {
		p.nextToken()
	}
	return lpexp
}

//parseFunction return parsed expression and check expression
func (p *Parser) parseFunction() ast.Expression {
	fnexp := &ast.Function{Token: p.nowToken}
	if p.nextTokenType(token.IDENT) {
		p.nextToken()
		fnexp.Name = &ast.Identifier{Token: p.nowToken, Value: p.nowToken.Literal}
	}
	if !p.expect(token.LPAREN) {
		return nil
	}
	fnexp.Parameters = p.parseFnParams()
	if !p.expect(token.LBRACE) {
		return nil
	}
	fnexp.Process = p.parseBlockstmt()

	return fnexp
}

//parse Function Parameters
func (p *Parser) parseFnParams() []*ast.Identifier {
	idents := []*ast.Identifier{}
	if p.nextTokenType(token.RPAREN) {
		p.nextToken()
		return idents
	}
	p.nextToken()
	ident := &ast.Identifier{Token: p.nowToken, Value: p.nowToken.Literal}
	idents = append(idents, ident)
	for p.nextTokenType(token.COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{Token: p.nowToken, Value: p.nowToken.Literal}
		idents = append(idents, ident)
	}
	if !p.expect(token.RPAREN) {
		return nil
	}
	return idents
}

//parseCall return parsed expression and check expression
func (p *Parser) parseCall(function ast.Expression) ast.Expression {
	callexp := &ast.Call{Token: p.nowToken, Function: function}
	callexp.Arguments = p.parseList(token.RPAREN)
	return callexp
}

//parseArray return parsed expression and check expression
func (p *Parser) parseArray() ast.Expression {
	ar := &ast.Array{Token: p.nowToken}
	ar.Elements = p.parseList(token.RBRACKET)
	return ar
}

//parse List(Argumetns,Elements) to []ast.Expression
func (p *Parser) parseList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}
	if p.nextTokenType(end) {
		p.nextToken()
		return list
	}
	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))
	for p.nextTokenType(token.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}
	if !p.expect(end) {
		return nil
	}
	return list
}

//parseIndex return parsed expression and check expression
func (p *Parser) parseIndex(left ast.Expression) ast.Expression {
	ix := &ast.Index{Token: p.nowToken, Left: left}
	p.nextToken()
	ix.Index = p.parseExpression(LOWEST)
	if !p.expect(token.RBRACKET) {
		return nil
	}
	return ix
}

//parse Block Statements
func (p *Parser) parseBlockstmt() *ast.BlockStmt {
	var bs *ast.BlockStmt
	bs = &ast.BlockStmt{Token: p.nowToken, Statements: []ast.Statement{}}
	p.nextToken()
	for !p.nowTokenType(token.RBRACE) && !p.nowTokenType(token.EOF) {
		stmt := p.parseStmt()
		if stmt != nil {
			bs.Statements = append(bs.Statements, stmt)
		}
		p.nextToken()
	}
	return bs
}

//parse Identifier
func (p *Parser) parseIdent() ast.Expression {
	if strings.Index(p.nowToken.Literal, "　") != -1 {
		p.setError(130, p.nowToken.Line)
	}
	if p.nextTokenType(token.ASSIGN) {
		tok := p.nowToken
		name := &ast.Identifier{Token: tok, Value: p.nowToken.Literal}
		p.nextToken()
		p.nextToken()
		value := p.parseExpression(LOWEST)
		if p.readToken.Type == token.SEMICOLON {
			p.nextToken()
		}
		return &ast.Assign{Token: tok, Name: name, Value: value}
	}
	return &ast.Identifier{Token: p.nowToken, Value: p.nowToken.Literal}
}

//parse Intenger
func (p *Parser) parseInt() ast.Expression {
	inted, err := strconv.ParseInt(p.nowToken.Literal, 10, 0)
	if err != nil {
		p.setError(111, p.nowToken.Line, p.nowToken.Literal)
		return nil
	}
	return &ast.Int{Token: p.nowToken, Value: inted}
}

//parse Intenger
func (p *Parser) parseFloat() ast.Expression {
	floated, err := strconv.ParseFloat(p.nowToken.Literal, 64)
	if err != nil {
		p.setError(112, p.nowToken.Line, p.nowToken.Literal)
		return nil
	}
	return &ast.Float{Token: p.nowToken, Value: floated}
}

//parse Boolean
func (p *Parser) parseBool() ast.Expression {
	return &ast.Bool{Token: p.nowToken, Value: p.nowTokenType(token.TRUE)}
}

//parse String Literal
func (p *Parser) parseString() ast.Expression {
	return &ast.String{Token: p.nowToken, Value: p.nowToken.Literal}
}

//parse Grouped Expression
func (p *Parser) parsegroup() ast.Expression {
	p.nextToken()
	ex := p.parseExpression(LOWEST)
	if !p.expect(token.RPAREN) {
		return nil
	}
	return ex
}

//parse Prefix Expression
func (p *Parser) parsePrefix() ast.Expression {
	ex := &ast.Prefix{Token: p.nowToken, Operator: p.nowToken.Literal}
	p.nextToken()
	ex.Value = p.parseExpression(PREFIX)
	return ex
}

//parse Infix Expression
func (p *Parser) parseInfix(left ast.Expression) ast.Expression {
	ex := &ast.Infix{
		Token:    p.nowToken,
		Operator: p.nowToken.Literal,
		Left:     left,
	}
	level := p.nowTokenPriority()
	p.nextToken()
	ex.Right = p.parseExpression(level)
	return ex
}

//expect next token and error check
func (p *Parser) expect(t token.TokenType) bool {
	if p.readToken.Type == t {
		p.nextToken()
		return true
	}
	p.setError(101, p.readToken.Line, t, p.readToken.Literal)
	return false
}

//setError set error
func (p *Parser) setError(code int, params ...interface{}) {
	message := fmt.Sprintf(errorwords.Err[code], params...)
	line := params[0].(int)
	p.errors = append(p.errors, Err{Message: message, Line: line})
}

//getError get error
func (p *Parser) GetError() []Err {
	return p.errors
}

//add fixfunction to functionsmap
func (p *Parser) addPrefix(tt token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tt] = fn
}
func (p *Parser) addInfix(tt token.TokenType, fn infixParseFn) {
	p.infixParseFns[tt] = fn
}

func (p *Parser) nowTokenType(t token.TokenType) bool {
	if t == p.nowToken.Type {
		return true
	}
	return false

}
func (p *Parser) nextTokenType(t token.TokenType) bool {
	if t == p.readToken.Type {
		return true
	}
	return false

}

//get token priority
func (p *Parser) nowTokenPriority() int {
	if pr, ok := priority[p.nowToken.Type]; ok {
		return pr
	}
	return LOWEST
}
func (p *Parser) nextTokenPriority() int {
	if pr, ok := priority[p.readToken.Type]; ok {
		return pr
	}
	return LOWEST
}
