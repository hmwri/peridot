# PeriDot

## 日本語

PeriDot は、テキストでのプログラミングを始める人を対象にした初心者・教育向けのプログラミング言語です。シンプルな手続き型のインタプリタ型言語で、日本語エラー表示と実行過程の確認に対応しています。

Web 版の Peree では、ブラウザ上で PeriDot を実行し、コードの保存や共有を行うことを目指しています。

参考:

- https://demo.honma.site/peridot/
- https://demo.honma.site/peridot/manual/

### 主な特徴

- Go 製のインタプリタ
- REPL による対話実行
- ファイル指定による実行
- 日本語のエラーメッセージ
- `-l` オプションによる実行過程ログ
- C、Go、JavaScript に近い手続き型の書き方

### PeriDot のサンプルコード

PeriDot は、`make` で変数を作り、`if`、`loop`、`func` で制御構造や関数を表現します。標準出力、入力、配列操作、数値変換などは大文字の組み込み関数で扱います。

Hello World:

```peridot
SAY("Hello, PeriDot!")
```

変数、計算、条件分岐:

```peridot
make score = 82
make bonus = 8
make total = score + bonus

if total >= 90 {
  SAY("excellent")
} else {
  SAY("keep going")
}
```

`loop` による繰り返し:

```peridot
make sum = 0
make i = 1

loop i <= 5 {
  sum = sum + i
  i = i + 1
}

SAY(sum)
```

関数:

```peridot
func double(x) {
  return x * 2
}

SAY(double(21))
```

配列と組み込み関数:

```peridot
make values = [1, 2, 3]

ADD(values, 4)
SAY(SIZE(values))
SAY(values[0])

DELETE(values, 1)
SAY(values)
```

入力、数値変換、平方根:

```peridot
SAY("数字を入力してください")
make input = GET()
make number = TONUM(input)

SAY(ROOT(number))
```

### 文法メモ

- 変数定義は `make name = value`、再代入は `name = value` と書きます。
- 条件分岐は `if condition { ... } else { ... }` です。
- 繰り返しは `loop condition { ... }`、または条件なしの `loop { ... }` です。
- 関数は `func name(arg1, arg2) { ... }` で定義し、`return` で値を返します。
- 配列は `[1, 2, 3]` と書き、`values[0]` のように参照します。
- `SAY`、`GET`、`GETNUM`、`TONUM`、`SIZE`、`ADD`、`DELETE`、`SLICE`、`RAND`、`ROOT`、`SLEEP` などの組み込み関数があります。

### 構成

- `main.go`: CLI の入口
- `lexer/`: 字句解析
- `parser/`: 構文解析
- `ast/`: 抽象構文木
- `eval/`: 評価器
- `object/`: 実行時オブジェクトと環境
- `repl/`: 対話実行環境
- `token/`: トークン定義
- `errorwords/`: 日本語エラー表示
- `log/`: 実行過程ログ
- `info/`: バージョン情報

### 必要環境

- Go 1.24.2 以降

### 実行方法

REPL を起動します。

```bash
go run .
```

ファイルを実行します。

```bash
go run . path/to/file.pri
```

実行過程を表示しながらファイルを実行します。

```bash
go run . -l path/to/file.pri
```

バージョンやヘルプを表示します。

```bash
go run . -v
go run . -h
```

### REPL の操作

- `Q!`: REPL を終了
- `LOG!`: 実行過程ログの表示を切り替え

### メモ

このリポジトリは言語処理系の実装です。Web 実行環境や教育向け UI まで含む Peree とは別に、PeriDot のコアとなるインタプリタ部分を管理しています。

## English

PeriDot is an educational programming language for people moving from visual programming or first programming lessons into text-based programming. It is a small procedural interpreted language with Japanese error messages and optional execution logs.

References:

- https://demo.honma.site/peridot/
- https://demo.honma.site/peridot/manual/

### Features

- Go-based interpreter
- Interactive REPL
- File execution
- Japanese error messages
- Execution logs with the `-l` option
- C/Go/JavaScript-like procedural syntax

### PeriDot Language Examples

PeriDot uses `make` for variable declarations, `if`, `loop`, and `func` for control flow and functions, and uppercase built-ins for output, input, conversion, math, and array operations.

Hello World:

```peridot
SAY("Hello, PeriDot!")
```

Variables, calculation, and branching:

```peridot
make score = 82
make bonus = 8
make total = score + bonus

if total >= 90 {
  SAY("excellent")
} else {
  SAY("keep going")
}
```

Loop:

```peridot
make sum = 0
make i = 1

loop i <= 5 {
  sum = sum + i
  i = i + 1
}

SAY(sum)
```

Function:

```peridot
func double(x) {
  return x * 2
}

SAY(double(21))
```

Arrays and built-ins:

```peridot
make values = [1, 2, 3]

ADD(values, 4)
SAY(SIZE(values))
SAY(values[0])

DELETE(values, 1)
SAY(values)
```

Input, conversion, and square root:

```peridot
SAY("Enter a number")
make input = GET()
make number = TONUM(input)

SAY(ROOT(number))
```

### Syntax Notes

- Declare variables with `make name = value`; reassign them with `name = value`.
- Branch with `if condition { ... } else { ... }`.
- Repeat with `loop condition { ... }`, or use `loop { ... }` for an unconditional loop.
- Define functions with `func name(arg1, arg2) { ... }` and return values with `return`.
- Create arrays with `[1, 2, 3]` and read items with expressions such as `values[0]`.
- Built-ins include `SAY`, `GET`, `GETNUM`, `TONUM`, `SIZE`, `ADD`, `DELETE`, `SLICE`, `RAND`, `ROOT`, and `SLEEP`.

### Project Structure

- `main.go`: CLI entry point
- `lexer/`: lexer
- `parser/`: parser
- `ast/`: abstract syntax tree
- `eval/`: evaluator
- `object/`: runtime objects and environment
- `repl/`: interactive shell
- `token/`: token definitions
- `errorwords/`: Japanese error messages
- `log/`: execution logging
- `info/`: version information

### Requirements

- Go 1.24.2 or later

### Usage

Start the REPL:

```bash
go run .
```

Run a file:

```bash
go run . path/to/file.pri
```

Run a file with execution logs:

```bash
go run . -l path/to/file.pri
```

Show version or help:

```bash
go run . -v
go run . -h
```

### REPL Commands

- `Q!`: quit the REPL
- `LOG!`: toggle execution logs
