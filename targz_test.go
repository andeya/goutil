package goutil

import (
	"fmt"
	"testing"
)

func TestTarGz(t *testing.T) {
	var src = "pool"
	var dst = fmt.Sprintf("%s.tar.gz", src)
	if err := TarGz(src, dst, false, t.Logf, ".git"); err != nil {
		t.Fatal(err)
	}
}
