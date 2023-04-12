package yacli

import (
	"fmt"
	"os"
	"testing"
)

func newCommand() *command {
	return NewRootCommand(
		WithFlags(
			NewFlag("amount", "n", "…", Integer),
		),
		WithMutualExclusiveFlags(
			NewFlag("lowercase", "l", "…", Bool),
			NewFlag("uppercase", "u", "…", Bool),
		),
		WithAlwaysTogetherFlags(
			NewFlag("separator", "s", "…", String),
			NewFlag("separator-amount", "a", "…", Integer),
		),
	)
}

func TestCommand_run(t *testing.T) {
	testCases := []struct {
		name   string
		args   []string
		c      *command
		f      []func(Context) error
		hasErr bool
	}{
		{
			name: "echo --amount 10",
			args: []string{"echo", "--amount", "10"},
			c:    newCommand(),
			f: []func(Context) error{
				func(ctx Context) error {
					n, isSet := ctx.Flags().Integer("amount")
					if !isSet {
						return fmt.Errorf("amount MUST be set")
					}
					if n != 10 {
						return fmt.Errorf("amount expected 10, got %d", n)
					}
					return nil
				},
			},
		},
		{
			name: "echo -n 10",
			args: []string{"echo", "-n", "10"},
			c:    newCommand(),
			f: []func(Context) error{
				func(ctx Context) error {
					n, isSet := ctx.Flags().Integer("amount")
					if !isSet {
						return fmt.Errorf("amount MUST be set")
					}
					if n != 10 {
						return fmt.Errorf("amount expected 10, got %d", n)
					}
					return nil
				},
			},
		},
		{
			name:   "echo --amount invalid-amount",
			args:   []string{"echo", "--amount", "invalid-amount"},
			c:      newCommand(),
			hasErr: true,
			f:      []func(Context) error{nil},
		},
		{
			name: "echo --amount 10 --lowercase true",
			args: []string{"echo", "--amount", "10", "--lowercase", "true"},
			c:    newCommand(),
			f: []func(Context) error{
				func(ctx Context) error {
					n, isSet := ctx.Flags().Integer("amount")
					if !isSet {
						return fmt.Errorf("amount MUST be set")
					}
					if n != 10 {
						return fmt.Errorf("amount expected 10, got %d", n)
					}
					if v, isSet := ctx.Flags().Bool("lowercase"); !isSet || !v {
						return fmt.Errorf("lowercase flag is not set or with incorrect value")
					}
					return nil
				},
			},
		},
		{
			name:   "echo --amount 10 --lowercase true --uppercase true",
			args:   []string{"echo", "--amount", "10", "--lowercase", "true", "--uppercase", "true"},
			c:      newCommand(),
			hasErr: true,
			f:      []func(Context) error{nil},
		},
		{
			name: "echo --separator smth --separator-amount 3",
			args: []string{"echo", "--separator", "smth", "--separator-amount", "3"},
			c:    newCommand(),
			f:    []func(Context) error{nil},
		},
		{
			name:   "echo --separator smth",
			args:   []string{"echo", "--separator", "smth"},
			c:      newCommand(),
			hasErr: true,
			f:      []func(Context) error{nil},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.args

			for _, f := range tt.f {
				tt.c.action = f
				err := tt.c.Run()

				if tt.hasErr && err == nil {
					t.Errorf("expected error, got %v", err)
				}

				if !tt.hasErr && err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			}
		})
	}
}
