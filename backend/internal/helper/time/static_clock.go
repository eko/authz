package time

import "time"

type StaticClock struct{}

func (c *StaticClock) Now() time.Time {
	return time.Date(2100, time.May, 16, 10, 0, 0, 0, time.UTC)
}

func NewStaticClock() *StaticClock {
	return &StaticClock{}
}
