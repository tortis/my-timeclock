package main

import "time"

type Week struct {
	Year      int
	WeekNum   int
	Days      []*Day
}

func NewWeek() *Week {
	year, week := time.Now().ISOWeek()
	return &Week{
		Year: year,
		WeekNum: week,
		Days: make([]*Day,7),
	}
}

func (w *Week) Set(dayOfWeek time.Weekday, day *Day) {
	w.Days[dayOfWeek] = day
}

func (w *Week) Today() *Day {
	return w.Days[time.Now().Weekday()]
}
