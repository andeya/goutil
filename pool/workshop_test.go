package pool

import (
	"sync"
	"testing"
	"time"
)

/*
go test -v -run=^TestChanPool$ -bench=^$ && go test -v -run=^TestWorkshop$ -bench=^$ && go test -v -run=^$ -bench=^Benchmark
*/

// // Scene A
// const (
// 	poolSize      = 50
// 	requests    = 1000000
// 	oneLogicCostTime = time.Millisecond
// )

// Scene B
const (
	poolSize         = 50
	requests         = 100000
	oneLogicCostTime = time.Millisecond * 10
)

type testWorker struct{ int }

func newTestWorker() (Worker, error) { return &testWorker{}, nil }
func (t *testWorker) Health() bool   { return true }
func (t *testWorker) Close() error   { return nil }
func (t *testWorker) Do()            { time.Sleep(oneLogicCostTime) }

func reportStat(t *testing.T, startNano int64) {
	totalNano := time.Now().UnixNano() - startNano
	t.Logf(
		"pool_size:%d, requests:%d, one_logic_cost:%v, total_cost:%v, QPS:%d",
		poolSize,
		requests,
		oneLogicCostTime,
		time.Duration(totalNano),
		int64(requests*time.Second)/totalNano,
	)
}

func TestChanPool(t *testing.T) {
	var workerPool = make(chan *testWorker, poolSize)
	for i := 0; i < poolSize; i++ {
		workerPool <- new(testWorker)
	}

	wg := new(sync.WaitGroup)
	wg.Add(requests)
	startNano := time.Now().UnixNano()
	for i := 0; i < requests; i++ {
		go func() {
			defer wg.Done()
			worker := <-workerPool
			worker.Do()
			workerPool <- worker
		}()
	}
	wg.Wait()
	reportStat(t, startNano)
}

func TestWorkshop(t *testing.T) {
	w := NewWorkshop(poolSize, time.Second, newTestWorker)
	defer w.Close()
	wg := new(sync.WaitGroup)
	wg.Add(requests)
	startNano := time.Now().UnixNano()
	for i := 0; i < requests; i++ {
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
	t.Logf("workshop stats: %+v", w.Stats())
	reportStat(t, startNano)
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
