package goutil

import (
	"testing"
)

func TestGetIP(t *testing.T) {
	extranet, err := ExtranetIP()
	if err != nil {
		t.Fatalf("call ExtranetIP() fail: %v", err)
	}
	t.Logf("your extranet ip: %s", extranet)

	intranet, err := IntranetIP()
	if err != nil {
		t.Fatalf("call IntranetIP() fail: %v", err)
	}
	t.Logf("your intranet ip: %s", intranet)
}
