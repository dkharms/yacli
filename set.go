package yacli

// commandset is a set of commands mapped by their name.
type commandset map[string]*command

// set is a method of the commandset struct that takes two arguments:
// a string representing the name of the command to add, and a pointer to a command struct.
// It returns a boolean value indicating whether the command was successfully added to the command set.
// If a command with the same name already exists in the command set, the function will return false.
func (cs commandset) set(name string, c *command) bool {
	if _, ok := cs[name]; ok {
		return false
	}
	cs[name] = c
	return true
}

// get returns the command associated with the given name and a bool indicating
// whether the command was found or not. Returns nil and false if the command is
// not present in the command set.
func (cs commandset) get(name string) (*command, bool) {
	c, ok := cs[name]
	return c, ok
}

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
	f, ok := fs.get(name)
	if !ok {
		return 0, false
	}
	v, ok := f.value.(int)
	return v, ok
}

func (fs flagset) Integer8(name string) (int8, bool) {
	f, ok := fs.get(name)
	if !ok {
		return 0, false
	}
	v, ok := f.value.(int8)
	return v, ok
}

func (fs flagset) Integer16(name string) (int16, bool) {
	f, ok := fs.get(name)
	if !ok {
		return 0, false
	}
	v, ok := f.value.(int16)
	return v, ok
}

func (fs flagset) Integer32(name string) (int32, bool) {
	f, ok := fs.get(name)
	if !ok {
		return 0, false
	}
	v, ok := f.value.(int32)
	return v, ok
}

func (fs flagset) Integer64(name string) (int64, bool) {
	f, ok := fs.get(name)
	if !ok {
		return 0, false
	}
	v, ok := f.value.(int64)
	return v, ok
}

// Float32 retrieves the value of a float32 flag.
// If the flag is not found, the second return value is false.
func (fs flagset) Float32(name string) (float32, bool) {
	f, ok := fs.get(name)
	if !ok {
		return 0, false
	}
	v, ok := f.value.(float32)
	return v, ok
}

// Float64 retrieves the value of a float64 flag.
// If the flag is not found, the second return value is false.
func (fs flagset) Float64(name string) (float64, bool) {
	f, ok := fs.get(name)
	if !ok {
		return 0, false
	}
	v, ok := f.value.(float64)
	return v, ok
}

// String retrieves the value of a string flag.
// If the flag is not found, the second return value is false.
func (fs flagset) String(name string) (string, bool) {
	f, ok := fs.get(name)
	if !ok {
		return "", false
	}
	v, ok := f.value.(string)
	return v, ok
}

// Bool retrieves the value of a boolean flag.
// If the flag is not found, the second return value is false.
func (fs flagset) Bool(name string) (bool, bool) {
	f, ok := fs.get(name)
	if !ok {
		return false, false
	}
	v, ok := f.value.(bool)
	return v, ok
}

// argset is a slice of pointers to argument objects.
type argset []*argument

// add adds an argument to the argset.
func (as argset) add(arg *argument) {
	as = append(as, arg)
}

// get returns the argument with the given name, if it exists in the argset.
func (as argset) get(name string) *argument {
	for _, arg := range as {
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

// Integer8 returns the value of the argument with the given name as an int8.
// Panics if the argument was not found or have different type.
func (as argset) Integer8(name string) int8 {
	return as.get(name).value.(int8)
}

// Integer16 returns the value of the argument with the given name as an int16.
// Panics if the argument was not found or have different type.
func (as argset) Integer16(name string) int16 {
	return as.get(name).value.(int16)
}

// Integer32 returns the value of the argument with the given name as an int32.
// Panics if the argument was not found or have different type.
func (as argset) Integer32(name string) int32 {
	return as.get(name).value.(int32)
}

// Integer64 returns the value of the argument with the given name as an int64.
// Panics if the argument was not found or have different type.
func (as argset) Integer64(name string) int64 {
	return as.get(name).value.(int64)
}

// Float32 returns the value of the argument with the given name as an float32.
// Panics if the argument was not found or have different type.
func (as argset) Float32(name string) float32 {
	return as.get(name).value.(float32)
}

// Float64 returns the value of the argument with the given name as an float64.
// Panics if the argument was not found or have different type.
func (as argset) Float64(name string) float64 {
	return as.get(name).value.(float64)
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
