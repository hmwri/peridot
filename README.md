# PeriDot

## 日本語

PeriDot は、テキストでのプログラミングを始める人を対象にした初心者・教育向けのプログラミング言語です。シンプルな手続き型のインタプリタ型言語で、日本語エラー表示と実行過程の確認に対応しています。

Web 版の Peree では、ブラウザ上で PeriDot を実行し、コードの保存や共有を行うことを目指しています。

参考: https://demo.honma.site/peridot/

### 主な特徴

- Go 製のインタプリタ
- REPL による対話実行
- ファイル指定による実行
- 日本語のエラーメッセージ
- `-l` オプションによる実行過程ログ
- C、Go、JavaScript に近い手続き型の書き方

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

Reference: https://demo.honma.site/peridot/

### Features

- Go-based interpreter
- Interactive REPL
- File execution
- Japanese error messages
- Execution logs with the `-l` option
- C/Go/JavaScript-like procedural syntax

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
