package pool

import (
	"fmt"
	"math"
	"sync"
)

type (
	// Worker woker interface
	Worker interface {
		Health() bool
		Close() error
	}
	// Workshop workshop
	Workshop struct {
		newFn    func() (Worker, error)
		maxQuota int32
		maxIdle  int32
		infos    map[*workerInfo]struct{}
		mostFree *workerInfo
		stats    *WorkshopStats
		lock     sync.Mutex
		wg       sync.WaitGroup
		closeCh  chan struct{}
	}
	workerInfo struct {
		worker Worker
		jobNum int32
		wg     sync.WaitGroup
	}
	// WorkshopStats workshop stats
	WorkshopStats struct {
		Worker    int32
		Idle      int32
		Created   uint64
		Closed    uint64
		Hire      uint64
		Fire      uint64
		Doing     int32
		MostUsed  int32
		LeastUsed int32
	}
)

const (
	defaultWorkerMaxQuota = 64
	defaultWorkerMaxIdle  = 8
)

// NewWorkshop creates a new workshop.
func NewWorkshop(maxQuota, maxIdle int32, newWorkerFunc func() (Worker, error)) *Workshop {
	if maxQuota <= 0 {
		maxQuota = defaultWorkerMaxQuota
	}
	if maxIdle <= 0 {
		maxIdle = defaultWorkerMaxIdle
	}
	if maxIdle > maxQuota {
		maxIdle = maxQuota
	}
	w := new(Workshop)
	w.stats = new(WorkshopStats)
	w.maxQuota = maxQuota
	w.maxIdle = maxIdle
	w.infos = make(map[*workerInfo]struct{}, maxQuota)
	w.closeCh = make(chan struct{})
	w.newFn = func() (worker Worker, err error) {
		defer func() {
			if p := recover(); p != nil {
				err = fmt.Errorf("%v", p)
			} else {
				w.stats.Created++
			}
		}()
		return newWorkerFunc()
	}
	return w
}

// Do assign a worker to execute the callback function.
func (w *Workshop) Do(callback func(Worker) error) error {
	select {
	case <-w.closeCh:
		return ErrWorkshopClosed
	default:
	}
	w.lock.Lock()
	info, err := w.hireLocked()
	w.lock.Unlock()
	if err != nil {
		return err
	}
	worker := info.worker
	defer func() {
		w.lock.Lock()
		w.fireLocked(info)
		w.lock.Unlock()
	}()
	return callback(worker)
}

var ErrWorkshopClosed = fmt.Errorf("%s", "workshop is closed")

func (w *Workshop) hireLocked() (*workerInfo, error) {
	var info = w.mostFree
	if len(w.infos) >= int(w.maxQuota) || (info != nil && info.jobNum == 0) {
		info.use()
		if info.jobNum == 1 {
			w.stats.Idle--
		}
		w.updateFreeLocked()
	} else {
		worker, err := w.newFn()
		if err != nil {
			return nil, err
		}
		info = &workerInfo{
			worker: worker,
			wg:     w.wg,
		}
		info.use()
		w.infos[info] = struct{}{}
		w.mostFree = info
	}

	w.stats.Hire++
	return info, nil
}

func (info *workerInfo) use() {
	info.jobNum++
	info.wg.Add(1)
}

func (info *workerInfo) free() {
	info.jobNum--
	info.wg.Add(-1)
}

func (w *Workshop) fireLocked(info *workerInfo) {
	w.stats.Fire++
	if !info.worker.Health() {
		delete(w.infos, info)
		w.stats.Closed++
		return
	}
	info.free()
	jobNum := info.jobNum
	if jobNum == 0 {
		if w.stats.Idle >= w.maxIdle {
			delete(w.infos, info)
			info.worker.Close()
			w.stats.Closed++
			return
		}
		w.stats.Idle++
	}
	if jobNum >= w.mostFree.jobNum {
		return
	}
	w.mostFree = info
}

func (w *Workshop) updateFreeLocked() {
	var mostFree = w.mostFree
	for info := range w.infos {
		if info.jobNum >= mostFree.jobNum {
			continue
		}
		mostFree = info
	}
	w.mostFree = mostFree
}

// Stats returns the current workshop stats.
func (w *Workshop) Stats() *WorkshopStats {
	w.lock.Lock()
	w.stats.Worker = int32(len(w.infos))
	w.stats.Doing = int32(w.stats.Hire - w.stats.Fire)
	var max, min int32
	if w.stats.Worker > 0 {
		var tmp int32
		min = math.MaxInt32
		for info := range w.infos {
			tmp = info.jobNum
			if tmp > max {
				max = tmp
			}
			if tmp < min {
				min = tmp
			}
		}
	}
	w.stats.LeastUsed = min
	w.stats.MostUsed = max
	w.lock.Unlock()
	return w.stats
}

// Close wait for all the work to be completed and close the workshop.
func (w *Workshop) Close() {
	select {
	case <-w.closeCh:
		return
	default:
		close(w.closeCh)
	}
	w.wg.Wait()
	w.lock.Lock()
	defer w.lock.Unlock()
	for info := range w.infos {
		info.worker.Close()
		w.stats.Closed++
	}
	w.infos = nil
	w.stats.Idle = 0
}
