package yacli

import "strconv"

type ytype int

const (
	Integer ytype = iota
	Float
	String
	Bool
)

type vfunc func(v any) (any, error)

var vfuncs map[ytype]vfunc = map[ytype]vfunc{
	Integer: func(v any) (any, error) { return strconv.ParseInt(v.(string), 10, 32) },
	Float:   func(v any) (any, error) { return strconv.ParseFloat(v.(string), 32) },
	String:  func(v any) (any, error) { return v.(string), nil },
	Bool:    func(v any) (any, error) { return strconv.ParseBool(v.(string)) },
}
