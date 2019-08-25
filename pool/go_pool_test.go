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
		err = gp.MustGo(func() {
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

func TestGoPool3(t *testing.T) {
	gp := NewGoPool(10, 0)
	wg := new(sync.WaitGroup)
	retryTimes := 0
	var err error
	for i := 0; i < 100; i++ {
		wg.Add(1)
		a := i
		err = gp.MustGo(func() {
			t.Log("done:", a)
			wg.Done()
		}, context.Background())
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

func BenchmarkGoPool_MustGo(b *testing.B) {
	gp := NewGoPool(10000000, 0)
	wg := new(sync.WaitGroup)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gp.MustGo(func() {})
	}
	wg.Wait()
}

func BenchmarkGoPool_MustGo_Background(b *testing.B) {
	gp := NewGoPool(10000000, 0)
	wg := new(sync.WaitGroup)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		gp.MustGo(func() {
			wg.Done()
		}, context.Background())
	}
	wg.Wait()
}

func BenchmarkGoPool_go(b *testing.B) {
	wg := new(sync.WaitGroup)
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestGoPool_Work(t *testing.T) {
	gp := NewGoPool(10, 0)
	defer gp.Stop()
	var goAdd = func(a, b int) <-chan int {
		ch := make(chan int, 1)
		gp.MustGo(func() {
			ch <- a + b
			close(ch)
		})
		return ch
	}
	ch := goAdd(1, 2)
	ret := <-ch
	if ret != 3 {
		t.Fatalf("except %d, but %d", 3, ret)
	}
}
