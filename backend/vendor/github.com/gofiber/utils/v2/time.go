package utils

import (
	"sync"
	"sync/atomic"
	"time"
)

var (
	timestamp   atomic.Uint32
	updaterMu   sync.Mutex
	stopUpdater chan struct{}
	updaterDone chan struct{}
)

// newTicker returns a channel and stop func because a fake *time.Ticker panics on Stop.
var (
	now       = time.Now
	newTicker = func(d time.Duration) (<-chan time.Time, func()) {
		t := time.NewTicker(d)
		return t.C, t.Stop
	}
)

// Timestamp returns the Unix time in seconds as a uint32, which wraps on 2106-02-07. While the updater
// runs it is the value it last published, at most about a second old, for the cost of one atomic load;
// otherwise it reads the clock.
func Timestamp() uint32 {
	if t := timestamp.Load(); t != 0 {
		return t
	}
	return clockTimestamp()
}

// clockTimestamp is Timestamp's out-of-line slow path.
//
//go:noinline
func clockTimestamp() uint32 {
	return uint32(now().Unix())
}

// StartTimeStampUpdater starts a goroutine that refreshes the cached timestamp every second.
// Repeated and concurrent calls are safe.
func StartTimeStampUpdater() {
	updaterMu.Lock()
	defer updaterMu.Unlock()
	if stopUpdater != nil {
		return
	}

	timestamp.Store(uint32(now().Unix()))
	stopUpdater = make(chan struct{})
	updaterDone = make(chan struct{})

	// Created before the goroutine so it exists when this returns.
	tick, stopTick := newTicker(time.Second)

	go func(stop, done chan struct{}) {
		// stopTick before close(done): Stop's wait then implies release.
		defer close(done)
		defer stopTick()

		for {
			select {
			case <-tick:
				timestamp.Store(uint32(now().Unix()))
			case <-stop:
				return
			}
		}
	}(stopUpdater, updaterDone)
}

// StopTimeStampUpdater stops the updater goroutine; call it on shutdown to avoid a leak.
// Timestamp reads the clock afterwards. Repeated and concurrent calls are safe.
func StopTimeStampUpdater() {
	updaterMu.Lock()
	defer updaterMu.Unlock()
	if stopUpdater == nil {
		return
	}

	close(stopUpdater)
	<-updaterDone
	stopUpdater = nil
	updaterDone = nil
	timestamp.Store(0) // nothing refreshes the value any more: read the clock instead
}
