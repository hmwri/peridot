package repl

import (
	"bufio"
	"fmt"
	"github.com/hmwri/peridot/errorwords"
	"github.com/hmwri/peridot/eval"
	"github.com/hmwri/peridot/info"
	"github.com/hmwri/peridot/lexer"
	"github.com/hmwri/peridot/log"
	"github.com/hmwri/peridot/object"
	"github.com/hmwri/peridot/parser"
	"os"
	"os/user"
	"strings"
)

func Open() {
	user, err := user.Current()
	log := false
	if err != nil {
		panic(err)
	}
	op := `+===+ +==== +=====  *** 
|   | |     |   /    |
+===+ +==== |=/      |      >>
|     |     | \      |     >>>
+     +==== +  \_   ***   >>>>`
	fmt.Printf("%v\n", op)
	fmt.Printf("PeriDot " + info.Version + info.CheckVersion() + "\n")
	fmt.Println("ぜひフィードバックにご協力ください!リンク:https://forms.gle/Cca4668Tah7x2o5YA\n")
	fmt.Printf("%vさんようこそ！ここでは対話式プログラム実行ができます！\n終了：Q!, ログ(実行過程)表示:LOG!\n>>", user.Username)
	scanner := bufio.NewScanner(os.Stdin)
	env := object.NewEnv()
	indent := 0
	statement := ""
	for scanner.Scan() {
		context := scanner.Text()
		if context == "Q!" {
			fmt.Printf("終了！(QUIT!)\n")
			break
		}
		if context == "LOG!" {
			if log {
				fmt.Printf("ログ表示機能をオフにしました！\n")
				log = false
			} else {
				fmt.Printf("ログ表示機能をオンにしました！\n")
				log = true
			}
			fmt.Printf(">>")
			continue
		}
		if plusInd := strings.Count(context, "{"); plusInd > 0 {
			if indent > -1 {
				indent += plusInd
			}
		}
		if indent > 0 {
			statement = statement + context + "\n"
		} else {
			statement = context
		}
		if minusInd := strings.Count(context, "}"); minusInd > 0 {
			if indent > -1 {
				indent -= minusInd
			}
		}
		for i := 0; i < indent; i++ {
			fmt.Printf("...")
		}
		if indent == 0 {
			writeMode(statement, env, log)
			fmt.Printf(">>")
			statement = ""
		}

	}

}

func writeMode(t string, env *object.Env, logswitch bool) {
	log.ResetLogs()
	l := lexer.New(t)
	p := parser.New(l)
	program := p.Parse()
	if !Checkerror(p) {
		eval := eval.Eval(program, env)
		if logswitch {
			for _, v := range log.Logs {
				fmt.Printf("[%d] %d行目: %v -> %v 「%v」\n", v.Num, v.Line, v.Evaled, v.Toeval, v.Message)
			}
		}
		if len(errorwords.Errp.Error) != 0 {
			fmt.Printf("\x1b[31m(´；ω；)<おっと!:%v\x1b[0m\n", deleteLine(errorwords.Errp.Error[0].Message))
		} else {
			if eval != nil {
				fmt.Printf("(≧▽≦)Answer:")
				fmt.Println(eval.Inspect())
			}
		}
	}
}

func Checkerror(p *parser.Parser) bool {
	e := p.GetError()
	n := len(e)
	if n == 0 {
		return false
	}
	fmt.Printf("(´；ω；)<文法のエラーが%v個あります\n", n)
	for _, w := range e {
		fmt.Printf("\x1b[31m%v\x1b[0m\n", deleteLine(w.Message))
	}
	return true
}
func deleteLine(str string) string {
	return string([]rune(str)[5:])
}
