package object

import (
	"bytes"
	"eduC/ast"
	"fmt"
	"strconv"
	"strings"
)

type (
	ObjectType  string
	BuiltInFunc func(line int, args ...Object) Object
)

const (
	//IntOBJ > intenger object
	IntOBJ = "INTENGER"
	//FloatOBJ > FloatLiteral object
	FloatOBJ = "FlOAT"
	//BoolOBJ > boolean object
	BoolOBJ = "BOOLEAN"
	//StringOBJ > StringLiteral object
	StringOBJ = "STRING"
	//ErrorOBJ > error object
	ErrorOBJ = "ERROR"
	//ReturnOBJ > return value object
	ReturnOBJ = "RETURN_VALUE"
	//StopOBJ > stop object
	StopOBJ = "STOP"
	//FunctionOBJ > return function object
	FunctionOBJ = "FUNCTION"
	//BuiltInOBJ > built in function object
	BuiltInOBJ = "BUILTIN"
	//ArrayOBJ > built in function object
	ArrayOBJ = "ARRAY"
)

//Object interface (Type(),Inspect(),GetVal(),GetLine())
type Object interface {
	Type() ObjectType
	Inspect() string
	GetVal() interface{}
	GetLine() int
}

//Int object
type Int struct {
	Value int64
	Line  int
}

//Inspect Get Int value(string)
func (i *Int) Inspect() string { return fmt.Sprintf("%d", i.Value) }

//Type Get Int type(ObjectType)
func (i *Int) Type() ObjectType { return IntOBJ }

//GetVal Get Int value (interface)
func (i *Int) GetVal() interface{} { return i.Value }

//GetLine Get Int Line (int)
func (i *Int) GetLine() int { return i.Line }

//Float object
type Float struct {
	Value float64
	Line  int
}

//Inspect Get Float value(string)
func (f *Float) Inspect() string { return strconv.FormatFloat(f.Value, 'f', -1, 64) }

//Type Get Float type(ObjectType)
func (f *Float) Type() ObjectType { return FloatOBJ }

//GetVal Get Float value (interface)
func (f *Float) GetVal() interface{} { return f.Value }

//GetLine Get Float Line (int)
func (f *Float) GetLine() int { return f.Line }

//Bool object
type Bool struct {
	Value bool
	Line  int
}

//Inspect Get Bool value(string)
func (b *Bool) Inspect() string { return fmt.Sprintf("%t", b.Value) }

//Type Get Bool type(ObjectType)
func (b *Bool) Type() ObjectType { return BoolOBJ }

//GetVal Get Bool value (interface)
func (b *Bool) GetVal() interface{} { return b.Value }

//GetLine Get Bool Line (int)
func (b *Bool) GetLine() int { return b.Line }

//String object
type String struct {
	Value string
	Line  int
}

//Inspect Get String value(string)
func (s *String) Inspect() string { return fmt.Sprintf(`"%v"`, s.Value) }

//Type Get String type(ObjectType)
func (s *String) Type() ObjectType { return StringOBJ }

//GetVal Get String value (interface)
func (s *String) GetVal() interface{} { return s.Value }

//GetLine Get String Line (int)
func (s *String) GetLine() int { return s.Line }

//Array object
type Array struct {
	Elements []Object
	Line     int
}

//Type Get Array type(ObjectType)
func (ar *Array) Type() ObjectType { return ArrayOBJ }

//Inspect Get Array value(string)
func (ar *Array) Inspect() string {
	var out bytes.Buffer
	elements := []string{}
	for _, e := range ar.Elements {
		elements = append(elements, e.Inspect())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

//GetVal Get Array value(interface)
func (ar *Array) GetVal() interface{} { return ar.Elements }

//GetLine Get Array Line(int)
func (ar *Array) GetLine() int { return ar.Line }

//ERROR object
type ERROR struct {
	Value string
	Line  int
}

//Inspect Get ERROR value(string)
func (er *ERROR) Inspect() string { return fmt.Sprintf("%s", er.Value) }

//Type Get ERROR type(ObjectType)
func (er *ERROR) Type() ObjectType { return ErrorOBJ }

//GetVal Get ERROR value (interface)
func (er *ERROR) GetVal() interface{} { return er.Value }

//GetLine Get ERROR Line (int)
func (er *ERROR) GetLine() int { return er.Line }

//ReturnValue object
type ReturnValue struct {
	Value Object
	Line  int
}

//Type Get ReturnValue type(ObjectType)
func (rv *ReturnValue) Type() ObjectType { return ReturnOBJ }

//Inspect Get ReturnValue value(string)
func (rv *ReturnValue) Inspect() string { return rv.Value.Inspect() }

//GetVal Get ReturnValue value(interface)
func (rv *ReturnValue) GetVal() interface{} { return rv.Value.GetVal() }

//GetLine Get ReturnValue Line(int)
func (rv *ReturnValue) GetLine() int { return rv.Line }

//Stop object
type Stop struct {
	Line int
}

//Type Get Stop type(ObjectType)
func (sp *Stop) Type() ObjectType { return StopOBJ }

//Inspect Get Stop value(string)
func (sp *Stop) Inspect() string { return "stop" }

//GetVal Get Stop value(interface)
func (sp *Stop) GetVal() interface{} { return "stop" }

//GetLine Get Stop Line(int)
func (sp *Stop) GetLine() int { return sp.Line }

//Function object
type Function struct {
	Name    *ast.Identifier
	Params  []*ast.Identifier
	Process *ast.BlockStmt
	Env     *Env
	Line    int
}

//Type Get Function type(ObjectType)
func (fn *Function) Type() ObjectType { return FunctionOBJ }

//Inspect Get Function value(string)
func (fn *Function) Inspect() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range fn.Params {
		params = append(params, p.String())
	}
	name := ""
	if fn.Name != nil {
		name = fn.Name.String()
	}
	out.WriteString("\nfunc")
	out.WriteString(" " + name + " ")
	out.WriteString("( ")
	out.WriteString(strings.Join(params, ","))
	out.WriteString(") {\n" + fn.Process.String() + "\n}")

	return out.String()
}

//GetVal Get Function value(interface)
func (fn *Function) GetVal() interface{} { return fn.Name.String() }

//GetLine Get Function Line(int)
func (fn *Function) GetLine() int { return fn.Line }

//BuiltIn object
type BuiltIn struct {
	Func BuiltInFunc
}

//Type Get BuiltIn type(ObjectType)
func (b *BuiltIn) Type() ObjectType { return BuiltInOBJ }

//Inspect Get BuiltIn value(string)
func (b *BuiltIn) Inspect() string {
	return "This is BuiltInFunction"
}

//GetVal Get BuiltIn value(interface)
func (b *BuiltIn) GetVal() interface{} { return "This is BuiltInFunction" }

//GetLine Get BuiltIn Line(int)
func (b *BuiltIn) GetLine() int { return 0 }
