package goutil

import (
	"testing"
)

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
