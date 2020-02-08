package password

// Flag denotes the expected or actual element of the password.
type Flag uint8

const (
	// N denotes the numerals.
	N Flag = 1 << 0
	// L_OR_U denotes the lowercase or uppercase letters.
	L_OR_U Flag = 1 << 1
	// L denotes the lowercase letters.
	L Flag = 1 << 2
	// U denotes the uppercase letters.
	U Flag = 1 << 3
	// S denotes the printable symbols found on the keyboard, except letters, numerals and spaces.
	S Flag = 1 << 4

	// mask is used for extracting the 5 least significant bits of the flag.
	mask Flag = 0x1f
)

// CheckPassword checks if the actual element of the password matches the expected flag.
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
	r &^= flag & (L|U) ^ (L|U)
	return r == flag
}
