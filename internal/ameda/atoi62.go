package ameda

import (
	"errors"
	"math"
	"strconv"
)

// ParseUint is like ParseInt but for unsigned numbers.
// NOTE:
//  Compatible with standard package strconv.
func ParseUint(s string, base int, bitSize int) (uint64, error) {
	// Ignore letter case
	if base <= 36 {
		return strconv.ParseUint(s, base, bitSize)
	}

	const fnParseUint = "ParseUint"

	if base > 62 {
		return 0, baseError(fnParseUint, s, base)
	}

	if s == "" || !underscoreOK(s) {
		return 0, syntaxError(fnParseUint, s)
	}

	if bitSize == 0 {
		bitSize = int(strconv.IntSize)
	} else if bitSize < 0 || bitSize > 64 {
		return 0, bitSizeError(fnParseUint, s, bitSize)
	}

	// Cutoff is the smallest number such that cutoff*base > maxUint64.
	// Use compile-time constants for common cases.
	cutoff := math.MaxUint64/uint64(base) + 1

	maxVal := uint64(1)<<uint(bitSize) - 1

	var n uint64
	for _, c := range []byte(s) {
		var d byte
		switch {
		case '0' <= c && c <= '9':
			d = c - '0'
		case 'a' <= c && c <= 'z':
			d = c - 'a' + 10
		case 'A' <= c && c <= 'Z':
			d = c - 'A' + 10 + 26
		case c == '_':
			continue
		default:
			return 0, syntaxError(fnParseUint, s)
		}

		if d >= byte(base) {
			return 0, syntaxError(fnParseUint, s)
		}

		if n >= cutoff {
			// n*base overflows
			return maxVal, rangeError(fnParseUint, s)
		}
		n *= uint64(base)

		n1 := n + uint64(d)
		if n1 < n || n1 > maxVal {
			// n+v overflows
			return maxVal, rangeError(fnParseUint, s)
		}
		n = n1
	}

	return n, nil
}

// ParseInt interprets a string s in the given base (0, 2 to 62) and
// bit size (0 to 64) and returns the corresponding value i.
//
// If base == 0, the base is implied by the string's prefix:
// base 2 for "0b", base 8 for "0" or "0o", base 16 for "0x",
// and base 10 otherwise. Also, for base == 0 only, underscore
// characters are permitted per the Go integer literal syntax.
// If base is below 0, is 1, or is above 62, an error is returned.
//
// The bitSize argument specifies the integer type
// that the result must fit into. Bit sizes 0, 8, 16, 32, and 64
// correspond to int, int8, int16, int32, and int64.
// If bitSize is below 0 or above 64, an error is returned.
//
// The errors that ParseInt returns have concrete type *NumError
// and include err.Num = s. If s is empty or contains invalid
// digits, err.Err = ErrSyntax and the returned value is 0;
// if the value corresponding to s cannot be represented by a
// signed integer of the given size, err.Err = ErrRange and the
// returned value is the maximum magnitude integer of the
// appropriate bitSize and sign.
// NOTE:
//  Compatible with standard package strconv.
func ParseInt(s string, base int, bitSize int) (i int64, err error) {
	// Ignore letter case
	if base <= 36 {
		return strconv.ParseInt(s, base, bitSize)
	}

	const fnParseInt = "ParseInt"

	if s == "" {
		return 0, syntaxError(fnParseInt, s)
	}

	// Pick off leading sign.
	s0 := s
	neg := false
	if s[0] == '+' {
		s = s[1:]
	} else if s[0] == '-' {
		neg = true
		s = s[1:]
	}

	// Convert unsigned and check range.
	var un uint64
	un, err = ParseUint(s, base, bitSize)
	if err != nil && err.(*strconv.NumError).Err != strconv.ErrRange {
		err.(*strconv.NumError).Func = fnParseInt
		err.(*strconv.NumError).Num = s0
		return 0, err
	}

	if bitSize == 0 {
		bitSize = int(strconv.IntSize)
	}

	cutoff := uint64(1 << uint(bitSize-1))
	if !neg && un >= cutoff {
		return int64(cutoff - 1), rangeError(fnParseInt, s0)
	}
	if neg && un > cutoff {
		return -int64(cutoff), rangeError(fnParseInt, s0)
	}
	n := int64(un)
	if neg {
		n = -n
	}
	return n, nil
}

// underscoreOK reports whether the underscores in s are allowed.
// Checking them in this one function lets all the parsers skip over them simply.
// Underscore must appear only between digits or between a base prefix and a digit.
func underscoreOK(s string) bool {
	// saw tracks the last character (class) we saw:
	// ^ for beginning of number,
	// 0 for a digit or base prefix,
	// _ for an underscore,
	// ! for none of the above.
	saw := '^'
	i := 0

	// Optional sign.
	if len(s) >= 1 && (s[0] == '-' || s[0] == '+') {
		s = s[1:]
	}

	// Optional base prefix.
	if len(s) >= 2 && s[0] == '0' && (lower(s[1]) == 'b' || lower(s[1]) == 'o' || lower(s[1]) == 'x') {
		i = 2
		saw = '0' // base prefix counts as a digit for "underscore as digit separator"
	}

	// Number proper.
	for ; i < len(s); i++ {
		// Digits are always okay.
		if '0' <= s[i] && s[i] <= '9' || 'a' <= lower(s[i]) && lower(s[i]) <= 'z' {
			saw = '0'
			continue
		}
		// Underscore must follow digit.
		if s[i] == '_' {
			if saw != '0' {
				return false
			}
			saw = '_'
			continue
		}
		// Saw non-digit, non-underscore.
		return false
	}
	return true
}

// Atoi is equivalent to ParseInt(s, 10, 0), converted to type int.
func Atoi(s string) (int, error) {
	const fnAtoi = "Atoi"

	sLen := len(s)
	if strconv.IntSize == 32 && (0 < sLen && sLen < 10) ||
		strconv.IntSize == 64 && (0 < sLen && sLen < 19) {
		// Fast path for small integers that fit int type.
		s0 := s
		if s[0] == '-' || s[0] == '+' {
			s = s[1:]
			if len(s) < 1 {
				return 0, &strconv.NumError{fnAtoi, s0, strconv.ErrSyntax}
			}
		}

		n := 0
		for _, ch := range []byte(s) {
			ch -= '0'
			if ch > 9 {
				return 0, &strconv.NumError{fnAtoi, s0, strconv.ErrSyntax}
			}
			n = n*10 + int(ch)
		}
		if s0[0] == '-' {
			n = -n
		}
		return n, nil
	}

	// Slow path for invalid, big, or underscored integers.
	i64, err := ParseInt(s, 10, 0)
	if nerr, ok := err.(*strconv.NumError); ok {
		nerr.Func = fnAtoi
	}
	return int(i64), err
}

// lower(c) is a lower-case letter if and only if
// c is either that lower-case letter or the equivalent upper-case letter.
// Instead of writing c == 'x' || c == 'X' one can write lower(c) == 'x'.
// Note that lower of non-letters can produce other non-letters.
func lower(c byte) byte {
	return c | ('x' - 'X')
}
func syntaxError(fn, str string) *strconv.NumError {
	return &strconv.NumError{fn, str, strconv.ErrSyntax}
}
func baseError(fn, str string, base int) *strconv.NumError {
	return &strconv.NumError{fn, str, errors.New("invalid base " + strconv.Itoa(base))}
}
func rangeError(fn, str string) *strconv.NumError {
	return &strconv.NumError{fn, str, strconv.ErrRange}
}
func bitSizeError(fn, str string, bitSize int) *strconv.NumError {
	return &strconv.NumError{fn, str, errors.New("invalid bit size " + strconv.Itoa(bitSize))}
}
