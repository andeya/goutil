package status

// Throw creates a status with stack, and panic.
func Throw(code int32, msg string, cause interface{}) {
	panic(NewWithStack(code, msg, cause))
}

// Panic panic with stack trace.
func Panic(stat *Status) {
	if stat == nil {
		stat = &Status{
			stack: callers(),
		}
	} else if stat.stack == nil {
		stat.stack = callers()
	}
	panic(stat)
}

// Check if err!=nil, create a status with stack, and panic.
func Check(err error, code int32, msg string) {
	if err == nil {
		return
	}
	panic(NewWithStack(code, msg, err))
}

// Catch recovers the panic and returns status.
// Example:
//  var stat *Status
//  defer Catch(&stat)
func Catch(statPtr **Status) {
	r := recover()
	if r == nil || statPtr == nil {
		return
	}
	switch v := r.(type) {
	case *Status:
		if v == nil {
			v = new(Status)
		}
		*statPtr = v
	case Status:
		*statPtr = &v
	default:
		*statPtr = New(UnknownError, "", v)
	}
}
