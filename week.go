package main

import "time"
import "encoding/json"

type Week struct {
	MondayDate time.Time
	Days       []*Day
	Hours      float64
}

func NewWeek() *Week {
	now := time.Now()
	monday := now.AddDate(0, 0, 1-int(now.Weekday()))
	r := &Week{
		MondayDate: monday,
		Days:       make([]*Day, 7),
		Hours:      0.0,
	}
	for i := 0; i < 7; i++ {
		r.Days[i] = NewDay(i)
	}
	return r
}

func NewSpecificWeek(t time.Time) *Week {
	monday := t.AddDate(0, 0, 1-int(t.Weekday()))
	r := &Week{
		MondayDate: monday,
		Days:       make([]*Day, 7),
		Hours:      0.0,
	}
	for i := 0; i < 7; i++ {
		r.Days[i] = NewDay(i)
	}
	return r
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

func (w *Week) weekKey() int {
	return w.MondayDate.Year() + int(w.MondayDate.Month())*3000 + (1+w.MondayDate.Day())*37000
}
