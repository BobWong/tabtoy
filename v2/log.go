package v2

import "github.com/bobwong89757/golog"

var log = golog.New("exportorv2")

func GetLog() *golog.Logger {
	return log
}
