package yacli

import "fmt"

type flagOption func(*flag)

type Flag interface {
	Name() string
	Short() string
	Value() any
	Type() ytype
	Description() string
	Deprecated() bool
}

var _ Flag = (*flag)(nil)

type flag struct {
	name        string
	short       string
	description string
	deprecated  bool
	value       any
	ttype       ytype
	cvalidators []func(f Flag) error
}

// NewFlag creates and returns a new Flag instance with the provided name, short name, description, and type.
func NewFlag(name, short, description string, ttype ytype, opts ...flagOption) *flag {
	f := &flag{
		name:        name,
		short:       short,
		description: description,
		ttype:       ttype,
	}

	for _, opt := range opts {
		opt(f)
	}

	return f
}

func WithFlagDeprecated(d bool) flagOption {
	return func(f *flag) {
		f.deprecated = d
	}
}

func WithFlagValidator(v func(Flag) error) flagOption {
	return func(f *flag) {
		f.cvalidators = append(f.cvalidators, v)
	}
}

func (f *flag) Name() string {
	return f.name
}

func (f *flag) Short() string {
	return f.short
}

func (f *flag) Value() any {
	return f.value
}

func (f *flag) Type() ytype {
	return f.ttype
}

func (f *flag) Description() string {
	return f.description
}

func (f *flag) Deprecated() bool {
	return f.deprecated
}

func (f *flag) String() string {
	return fmt.Sprintf("%s %s", f.Name(), f.Type())
}

func (f *flag) validate() error {
	if f.value == nil {
		return nil
	}

	v, err := vfuncs[f.ttype](f.value)
	if err != nil {
		return err
	}
	f.value = v

	for _, cvalidator := range f.cvalidators {
		if err := cvalidator(f); err != nil {
			return err
		}
	}

	return nil
}
