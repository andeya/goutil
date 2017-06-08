// Shutdown closes current program gracefully.
package shutdown

import (
	"context"
	logPkg "log"
	"os"
	"os/signal"
	"sync/atomic"
	"time"
)

func init() {
	go func() {
		// subscribe to SIGINT signals
		stopChan := make(chan os.Signal)
		signal.Notify(stopChan, os.Interrupt, os.Kill)
		<-stopChan // wait for SIGINT
		signal.Stop(stopChan)
		Shutdown()
	}()
}

var (
	finalizers      []func() error
	defaultTimeout         = time.Minute
	shutdownTimeout        = defaultTimeout
	log             Logger = new(logger)
)

type (
	// Logger logger interface
	Logger interface {
		Infof(format string, v ...interface{})
		Errorf(format string, v ...interface{})
	}
	logger struct{}
)

func (l *logger) Infof(format string, v ...interface{}) {
	logPkg.Printf("[I] "+format, v...)
}

func (l *logger) Errorf(format string, v ...interface{}) {
	logPkg.Printf("[E] "+format, v...)
}

// SetLog resets logger.
func SetLog(logger Logger) {
	log = logger
}

// SetShutdown sets the function which is called after current program shutdown,
// and the time-out period for current program shutdown.
// If parameter timeout is 0, automatically use default `defaultTimeout`(60s).
// If parameter timeout less than 0, it is indefinite period.
// The finalizer function is executed before the shutdown deadline, but it is not guaranteed to be completed.
func SetShutdown(timeout time.Duration, fn ...func() error) {
	if timeout == 0 {
		timeout = defaultTimeout
	} else if timeout < 0 {
		timeout = 1<<63 - 1
	}
	shutdownTimeout = timeout
	finalizers = fn
}

// Shutdown closes current program gracefully.
// Parameter timeout is used to reset time-out period for current program shutdown.
func Shutdown(timeout ...time.Duration) {
	log.Infof("shutting down servers...\n")
	if len(timeout) > 0 {
		SetShutdown(timeout[0], finalizers...)
	}
	graceful := shutdown()
	if graceful {
		log.Infof("servers are shutted down gracefully.\n")
		os.Exit(0)
	} else {
		log.Errorf("servers are shutted down, but not gracefully.\n")
		os.Exit(-1)
	}
}

func shutdown() bool {
	ctxTimeout, _ := context.WithTimeout(context.Background(), shutdownTimeout)
	var flag int32 = 1
	fchan := make(chan bool)
	var idx int32
	go func() {
		for i, finalizer := range finalizers {
			atomic.StoreInt32(&idx, int32(i))
			if finalizer == nil {
				continue
			}
			select {
			case <-ctxTimeout.Done():
				break
			default:
				if err := finalizer(); err != nil {
					atomic.StoreInt32(&flag, 0)
					log.Errorf("[shutdown-finalizer%d] %s\n", i, err.Error())
				}
			}
		}
		close(fchan)
	}()
	select {
	case <-ctxTimeout.Done():
		if err := ctxTimeout.Err(); err != nil {
			atomic.StoreInt32(&flag, 0)
			log.Errorf("[shutdown-finalizer%d] %s\n", atomic.LoadInt32(&idx), err.Error())
		}
	case <-fchan:
	}
	return flag == 1
}
