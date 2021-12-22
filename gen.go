package gen

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

var (
	tpl       = template.New("")
	parserMap = map[string]Parser{}
)

func init() {
	for _, p := range []Parser{&dockerParser{}} {
		parserMap[p.Name()] = p
	}

	if _, err := tpl.ParseGlob("tmpls/*.tmpl"); err != nil {
		log.Panicf("gen: parse glob failed: %w", err)
	}
}

type Parser interface {
	Name() string
	Parse([]byte) (*Result, error)
}

type Result struct {
	Name        string
	Desc        string
	Options     []*Option
	SubCommands []*Command
}

type Command struct {
	Name        string
	Desc        string
	Args        []*Arg
	SubCommands []*Command
	Options     []*Option
}

type Option struct {
	Name string
	Desc string
	Args []*Arg
}

type Arg struct {
	Name       string
	Template   string
	IsVariadic bool
}

func Gen(cmd string, w io.Writer) error {
	p, ok := parserMap[cmd]
	if !ok {
		return fmt.Errorf("gen: unsupported command: %s", cmd)
	}

	output, err := runCommand(cmd)
	if err != nil {
		return fmt.Errorf("gen: run command error, command: %s, error: %w", cmd, err)
	}
	r, err := p.Parse(output)
	if err != nil {
		return err
	}

	tmplName := cmd + ".tmpl"
	if err := tpl.ExecuteTemplate(w, tmplName, r); err != nil {
		return fmt.Errorf("gen: execute template error, name: %s, error: %w", tmplName, err)
	}
	return nil
}

func runCommand(cmd string) ([]byte, error) {
	cs := strings.Split(cmd, " ")
	args := append(cs[1:], "--help")

	c := exec.Command(cs[0], args...)
	c.Stdin = os.Stdin
	data, err := c.Output()
	if err != nil {
		return nil, fmt.Errorf("gen: run command error, command: %s, error: %w", cmd, err)

	}
	return data, nil
}
