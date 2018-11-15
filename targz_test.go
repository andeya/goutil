package goutil

import (
	"bytes"
	"testing"
)

func TestTarGzTo(t *testing.T) {
	var src = "pool"
	var dstWriter = bytes.NewBuffer(nil)
	if err := TarGzTo(src, dstWriter, false, t.Logf, ".git"); err != nil {
		t.Fatal(err)
	}
}
