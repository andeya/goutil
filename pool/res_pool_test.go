package pool

import (
	"context"
	"errors"
	"testing"
)

func TestResPool(t *testing.T) {
	p := NewResPool("name", func(context.Context) (Resource, error) {
		return nil, errors.New("new error")
	})
	res, err := p.Get()
	t.Logf("res: %#v, err: %v", res, err)
}
