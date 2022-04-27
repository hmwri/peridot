package lexer

//Oridinal struct
type Oridinal struct {
	//input Target source
	input string
	//runeinput Input change to rune(Unicode support)
	runeinput []rune
	//position Current position
	position int
	//readPosition Reading position
	readPosition int
	//ch Testing character
	ch rune
}

//New Make Oridinal struct and call readChar
func newC(input string) *Oridinal {
	o := &Oridinal{input: input, runeinput: []rune(input)}
	o.runeinput = append(o.runeinput, '\n')
	o.readChar()
	return o
}

//Delete Comment from Oridinal
func DeleteComment(s string) string {
	o := newC(s)
	deleted := ""
	readnum := 0
	inflag := false
	for readnum < len(o.runeinput) {
		if o.ch == '"' && inflag == false {
			inflag = true
		} else {
			if o.ch == '"' {
				inflag = false
			}
		}
		if o.ch == '<' && o.nextRead() == '<' && inflag == false {
			for !(o.ch == '>' && o.nextRead() == '>') {
				o.readChar()
				readnum++

			}
			o.readChar()
			readnum++
			o.readChar()
			readnum++
		}
		deleted += string(o.ch)
		o.readChar()

		readnum++
	}

	return deleted

}

//readChar Read next character
func (o *Oridinal) readChar() {
	if o.readPosition >= len(o.runeinput) {
		o.ch = 0 //EOF
	} else {
		o.ch = o.runeinput[o.readPosition]
	}
	o.position = o.readPosition
	o.readPosition++
}

//nextRead read next character
func (o *Oridinal) nextRead() rune {
	if o.readPosition >= len(o.runeinput) {
		return 0 //EOF
	}
	return o.runeinput[o.readPosition]
}
