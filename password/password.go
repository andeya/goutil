package password

// Flag password material requirement
type Flag uint8

const (
	// N Numbers
	N Flag = 1 << 0
	// L_OR_U Uppercase or lowercase letter
	L_OR_U Flag = 1 << 1
	// L Lowercase letters
	L Flag = 1 << 2
	// U Uppercase letter
	U Flag = 1 << 3
	// S Symbols found on the keyboard (all keyboard characters not defined as letters or numerals) and spaces
	S Flag = 1 << 4

	mask Flag = 0x1f
)

// CheckPassword checks if the password matches the format requirements.
func CheckPassword(pw string, flag Flag, minLen int, maxLen ...int) bool {
	if len(pw) < minLen ||
		(len(maxLen) > 0 && len(pw) > maxLen[0]) {
		return false
	}
	flag &= mask
	if flag == flag|L || flag == flag|U {
		flag |= L_OR_U
	}
	var r Flag
	for _, c := range pw {
		if c >= '0' && c <= '9' {
			r |= N
		} else if c >= 'a' && c <= 'z' {
			r |= L
			r |= L_OR_U
		} else if c >= 'A' && c <= 'Z' {
			r |= U
			r |= L_OR_U
		} else if (c >= '!' && c <= '/') ||
			(c >= ':' && c <= '@') ||
			(c >= '[' && c <= '`') ||
			(c >= '{' && c <= '~') {
			r |= S
		} else {
			return false
		}
	}
	return r == flag
}
