package main

import (
	"eduC/errorwords"
	"eduC/eval"
	"eduC/info"
	"eduC/lexer"
	"eduC/log"
	"eduC/object"
	"eduC/parser"
	"eduC/repl"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	//日本語エラー
	errorwords.Jerror()
	//引数
	arglen := len(os.Args)
	if arglen == 1 {
		repl.Open()
		return
	}
	if arglen > 3 {
		fmt.Println("priの引数の数は最大2個です")
		return
	}
	if arglen == 2 {
		if os.Args[1][0] == '-' {
			selectOption(arglen)
			return
		}
		env := object.NewEnv()
		readMode(os.Args[1], env, false)
		return
	}
	if os.Args[1][0] != '-' {
		fmt.Println("オプションは - から初めてください")
		options()
	}
	selectOption(arglen)

}
func selectOption(arglen int) {
	switch os.Args[1][1] {
	case 'l':
		if arglen == 2 {
			fmt.Println("ファイル名を指定してください")
			return
		}
		env := object.NewEnv()
		readMode(os.Args[2], env, true)
	case 'v':
		fmt.Println("PeriDot " + info.Version + info.CheckVersion())
	case 'h':
		help()
	default:
		fmt.Println(os.Args[1] + "というオプションが見つかりません")
		options()
	}
}
func options() {
	fmt.Println("オプション一覧")
	fmt.Println("[-l] ログ(実行過程)を表示")
	fmt.Println("[-v] バージョンを表示")
	fmt.Println("[-h] ヘルプを表示")
}
func help() {
	fmt.Println("ヘルプ")
	help := `
対話形式で実行したい
-> "pri"
ファイルを指定して実行したい
-> "pri ファイル名"
実行過程をみたい
-> "pri -l ファイル名"

詳しくは>hmwri.com/pridot/manual/`
	fmt.Println(help)
	options()
}
func readMode(t string, env *object.Env, logswitch bool) {
	if filepath.Ext(t) != ".pri" {
		fmt.Println("priファイルを指定してください")
		return
	}
	f, err := os.Open(t)
	if err != nil {
		fmt.Println("そのようなファイルはみつかりません")
		return
	}
	defer f.Close()
	w, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("ファイルの読み込み中に問題発生しました")
		return
	}
	log.ResetLogs()
	l := lexer.New(string(w))
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
			fmt.Printf("(´；ω；)<おっと!:%v\n", errorwords.Errp.Error[0].Message)
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
		fmt.Printf("%v\n", w.Message)
	}
	return true
}
