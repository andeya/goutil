package status

// Check if err!=nil, create a status with stack, and panic.
// NOTE:
//  If err!=nil and msg=="", error text is set to msg
func Check(err error, code int32, msg string, whenError ...func()) {
	if err == nil {
		return
	}
	if len(whenError) > 0 && whenError[0] != nil {
		whenError[0]()
	}
	panic(New(code, msg, err).TagStack(1))
}

// Throw creates a status with stack, and panic.
func Throw(code int32, msg string, cause ...interface{}) {
	panic(New(code, msg, cause...).TagStack(1))
}

// Panic panic.
// TODO: remove
func Panic(stat *Status) {
	panic(stat)
}

// CatchWithStack recovers the panic and returns status copy with stack.
// NOTE:
//  Set `realStat` to true if a `Status` type is recovered
// Example:
//  var stat *Status
//  defer CatchWithStack(&stat)
func CatchWithStack(statPtr **Status, realStat ...*bool) {
	r := recover()

	if statPtr == nil {
		switch r.(type) {
		case *Status, Status:
			trySetBool(realStat, true)
		default:
			trySetBool(realStat, false)
		}
		return
	}
	stack := findPanicStack()
	switch v := r.(type) {
	case nil:
		// Keep the original abnormal status
		if !(*statPtr).OK() {
			trySetBool(realStat, true)
			(*statPtr).stack = stack
		} else {
			trySetBool(realStat, false)
			*statPtr = &Status{stack: stack}
		}
	case *Status:
		trySetBool(realStat, true)
		if v == nil {
			*statPtr = &Status{stack: stack}
		} else {
			*statPtr = &Status{code: v.code, msg: v.msg, cause: v.cause, stack: stack}
		}
	case Status:
		trySetBool(realStat, true)
		*statPtr = &Status{code: v.code, msg: v.msg, cause: v.cause, stack: stack}
	default:
		trySetBool(realStat, false)
		*statPtr = &Status{code: UnknownError, cause: toErr(v), stack: stack}
	}
}

// Catch recovers the panic and returns status.
// NOTE:
//  Set `realStat` to true if a `Status` type is recovered
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
		return
	}

	switch v := r.(type) {
	case nil:
		// Keep the original abnormal status
		if !(*statPtr).OK() {
			trySetBool(realStat, true)
			return
		}
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
		*statPtr = New(UnknownError, "", v).TagStack(2)
	}
}

func trySetBool(a []*bool, v bool) {
	if len(a) > 0 && a[0] != nil {
		*(a[0]) = v
	}
}
