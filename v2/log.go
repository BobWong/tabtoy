package v2

import "github.com/BobWong/golog"

var log = golog.New("exportorv2")

func GetLog() *golog.Logger {
	return log
}
