package srcpool

import (
	"context"
	"errors"
	"testing"
)

func TestPool(t *testing.T) {
	p := New("name", func(context.Context) (Resource, error) {
		return nil, errors.New("new error")
	})
	src, err := p.Get()
	t.Logf("src: %#v, err: %v", src, err)
}
