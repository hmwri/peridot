package eval

import (
	"bufio"
	"fmt"
	"github.com/hmwri/peridot/errorwords"
	"github.com/hmwri/peridot/log"
	"github.com/hmwri/peridot/object"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"
	"unicode/utf8"

	"golang.org/x/exp/utf8string"
)

var builtIns = map[string]*object.BuiltIn{
	//SIZE return Array or String length
	"SIZE": &object.BuiltIn{
		Func: func(line int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return errorwords.SetError(300, line, "SIZE", 1)
			}
			switch arg := args[0].(type) {
			case *object.String:
				log.SetLog(line, "SIZE("+arg.Inspect()+")", strconv.Itoa(len(arg.Value)), "組み込み関数SIZEを実行")
				return &object.Int{Value: int64(utf8.RuneCountInString(arg.Value)), Line: line}
			case *object.Array:
				log.SetLog(line, "SIZE("+arg.Inspect()+")", strconv.Itoa(len(arg.Elements)), "組み込み関数SIZEを実行")
				return &object.Int{Value: int64(len(arg.Elements)), Line: line}
			default:
				return errorwords.SetError(301, line, "SIZE", "文字列,配列")
			}
		},
	},
	//ADD add object into array
	"ADD": &object.BuiltIn{
		Func: func(line int, args ...object.Object) object.Object {
			if len(args) != 2 {
				return errorwords.SetError(300, line, "ADD", 2)
			}
			switch arg := args[0].(type) {
			case *object.Array:
				if arg == args[1] {
					return errorwords.SetError(303, line)
				}
				before := arg.Inspect()
				arg.Elements = append(arg.Elements, args[1])
				log.SetLog(line, "ADD("+before+")", arg.Inspect(), "組み込み関数ADDを実行")
				return nil
			default:
				return errorwords.SetError(302, line)
			}
		},
	},
	//DELETE delete object from array
	"DELETE": &object.BuiltIn{
		Func: func(line int, args ...object.Object) object.Object {
			if len(args) != 2 {
				return errorwords.SetError(300, line, "DELETE", 2)
			}
			switch arg := args[0].(type) {
			case *object.Array:
				arg2, ok := args[1].(*object.Int)
				if !ok {
					return errorwords.SetError(304, line, "整数")
				}
				num := int(arg2.Value)
				if num < 0 {
					return errorwords.SetError(304, line, "0以上")
				}
				if num >= len(arg.Elements) {
					return errorwords.SetError(304, line, strconv.Itoa(len(arg.Elements)-1)+"以下")
				}
				before := arg.Inspect()
				arg.Elements = delete(arg.Elements, num)
				log.SetLog(line, "ADD("+before+")", arg.Inspect(), "組み込み関数DELETEを実行")
				return nil
			default:
				return errorwords.SetError(302, line)
			}
		},
	},
	//SLICE string,array
	"SLICE": &object.BuiltIn{
		Func: func(line int, args ...object.Object) object.Object {
			if len(args) > 3 || len(args) < 2 {
				return errorwords.SetError(300, line, "SLICE", "2個または３個")
			}
			switch arg := args[0].(type) {
			case *object.String:
				arg2, ok := args[1].(*object.Int)
				if !ok {
					return errorwords.SetError(305, line)
				}
				max := utf8.RuneCountInString(arg.Value)
				start := int(arg2.Value)
				if start < 1 {
					return errorwords.SetError(306, line, 1, "1以上")
				}
				if start > max {
					return errorwords.SetError(306, line, 1, fmt.Sprintf("%v以下", max))
				}
				end := max
				if len(args) == 3 {
					arg3, ok := args[2].(*object.Int)
					num := int(arg3.Value)
					if !ok {
						return errorwords.SetError(305, line)
					}
					if num > max {
						return errorwords.SetError(306, line, 2, fmt.Sprintf("%v以下", max))
					}
					if num < 1 {
						return errorwords.SetError(306, line, 1, "1以上")
					}
					end = num
				}
				utfStr := utf8string.NewString(arg.Value)
				result := utfStr.Slice(start-1, end)
				log.SetLog(line, "SLICE("+arg.Inspect()+")", result, "組み込み関数SLICEを実行")
				return &object.String{Value: result, Line: line}
			case *object.Array:
				arg2, ok := args[1].(*object.Int)

				if !ok {
					return errorwords.SetError(305, line)
				}
				max := len(arg.Elements)
				start := int(arg2.Value)
				if start < 0 {
					return errorwords.SetError(306, line, 1, "0以上")
				}
				if start > max {
					return errorwords.SetError(306, line, 1, fmt.Sprintf("%v以下", max))
				}
				end := max
				if len(args) == 3 {
					arg3, ok := args[2].(*object.Int)
					num := int(arg3.Value)
					if !ok {
						return errorwords.SetError(305, line)
					}
					if num > max {
						return errorwords.SetError(306, line, 2, fmt.Sprintf("%v以下", max))
					}
					if num < 0 {
						return errorwords.SetError(306, line, 1, "0以上")
					}
					end = num
				}
				result := &object.Array{Elements: arg.Elements[start:end], Line: line}
				log.SetLog(line, "SLICE("+arg.Inspect()+")", result.Inspect(), "組み込み関数SLICEを実行")
				return result
			default:
				return errorwords.SetError(301, line, "SLICE", "文字列,配列")
			}
		},
	},
	//stdin
	"GET": &object.BuiltIn{
		Func: func(line int, args ...object.Object) object.Object {
			if len(args) != 0 {
				return errorwords.SetError(300, line, "GET", 0)
			}
			in := bufio.NewScanner(os.Stdin)
			in.Scan()
			val := in.Text()
			log.SetLog(line, "", `"`+val+`"`, "組み込み関数GETを実行")
			return &object.String{Value: val, Line: line}
		},
	},
	//stdin - int
	"GETNUM": &object.BuiltIn{
		Func: func(line int, args ...object.Object) object.Object {
			if len(args) != 0 {
				return errorwords.SetError(300, line, "GETNUM", 0)
			}
			in := bufio.NewScanner(os.Stdin)
			in.Scan()
			val := in.Text()
			for isNum([]rune(val)) == "ERROR" {
				fmt.Println("数字を入力してください")
				in.Scan()
				val = in.Text()
			}
			isnum := isNum([]rune(val))
			if isnum == "FLOAT" {
				float, err := strconv.ParseFloat(val, 64)
				if err != nil {
					return errorwords.SetError(112, line, float)
				}
				log.SetLog(line, "", fmt.Sprintf("%v", float), "組み込み関数GETNUMを実行")
				return &object.Float{Value: float, Line: line}
			}
			intnum, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return errorwords.SetError(111, line, val)
			}
			log.SetLog(line, "", val, "組み込み関数GETNUMを実行")
			return &object.Int{Value: intnum, Line: line}
		},
	},
	//calc root
	"ROOT": &object.BuiltIn{
		Func: func(line int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return errorwords.SetError(300, line, "ROOT", 1)
			}
			switch arg := args[0].(type) {
			case *object.Int:
				if arg.Value < 0 {
					return errorwords.SetError(301, line, "ROOT", "0以上")
				}
				log.SetLog(line, "ROOT("+arg.Inspect()+")", fmt.Sprintf("%v", math.Sqrt(float64(arg.Value))), "組み込み関数ROOTを実行")
				return &object.Float{Value: math.Sqrt(float64(arg.Value)), Line: line}
			case *object.Float:
				if arg.Value < 0 {
					return errorwords.SetError(301, line, "ROOT", "0以上")
				}
				log.SetLog(line, "ROOT("+arg.Inspect()+")", fmt.Sprintf("%v", math.Sqrt(arg.Value)), "組み込み関数ROOTを実行")
				return &object.Float{Value: math.Sqrt(arg.Value), Line: line}
			default:
				return errorwords.SetError(301, line, "ROOT", "数値")
			}
		},
	},
	//string convert to int or float
	"TONUM": &object.BuiltIn{
		Func: func(line int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return errorwords.SetError(300, line, "TONUM", 1)
			}
			switch arg := args[0].(type) {
			case *object.String:
				str := []rune(arg.Value)
				isnum := isNum(str)
				if isnum == "ERROR" {
					return errorwords.SetError(301, line, "TONUM", "整数または少数を表す文字列")
				}
				if isnum == "FLOAT" {
					val, err := strconv.ParseFloat(arg.Value, 64)
					if err != nil {
						return errorwords.SetError(112, line, arg.Value)
					}
					log.SetLog(line, `"`+arg.Value+`"`, fmt.Sprintf("%v", val), "組み込み関数TONUMを実行")
					return &object.Float{Value: val, Line: line}
				}
				val, err := strconv.ParseInt(arg.Value, 10, 64)
				if err != nil {
					return errorwords.SetError(111, line, arg.Value)
				}
				log.SetLog(line, `"`+arg.Value+`"`, fmt.Sprintf("%v", val), "組み込み関数TONUMを実行")
				return &object.Int{Value: val, Line: line}
			default:
				return errorwords.SetError(301, line, "TONUM", "文字列")
			}
		},
	},
	//rand(min,max)
	"RAND": &object.BuiltIn{
		Func: func(line int, args ...object.Object) object.Object {
			if len(args) != 2 {
				return errorwords.SetError(300, line, "RAND", 2)
			}
			switch min := args[0].(type) {
			case *object.Int:
				switch max := args[1].(type) {
				case *object.Int:
					return &object.Int{Value: randInt(min.Value, max.Value, line), Line: line}
				case *object.Float:
					return &object.Float{Value: randFloat(float64(min.Value), max.Value, line), Line: line}
				default:
					return errorwords.SetError(301, line, "RAND", "数値")
				}
			case *object.Float:
				switch max := args[1].(type) {
				case *object.Int:
					return &object.Float{Value: randFloat(min.Value, float64(max.Value), line), Line: line}
				case *object.Float:
					return &object.Float{Value: randFloat(float64(min.Value), max.Value, line), Line: line}
				default:
					return errorwords.SetError(301, line, "RAND", "数値")
				}
			default:
				return errorwords.SetError(301, line, "RAND", "数値")
			}
		},
	},
	//Print
	"SAY": &object.BuiltIn{
		Func: func(line int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return errorwords.SetError(300, line, "SAY", 1)
			}
			if str, ok := args[0].(*object.String); ok {
				fmt.Println(str.Value)
				return nil
			}
			log.SetLog(line, args[0].Inspect(), "出力", "組み込み関数SAYを実行")
			fmt.Println(args[0].Inspect())
			return nil
		},
	},
	//Wait some time
	"SLEEP": &object.BuiltIn{
		Func: func(line int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return errorwords.SetError(300, line, "SAY", 1)
			}
			s, ok := args[0].(*object.Int)
			if !ok {
				return errorwords.SetError(301, line, "SLEEP", "整数")
			}
			second := int(s.Value)
			log.SetLog(line, args[0].Inspect()+"秒", "待つ", "組み込み関数SLEEPを実行")
			time.Sleep(time.Duration(second) * time.Second)
			return nil
		},
	},
}

//isNumber If a character is number return true
func isNumber(ch rune) bool {
	if '0' <= ch && ch <= '9' {
		return true
	}
	return false
}

//isNumber If a character is number(int ot float) return (float,int)
func isNum(str []rune) string {
	isFloat := false
	for i := 0; i < len(str); i++ {
		if !isNumber(str[i]) && str[i] != '.' {
			return "ERROR"
		}
		if str[i] == '.' {
			if isFloat {
				return "ERROR"
			}
			isFloat = true
		}
	}
	if isFloat {
		return "FLOAT"
	}
	return "INT"
}

//remove object from slice
func delete(els []object.Object, i int) []object.Object {
	return append(els[:i], els[i+1:]...)
}

//random Int number
func randInt(min int64, max int64, line int) int64 {
	rand.Seed(time.Now().UnixNano())
	result := rand.Int63n(max-min) + min
	log.SetLog(line, "乱数("+fmt.Sprintf("%v", min)+`~`+fmt.Sprintf("%v", max)+")", fmt.Sprintf("%v", result), "組み込み関数RANDを実行")
	return result
}

//random Int number
func randFloat(min float64, max float64, line int) float64 {
	rand.Seed(time.Now().UnixNano())
	result := rand.Float64()*(max-min) + min
	log.SetLog(line, "乱数("+fmt.Sprintf("%v", min)+`~`+fmt.Sprintf("%v", max)+")", fmt.Sprintf("%v", result), "組み込み関数RANDを実行")
	return result
}
