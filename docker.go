package gen

import (
	"fmt"
	"os"
	"strings"
)

type dockerParser struct{}

func (p *dockerParser) Name() string {
	return "docker"
}

func (p *dockerParser) Parse(output []byte) (*Result, error) {
	r, err := p.parse(p.Name(), output)
	if err != nil {
		return nil, err
	}
	return &Result{
		Name:        "docker",
		Desc:        "A self-sufficient runtime for containers",
		Options:     r.Options,
		SubCommands: r.SubCommands,
	}, nil
}

func (p *dockerParser) parse(command string, output []byte) (*Result, error) {
	isOptions, isCommand := false, false
	optionsLines, commandLines := []string{}, []string{}
	for _, line := range strings.Split(string(output), "\n") {
		if strings.TrimSpace(line) == "" {
			isOptions, isCommand = false, false
		}

		if ok := strings.HasPrefix(line, "Options"); ok || isOptions {
			isOptions = true
			if !ok {
				optionsLines = append(optionsLines, line)
			}
		}

		if ok := (strings.HasPrefix(line, "Management Commands") || strings.HasPrefix(line, "Commands")); ok || isCommand {
			isCommand = true
			if !ok {
				commandLines = append(commandLines, line)
			}
		}
	}
	commands, err := p.parseCommands(command, commandLines)
	if err != nil {
		return nil, err
	}

	return &Result{
		Options:     p.parseOptions(optionsLines),
		SubCommands: commands,
	}, nil
}

func (p *dockerParser) parseCommands(command string, lines []string) ([]*Command, error) {
	cmds := []*Command{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		arr := strings.SplitN(line, " ", 2)
		cmd := &Command{
			Name: strings.TrimSuffix(strings.TrimSpace(arr[0]), "*"),
			Desc: strings.TrimSpace(arr[1]),
		}
		subCommand := fmt.Sprintf("%s %s", command, cmd.Name)
		subCommandOutput, err := runCommand(subCommand)
		if err != nil {
			return nil, err
		}

		r, err := p.parse(subCommand, subCommandOutput)
		if err != nil {
			return nil, err
		}

		cmd.SubCommands = r.SubCommands
		cmd.Options = r.Options
		cmds = append(cmds, cmd)
	}
	return cmds, nil
}

func (p *dockerParser) parseOptions(lines []string) []*Option {
	opts := []*Option{}
	hasHelp := false

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "-") {
			opts[len(opts)-1].Desc += "\n" + convertDesc(line)
			continue
		}

		r := strings.SplitN(line, "  ", 2)
		str := r[0]
		desc := convertDesc(r[1])
		arr := strings.Split(str, " ")

		opt := &Option{Desc: desc}
		var arg *Arg
		if strings.HasSuffix(arr[0], ",") {
			opt.Name = fmt.Sprintf(`["%s", "%s"]`, arr[0][0:len(arr[0])-1], arr[1])
			if len(arr) > 2 {
				arg = &Arg{
					Name: arr[2],
				}
			}
		} else {
			if arr[0] == "--help" {
				hasHelp = true
			}
			opt.Name = fmt.Sprintf(`["%s"]`, arr[0])
			if len(arr) > 1 {
				arg = &Arg{
					Name: arr[1],
				}
			}
		}
		if arg != nil {
			opt.Args = append(opt.Args, arg)
		}
		opts = append(opts, opt)
	}
	if !hasHelp {
		opts = append(opts, &Option{
			Name: `"--help"`,
			Desc: "Print usage",
		})
	}

	return opts
}

func convertDesc(desc string) string {
	tmp := strings.ReplaceAll(strings.TrimSpace(desc), `"`, `'`)
	return strings.ReplaceAll(tmp, os.Getenv("HOME"), "$HOME")
}
