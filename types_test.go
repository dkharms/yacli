package yacli

import (
	"reflect"
	"testing"
)

func TestValidateInteger(t *testing.T) {
	tests := []struct {
		name   string
		v      any
		want   any
		hasErr bool
	}{
		{"valid positive int", "123", int(123), false},
		{"valid negative int", "-123", int(-123), false},
		{"invalid int", "123a", 0, true},
		{"valid positive int8", "123", int8(123), false},
		{"valid negative int8", "-123", int8(-123), false},
		{"invalid int8", "123a", 0, true},
		{"valid positive int16", "123", int16(123), false},
		{"valid negative int16", "-123", int16(-123), false},
		{"invalid int16", "123a", 0, true},
		{"valid positive int32", "123", int32(123), false},
		{"valid negative int32", "-123", int32(-123), false},
		{"invalid int32", "123a", 0, true},
		{"valid positive int64", "123", int64(123), false},
		{"valid negative int64", "-123", int64(-123), false},
		{"invalid int64", "123a", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				got any
				err error
			)

			switch reflect.TypeOf(tt.want).Kind() {
			case reflect.Int:
				got, err = validateInteger[int](tt.v)
			case reflect.Int8:
				got, err = validateInteger[int8](tt.v)
			case reflect.Int16:
				got, err = validateInteger[int16](tt.v)
			case reflect.Int32:
				got, err = validateInteger[int32](tt.v)
			case reflect.Int64:
				got, err = validateInteger[int64](tt.v)
			}

			if (err != nil) != tt.hasErr {
				t.Errorf(
					"validateInteger('%v')=%v, error= %v, wantErr=%v",
					tt.v, got, err, tt.hasErr,
				)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("validateInteger('%v')=%v, want=%v", tt.v, got, tt.want)
			}
		})
	}
}

func TestValidateFloat(t *testing.T) {
	tests := []struct {
		name   string
		v      any
		want   any
		hasErr bool
	}{
		{"valid positive float32", "123.45", float32(123.45), false},
		{"valid negative float32", "-123.45", float32(-123.45), false},
		{"invalid float32", "123.45o", float32(0), true},
		{"valid positive float64", "123.45", float64(123.45), false},
		{"valid negative float64", "-123.45", float64(-123.45), false},
		{"invalid float64", "123.45o", float64(0), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				got any
				err error
			)

			switch reflect.TypeOf(tt.want).Kind() {
			case reflect.Float32:
				got, err = validateFloat[float32](tt.v)
			case reflect.Float64:
				got, err = validateFloat[float64](tt.v)
			}

			if (err != nil) != tt.hasErr {
				t.Errorf(
					"validateFloat('%v')=%v error=%v, wantErr %v",
					tt.v, got, err, tt.hasErr,
				)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("validateFloat('%v')=%v, want %v", tt.v, got, tt.want)
			}
		})
	}
}
