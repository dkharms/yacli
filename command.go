package yacli

import (
	"fmt"
	"os"
	"strings"
)

type commandOption func(*command)

// Command represents a command definition which has a name,
// a description, flags, arguments, and an action.
type Command interface {
	// Name returns the name of the command.
	Name() string

	// Description returns a brief description of the command.
	Description() string

	// Deprecated returns a boolean indicating whether or not the command is deprecated.
	Deprecated() bool

	// Subcommands returns a list of subcommands that are associated with this command.
	Subcommands() []Command

	// Flags returns a list of flags of command.
	Flags() []Flag

	// Arguments returns a list of arguments that were passed to command.
	Arguments() []Argument

	// Usage returns a string describing how the command should be used in POSIX format.
	Usage() string

	// Help returns a string containing a more detailed explanation of how to use the command.
	Help() string
}

var _ Command = (*command)(nil)

// command represents a command definition which has a name,
// a description, flags, arguments, and an action.
type command struct {
	// name is the name of the command.
	name string

	// description is a brief description of the command.
	description string

	// deprecated is a flag that indicates if the command is deprecated or not
	deprecated bool

	// cs is the sub-commands under this command.
	cs commandset

	// fsl is the flags associated with this command.
	fsl flagset

	// fss is the flags associated with this command.
	fss flagset

	// fg is the flag groups associated with this command.
	fg flaggroup

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
//
// The command is initialized with an empty flagset and commandset.
// It then applies the given options to the command.
// Returns a pointer to the created command.
func NewCommand(name string, opts ...commandOption) *command {
	hFlag := &flag{name: "help", short: "h", description: "Print this message", ttype: Bool}

	c := &command{
		name: name,
		fsl:  flagset{"help": hFlag},
		fss:  flagset{"h": hFlag},
		fg:   make(flaggroup),
		cs:   make(commandset),
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
		if _, ok := c.cs.get(subc.name); ok {
			panic(fmt.Errorf(
				"invalid command: subcommand with name '%s' is already defined on command '%s'",
				subc.name, c.name,
			))
		}
		c.cs.set(subc.name, subc)
	}
}

// WithCommandDescription sets the description for a command.
//
// The description is typically used for displaying help or usage information.
// It should be a brief summary of what the command does.
func WithCommandDescription(desc string) commandOption {
	return func(c *command) {
		c.description = desc
	}
}

func WithCommandDeprecated(d bool) commandOption {
	return func(c *command) {
		c.deprecated = d
	}
}

// WithFlags sets the provided flags as options for the command.
//
// This commandOption takes a variable number of flag pointers as input
// and adds each flag to the command's flagset.
// If a flag with the same name already exists in the flagset, it will be replaced.
func WithFlags(flags ...*flag) commandOption {
	return func(c *command) {
		g := c.fg.new(groupDefault)
		for _, f := range flags {
			panicIfFlagAlreadyDefined(c, f)
			c.fsl.set(f.name, f)
			c.fss.set(f.short, f)
			g.add(f)
		}
	}
}

// WithMutualExclusiveFlags is a commandOption that creates a group
// of flags that are mutually exclusive.

// This function takes a variable number of flags as its arguments.
// The flags passed to this function will be added to the group,
// and an error will be thrown if more than one flag in the group is set.
func WithMutualExclusiveFlags(flags ...*flag) commandOption {
	return func(c *command) {
		g := c.fg.new(groupMutex)
		for _, f := range flags {
			panicIfFlagAlreadyDefined(c, f)
			c.fsl.set(f.name, f)
			c.fss.set(f.short, f)
			g.add(f)
			f.cvalidators = append(f.cvalidators,
				func(f Flag) error {
					g.met++
					if g.met > 1 {
						var names []string
						for _, ff := range g.flags {
							names = append(names, ff.Name())
						}
						return fmt.Errorf(
							"invalid flags: flags %s are mutual exclusive flags", flags,
						)
					}
					return nil
				},
			)
		}
	}

}

// WithAlwaysTogetherFlags is a commandOption that creates a group
// of flags that must always be set together.
//
// This function takes a variable number of flags as its arguments.
// The flags passed to this function will be added to the group,
// and an error will not be thrown even if one of the flags is not set.
func WithAlwaysTogetherFlags(flags ...*flag) commandOption {
	return func(c *command) {
		g := c.fg.new(groupTogether)
		for _, f := range flags {
			panicIfFlagAlreadyDefined(c, f)
			c.fsl.set(f.name, f)
			c.fss.set(f.short, f)
			g.add(f)
			f.cvalidators = append(f.cvalidators,
				func(_ Flag) error {
					g.met++
					return nil
				},
			)
		}
	}

}

// WithArguments is a commandOption that adds one or more arguments to a command object.
//
// This function takes a variable number of argument pointers as its arguments.
// Each argument passed to this function will be added to the command object's argument slice.
// If an argument with the same name has already been defined, an error will be thrown.
func WithArguments(args ...*argument) commandOption {
	return func(c *command) {
		for _, arg := range args {
			panicIfArgumentAlreadyDefined(c, arg)
			c.as = append(c.as, arg)
		}
	}
}

// WithAction is a commandOption that sets the action function for a command object.
//
// This function takes a single argument, a function that accepts a Context and returns an error.
func WithAction(f func(Context) error) commandOption {
	return func(c *command) {
		c.action = f
	}
}

// Name method returns the name of the command.
func (c *command) Name() string {
	return c.name
}

// Description method returns the description of the command.
func (c *command) Description() string {
	return c.description
}

// Deprecated method returns a boolean indicating whether the command is deprecated or not.
func (c *command) Deprecated() bool {
	return c.deprecated
}

// Subcommands method returns a slice of Command objects representing the subcommands of this command.
func (c *command) Subcommands() []Command {
	var subcommands []Command
	for _, subcommand := range c.cs {
		subcommands = append(subcommands, subcommand)
	}
	return subcommands
}

// Flags method returns a slice of Flag objects representing the flags of this command.
func (c *command) Flags() []Flag {
	var flags []Flag
	for _, flag := range c.fsl {
		flags = append(flags, flag)
	}
	return flags
}

// Arguments method returns a slice of Argument objects representing the arguments of this command.
func (c *command) Arguments() []Argument {
	var arguments []Argument
	for _, argument := range c.as {
		arguments = append(arguments, argument)
	}
	return arguments
}

// Help method returns a string representing the help information for the command.
func (c *command) Help() string {
	var s strings.Builder
	if err := helpTemplate.Execute(&s, c); err != nil {
		panic(err)
	}
	return s.String()
}

// The Run method is responsible for parsing the command line arguments and executing the command.
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

	for _, g := range c.fg {
		if g.ttype == groupTogether && 0 < g.met && g.met < len(g.flags) {
			var flags []string
			for _, f := range g.flags {
				flags = append(flags, f.name)
			}
			return fmt.Errorf(
				"invalid flags: you have to pass flags %s together", flags,
			)
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

	if c.help() {
		return nil
	}

	if err := c.validate(); err != nil {
		return err
	}

	if c.action == nil {
		return nil
	}

	return c.action(&context{fs: c.fsl, as: c.as})
}

func (c *command) help() bool {
	if v, _ := c.fsl.get("help"); v.value != nil {
		fmt.Print(c.Help())
		return true
	}
	return false
}

func panicIfFlagAlreadyDefined(c *command, f *flag) {
	if _, ok := c.fsl.get(f.name); ok {
		panic(fmt.Errorf(
			"invalid command: long flag '%s' is alredy defined for command '%s'",
			f.name, c.name,
		))
	}

	if _, ok := c.fss.get(f.short); ok {
		panic(fmt.Errorf(
			"invalid command: short flag '%s' is alredy defined for command '%s'",
			f.short, c.name,
		))
	}
}

func panicIfArgumentAlreadyDefined(c *command, arg *argument) {
	if v := c.as.get(arg.name); v != nil {
		panic(fmt.Errorf(
			"invalid command: argument '%s' is already defined for command '%s'",
			arg.name, c.name,
		))
	}
}
