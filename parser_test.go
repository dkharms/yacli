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
				mf: map[string]string{"key": "value"},
			},
		},
		{
			name:  "--key=value",
			flags: []string{"--key=value"},
			expected: repository{
				mf: map[string]string{"key": "value"},
			},
		},
		{
			name:  "--key='value'",
			flags: []string{"--key='value'"},
			expected: repository{
				mf: map[string]string{"key": "'value'"},
			},
		},
		{
			name:  "--key=",
			flags: []string{"--key="},
			expected: repository{
				mf: map[string]string{"key": ""},
			},
		},
		{
			name:  "--key ikey=ivalue",
			flags: []string{"--key", "ikey=ivalue"},
			expected: repository{
				mf: map[string]string{"key": "ikey=ivalue"},
			},
		},
		{
			name:  "--key=ikey=ivalue",
			flags: []string{"--key=ikey=ivalue"},
			expected: repository{
				mf: map[string]string{"key": "ikey=ivalue"},
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
				mf: map[string]string{"akey": "", "bkey": "value"},
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

			if !reflect.DeepEqual(tt.expected.mf, r.mf) {
				t.Fatalf("expected %v, got %v", tt.expected.mf, r.mf)
			}
		})
	}
}