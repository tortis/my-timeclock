package main

import "time"

type Day struct {
	Blocks []*TimeBlock
	Dotw   int
	Hours  float64
}

func NewDay(dotw int) *Day {
	return &Day{
		Blocks: make([]*TimeBlock, 0),
		Dotw:   dotw,
		Hours:  0.0,
	}
}

func (day *Day) TimeWorked() time.Duration {
	var r time.Duration = 0
	for _, block := range day.Blocks {
		r += block.GetDuration()
	}
	return r
}

func (day *Day) AddTimeBlock(b *TimeBlock) {
	day.Blocks = append(day.Blocks, b)
	day.Hours = day.TimeWorked().Hours()
}
