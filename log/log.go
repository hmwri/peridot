package log

var Logs []Log
var num = 0

//Log evaluate log
type Log struct {
	Num     int
	Line    int
	Evaled  string
	Toeval  string
	Message string
}

//ResetLog
func ResetLogs() {
	Logs = []Log{}
	num = 0
}

//SetLog
func SetLog(line int, from string, to string, message string) {
	num++
	l := Log{}
	l.Num = num
	l.Line = line
	l.Evaled = from
	l.Toeval = to
	l.Message = message
	Logs = append(Logs, l)
}
