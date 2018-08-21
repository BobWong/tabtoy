package printer

import (
	"github.com/BobWong/golog"
)

var log *golog.Logger = golog.New("printer")

func GetLog() *golog.Logger {
	return log
}
