package yacli

import "testing"

func TestFlagSet(t *testing.T) {
	fs := make(flagset)

	// Table tests
	tests := []struct {
		name  string
		value any
	}{
		{"int", 42},
		{"float", float32(3.14)},
		{"string", "hello world"},
		{"bool", true},
	}

	// Set flags
	for _, tt := range tests {
		f := &flag{value: tt.value}
		if !fs.set(tt.name, f) {
			t.Errorf("failed to set flag %s", tt.name)
		}
	}

	// Test Integer
	if v, ok := fs.Integer("int"); !ok || v != 42 {
		t.Errorf("expected Integer('int') to return 42, got %v (ok=%v)", v, ok)
	}
	if _, ok := fs.Integer("missing"); ok {
		t.Error("expected Integer('missing') to return false")
	}

	// Test Float32
	if v, ok := fs.Float32("float"); !ok || v != 3.14 {
		t.Errorf("expected Float32('float') to return 3.14, got %v (ok=%v)", v, ok)
	}
	if _, ok := fs.Float32("missing"); ok {
		t.Error("expected Float32('missing') to return false")
	}

	// Test String
	if v, ok := fs.String("string"); !ok || v != "hello world" {
		t.Errorf("expected String('string') to return 'hello world', got %v (ok=%v)", v, ok)
	}
	if _, ok := fs.String("missing"); ok {
		t.Error("expected String('missing') to return false")
	}

	// Test Bool
	if v, ok := fs.Bool("bool"); !ok || !v {
		t.Errorf("expected Bool('bool') to return true, got %v (ok=%v)", v, ok)
	}
	if _, ok := fs.Bool("missing"); ok {
		t.Error("expected Bool('missing') to return false")
	}
}
