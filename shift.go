package main

import (
	"errors"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Shift struct {
	Id     bson.ObjectId `bson:"_id,omitempty"`
	On     int64         `bson:"on,omitempty"`
	Off    int64         `bson:"off"`
	Active bool          `bson:"active"`
}

func NewShift() *Shift {
	return &Shift{
		On:     time.Now().Unix(),
		Off:    time.Now().Unix(),
		Active: true,
	}
}

func FullShift(on, off int64) (*Shift, error) {
	if off < on {
		return nil, errors.New("Off time is before on time.")
	}
	return &Shift{
		On:     on,
		Off:    off,
		Active: false,
	}, nil
}

func (s *Shift) DayOverlap(day_start int64) int64 {
	day_end := day_start + 86400 // Seconds in a day
	// Handle cases where shift does not overlap with day
	if s.On >= day_end {
		return 0
	}
	if s.Off < day_start {
		return 0
	}
	if s.Active {
		return min(day_end, time.Now().Unix()) - max(day_start, s.On)
	} else {
		return min(day_end, s.Off) - max(day_start, s.On)
	}
}

func (s *Shift) OnDay(day_start int64) bool {
	day_end := day_start + 86400 // Seconds in a day
	if s.Off > day_start && s.Off <= day_end {
		return true
	} else {
		return false
	}
}

func max(a, b int64) int64 {
	if a > b {
		return a
	} else {
		return b
	}
}

func min(a, b int64) int64 {
	if a < b {
		return a
	} else {
		return b
	}
}
