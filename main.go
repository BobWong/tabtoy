package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/bobwong89757/golog"
	"tabtoy/v2"
	"tabtoy/v2/printer"
)

var log = golog.New("main")

const (
	Version_v2 = "2.9.1"
	Version_v3 = "3.0.0"
)

func main() {

	flag.Parse()

	// 版本
	if *paramVersion {
		fmt.Printf("%s, %s", Version_v2, Version_v3)
		return
	}
	printer.GetLog().SetEnable(*paramLogEnable)
	v2.GetLog().SetEnable(*paramLogEnable)

	switch *paramMode {
	case "v2":
		V2Entry()
	default:
		fmt.Println("--mode not specify")
		os.Exit(1)
	}

}
