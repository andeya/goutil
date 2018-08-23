package pool

import (
	"sync"
	"testing"
	"time"
)

const (
	poolSize      = 50
	operations    = 1000000
	logicCostTime = time.Millisecond
)

type testWorker struct{ int }

func newTestWorker() (Worker, error) { return &testWorker{}, nil }
func (t *testWorker) Health() bool   { return true }
func (t *testWorker) Close() error   { return nil }
func (t *testWorker) Do()            { time.Sleep(logicCostTime) }

func TestWorkshop(t *testing.T) {
	w := NewWorkshop(poolSize, time.Second, newTestWorker)
	defer w.Close()
	wg := new(sync.WaitGroup)
	wg.Add(operations)
	startNano := time.Now().UnixNano()
	for i := 0; i < operations; i++ {
		go func() {
			defer wg.Done()
			err := w.Callback(func(worker Worker) error {
				worker.(*testWorker).Do()
				return nil
			})
			if err != nil {
				t.Fatalf("%v", err)
			}
		}()
	}
	wg.Wait()
	totalNano := time.Now().UnixNano() - startNano
	t.Logf(
		"stats: %+v, cost: %v ms for %d operations, TPS: %v",
		w.Stats(),
		totalNano/int64(time.Millisecond),
		operations,
		int64(operations)/(totalNano/int64(time.Second)),
	)
}

func TestChanPool(t *testing.T) {
	var workerPool = make(chan *testWorker, poolSize)
	for i := 0; i < poolSize; i++ {
		workerPool <- new(testWorker)
	}

	wg := new(sync.WaitGroup)
	wg.Add(operations)
	startNano := time.Now().UnixNano()
	for i := 0; i < operations; i++ {
		go func() {
			defer wg.Done()
			worker := <-workerPool
			worker.Do()
			workerPool <- worker
		}()
	}
	wg.Wait()
	totalNano := time.Now().UnixNano() - startNano
	t.Logf(
		"Worker: %d, Created: %[1]d, cost: %v ms for %d operations, TPS: %v",
		poolSize,
		totalNano/int64(time.Millisecond),
		operations,
		int64(operations)/(totalNano/int64(time.Second)),
	)
}

func BenchmarkWorkshop(b *testing.B) {
	w := NewWorkshop(poolSize, time.Second, newTestWorker)
	defer w.Close()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		err := w.Callback(func(worker Worker) error {
			worker.(*testWorker).Do()
			return nil
		})
		if err != nil {
			b.Fatalf("%v", err)
		}
	}
}

func BenchmarkChanPool(b *testing.B) {
	var workerPool = make(chan *testWorker, poolSize)
	for i := 0; i < poolSize; i++ {
		workerPool <- new(testWorker)
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		worker := <-workerPool
		worker.Do()
		workerPool <- worker
	}
}
