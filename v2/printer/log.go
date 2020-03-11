package printer

import (
	"github.com/bobwong89757/golog"
)

var log *golog.Logger = golog.New("printer")

func GetLog() *golog.Logger {
	return log
}
