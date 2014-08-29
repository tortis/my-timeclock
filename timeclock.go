package main

import "time"

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

func (c *TimeClock) ClockIn() bool {
	if c.onClock {
		return false
	} else {
		if _, err := c.store.Get(time.Now().ISOWeek()); err == nil {
			c.block = NewTimeBlock()
			c.onClock = true
		} else {
			c.store.Put(NewWeek())
			println("Adding this week to week store.")
			c.block = NewTimeBlock()
			c.onClock = true
		}
		return true
	}

}

func (c *TimeClock) ClockOut() bool {
	if c.onClock {
		c.block.End()
		week, _ := c.store.Get(time.Now().ISOWeek())
		week.Today().AddTimeBlock(c.block)
		c.onClock = false
		return true
	}
	return false
}

func (c *TimeClock) TimeThisWeek() time.Duration {
	if week, err := c.store.Get(time.Now().ISOWeek()); err == nil {
		var r time.Duration = 0
		for _, day := range week.Days {
			r += day.TimeWorked()
		}
		return r + c.block.GetDuration()
	} else {
		c.store.Put(NewWeek())
		return 0
	}
}

func (c *TimeClock) TimeToday() time.Duration {
	if week, err := c.store.Get(time.Now().ISOWeek()); err == nil {
		return week.Today().TimeWorked() + c.block.GetDuration()
	} else {
		c.store.Put(NewWeek())
		return 0
	}
}

func (c *TimeClock) GetWeek(year, week int) *Week {
	if w, err := c.store.Get(year, week); err == nil {
		return w
	} else {
		c.store.Put(NewSpecificWeek(year, week))
		w, _ = c.store.Get(year, week)
		return w
	}
}
