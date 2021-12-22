package main

import (
	"os"

	gen "github.com/xwjdsh/completion-gen"
)

func main() {
	if err := gen.Gen("docker", os.Stdout); err != nil {
		panic(err)
	}
}
