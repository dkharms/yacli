package yacli

import (
	"fmt"
	"strings"
)

// repository is a struct that represents a collection of arguments passed to CLI tool.
// It contains the positional arguments and flags provided by the user.
type repository struct {
	// A slice of strings representing the arguments passed before the flags.
	// These are the arguments that do not start with a dash (-) or double-dash (--).
	beforeFlags []string

	// A map that contains the flags and their values.
	// The keys represent the flag names, and the values represent their values.
	// Flags are the arguments that start with a dash (-) or double-dash (--).
	flagSet map[string]string

	// A slice of strings representing the positional arguments.
	// These are the arguments that are not flags, and their order matters.
	// They are provided after the beforeFlags and the flags.
	positionalArgs []string
}

// parser represents a command-line argument parser
// which holds the list of unparsed arguments and a boolean
// flag indicating whether the flags have been parsed or not.
type parser struct {
	// osargs is a slice of strings representing the command-line arguments
	// to be parsed.
	osargs []string

	// parsedFlags is a boolean value indicating whether or not the flags
	// have been parsed yet. If this value is true, then all remaining command-line
	// arguments are treated as positional arguments and are not parsed as flags.
	parsedFlags bool
}

func newParser(args []string) *parser {
	return &parser{osargs: args}
}

// parse is a method of the parser type which parses the
// command-line arguments and returns a repository and error.
// The repository contains the parsed flags and positional
// arguments, while the error will be non-nil if there was an error
// encountered during parsing.
func (p *parser) parse() (repository, error) {
	beforeFlags, _ := p.parseTillFlags()

	r := repository{
		beforeFlags:    beforeFlags,
		flagSet:        make(map[string]string),
		positionalArgs: []string{},
	}

	var i int
	for i < len(p.osargs) {
		arg := p.osargs[i]

		switch {
		case isLongFlag(arg):
			if p.parsedFlags {
				return repository{}, fmt.Errorf("already parsed flags: invalid argument '%s'", arg)
			}

			targ := strings.TrimPrefix(arg, "--")
			if len(arg) == 0 {
				return repository{}, fmt.Errorf("invalid flag: '%s'", arg)
			}

			var name, value string
			parts := strings.SplitN(targ, "=", 2)
			name = parts[0]

			if len(name) == 0 {
				return repository{}, fmt.Errorf("flag name can not be empty: '%s'", arg)
			}

			if len(parts) == 2 {
				value = parts[1]
			} else if i+1 < len(p.osargs) && !(isShortFlag(p.osargs[i+1]) || isLongFlag(p.osargs[i+1])) {
				value = p.osargs[i+1]
				i++
			}

			r.flagSet[name] = value
		case isShortFlag(arg):
			if p.parsedFlags {
				return repository{}, fmt.Errorf("already parsed flags: invalid argument '%s'", arg)
			}

			targ := strings.TrimPrefix(arg, "-")
			if len(targ) == 0 {
				return repository{}, fmt.Errorf("invalid flag: '%s'", arg)
			}

			if len(targ) > 1 {
				for _, l := range targ {
					r.flagSet[string(l)] = ""
				}
			} else if i+1 < len(p.osargs) && !(isShortFlag(p.osargs[i+1]) || isLongFlag(p.osargs[i+1])) {
				r.flagSet[targ] = p.osargs[i+1]
				i++
			} else {
				r.flagSet[targ] = ""
			}
		default:
			r.positionalArgs = append(r.positionalArgs, arg)
			p.parsedFlags = true
		}

		i++
	}

	return r, nil
}

// parseTillFlags is a method of the parser struct
// that parses the command line arguments passed to a GoLang program till it encounters a flag argument.
// It returns a slice of strings containing the arguments before the flags
// and a boolean value indicating whether it encountered a flag argument or not.
func (p *parser) parseTillFlags() ([]string, bool) {
	var beforeFlags []string

	for len(p.osargs) > 0 {
		arg := p.osargs[0]
		if isShortFlag(arg) || isLongFlag(arg) {
			return beforeFlags, true
		}
		p.osargs = p.osargs[1:]
		beforeFlags = append(beforeFlags, arg)
	}

	return beforeFlags, false
}

// isLongFlag is a function that takes a single argument
// of type string and returns a boolean value indicating
// whether the given string is a long flag or not.
// A long flag is a command-line argument that starts with a double dash (--).
func isLongFlag(arg string) bool {
	return strings.HasPrefix(arg, "--")
}

// isShortFlag is a function that takes a single argument
// of type string and returns a boolean value indicating
// whether the given string is a short flag or not.
// A short flag is a command-line argument that starts with a single dash (-).
func isShortFlag(arg string) bool {
	return strings.HasPrefix(arg, "-")
}
