package main

import "time"

type Day struct {
	Blocks []TimeBlock
}

func NewDay() *Day {
	return &Day{Blocks: make([]TimeBlock, 0)}
}

func (day *Day) TimeWorked() time.Duration {
	var r time.Duration = 0
	for _, block := range day.Blocks {
		r += block.GetDuration()
	}
	return r
}

func (day *Day) AddTimeBlock(b *TimeBlock) {
	day.Blocks = append(day.Blocks, *b)
}
