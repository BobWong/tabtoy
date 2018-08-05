package v2

import "github.com/davyxu/golog"

var log = golog.New("exportorv2")

func GetLog() *golog.Logger {
	return log
}
