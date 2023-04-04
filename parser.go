package yacli

import (
	"fmt"
	"strings"
)

type repository struct {
	// Positional argument.
	posarg []string
	// Flags.
	mf map[string]string
}

type parser struct {
	fm   int
	args []string
}

func newParser(args []string) *parser {
	return &parser{args: args}
}

func (p *parser) parse() (repository, error) {
	var posargs []string

	r := repository{mf: make(map[string]string)}

	var i int
	for i < len(p.args) {
		arg := p.args[i]

		switch {
		case isLongFlag(arg):
			targ := strings.TrimPrefix(arg, "--")
			if len(arg) == 0 {
				return repository{}, fmt.Errorf("'%s' is not valid flag", arg)
			}

			var name, value string
			parts := strings.SplitN(targ, "=", 2)
			name = parts[0]

			if len(name) == 0 {
				return repository{}, fmt.Errorf("flag name can not be empty: '%s'", arg)
			}

			if len(parts) == 2 {
				value = parts[1]
			} else if i+1 < len(p.args) && !(isShortFlag(p.args[i+1]) || isLongFlag(p.args[i+1])) {
				value = p.args[i+1]
				i++
			}

			r.mf[name] = value
			i++
		case isShortFlag(arg):
		default:
			posargs = append(posargs, arg)
			i++
		}
	}

	return r, nil
}

func isLongFlag(arg string) bool {
	return strings.HasPrefix(arg, "--")
}

func isShortFlag(arg string) bool {
	return strings.HasPrefix(arg, "-")
}
