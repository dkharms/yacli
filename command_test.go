package yacli

import "testing"

func TestCommand_run(t *testing.T) {
	_ = []struct {
		name string
		r    repository
		c    command
	}{
		{
			name: "echo",
			r:    repository{},
			c:    command{},
		},
	}
}
