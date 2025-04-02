package eval

import (
	"fmt"
	"github.com/hmwri/peridot/ast"
	"github.com/hmwri/peridot/errorwords"
	"github.com/hmwri/peridot/log"
	"github.com/hmwri/peridot/object"
	"strconv"
	"unicode/utf8"

	"golang.org/x/exp/utf8string"
)

//Eval evaluator
func Eval(node ast.Node, env *object.Env) object.Object {

	switch node := node.(type) {
	case *ast.Root:
		return evalRoot(node.Statements, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.Int:
		return &object.Int{Value: node.Value, Line: node.Token.Line}
	case *ast.Float:
		return &object.Float{Value: node.Value, Line: node.Token.Line}
	case *ast.Bool:
		return &object.Bool{Value: node.Value, Line: node.Token.Line}
	case *ast.String:
		return &object.String{Value: node.Value, Line: node.Token.Line}
	case *ast.Array:
		els := evalExps(node.Elements, env)

		if len(els) == 1 && isError(els[0]) {
			return els[0]
		}
		return &object.Array{Elements: els}
	case *ast.Index:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndex(left, index, node.Token.Line)
	case *ast.Prefix:
		value := Eval(node.Value, env)
		if isError(value) {
			return value
		}
		return evalPrefix(node.Operator, value, node.Token.Line)
	case *ast.Infix:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfix(node.Operator, left, right, node.Token.Line)
	case *ast.BlockStmt:
		return evalStmt(node.Statements, env)
	case *ast.If:
		return evalIf(node, env, node.Token.Line)
	case *ast.Loop:
		return evalLoop(node, env, node.Token.Line)
	case *ast.Return:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val, Line: node.Token.Line}
	case *ast.Stop:
		log.SetLog(node.Token.Line, "stop", "stop", "ループ離脱")
		return &object.Stop{Line: node.Token.Line}
	case *ast.Make:
		if node == nil {
			return errorwords.SetError(200, 0)
		}
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		if v, ok := node.Value.(*ast.Function); ok {
			//If you try to assign an anonymous function
			if v.Name != nil {
				return errorwords.SetError(230, node.Token.Line, v.Name.String())
			}
		}
		if val == nil {
			return errorwords.SetError(211, node.Token.Line)
		}
		log.SetLog(node.Token.Line, node.Name.String(), val.Inspect(), "変数("+node.Name.String()+")を定義し("+val.Inspect()+")を代入")
		env.SetEnv(node.Name.Value, val)
	case *ast.Assign:
		val := Eval(node.Value, env)
		if v, ok := node.Value.(*ast.Function); ok {
			//If you try to assign an anonymous function
			if v.Name != nil {
				return errorwords.SetError(230, node.Token.Line, v.Name.String())
			}
		}
		if val == nil {
			return errorwords.SetError(211, node.Token.Line)
		}
		if _, found := env.GetEnv(node.Name.Value); found {
			env.SetEnv(node.Name.Value, val)
			log.SetLog(node.Token.Line, node.Name.String(), val.Inspect(), "変数("+node.Name.String()+")に"+val.Inspect()+"を代入")
		} else {
			return errorwords.SetError(210, node.Token.Line, node.Name.Value)
		}

	case *ast.Identifier:
		return evalIdent(node, node.Token.Line, env)
	case *ast.Function:
		params := node.Parameters
		process := node.Process
		name := node.Name
		obj := &object.Function{Params: params, Env: env, Process: process, Name: name, Line: node.Token.Line}
		if name != nil {
			log.SetLog(node.Token.Line, node.Name.String(), obj.Inspect(), "関数を定義(名前:"+node.Name.String()+")")
			env.SetEnv(node.Name.Value, obj)
		} else {
			return obj
		}
	case *ast.Call:
		function := Eval(node.Function, env)
		log.SetLog(node.Token.Line, node.Function.String(), "実行", "関数を実行(関数:"+node.Function.String()+")")
		if function.Type() == object.ErrorOBJ {
			return function
		}
		args := evalExps(node.Arguments, env)
		if len(args) == 1 && args[0].Type() == object.ErrorOBJ {
			return args[0]
		}
		return exeFunction(function, args, node.Token.Line)
	default:
		return errorwords.SetError(200, 0)

	}
	return nil
}

//Root Statement Evaluator
func evalRoot(stmts []ast.Statement, env *object.Env) object.Object {
	var result object.Object
	for _, stmt := range stmts {
		result = Eval(stmt, env)
		switch v := result.(type) {
		case *object.ReturnValue:
			return v.Value
		case *object.ERROR:
			return v
		}
	}
	return result
}

//Block Statement Evaluator
func evalStmt(stmts []ast.Statement, env *object.Env) object.Object {
	var result object.Object
	for _, stmt := range stmts {
		result = Eval(stmt, env)
		if result != nil && (result.Type() == object.ReturnOBJ || result.Type() == object.StopOBJ || result.Type() == object.ErrorOBJ) {
			return result
		}

	}
	return result
}

//Expressions Evaluator
func evalExps(exps []ast.Expression, env *object.Env) []object.Object {
	var result []object.Object
	for _, exp := range exps {
		evaled := Eval(exp, env)
		if isError(evaled) {
			err := []object.Object{evaled}
			return err
		}
		result = append(result, evaled)
	}
	return result
}

//Prefix Expression Evaluator
func evalPrefix(operator string, value object.Object, line int) object.Object {
	switch operator {
	case "!":
		return evalBang(value, line)
	case "-":
		return evalMinus(value, line)
	default:
		return errorwords.SetError(201, line, operator)
	}
}

//Bang Expression Evaluator
func evalBang(value object.Object, line int) object.Object {

	switch value.GetVal() {
	case true:
		log.SetLog(line, "true", "false", "")
		return makeBoolObj(false, line)
	case false:
		log.SetLog(line, "false", "true", "")
		return makeBoolObj(true, line)
	default:
		log.SetLog(line, "Not Boolean(真偽値以外)", "false", "")
		return makeBoolObj(false, line)
	}
}

//Minus Expression Evaluator
func evalMinus(value object.Object, line int) object.Object {
	if value.Type() != object.IntOBJ && value.Type() != object.FloatOBJ {
		return errorwords.SetError(202, line)
	}
	intval, ok := value.GetVal().(int64)
	if ok {
		return &object.Int{Value: -intval, Line: line}
	}
	floatval, ok := value.GetVal().(float64)
	if ok {
		return &object.Float{Value: -floatval, Line: line}
	}
	return errorwords.SetError(202, line)
}

//Infix Expression Evaluator
func evalInfix(operator string, left object.Object, right object.Object, line int) object.Object {
	switch v := left.GetVal().(type) {
	case int64:
		lval := v
		rval, ok := right.GetVal().(int64)
		if !ok {
			rstr, ok := right.GetVal().(string)
			if ok {
				//if right is string ,left value convert to string
				lstr := strconv.FormatInt(lval, 10)
				if operator != "+" {
					return errorwords.SetError(207, line)
				}
				log.SetLog(line, lstr+" + "+rstr, lstr+rstr, "文字列結合")
				return &object.String{Value: lstr + rstr, Line: line}
			}
			floatRval, ok := right.GetVal().(float64)
			if !ok {
				return errorwords.SetError(203, line, "整数", "数,文字列以外")
			}
			//If right is float64,left value convert to float64
			floatLval := float64(lval)
			return floatCalc(floatLval, operator, floatRval, line)
		}
		switch operator {
		case "+":
			log.SetLog(line, infixIntString(lval, " + ", rval), strconv.FormatInt(lval+rval, 10), "計算")
			return &object.Int{Value: lval + rval, Line: line}
		case "-":
			log.SetLog(line, infixIntString(lval, " - ", rval), strconv.FormatInt(lval-rval, 10), "計算")
			return &object.Int{Value: lval - rval, Line: line}
		case "*":
			log.SetLog(line, infixIntString(lval, " * ", rval), strconv.FormatInt(lval*rval, 10), "計算")
			return &object.Int{Value: lval * rval, Line: line}
		case "/":
			if lval%rval != 0 {
				result := float64(lval) / float64(rval)
				log.SetLog(line, infixIntString(lval, " / ", rval), fmt.Sprintf("%v", result), "計算")
				return &object.Float{Value: result, Line: line}
			}
			log.SetLog(line, infixIntString(lval, " / ", rval), strconv.FormatInt(lval/rval, 10), "計算")
			return &object.Int{Value: lval / rval, Line: line}
		case "%":
			log.SetLog(line, infixIntString(lval, " % ", rval), strconv.FormatInt(lval%rval, 10), "計算")
			return &object.Int{Value: lval % rval, Line: line}
		case "<":
			log.SetLog(line, infixIntString(lval, " < ", rval), booltoString(lval < rval), "評価(true or false)")
			return makeBoolObj(lval < rval, line)
		case ">":
			log.SetLog(line, infixIntString(lval, " > ", rval), booltoString(lval > rval), "評価(true or false)")
			return makeBoolObj(lval > rval, line)
		case "<=":
			log.SetLog(line, infixIntString(lval, " <= ", rval), booltoString(lval <= rval), "評価(true or false)")
			return makeBoolObj(lval <= rval, line)
		case ">=":
			log.SetLog(line, infixIntString(lval, " >= ", rval), booltoString(lval >= rval), "評価(true or false)")
			return makeBoolObj(lval >= rval, line)
		case "!=":
			log.SetLog(line, infixIntString(lval, " != ", rval), booltoString(lval != rval), "評価(true or false)")
			return makeBoolObj(lval != rval, line)
		case "==":
			log.SetLog(line, infixIntString(lval, " == ", rval), booltoString(lval == rval), "評価(true or false)")
			return makeBoolObj(lval == rval, line)

		default:
			return errorwords.SetError(204, line, operator)
		}
	case float64:
		lval := v
		rval, ok := right.GetVal().(float64)
		if !ok {
			rstr, ok := right.GetVal().(string)
			if ok {
				//if right is string ,left value convert to string
				lstr := strconv.FormatFloat(lval, 'f', -1, 64)
				if operator != "+" {
					return errorwords.SetError(207, line)
				}
				log.SetLog(line, lstr+" + "+rstr, lstr+rstr, "文字列結合")
				return &object.String{Value: lstr + rstr, Line: line}
			}
			rval, ok := right.GetVal().(int64)
			if !ok {
				return errorwords.SetError(203, line, "整数", "数,文字列以外")
			}
			//If right is int64,right value convert to float64
			floatRval := float64(rval)
			return floatCalc(lval, operator, floatRval, line)
		}

		return floatCalc(lval, operator, rval, line)

	case bool:
		lval := v
		rval, ok := right.GetVal().(bool)
		if !ok {
			return errorwords.SetError(203, line, "真偽値or式", "真偽値or式以外")
		}
		switch operator {
		case "!=":
			log.SetLog(line, booltoString(lval)+" != "+booltoString(rval), booltoString(lval != rval), "比較")
			return makeBoolObj(lval != rval, line)
		case "==":
			log.SetLog(line, booltoString(lval)+" == "+booltoString(rval), booltoString(lval == rval), "比較")
			return makeBoolObj(lval == rval, line)
		case "and":
			log.SetLog(line, booltoString(lval)+" and "+booltoString(rval), booltoString(lval && rval), "論理積")
			return makeBoolObj(lval && rval, line)
		case "or":
			log.SetLog(line, booltoString(lval)+" or "+booltoString(rval), booltoString(lval || rval), "論理和")
			return makeBoolObj(lval || rval, line)
		default:
			return errorwords.SetError(205, line, operator)
		}
	case string:
		lstr := v
		rstr, ok := right.GetVal().(string)
		if !ok {
			irval, ok := right.GetVal().(int64)
			if ok {
				rstr = strconv.FormatInt(irval, 10)
			}
			frval, ok2 := right.GetVal().(float64)
			if ok2 {
				rstr = strconv.FormatFloat(frval, 'f', -1, 64)
			}
			if !ok && !ok2 {
				return errorwords.SetError(203, line, "文字列", "文字列,数以外")
			}

		}
		if operator == "+" {
			log.SetLog(line, lstr+" + "+rstr, lstr+rstr, "文字列結合")
			return &object.String{Value: lstr + rstr, Line: line}
		}
		if operator == "==" {
			log.SetLog(line, lstr+" == "+rstr, booltoString(lstr == rstr), "文字列比較")
			return makeBoolObj(lstr == rstr, line)
		}
		if operator == "!=" {
			log.SetLog(line, lstr+" != "+rstr, booltoString(lstr != rstr), "文字列比較")
			return makeBoolObj(lstr != rstr, line)
		}
		return errorwords.SetError(207, line)

	default:
		return errorwords.SetError(206, line, operator)
	}

}

//return bool object
func makeBoolObj(b bool, line int) object.Object {
	return &object.Bool{Value: b, Line: line}
}

//evaluate "If Expression"
func evalIf(i *ast.If, env *object.Env, line int) object.Object {
	condition := Eval(i.Condition, env)
	if isError(condition) {
		return condition
	}
	if condition == nil {
		return errorwords.SetError(220, line, "if")
	}
	if isTrue(condition) {
		log.SetLog(line, i.Condition.String(), "true", "条件がtrueであったためif内を実行")
		return Eval(i.Consequence, env)
	} else if i.Alternative != nil {
		log.SetLog(line, i.Condition.String(), "false", "条件がfalseであったためelse内を実行")
		return Eval(i.Alternative, env)
	} else {
		log.SetLog(line, i.Condition.String(), "false", "条件がfalseであったためif内をスキップ")
		return nil
	}
}

//evaluate "Loop Expression"
func evalLoop(l *ast.Loop, env *object.Env, line int) object.Object {
	log.SetLog(line, "loop", "start", "ループ開始")
	condition := Eval(l.Condition, env)
	if isError(condition) {
		return condition
	}
	_, isStr := condition.(*object.String)
	if isStr {
		return errorwords.SetError(500, line, "文字列")
	}
	var obj object.Object
	fnum, ok := condition.(*object.Float)
	//loop(number-float){}
	if ok {
		for i := 0.0; i < float64(fnum.Value); i++ {
			log.SetLog(line, fmt.Sprintf("%v <= %v", i+1, fnum.Value), "true", fmt.Sprintf("条件がtrueであったためloop内を実行(%v回目)", i+1))
			obj = Eval(l.Process, env)
			if isStopType(obj) {
				break
			}
		}
		log.SetLog(line, "loop", "end", "ループ終了")
		return obj
	}
	num, ok := condition.(*object.Int)
	//loop(number-int){}
	if ok {
		for i := 0; i < int(num.Value); i++ {
			log.SetLog(line, fmt.Sprintf("%v <= %v", i+1, num.Value), "true", fmt.Sprintf("条件がtrueであったためloop内を実行(%v回目)", i+1))
			obj = Eval(l.Process, env)
			if isStopType(obj) {
				break
			}
		}
		log.SetLog(line, "loop", "end", "ループ終了")
		return obj
	}

	i := 0
	for ; isTrue(condition); condition = Eval(l.Condition, env) {
		if isError(condition) {
			return condition
		}
		if condition == nil {
			return errorwords.SetError(220, line, "loop")
		}
		i++
		log.SetLog(line, l.Condition.String(), "true", fmt.Sprintf("条件がtrueであったためloop内を実行(%v回目)", i))
		obj := Eval(l.Process, env)
		if isStopType(obj) {
			break
		}
	}
	log.SetLog(line, "loop", "end", "ループ終了")
	return obj
}

//evaluate identifier expression
func evalIdent(id *ast.Identifier, line int, env *object.Env) object.Object {
	if val, found := env.GetEnv(id.Value); found {
		log.SetLog(line, id.Value, val.Inspect(), "変数("+id.Value+")を参照")
		return val
	}
	if blt, found := builtIns[id.Value]; found {
		return blt
	}
	return errorwords.SetError(210, line, id.Value)
}

//execute function object and return ReturnValue
func exeFunction(fn object.Object, args []object.Object, line int) object.Object {
	switch funcObj := fn.(type) {
	case *object.Function:
		newEnv := addFuncEnv(funcObj, args, line)
		if newEnv == nil {
			return nil
		}
		evaled := Eval(funcObj.Process, newEnv)
		return getRV(evaled)
	case *object.BuiltIn:
		return funcObj.Func(line, args...)
	}
	return errorwords.SetError(231, line, fn.Type())
}

//Add new Environment(in function) and set function parameters ,args in this env
func addFuncEnv(fn *object.Function, args []object.Object, line int) *object.Env {
	nenv := object.AddEnv(fn.Env)
	if len(args) != len(fn.Params) {
		errorwords.SetError(232, line, len(args), len(fn.Params))
		return nil
	}
	for i, param := range fn.Params {
		log.SetLog(param.Token.Line, param.String(), args[i].Inspect(), "関数のパラメータに引数を代入")
		nenv.SetEnv(param.Value, args[i])
	}
	return nenv
}
func evalIndex(left, index object.Object, line int) object.Object {
	switch {
	case left.Type() == object.ArrayOBJ && index.Type() == object.IntOBJ:
		return evalArrayIndex(left, index, line)
	case left.Type() == object.StringOBJ && index.Type() == object.IntOBJ:
		return evalStringIndex(left, index, line)
	default:
		if index.Type() != object.IntOBJ {
			return errorwords.SetError(401, line)
		}
		return errorwords.SetError(400, line, left.Type())
	}
}
func evalArrayIndex(left, index object.Object, line int) object.Object {
	array := left.(*object.Array)
	ix := index.(*object.Int).Value
	end := int64(len(array.Elements) - 1)
	if ix < 0 {
		return errorwords.SetError(402, line, left.Type())
	}
	if ix > end {
		return errorwords.SetError(403, line, ix, end)
	}
	log.SetLog(line, fmt.Sprintf("%v[%v]", array.Inspect(), ix), array.Elements[ix].Inspect(), "配列から値をとりだす")
	return array.Elements[ix]
}
func evalStringIndex(left, index object.Object, line int) object.Object {
	str := left.(*object.String)
	ix := index.(*object.Int).Value
	end := int64(utf8.RuneCountInString(str.Value))
	if ix <= 0 {
		return errorwords.SetError(502, line, left.Type())
	}
	if ix > end {
		return errorwords.SetError(503, line, ix, end)
	}
	utfStr := utf8string.NewString(str.Value)
	result := utfStr.Slice(int(ix)-1, int(ix))
	log.SetLog(line, fmt.Sprintf("%v[%v]", str.Inspect(), ix), result, "文字列から文字をとりだす")
	return &object.String{Value: result, Line: line}
}

//get return value from return object
func getRV(obj object.Object) object.Object {
	returnObj, ok := obj.(*object.ReturnValue)
	if ok {
		log.SetLog(obj.GetLine(), returnObj.Value.Inspect(), "呼び出し元へ", "関数の実行結果を返す")
		return returnObj.Value
	}
	log.SetLog(obj.GetLine(), "", "", "関数の実行完了")

	return obj
}

//object is true?
func isTrue(c object.Object) bool {
	v := c.GetVal()
	switch v {
	case true:
		return true
	case false:
		return false
	default:
		//WARNING
		return false
	}
}

//(left(int64)) operator (right(int64)) >> (left(string)) operator (right(string))
func infixIntString(left int64, operator string, right int64) string {
	return strconv.FormatInt(left, 10) + operator + strconv.FormatInt(right, 10)
}

//boolean > "boolean"
func booltoString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

//check to if target is errorobject
func isError(t object.Object) bool {
	if t == nil {
		return false
	}
	return t.Type() == object.ErrorOBJ
}

//calculate float
func floatCalc(lval float64, operator string, rval float64, line int) object.Object {
	switch operator {
	case "+":
		log.SetLog(line, infixFloatString(lval, " + ", rval), strconv.FormatFloat(lval+rval, 'f', -1, 64), "計算")
		return &object.Float{Value: lval + rval, Line: line}
	case "-":
		log.SetLog(line, infixFloatString(lval, " - ", rval), strconv.FormatFloat(lval-rval, 'f', -1, 64), "計算")
		return &object.Float{Value: lval - rval, Line: line}
	case "*":
		log.SetLog(line, infixFloatString(lval, " * ", rval), strconv.FormatFloat(lval*rval, 'f', -1, 64), "計算")
		return &object.Float{Value: lval * rval, Line: line}
	case "/":
		log.SetLog(line, infixFloatString(lval, " / ", rval), strconv.FormatFloat(lval/rval, 'f', -1, 64), "計算")
		return &object.Float{Value: lval / rval, Line: line}
	case "%":
		return errorwords.SetError(208, line, "%")
	case "<":
		log.SetLog(line, infixFloatString(lval, " < ", rval), booltoString(lval < rval), "評価(true or false)")
		return makeBoolObj(lval < rval, line)
	case ">":
		log.SetLog(line, infixFloatString(lval, " > ", rval), booltoString(lval > rval), "評価(true or false)")
		return makeBoolObj(lval > rval, line)
	case "<=":
		log.SetLog(line, infixFloatString(lval, " <= ", rval), booltoString(lval <= rval), "評価(true or false)")
		return makeBoolObj(lval <= rval, line)
	case ">=":
		log.SetLog(line, infixFloatString(lval, " >= ", rval), booltoString(lval >= rval), "評価(true or false)")
		return makeBoolObj(lval >= rval, line)
	case "!=":
		log.SetLog(line, infixFloatString(lval, " != ", rval), booltoString(lval != rval), "評価(true or false)")
		return makeBoolObj(lval != rval, line)
	case "==":
		log.SetLog(line, infixFloatString(lval, " == ", rval), booltoString(lval == rval), "評価(true or false)")
		return makeBoolObj(lval == rval, line)
	default:
		return errorwords.SetError(204, line, operator)
	}
}

//(left(float64)) operator (right(float64)) >> (left(string)) operator (right(string))
func infixFloatString(left float64, operator string, right float64) string {
	return strconv.FormatFloat(left, 'f', -1, 64) + operator + strconv.FormatFloat(right, 'f', -1, 64)
}
func isStopType(obj object.Object) bool {
	if obj != nil && (obj.Type() == object.StopOBJ || obj.Type() == object.ReturnOBJ || obj.Type() == object.ErrorOBJ) {
		return true
	}
	return false
}
