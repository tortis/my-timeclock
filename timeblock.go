package main

import "time"

type TimeBlock struct {
	In  time.Time
	Out time.Time
}

func NewTimeBlock() *TimeBlock {
	return &TimeBlock{In: time.Now()}
}

func (b *TimeBlock) GetDuration() time.Duration {
	if b.Out.Unix() < 0 {
		return time.Now().Sub(b.In)
	} else {
		return b.Out.Sub(b.In)
	}
}

func (b *TimeBlock) End() {
	b.Out = time.Now()
}
