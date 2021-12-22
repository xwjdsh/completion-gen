package main

import (
	"flag"
	"os"

	gen "github.com/xwjdsh/completion-gen"
)

var (
	fp      = flag.String("w", "", "write file path")
	command = flag.String("c", "", "command")
)

func main() {
	flag.Parse()

	if *command == "" {
		panic("command not set")
	}
	writer := os.Stdout
	if *fp != "" {
		f, err := os.Create(*fp)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		writer = f
	}

	if err := gen.Gen(*command, writer); err != nil {
		panic(err)
	}
}
