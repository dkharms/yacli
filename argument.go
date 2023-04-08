package yacli

import "fmt"

type argumentOption func(*argument)

// Argument represents a command argument.
type Argument interface {
	// Name returns the name of the argument.
	Name() string

	// Value returns the value of the argument.
	Value() any

	// Type returns the type of the argument.
	Type() ytype
}

var _ (Argument) = (*argument)(nil)

// argument represents a command-line argument with a name, description,
// type, value, and custom validators.
type argument struct {
	name        string
	description string
	ttype       ytype
	value       any
	cvalidators []func(Argument) error
}

// NewArgument creates a new argument with the given name, description,
// and type, and applies the specified options to it.
func NewArgument(
	name, description string, ttype ytype, opts ...argumentOption,
) *argument {
	a := &argument{
		name:        name,
		description: description,
		ttype:       ttype,
	}

	for _, opt := range opts {
		opt(a)
	}

	return a
}

// WithArgumentValidator is an argumentOption function that allows adding custom validators to an argument.
// It takes a function with an Argument parameter that returns an error, and appends it to the argument's
// list of custom validators.
func WithArgumentValidator(v func(Argument) error) argumentOption {
	return func(a *argument) {
		a.cvalidators = append(a.cvalidators, v)
	}
}

// Name returns the name of the argument.
func (a *argument) Name() string {
	return a.name
}

// Type returns the type of the argument.
func (a *argument) Type() ytype {
	return a.ttype
}

// Value returns the value of the argument.
func (a *argument) Value() any {
	return a.value
}

// validate validates the argument value against its type and custom validators.
// It returns an error if the validation fails.
func (a *argument) validate() error {
	if a.value == nil {
		return fmt.Errorf("invalid argument: missing value for argument '%s'", a.name)
	}

	v, err := vfuncs[a.ttype](a.value)
	if err != nil {
		return err
	}
	a.value = v

	for _, cvalidator := range a.cvalidators {
		if err := cvalidator(a); err != nil {
			return err
		}
	}

	return nil
}
