package goutil

import (
	"fmt"
	"testing"
)

func TestBytesToString(t *testing.T) {
	bb := []byte("testing: BytesToString")
	ss := BytesToString(bb)
	t.Logf("type: %T, value: %v", ss, ss)
}

func TestStringToBytes(t *testing.T) {
	s := "testing: StringToBytes"
	b := StringToBytes(s)
	t.Logf("type: %T, value: %v, val-string: %s\n", b, b, b)
	b = append(b, '!')
	t.Logf("after append:\ntype: %T, value: %v, val-string: %s\n", b, b, b)
}

func TestSnakeString(t *testing.T) {
	data := [][2]string{
		{"XxYy", "xx_yy"},
		{"_XxYy", "_xx_yy"},
		{"TcpRpc", "tcp_rpc"},
		{"ID", "id"},
		{"UserID", "user_id"},
		{"RPC", "rpc"},
		{"TCP_RPC", "tcp_rpc"},
		{"wakeRPC", "wake_rpc"},
		{"_TCP__RPC", "_tcp__rpc"},
		{"_TcP__RpC_", "_tc_p__rp_c_"},
	}
	for _, p := range data {
		r := SnakeString(p[0])
		if r != p[1] {
			t.Fatalf("[SnakeString] %s: expect: %s, but get %s", p[0], p[1], r)
		}
		r = SnakeString(p[1])
		if r != p[1] {
			t.Fatalf("[SnakeString] %s: expect: %s, but get %s", p[1], p[1], r)
		}
	}
}

func TestCamelString(t *testing.T) {
	data := [][2]string{
		{"xx_yy", "XxYy"},
		{"_xx_yy", "_XxYy"},
		{"id", "Id"},
		{"user_id", "UserId"},
		{"rpc", "Rpc"},
		{"tcp_rpc", "TcpRpc"},
		{"wake_rpc", "WakeRpc"},
		{"_tcp___rpc", "_Tcp__Rpc"},
		{"_tc_p__rp_c__", "_TcP_RpC__"},
	}
	for _, p := range data {
		r := CamelString(p[0])
		if r != p[1] {
			t.Fatalf("[CamelString] %s: expect: %s, but get %s", p[0], p[1], r)
		}
		r = CamelString(p[1])
		if r != p[1] {
			t.Fatalf("[CamelString] %s: expect: %s, but get %s", p[1], p[1], r)
		}
	}
}

func TestHTMLEntityToUTF8(t *testing.T) {
	want := `{"info":[["color","咖啡色|绿色"]]｝`
	got := HTMLEntityToUTF8(`{"info":[["color","&#5496;&#5561;&#8272;&#7c;&#7eff;&#8272;"]]｝`, 16)
	if got != want {
		t.Fatalf("want: %q, got: %q", want, got)
	}
}

func TestCodePointToUTF8(t *testing.T) {
	got := CodePointToUTF8(`{"info":[["color","\u5496\u5561\u8272\u7c\u7eff\u8272"]]｝`, 16)
	want := `{"info":[["color","咖啡色|绿色"]]｝`
	if got != want {
		t.Fatalf("want: %q, got: %q", want, got)
	}
}

func TestSpaceInOne(t *testing.T) {
	a := struct {
		input  string
		output string
	}{
		input: `# authenticate method 

		//  comment2	

		/*  some other 
			  comments */
		`,
		output: `# authenticate method
	// comment2
	/* some other
	comments */
	`,
	}
	r := SpaceInOne(a.input)
	if r != a.output {
		t.Fatalf("want: %q, got: %q", a.output, r)
	}
}

func ExampleStringMarshalJSON() {
	s := `<>&{}""`
	fmt.Printf("%s\n", StringMarshalJSON(s, true))
	fmt.Printf("%s\n", StringMarshalJSON(s, false))
	// Output:
	// "\u003c\u003e\u0026{}\"\""
	// "<>&{}\"\""
}
