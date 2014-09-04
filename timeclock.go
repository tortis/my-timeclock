package main

import "time"

type TimeClock struct {
	store   *WeekStore
	onClock bool
	block   *TimeBlock
}

func NewTimeClock(storefile string) *TimeClock {
	return &TimeClock{
		store:   NewWeekStore(storefile),
		onClock: false,
	}
}

func (c *TimeClock) ClockIn() bool {
	if c.onClock {
		return false
	} else {
		if _, err := c.store.Get(time.Now()); err == nil {
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
		week, _ := c.store.Get(time.Now())
		week.Today().AddTimeBlock(c.block)
		c.store.SaveWeeks();
		c.onClock = false
		return true
	}
	return false
}

func (c *TimeClock) TimeThisWeek() time.Duration {
	if week, err := c.store.Get(time.Now()); err == nil {
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
	if week, err := c.store.Get(time.Now()); err == nil {
		return week.Today().TimeWorked() + c.block.GetDuration()
	} else {
		c.store.Put(NewWeek())
		return 0
	}
}

func (c *TimeClock) TimeOn() time.Duration {
	if c.onClock {
		return c.block.GetDuration()
	} else {
		return 0
	}
}

func (c *TimeClock) GetWeek(t time.Time) *Week {
	if w, err := c.store.Get(t); err == nil {
		return w
	} else {
		c.store.Put(NewSpecificWeek(t))
		w, _ = c.store.Get(t)
		return w
	}
}
