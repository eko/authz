package time

import "time"

type Clock interface {
	Now() time.Time
}

type systemClock struct{}

func NewClock() Clock {
	return &systemClock{}
}

func (c *systemClock) Now() time.Time {
	return time.Now()
}
