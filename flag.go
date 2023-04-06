package yacli

type flagOption func(*flag)

type Flag interface {
	Name() string
	Value() any
	Type() ytype
}

var _ (Flag) = (*flag)(nil)

type flag struct {
	name        string
	short       string
	description string
	deprecated  bool
	value       any
	ttype       ytype
	cvalidators []func(f Flag) error
}

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

func WithFlagValidator(v func(Flag) error) flagOption {
	return func(f *flag) {
		f.cvalidators = append(f.cvalidators, v)
	}
}

func (f *flag) Name() string {
	return f.name
}

func (f *flag) Value() any {
	return f.value
}

func (f *flag) Type() ytype {
	return f.ttype
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