package yacli

// Context is an interface that defines the methods for accessing command-line arguments
// and flags that were parsed by a command-line parser.
type Context interface {
	// Flags returns a flagset object that contains
	// all the flag values that were parsed.
	Flags() flagset

	// Arguments returns an argset object that contains
	// all the positional arguments that were parsed.
	Arguments() argset
}

var _ Context = (*context)(nil)

// context represents a parsed command-line context, containing the
// parsed flags and arguments.
type context struct {
	// fs is a flagset containing all the parsed flags.
	fs flagset

	// as is an argset containing all the parsed arguments.
	as argset
}

// Flags returns the flagset associated with the context.
func (c *context) Flags() flagset {
	return c.fs
}

// Arguments returns the argument set associated with this context.
func (c *context) Arguments() argset {
	return c.as
}
