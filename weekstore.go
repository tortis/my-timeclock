package main

import "os"
import "encoding/gob"
import "log"
import "errors"
import "strconv"
import "time"

type WeekStore struct {
	// The integer key to get a Week is (year + month(1-12) * 3000 + day(1-31) * 37000)
	weeks map[int]*Week
	fname string
}

func NewWeekStore(fname string) *WeekStore {
	r := &WeekStore{
		weeks: make(map[int]*Week),
		fname: fname,
	}
	if err := r.LoadWeeks(); err != nil {
		println("Could not load data from given weekstore. Proceeding with empty store.", err)
	}
	println("Loaded " + strconv.Itoa(len(r.weeks)) + " weeks from file.")
	return r
}

func (ws *WeekStore) LoadWeeks() error {
	f, err := os.OpenFile(ws.fname, os.O_RDONLY, 0644)
	defer f.Close()
	if err != nil {
		return err
	}

	if _, err := f.Seek(0, 0); err != nil {
		return err
	}

	d := gob.NewDecoder(f)
	if err = d.Decode(&ws.weeks); err != nil {
		return err
	}
	return nil
}

func (ws *WeekStore) SaveWeeks() error {
	f, err := os.OpenFile(ws.fname, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Println("WARNING: Unable to save weeks to store. Some data may have been lost.", err)
		return err
	}
	defer f.Close()
	e := gob.NewEncoder(f)
	err = e.Encode(ws.weeks)
	if err != nil {
		log.Println("Falied to save week records to the store. Some data may have been lost.", err)
		return err
	}
	return nil
}

func (ws *WeekStore) Get(t time.Time) (*Week, error) {
	// Convert the given time to Monday
	m := t.AddDate(0, 0, 1-int(t.Weekday()))
	if w, present := ws.weeks[m.Year()+int(m.Month())*3000+(1+m.Day())*37000]; present {
		return w, nil
	} else {
		return nil, errors.New("The requested week does not exist in the store.")
	}
}

func (ws *WeekStore) Put(week *Week) error {
	if _, present := ws.weeks[week.weekKey()]; present {
		return errors.New("There is already a value in the week store with the provided key.")
	}
	ws.weeks[week.weekKey()] = week
	return nil
}
