package errorwords

import (
	"fmt"
	"github.com/hmwri/peridot/log"
	"github.com/hmwri/peridot/object"
)

//Errors struct
type errors struct {
	Error []Error
}
type Error struct {
	Message string
	Line    int
}

var Errp *errors

//Err Error Words
var Err = map[int]string{}

//Jerror Japanese error words
func Jerror() {
	Errp = &errors{}
	Err[101] = "[%d行目]'%v'となるべきところが'%v'になっています！"

	Err[111] = "[%d行目]'%v'を整数に変換できません！桁数が大きすぎるかも！"
	Err[112] = "[%d行目]'%v'を少数に変換できません！桁数が大きすぎるかも！"
	Err[121] = "[%d行目]'%v'は適切な演算子ではありません！"
	Err[122] = "[%d行目]'%v'は適切な演算子ではありません！\n>>もしかして'=='?"

	Err[130] = "[%d行目]全角スペース(　)は使わないでください！"
	Err[200] = "[＞%d＜][エラー]"
	Err[201] = "[%d行目]前置演算子エラー。'%v'は不適です。値の前に置けるのは-と!のみです"
	Err[202] = "[%d行目]マイナスの後に整数、少数以外を置くことはできません。"
	Err[203] = "[%d行目]式の両側は基本的に同じ型である必要があります。左＝%v,右＝%vになっています"
	Err[204] = "[%d行目]式の記号が不適切です。'%v'は不適です。数値の間で使えるのは[+,-,*,/,<,<=,>,>=,!=,==,and,or]のみです"
	Err[205] = "[%d行目]式の記号が不適切です。'%v'は不適です。真偽値や式の間で使えるのは[!=,==,and,or]のみです"
	Err[206] = "[%d行目]%vの左側が不適切な値です。"
	Err[207] = "[%d行目]文字列の操作,比較につかえるのは'+','==','!='のみです"
	Err[208] = "[%d行目]'%v'は少数には使用できません"
	Err[210] = "[%d行目]%vはまだ定義されていません。変数:make 名前 = 値,関数:func 名前 (引数){}という形で定義してください"
	Err[211] = "[%d行目]変数にnil(空)を代入できません"
	Err[220] = "[%d行目]%vの条件式が不適切です"
	Err[230] = "[%d行目]名前がついた関数は変数に代入できません(関数名:%v)"
	Err[231] = "[%d行目]呼び出したものが関数ではありません(関数ではない:%v)"
	Err[232] = "[%d行目]引数の数が不適切です(呼び出し側:%v個,関数側:%v個)"
	Err[300] = "[%d行目]組み込み関数:%vの引数は%v個である必要があります"
	Err[301] = "[%d行目]組み込み関数:%vの第一引数は%vである必要があります"
	Err[302] = "[%d行目]組み込み関数:ADDの第一引数は配列である必要があります"
	Err[303] = "[%d行目]組み込み関数:ADD>同じ配列を追加することはできません"
	Err[304] = "[%d行目]組み込み関数:DELETE>対応する値が見つかりません。DELETEの第２引数は%vである必要があります"
	Err[305] = "[%d行目]組み込み関数:SLICEの第1,第2引数は整数である必要があります"
	Err[306] = "[%d行目]組み込み関数:SLICE>対応する値が見つかりません。SLICEの第%v引数は%vである必要があります"
	Err[307] = "[%d行目]組み込み関数:%vの色指定は整数の配列[赤の量,緑の量,青の量,透明度(任意)]で行ってください"
	Err[308] = "[%d行目]組み込み関数:%vの色指定は整数値(255以下)で行ってください"
	Err[309] = "[%d行目]組み込み関数:%vの第%v引数は%vである必要があります"
	Err[400] = "[%d行目]配列から値を取り出せませんでした。%vはサポートしていません"
	Err[401] = "[%d行目]配列から値を取り出せませんでした。添字は整数にしてください。例:Array[1]"
	Err[402] = "[%d行目]配列から値を取り出せませんでした。添字は0以上にしてください。例:Array[1]"
	Err[403] = "[%d行目]配列から値を取り出せませんでした。[ %v ]に対応する値がみつかりません。(添字は%v以下である必要があります)"
	Err[502] = "[%d行目]文字列から文字を取り出せませんでした。添字は1以上にしてください。例:'HELLO'[1]"
	Err[503] = "[%d行目]文字列から文字を取り出せませんでした。[ %v ]に対応する文字がみつかりません。(添字は%v以下である必要があります)"

	Err[500] = "[%d行目]ループの条件に%vは対応していません"

}

//Eerror Engilesh error words
func Eerror() {
	Errp = &errors{}
	Err[101] = "<Woops!!>[line:%v]I expected next token to be %v,but I got '%v'!"

}
func SetError(code int, params ...interface{}) object.Object {
	message := fmt.Sprintf(Err[code], params...)
	line := params[0].(int)
	Errp.Error = append(Errp.Error, Error{Message: message, Line: line})
	log.SetLog(line, "", "ERROR", message)
	return &object.ERROR{Value: message, Line: params[0].(int)}
}
