package gograceful

import (
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
)

//counts running operations, only allows shutdown when = 0
var runningOperations uint64

//when shutdown signal is initiated this blocks running operations from being incremented
var runningBlocked uint64

func shouldShutDown() bool {
	atomic.AddUint64(&runningBlocked, 1)
	return atomic.LoadUint64(&runningOperations) == 0
}

//AddRunningOperation - call when you want to stop any operations from being terminated by shutdown, if false don't proceed with your operation.
func AddRunningOperation() bool {
	if atomic.LoadUint64(&runningBlocked) != 0 {
		return false
	}
	atomic.AddUint64(&runningOperations, 1)
	return true
}

//FinishRunningOperation - call when you finish the important operation to proceed with shutdown.
func FinishRunningOperation() {
	atomic.AddUint64(&runningOperations, ^uint64(0))
}

func HandleGracefulShutdown(onShutdown func()) {
	gracefulStop := make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	go func() {
		sig := <-gracefulStop
		if sig == syscall.SIGTERM || sig == syscall.SIGINT {
			//wait until all important operations are finished.
			for !shouldShutDown() {
			}
			if onShutdown != nil {
				onShutdown()
			}
			os.Exit(0)
		}
	}()
}
