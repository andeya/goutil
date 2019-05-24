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
// NOTE:
//  Set "realStat" to true if a "state" type is recovered
// Example:
//  var stat *Status
//  defer Catch(&stat)
func Catch(statPtr **Status, realStat ...*bool) {
	r := recover()

	if statPtr == nil {
		switch r.(type) {
		case *Status, Status:
			trySetBool(realStat, true)
		default:
			trySetBool(realStat, false)
		}
	}

	switch v := r.(type) {
	case nil:
		trySetBool(realStat, false)
		*statPtr = new(Status)
	case *Status:
		trySetBool(realStat, true)
		if v == nil {
			v = new(Status)
		}
		*statPtr = v
	case Status:
		trySetBool(realStat, true)
		*statPtr = &v
	default:
		trySetBool(realStat, false)
		*statPtr = New(UnknownError, "", v)
	}
}

func trySetBool(a []*bool, v bool) {
	if len(a) > 0 && a[0] != nil {
		*(a[0]) = v
	}
}
