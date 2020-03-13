package password

import "fmt"

func ExampleCheckPassword() {
	fmt.Println(CheckPassword("A", N|L_OR_U, 2))
	fmt.Println(CheckPassword("a", N|L_OR_U, 2))
	fmt.Println(CheckPassword("1", N|L_OR_U, 2))
	fmt.Println(CheckPassword("Aa", N|L_OR_U, 2))
	fmt.Println(CheckPassword("A1", N|L_OR_U, 2))
	fmt.Println(CheckPassword("a1", N|L_OR_U, 2))
	fmt.Println(CheckPassword("Aa1", N|L_OR_U, 2))
	fmt.Println(CheckPassword("AaBbCcDd", L|U, 6))
	fmt.Println(CheckPassword("ABCD1234", N|U, 6))
	fmt.Println(CheckPassword("ABCD1234", N|L_OR_U|U, 6))
	fmt.Println(CheckPassword("ABCD1234", N|L_OR_U|L, 6))
	fmt.Println(CheckPassword("abcd1234", N|L, 6, 7))
	fmt.Println(CheckPassword("Aa123456", N|L_OR_U, 6, 8))
	fmt.Println(CheckPassword("Aa@123456", Flag(1<<6)|N|L|U|S, 6, 16))
	fmt.Println(CheckPassword("Aa123456语言", N|L|U, 6))
	// Output:
	// false
	// false
	// false
	// false
	// true
	// true
	// true
	// true
	// true
	// true
	// false
	// false
	// true
	// true
	// false
}
