package yacli

// flagset is a map of string flag names to pointers to flag objects.
// It is used to store and retrieve flag values by name.
type flagset map[string]*flag

// set adds the given flag to the flagset, if the flag with the same name doesn't already exist.
// If the flag with the same name already exists, set returns false, otherwise it returns true.
func (fs flagset) set(name string, f *flag) bool {
	if _, ok := fs[name]; ok {
		return false
	}
	fs[name] = f
	return true
}

// get retrieves the flag with the given name from the flagset, if it exists.
// It returns a pointer to the flag and a boolean indicating whether the flag was found.
func (fs flagset) get(name string) (*flag, bool) {
	f, ok := fs[name]
	return f, ok
}

// Integer retrieves the value of an integer flag.
// If the flag is not found, the second return value is false.
func (fs flagset) Integer(name string) (int, bool) {
	v, ok := fs[name].value.(int)
	return v, ok
}

// Float retrieves the value of a float flag.
// If the flag is not found, the second return value is false.
func (fs flagset) Float(name string) (float32, bool) {
	v, ok := fs[name].value.(float32)
	return v, ok
}

// String retrieves the value of a string flag.
// If the flag is not found, the second return value is false.
func (fs flagset) String(name string) (string, bool) {
	v, ok := fs[name].value.(string)
	return v, ok
}

// Bool retrieves the value of a boolean flag.
// If the flag is not found, the second return value is false.
func (fs flagset) Bool(name string) (bool, bool) {
	v, ok := fs[name].value.(bool)
	return v, ok
}

// argset is a slice of pointers to argument objects.
type argset []*argument

// add adds an argument to the argset.
func (a argset) add(arg *argument) {
	a = append(a, arg)
}

// get returns the argument with the given name, if it exists in the argset.
func (a argset) get(name string) *argument {
	for _, arg := range a {
		if name == arg.name {
			return arg
		}
	}
	return nil
}

// Integer returns the value of the argument with the given name as an int.
// Panics if the argument was not found or have different type.
func (as argset) Integer(name string) int {
	return as.get(name).value.(int)
}

// Float returns the value of the argument with the given name as an float32.
// Panics if the argument was not found or have different type.
func (as argset) Float(name string) float32 {
	return as.get(name).value.(float32)
}

// String returns the value of the argument with the given name as an string.
// Panics if the argument was not found or have different type.
func (as argset) String(name string) string {
	return as.get(name).value.(string)
}

// Bool returns the value of the argument with the given name as an bool.
// Panics if the argument was not found or have different type.
func (as argset) Bool(name string) bool {
	return as.get(name).value.(bool)
}
