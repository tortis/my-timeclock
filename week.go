package main

import "time"
import "encoding/json"

type Week struct {
	Year    int
	WeekNum int
	Days    []*Day
	Hours   float64
}

func NewWeek() *Week {
	year, week := time.Now().ISOWeek()
	r := &Week{
		Year:    year,
		WeekNum: week,
		Days:    make([]*Day, 7),
		Hours:   0.0,
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
		Hours:   0.0,
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

func (w *Week) sumHours() {
	var h float64 = 0.0
	for _, d := range w.Days {
		h += d.Hours
	}
	w.Hours = h
}

func (w *Week) ToJSON() string {
	w.sumHours()
	jbyte, err := json.Marshal(w)
	if err == nil {
		return string(jbyte)
	} else {
		return ""
	}
}
