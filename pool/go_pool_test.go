package pool

import (
	"context"
	"sync"
	"testing"
)

func TestGoPool(t *testing.T) {
	gp := NewGoPool(10, 0)
	wg := new(sync.WaitGroup)
	retryTimes := 0
	var err error
	for i := 0; i < 100; i++ {
		wg.Add(1)
		a := i
		err = gp.Go(func() {
			t.Log("done:", a)
			wg.Done()
		})
		if err != nil {
			retryTimes++
			i--
			t.Log(err)
			wg.Done()
		}
	}
	wg.Wait()
	gp.Stop()
	t.Logf("retryTimes: %d", retryTimes)
}

func TestGoPool2(t *testing.T) {
	gp := NewGoPool(10, 0)
	wg := new(sync.WaitGroup)
	retryTimes := 0
	var err error
	for i := 0; i < 100; i++ {
		wg.Add(1)
		a := i
		err = gp.MustGo(context.Background(), func() {
			t.Log("done:", a)
			wg.Done()
		})
		if err != nil {
			retryTimes++
			i--
			t.Log(err)
			wg.Done()
		}
	}
	wg.Wait()
	gp.Stop()
	if retryTimes > 0 {
		t.Fatalf("except 0, but got %d", retryTimes)
	}
}
