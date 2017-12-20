package pool

import (
	"sync"
	"testing"
	"time"
)

type testWorker struct{}

func newTestWorker() (Worker, error) { return &testWorker{}, nil }
func (t *testWorker) Health() bool   { return true }
func (t *testWorker) Close() error   { return nil }
func (t *testWorker) Do()            {}

func TestWorkshop(t *testing.T) {
	w := NewWorkshop(0, 0, newTestWorker)
	defer w.Close()
	n := 100000
	wg := new(sync.WaitGroup)
	wg.Add(n)
	t1 := time.Now()
	go func() {
		for {
			t.Logf("%+v", w.Stats())
			time.Sleep(time.Microsecond)
		}
	}()
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Add(-1)
			err := w.Do(func(w Worker) error {
				w.(*testWorker).Do()
				return nil
			})
			if err != nil {
				t.Fatalf("%v", err)
			}
		}()
	}
	wg.Wait()
	d := time.Since(t1)
	t.Logf("stats: %+v, cost: %v, TPS: %v", w.Stats(), d, d/time.Duration(n))
}
