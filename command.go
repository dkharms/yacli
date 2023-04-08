package yacli

import (
	"fmt"
	"os"
	"strings"
)

type commandOption func(*command)

// commandset is a set of commands mapped by their name.
type commandset map[string]*command

// set is a method of the commandset struct that takes two arguments:
// a string representing the name of the command to add, and a pointer to a command struct.
// It returns a boolean value indicating whether the command was successfully added to the command set.
// If a command with the same name already exists in the command set, the function will return false.
func (cs commandset) set(name string, c *command) bool {
	if _, ok := cs[name]; ok {
		return false
	}
	cs[name] = c
	return true
}

// get returns the command associated with the given name and a bool indicating
// whether the command was found or not. Returns nil and false if the command is
// not present in the command set.
func (cs commandset) get(name string) (*command, bool) {
	c, ok := cs[name]
	return c, ok
}

// command represents a command definition which has a name,
// a description, flags, arguments, and an action.
type command struct {
	// name is the name of the command.
	name string

	// desc is a brief description of the command.
	desc string

	// deprecated is a flag that indicates if the command is deprecated or not
	deprecated bool

	// cs is the sub-commands under this command.
	cs commandset

	// fsl is the flags associated with this command.
	fsl flagset

	fss flagset

	// as is the arguments associated with this command.
	as argset

	// action is the function to execute when the command is invoked.
	action func(Context) error
}

// NewRootCommand returns a new instance of a `command` struct with the name of the command
// set to the base name of the program.
// It accepts an optional argument `opts`, which is a
// variadic parameter of `commandOption` type, representing functional options for customizing the command.
func NewRootCommand(opts ...commandOption) *command {
	return NewCommand(strings.TrimPrefix(os.Args[0], "./"), opts...)
}

// NewCommand creates a new command with the specified name and options.
// The command is initialized with an empty flagset and commandset.
// It then applies the given options to the command.
// Returns a pointer to the created command.
func NewCommand(name string, opts ...commandOption) *command {
	c := &command{
		name: name, fsl: make(flagset), fss: make(flagset), cs: make(commandset),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// WithSubcommand returns a commandOption that adds a subcommand to a command.
//
// The subcommand is added to the parent command's commandset using the subcommand's name as the key.
// The subcommand can be invoked by calling the parent command with the subcommand's name as an argument.
func WithSubcommand(subc *command) commandOption {
	return func(c *command) {
		c.cs.set(subc.name, subc)
	}
}

// WithCommandDescription sets the description for a command.
// The description is typically used for displaying help or usage information.
// It should be a brief summary of what the command does.
func WithCommandDescription(desc string) commandOption {
	return func(c *command) {
		c.desc = desc
	}
}

func WithCommandDeprecated(d bool) commandOption {
	return func(c *command) {
		c.deprecated = d
	}
}

// WithFlags sets the provided flags as options for the command.
// This commandOption takes a variable number of flag pointers as input
// and adds each flag to the command's flagset.
// If a flag with the same name already exists in the flagset, it will be replaced.
func WithFlags(flags ...*flag) commandOption {
	return func(c *command) {
		for _, f := range flags {
			c.fsl.set(f.name, f)
			c.fss.set(f.short, f)
		}
	}
}

func WithMutualExclusiveFlags(flags ...*flag) commandOption {
	var cnt int
	var metFlags []string

	return func(c *command) {
		for _, f := range flags {
			c.fsl.set(f.name, f)
			c.fss.set(f.short, f)
			f.cvalidators = append(f.cvalidators,
				func(f Flag) error {
					if cnt > 0 {
						return fmt.Errorf("mutual exclusive flags: beside flags %s met '%s' flag", metFlags, f.Name())
					}
					metFlags = append(metFlags, f.Name())
					cnt++
					return nil
				},
			)
		}
	}

}

func WithAlwaysTogetherFlags(flags ...*flag) commandOption {
	var cnt int
	return func(c *command) {
		for i, f := range flags {
			i := i
			c.fsl.set(f.name, f)
			c.fss.set(f.short, f)
			f.cvalidators = append(f.cvalidators,
				func(f Flag) error {
					if i != cnt {
						return fmt.Errorf("together flags: met not all specified flag")
					}
					cnt++
					return nil
				},
			)
		}
	}

}

func WithAguments(args ...*argument) commandOption {
	return func(c *command) {
		for _, arg := range args {
			c.as = append(c.as, arg)
		}
	}
}

func WithAction(f func(Context) error) commandOption {
	return func(c *command) {
		c.action = f
	}
}

func (c *command) Run() error {
	p := newParser(os.Args[1:])

	r, err := p.parse()
	if err != nil {
		return err
	}

	currc := c
	for len(r.beforeFlags) > 0 {
		cname := r.beforeFlags[0]

		sc, ok := currc.cs.get(cname)
		if !ok {
			break
		}

		currc = sc
		r.beforeFlags = r.beforeFlags[1:]
	}

	if len(r.beforeFlags) > 0 && len(r.flagSet) > 0 {
		return fmt.Errorf(
			"invalid syntax: met positional args %s before flags %v", r.beforeFlags, r.flagSet,
		)
	}

	return currc.run(r)
}

func (c *command) prepare() error {
	return nil
}

func (c *command) init(r repository) error {
	var argi int

	for _, arg := range r.beforeFlags {
		if len(c.as) < argi+1 {
			break
		}

		c.as[argi].value = arg
		argi++
	}

	var (
		f  *flag
		ok bool
	)

	for fname, fentry := range r.flagSet {
		switch {
		case fentry.isLong:
			f, ok = c.fsl.get(fname)
		case !fentry.isLong:
			f, ok = c.fss.get(fname)
		}

		if !ok {
			return fmt.Errorf(
				"invalid flag: met unexpected flag '%s' for command '%s'", fname, c.name,
			)
		}

		f.value = fentry.value
	}

	for _, arg := range r.positionalArgs {
		if len(c.as) < argi+1 {
			break
		}

		c.as[argi].value = arg
		argi++
	}

	return nil
}

func (c *command) validate() error {
	for _, f := range c.fsl {
		if err := f.validate(); err != nil {
			return err
		}
	}

	for _, arg := range c.as {
		if err := arg.validate(); err != nil {
			return err
		}
	}

	return nil
}

func (c *command) run(r repository) error {
	if err := c.init(r); err != nil {
		return err
	}

	if err := c.validate(); err != nil {
		return err
	}

	if c.action == nil {
		return nil
	}

	return c.action(&context{fs: c.fsl, as: c.as})
}
