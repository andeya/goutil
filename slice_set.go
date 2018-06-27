package goutil

// SetStrings sets a element to the string set.
func SetStrings(set []string, a string) []string {
	for _, s := range set {
		if s == a {
			return set
		}
	}
	return append(set, a)
}

// SetInts sets a element to the int set.
func SetInts(set []int, a int) []int {
	for _, s := range set {
		if s == a {
			return set
		}
	}
	return append(set, a)
}

// SetInt32s sets a element to the int32 set.
func SetInt32s(set []int32, a int32) []int32 {
	for _, s := range set {
		if s == a {
			return set
		}
	}
	return append(set, a)
}

// SetInt64s sets a element to the int64 set.
func SetInt64s(set []int64, a int64) []int64 {
	for _, s := range set {
		if s == a {
			return set
		}
	}
	return append(set, a)
}

// SetInterfaces sets a element to the interface{} set.
func SetInterfaces(set []interface{}, a interface{}) []interface{} {
	for _, s := range set {
		if s == a {
			return set
		}
	}
	return append(set, a)
}
