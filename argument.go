package yacli

import "fmt"

type argumentOption func(*argument)

type Argument interface {
	Name() string
	Value() any
	Type() ytype
}

var _ (Argument) = (*argument)(nil)

type argument struct {
	name        string
	description string
	ttype       ytype
	value       any
	cvalidators []func(Argument) error
}

func NewArgument(name, description string, ttype ytype, opts ...argumentOption) *argument {
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

func WithArgumentValidator(v func(Argument) error) argumentOption {
	return func(a *argument) {
		a.cvalidators = append(a.cvalidators, v)
	}
}

func (a *argument) Name() string {
	return a.name
}

func (a *argument) Type() ytype {
	return a.ttype
}

func (a *argument) Value() any {
	return a.value
}

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
