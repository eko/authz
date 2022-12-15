package time

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClock_Now(t *testing.T) {
	// Given
	clock := NewClock()

	// When - Then
	assert.Implements(t, new(Clock), clock)
	assert.IsType(t, time.Time{}, clock.Now())
}
