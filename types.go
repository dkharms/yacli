package yacli

import (
	"reflect"
	"strconv"

	"golang.org/x/exp/constraints"
)

type ytype int

func (y ytype) String() string {
	switch y {
	case Integer, Integer8, Integer16, Integer32, Integer64:
		return "INTEGER"
	case Float32, Float64:
		return "FLOAT"
	case String:
		return "STRING"
	case Bool:
		return "BOOL"
	default:
		return "UNKNOWN"
	}
}

const (
	Integer ytype = iota
	Integer8
	Integer16
	Integer32
	Integer64
	Float32
	Float64
	String
	Bool
)

type vfunc func(v any) (any, error)

var vfuncs map[ytype]vfunc = map[ytype]vfunc{
	Integer:   func(v any) (any, error) { return validateInteger[int](v) },
	Integer8:  func(v any) (any, error) { return validateInteger[int8](v) },
	Integer16: func(v any) (any, error) { return validateInteger[int16](v) },
	Integer32: func(v any) (any, error) { return validateInteger[int32](v) },
	Integer64: func(v any) (any, error) { return validateInteger[int64](v) },
	Float32:   func(v any) (any, error) { return validateFloat[float32](v) },
	Float64:   func(v any) (any, error) { return validateFloat[float64](v) },
	String:    func(v any) (any, error) { return v.(string), nil },
	Bool:      func(v any) (any, error) { return validateBool(v) },
}

func validateInteger[T constraints.Integer](v any) (T, error) {
	var t T
	i, err := strconv.ParseInt(v.(string), 10, reflect.TypeOf(t).Bits())
	if err != nil {
		return t, err
	}
	return T(i), err
}

func validateFloat[T constraints.Float](v any) (T, error) {
	var t T
	i, err := strconv.ParseFloat(v.(string), reflect.TypeOf(t).Bits())
	if err != nil {
		return t, err
	}
	return T(i), err
}

func validateBool(v any) (bool, error) {
	if v.(string) == "" {
		return true, nil
	}
	return strconv.ParseBool(v.(string))
}
