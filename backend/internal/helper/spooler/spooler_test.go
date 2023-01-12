package spooler

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewSpooler(t *testing.T) {
	// Given
	flushFunc := func(values []int) {
		// Do nothing
	}

	s := New(flushFunc)

	// When - Then
	assert := assert.New(t)
	assert.IsType(new(spooler[int]), s)
	assert.Equal(defaultFlushInterval, s.flushInterval)
	assert.Equal(defaultFlushLimit, s.flushLimit)
	assert.Len(s.values, 0)
}

func TestNewSpooler_Generic(t *testing.T) {
	// Given
	flushFunc1 := func(values []int) {}
	flushFunc2 := func(values []string) {}
	flushFunc3 := func(values []any) {}

	spooler1 := New(flushFunc1)
	spooler2 := New(flushFunc2)
	spooler3 := New(flushFunc3)

	// When - Then
	assert := assert.New(t)

	assert.IsType(new(spooler[int]), spooler1)
	assert.IsType(new(spooler[string]), spooler2)
	assert.IsType(new(spooler[any]), spooler3)
}

func TestSpooler_Add(t *testing.T) {
	// Given
	s := New(func(values []int) {})

	// When
	s.Add(1)
	s.Add(2)
	s.Add(3)

	// Then
	assert := assert.New(t)
	assert.Len(s.values, 3)
}

func TestSpooler_Add_FlushWhenMaxReached(t *testing.T) {
	assert := assert.New(t)

	// Given
	var done = make(chan struct{}, 1)

	s := New(func(values []int) {
		assert.Len(values, 3) // 3 items should be flushed
		done <- struct{}{}
	}, WithFlushLimit(3))

	// When
	s.Add(1)
	s.Add(2)
	s.Add(3)
	s.Add(4)

	// Then
	var timeout = time.After(2 * time.Second)

	select {
	case <-timeout:
		assert.Fail("timeout reached, flush never happened.")
	case <-done:
		assert.Len(s.values, 1) // Item 4 should be in channel
	}
}

func TestSpooler_Add_FlushWhenTickerReached(t *testing.T) {
	assert := assert.New(t)

	// Given
	var done = make(chan struct{}, 1)

	s := New(func(values []int) {
		assert.Len(values, 2) // 2 items should be flushed
		done <- struct{}{}
	},
		WithFlushLimit(3),
		WithFlushInterval(250*time.Millisecond))

	// When
	s.Add(1)
	s.Add(2)

	// Then
	var timeout = time.After(2 * time.Second)

	select {
	case <-timeout:
		assert.Fail("timeout reached, flush never happened.")
	case <-done:
		assert.Len(s.values, 0) // All items should be flushed after 250ms
	}
}

func TestSpooler_Flush(t *testing.T) {
	assert := assert.New(t)

	// Given
	var done = make(chan struct{}, 1)

	s := New(func(values []int) {
		assert.Len(values, 2) // 2 items should be flushed
		done <- struct{}{}
	},
		WithFlushLimit(3),
		WithFlushInterval(500*time.Millisecond))

	s.Add(1)
	s.Add(2)

	// When
	s.Flush()

	// Then
	var timeout = time.After(2 * time.Second)

	select {
	case <-timeout:
		assert.Fail("timeout reached, flush never happened.")
	case <-done:
		assert.Len(s.values, 0) // All items should be flushed on Flush() call
	}
}

func TestSpooler_Flush_WhenNoItems(t *testing.T) {
	assert := assert.New(t)

	// Given
	var done = make(chan struct{}, 1)

	s := New(func(values []int) {
		assert.Len(values, 0) // 0 items should be flushed
		done <- struct{}{}
	},
		WithFlushLimit(3),
		WithFlushInterval(500*time.Millisecond))

	// When
	s.Flush()

	// Then
	var timeout = time.After(2 * time.Second)

	select {
	case <-timeout:
		assert.Fail("timeout reached, flush never happened.")
	case <-done:
		assert.Len(s.values, 0) // No items added so should be kept to 0
	}
}
