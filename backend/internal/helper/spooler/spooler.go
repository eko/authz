package spooler

import (
	"time"
)

const (
	// Default temporary channel size value.
	defaultChanSize uint32 = 10000

	// Default flush interval value: all items in channel will be flushed.
	defaultFlushInterval = 3 * time.Second

	// Default flush limit: when channel hits this number of items, they are all flushed.
	defaultFlushLimit uint32 = 100
)

// SpoolerOption represents the spooler available options.
type SpoolerOption func(*spoolerOptions)

type spoolerOptions struct {
	chanSize      uint32
	flushLimit    uint32
	flushInterval time.Duration
}

// WithChanSize defines a custom channel maximum size (capacity).
//
// When the maximum size of the channel is hit, the spooler automatically flushes items
// before allowing you to add new ones.
func WithChanSize(chanSize uint32) SpoolerOption {
	return func(o *spoolerOptions) {
		o.chanSize = chanSize
	}
}

// WithFlushInterval defines the flush interval delay (every 3 seconds by default).
//
// Every items that will be in the channel when this delay occurs will be flushed.
func WithFlushInterval(flushInterval time.Duration) SpoolerOption {
	return func(o *spoolerOptions) {
		o.flushInterval = flushInterval
	}
}

// WithFlushLimit defines when the flush should occur.
//
// Item can still be added to the channel if the maximum chan size is not hit.
func WithFlushLimit(flushLimit uint32) SpoolerOption {
	return func(o *spoolerOptions) {
		o.flushLimit = flushLimit
	}
}

type spooler[T any] struct {
	flushInterval time.Duration
	flushLimit    uint32
	values        chan T
	flusher       func(values []T)
}

// New instanciates a new spooler generic instance.
func New[T any](flusher func(values []T), options ...SpoolerOption) *spooler[T] {
	var opts = &spoolerOptions{
		chanSize:      defaultChanSize,
		flushInterval: defaultFlushInterval,
		flushLimit:    defaultFlushLimit,
	}

	for _, option := range options {
		option(opts)
	}

	s := &spooler[T]{
		flushInterval: opts.flushInterval,
		flushLimit:    opts.flushLimit,
		values:        make(chan T, opts.chanSize),
		flusher:       flusher,
	}

	go s.loop()

	return s
}

func (s *spooler[T]) loop() {
	ticker := time.NewTicker(s.flushInterval)

	for range ticker.C {
		s.Flush()
	}
}

// Add allows to add an item to the spooler queue.
func (s *spooler[T]) Add(item T) {
	s.values <- item

	if len(s.values) >= cap(s.values) || len(s.values) >= int(s.flushLimit) {
		s.Flush()
	}
}

// Flush allows to manually flush the queue.
func (s *spooler[T]) Flush() {
	var values = make([]T, 0)

	total := len(s.values)

	for i := 0; i <= total-1; i++ {
		values = append(values, <-s.values)
	}

	s.flusher(values)
}
