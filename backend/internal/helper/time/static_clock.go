package time

import (
	"time"
)

type StaticClock struct{}

func (c *StaticClock) Now() time.Time {
	time.Local = time.UTC

	return time.Date(2100, time.January, 1, 1, 0, 0, 0, time.UTC)
}

func NewStaticClock() *StaticClock {
	return &StaticClock{}
}
