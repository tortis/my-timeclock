package main

import "time"

const (
	Sunday    = 0
	Monday    = 1
	Tuesday   = 2
	Wednesday = 3
	Thursday  = 4
	Friday    = 5
	Saturday  = 6
)

type Week struct {
	Year      int
	WeekNum   int
	Days      []Day
}

func NewWeek() *Week {
	year, week := time.Now().ISOWeek()
	return &Week{
		Year: year,
		WeekNum: week,
		Days: make([]Day,7),
	}
}
