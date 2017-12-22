package pool

import (
	"sync"
	"testing"
	"time"
)

type testWorker struct{ int }

func newTestWorker() (Worker, error) { return &testWorker{}, nil }
func (t *testWorker) Health() bool   { return true }
func (t *testWorker) Close() error   { return nil }
func (t *testWorker) Do()            {}

func TestWorkshop(t *testing.T) {
	w := NewWorkshop(100, time.Second, newTestWorker)
	defer w.Close()
	n := 100000
	wg := new(sync.WaitGroup)
	wg.Add(n * 2)
	t1 := time.Now()
	var closeCh = make(chan struct{})
	go func() {
		for {
			select {
			case <-closeCh:
				stats := w.Stats()
				if uint64(stats.Worker) > stats.Created ||
					stats.Idle > stats.Worker ||
					stats.MinLoad > stats.MaxLoad ||
					stats.MaxLoad > stats.Doing {
					t.Fatalf("stats has bug: %+v", stats)
				} else {
					t.Logf("%+v", stats)
				}
				return
			default:
				stats := w.Stats()
				if uint64(stats.Worker) > stats.Created ||
					stats.Idle > stats.Worker ||
					stats.MinLoad > stats.MaxLoad ||
					stats.MaxLoad > stats.Doing {
					t.Fatalf("stats has bug: %+v", stats)
				} else {
					t.Logf("%+v", stats)
				}
				time.Sleep(time.Microsecond * 100)
			}
		}
	}()
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Add(-1)
			err := w.Callback(func(worker Worker) error {
				worker.(*testWorker).Do()
				return nil
			})
			if err != nil {
				t.Fatalf("%v", err)
			}
		}()
		go func() {
			defer wg.Add(-1)
			worker, err := w.Hire()
			if err != nil {
				t.Fatalf("%v", err)
			}
			worker.(*testWorker).Do()
			w.Fire(worker)
		}()
	}
	wg.Wait()
	d := time.Since(t1)
	close(closeCh)
	time.Sleep(time.Millisecond * 2500)
	t.Logf("stats: %+v, cost: %v, TPS: %v", w.Stats(), d, d/time.Duration(n*2))
}
