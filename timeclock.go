package main

import "time"
import "encoding/json"

const (
	STOREFILE = "weeks.gob"
)

type TimeClock struct {
	store   *WeekStore
	onClock bool
	block   *TimeBlock
}

func NewTimeClock() *TimeClock {
	return &TimeClock{
		store:   NewWeekStore(STOREFILE),
		onClock: false,
	}
}

func (c *TimeClock) ClockIn() {
	if _, err := c.store.Get(time.Now().ISOWeek()); err == nil {
		if !c.onClock {
			c.block = NewTimeBlock()
			c.onClock = true
		}
	} else {
		c.store.Put(NewWeek())
		println("Adding this week to week store.")
		if !c.onClock {
			c.block = NewTimeBlock()
			c.onClock = true
		}
	}
}

func (c *TimeClock) ClockOut() {
	if c.onClock {
		c.block.End()
		week, _ := c.store.Get(time.Now().ISOWeek())
		week.Today().AddTimeBlock(c.block)
		c.store.SaveWeeks()
	}
}

func (c *TimeClock) TimeThisWeek() time.Duration {
	if week, err := c.store.Get(time.Now().ISOWeek()); err == nil {
		var r time.Duration = 0
		for _, day := range week.Days {
			r += day.TimeWorked()
		}
		return r
	} else {
		c.store.Put(NewWeek())
		return 0
	}
}

func (c *TimeClock) TimeToday() time.Duration {
	if week, err := c.store.Get(time.Now().ISOWeek()); err == nil {
		return week.Today().TimeWorked()
	} else {
		c.store.Put(NewWeek())
		return 0
	}
}

func (c *TimeClock) JSONWeek(year, week int) string {
	if w, err := c.store.Get(year, week); err == nil {
		if jbyte, err := json.Marshal(w); err == nil {
			return string(jbyte)
		} else {
			return ""
		}
	} else {
		c.store.Put(NewWeek())
		w, _ = c.store.Get(year, week)
		if jbyte, err := json.Marshal(w); err == nil {
			return string(jbyte)
		} else {
			return ""
		}
	}
}

func (c *TimeClock) GetWeek(year, week int) *Week {
	if w, err := c.store.Get(year, week); err == nil {
		return w
	} else {
		c.store.Put(NewWeek())
		w, _ = c.store.Get(year, week)
		return w
	}
}
