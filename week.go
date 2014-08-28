package main

import "time"
import "encoding/json"

type Week struct {
	Year    int
	WeekNum int
	Days    []*Day
}

func NewWeek() *Week {
	year, week := time.Now().ISOWeek()
	r := &Week{
		Year:    year,
		WeekNum: week,
		Days:    make([]*Day, 7),
	}
	for i := 0; i < 7; i++ {
		r.Days[i] = NewDay(i)
	}
	return r
}

func NewSpecificWeek(year, week int) *Week {
	r := &Week{
		Year:    year,
		WeekNum: week,
		Days:    make([]*Day, 7),
	}
	for i := 0; i < 7; i++ {
		r.Days[i] = NewDay(i)
	}
	return r
}

func (w *Week) Set(dayOfWeek time.Weekday, day *Day) {
	w.Days[dayOfWeek] = day
}

func (w *Week) Today() *Day {
	return w.Days[time.Now().Weekday()]
}

func (w *Week) ToJSON() string {
	jbyte, err := json.Marshal(w)
	if err == nil {
		return string(jbyte)
	} else {
		return ""
	}
}
