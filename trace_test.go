package goutil

import (
	"testing"
)

func TestGetCallLine(t *testing.T) {
	line := GetCallLine(0)
	t.Log("====", line)
}
