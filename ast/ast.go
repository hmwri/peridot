package ast

import (
	"bytes"
	"eduC/token"
	"strings"
)

//ノード
type Node interface {
	TokenLiteral() string
	String() string
}

//文
type Statement interface {
	Node
	statementNode()
}

//式
type Expression interface {
	Node
	expressionNode()
}

//Root Root node(Statements)
type Root struct {
	Statements []Statement
}

//select Make Buffer and select String() each Statement
func (r *Root) String() string {
	var out bytes.Buffer

	for _, s := range r.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

func (r *Root) TokenLiteral() string {
	if len(r.Statements) > 0 {
		return r.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

//Make Make node(Statements)
type Make struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (m *Make) statementNode() {}
func (m *Make) TokenLiteral() string {
	if m != nil {
		return m.Token.Literal
	}
	return "nil"
}

func (m *Make) String() string {
	var out bytes.Buffer
	if m != nil {
		out.WriteString(m.TokenLiteral() + " ")
		out.WriteString(m.Name.String())
		out.WriteString(" = ")
		if m.Value != nil {
			out.WriteString(m.Value.String())
		}
	}
	return out.String()
}

//Return node(Statements)
type Return struct {
	Token token.Token
	Value Expression //return value
}

func (r *Return) statementNode()       {}
func (r *Return) TokenLiteral() string { return r.Token.Literal }

func (r *Return) String() string {
	var out bytes.Buffer
	out.WriteString(r.TokenLiteral() + " ")
	if r.Value != nil {
		out.WriteString(r.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

//Stop node(Statements)
type Stop struct {
	Token token.Token
}

func (s *Stop) statementNode()       {}
func (s *Stop) TokenLiteral() string { return s.Token.Literal }

func (s *Stop) String() string {
	return "Stop"
}

//Identifier Identifier node(Expressions)
type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string {
	return i.Value
}

//ExpressionStatement node
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (e *ExpressionStatement) statementNode()       {}
func (e *ExpressionStatement) TokenLiteral() string { return e.Token.Literal }
func (e *ExpressionStatement) String() string {
	if e.Expression != nil {
		return e.Expression.String()
	}
	return ""
}

//Int node
type Int struct {
	Token token.Token
	Value int64
}

func (i *Int) expressionNode()      {}
func (i *Int) TokenLiteral() string { return i.Token.Literal }
func (i *Int) String() string {
	return i.Token.Literal
}

//Float node
type Float struct {
	Token token.Token
	Value float64
}

func (f *Float) expressionNode()      {}
func (f *Float) TokenLiteral() string { return f.Token.Literal }
func (f *Float) String() string {
	return f.Token.Literal
}

//Bool node
type Bool struct {
	Token token.Token
	Value bool
}

func (b *Bool) expressionNode()      {}
func (b *Bool) TokenLiteral() string { return b.Token.Literal }
func (b *Bool) String() string {
	return b.Token.Literal
}

//Prefix node
type Prefix struct {
	Token    token.Token
	Operator string
	Value    Expression
}

func (pf *Prefix) expressionNode()      {}
func (pf *Prefix) TokenLiteral() string { return pf.Token.Literal }
func (pf *Prefix) String() string {
	var out bytes.Buffer
	if pf != nil && pf.Value != nil {
		out.WriteString("(" + pf.Operator + pf.Value.String() + ")")
	}
	return out.String()
}

//Infix node
type Infix struct {
	Token    token.Token
	Operator string
	Left     Expression
	Right    Expression
}

func (If *Infix) expressionNode()      {}
func (If *Infix) TokenLiteral() string { return If.Token.Literal }
func (If *Infix) String() string {
	var out bytes.Buffer
	out.WriteString("(" + If.Left.String() + " " + If.Operator + " " + If.Right.String() + ")")
	return out.String()
}

//If Else expression Node
type If struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStmt
	Alternative *BlockStmt
}

func (If *If) expressionNode()      {}
func (If *If) TokenLiteral() string { return If.Token.Literal }
func (If *If) String() string {
	var out bytes.Buffer
	out.WriteString("if" + If.Condition.String() + " " + If.Consequence.String())
	if If.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(If.Alternative.String())
	}
	return out.String()
}

//Loop expression Node
type Loop struct {
	Token     token.Token
	Condition Expression
	Process   *BlockStmt
}

func (lp *Loop) statementNode()       {}
func (lp *Loop) TokenLiteral() string { return lp.Token.Literal }
func (lp *Loop) String() string {
	var out bytes.Buffer
	out.WriteString("loop" + lp.Condition.String() + " " + lp.Process.String())
	return out.String()
}

//Function expression Node
type Function struct {
	Token      token.Token
	Name       *Identifier
	Parameters []*Identifier
	Process    *BlockStmt
}

func (fn *Function) expressionNode()      {}
func (fn *Function) TokenLiteral() string { return fn.Token.Literal }
func (fn *Function) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range fn.Parameters {
		params = append(params, p.String())
	}
	name := ""
	if fn.Name != nil {
		name = fn.Name.String()
	}
	out.WriteString(fn.TokenLiteral())
	out.WriteString(" " + name + " ")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ","))
	out.WriteString(") " + fn.Process.String())

	return out.String()
}

//Call expression Node
type Call struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (cl *Call) expressionNode()      {}
func (cl *Call) TokenLiteral() string { return cl.Token.Literal }
func (cl *Call) String() string {
	var out bytes.Buffer
	args := []string{}

	for _, a := range cl.Arguments {
		if a != nil {
			args = append(args, a.String())

		}
	}
	_, ok := cl.Function.(*Function)
	_, ok2 := cl.Function.(*Identifier)
	if ok || ok2 {
		out.WriteString(cl.Function.String())
	}

	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

//Assign Assign node(expressions)
type Assign struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (a *Assign) expressionNode()      {}
func (a *Assign) TokenLiteral() string { return a.Token.Literal }

func (a *Assign) String() string {
	var out bytes.Buffer
	out.WriteString(a.Name.String())
	out.WriteString(" = ")
	if a.Value != nil {
		out.WriteString(a.Value.String())
	}
	return out.String()
}

//Array expression Node
type Array struct {
	Token    token.Token
	Elements []Expression
}

func (ar *Array) expressionNode()      {}
func (ar *Array) TokenLiteral() string { return ar.Token.Literal }
func (ar *Array) String() string {
	var out bytes.Buffer
	els := []string{}

	for _, a := range ar.Elements {
		if a != nil {
			els = append(els, a.String())
		}
	}
	out.WriteString("[")
	out.WriteString(strings.Join(els, ", "))
	out.WriteString("]")

	return out.String()
}

//Index expression Node
type Index struct {
	Token token.Token
	Left  Expression
	Index Expression
}

func (ix *Index) expressionNode()      {}
func (ix *Index) TokenLiteral() string { return ix.Token.Literal }
func (ix *Index) String() string {
	var out bytes.Buffer
	out.WriteString("(" + ix.Left.String())
	out.WriteString("[" + ix.Index.String() + "])")
	return out.String()
}

type BlockStmt struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStmt) statementNode()       {}
func (bs *BlockStmt) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStmt) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

//String String LIteral node
type String struct {
	Token token.Token
	Value string
}

func (st *String) expressionNode()      {}
func (st *String) TokenLiteral() string { return st.Token.Literal }
func (st *String) String() string {
	return st.Value
}
