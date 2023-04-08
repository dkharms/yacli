package yacli

import (
	"reflect"
	"testing"
)

func TestParser_long(t *testing.T) {
	testCase := []struct {
		name     string
		flags    []string
		expected repository
		hasErr   bool
	}{
		{
			name:  "--key value",
			flags: []string{"--key", "value"},
			expected: repository{
				flagSet: map[string]entry{"key": {"value", true}},
			},
		},
		{
			name:  "--key=value",
			flags: []string{"--key=value"},
			expected: repository{
				flagSet: map[string]entry{"key": {"value", true}},
			},
		},
		{
			name:  "--key='value'",
			flags: []string{"--key='value'"},
			expected: repository{
				flagSet: map[string]entry{"key": {"'value'", true}},
			},
		},
		{
			name:  "--key=",
			flags: []string{"--key="},
			expected: repository{
				flagSet: map[string]entry{"key": {"", true}},
			},
		},
		{
			name:  "--key ikey=ivalue",
			flags: []string{"--key", "ikey=ivalue"},
			expected: repository{
				flagSet: map[string]entry{"key": {"ikey=ivalue", true}},
			},
		},
		{
			name:  "--key=ikey=ivalue",
			flags: []string{"--key=ikey=ivalue"},
			expected: repository{
				flagSet: map[string]entry{"key": {"ikey=ivalue", true}},
			},
		},
		{
			name:   "--=ikey=ivalue",
			flags:  []string{"--=ikey=ivalue"},
			hasErr: true,
		},
		{
			name:   "-- ikey=ivalue",
			flags:  []string{"--", "ikey=ivalue"},
			hasErr: true,
		},
		{
			name:  "--akey --bkey bvalue",
			flags: []string{"--akey", "--bkey", "value"},
			expected: repository{
				flagSet: map[string]entry{"akey": {"", true}, "bkey": {"value", true}},
			},
		},
	}

	for _, tt := range testCase {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			p := newParser(tt.flags)

			r, err := p.parse()
			if err != nil && !tt.hasErr {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(tt.expected.flagSet, r.flagSet) {
				t.Fatalf("expected %v, got %v", tt.expected.flagSet, r.flagSet)
			}
		})
	}
}

func TestParser_short(t *testing.T) {
	testCase := []struct {
		name     string
		flags    []string
		expected repository
		hasErr   bool
	}{
		{
			name:  "-k value",
			flags: []string{"-k", "value"},
			expected: repository{
				flagSet: map[string]entry{"k": {"value", false}},
			},
		},
		{
			name:  "-abc",
			flags: []string{"-abc"},
			expected: repository{
				flagSet: map[string]entry{"a": {"", false}, "b": {"", false}, "c": {"", false}},
			},
		},
		{
			name:  "-a -b -c",
			flags: []string{"-a", "-b", "-c"},
			expected: repository{
				flagSet: map[string]entry{"a": {"", false}, "b": {"", false}, "c": {"", false}},
			},
		},
		{
			name:  "-a -b -c value",
			flags: []string{"-a", "-b", "-c", "value"},
			expected: repository{
				flagSet: map[string]entry{"a": {"", false}, "b": {"", false}, "c": {"value", false}},
			},
		},
		{
			name:  "-a value",
			flags: []string{"-a", "-b", "-c", "value"},
			expected: repository{
				flagSet: map[string]entry{"a": {"", false}, "b": {"", false}, "c": {"value", false}},
			},
		},
	}

	for _, tt := range testCase {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			p := newParser(tt.flags)

			r, err := p.parse()
			if err != nil && !tt.hasErr {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(tt.expected.flagSet, r.flagSet) {
				t.Fatalf("expected %v, got %v", tt.expected.flagSet, r.flagSet)
			}
		})
	}
}

func TestParser_mixed(t *testing.T) {
	testCase := []struct {
		name     string
		args     []string
		expected repository
		hasErr   bool
	}{
		{
			name: "command subcommand argument",
			args: []string{"command", "subcommand", "argument"},
			expected: repository{
				beforeFlags:    []string{"command", "subcommand", "argument"},
				flagSet:        map[string]entry{},
				positionalArgs: []string{},
			},
		},
		{
			name: "command subcommand --flag flag-value argument",
			args: []string{"command", "subcommand", "--flag", "flag-value", "argument"},
			expected: repository{
				beforeFlags:    []string{"command", "subcommand"},
				flagSet:        map[string]entry{"flag": {"flag-value", true}},
				positionalArgs: []string{"argument"},
			},
		},
		{
			name: "command subcommand --flag flag-value --another-flag=argument",
			args: []string{"command", "subcommand", "--flag", "flag-value", "--another-flag=argument"},
			expected: repository{
				beforeFlags:    []string{"command", "subcommand"},
				flagSet:        map[string]entry{"flag": {"flag-value", true}, "another-flag": {"argument", true}},
				positionalArgs: []string{},
			},
		},
		{
			name:   "command subcommand --flag flag-value incorrect-argument --another-flag=argument",
			args:   []string{"command", "subcommand", "--flag", "flag-value", "incorrect-argument", "--another-flag=argument"},
			hasErr: true,
		},
	}

	for _, tt := range testCase {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			p := newParser(tt.args)

			r, err := p.parse()
			if err != nil && !tt.hasErr {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(tt.expected.beforeFlags, r.beforeFlags) {
				t.Fatalf("before flags: expected %v, got %v", tt.expected.beforeFlags, r.beforeFlags)
			}

			if !reflect.DeepEqual(tt.expected.flagSet, r.flagSet) {
				t.Fatalf("flag set: expected %v, got %v", tt.expected.flagSet, r.flagSet)
			}

			if !reflect.DeepEqual(tt.expected.positionalArgs, r.positionalArgs) {
				t.Fatalf("positional args: expected %v, got %v", tt.expected.positionalArgs, r.positionalArgs)
			}
		})
	}
}