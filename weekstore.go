package main

import "os"
import "encoding/gob"
import "io"
import "log"
import "errors"

type WeekStore struct {
	// The integer key to get a Week is (year*54+ISOWeek)
	weeks map[int]*Week
	fname string
}

type weekRecord struct {
	Key   int
	Value *Week
}

func NewWeekStore(fname string) *WeekStore {
	r := &WeekStore{
		weeks: make(map[int]*Week),
		fname: fname,
	}
	if r.LoadWeeks() != nil {
		println("Could not load data from given weekstore. Proceeding with empty store.")
	}
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
	err = nil
	for err == nil {
		var r weekRecord
		if err = d.Decode(&r); err == nil {
			ws.weeks[r.Key] = r.Value
		}
	}
	if err == io.EOF {
		return nil
	}

	return err
}

func (ws *WeekStore) SaveWeeks() error {
	f, err := os.OpenFile(ws.fname, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Println("WARNING: Unable to save weeks to store. Some data may have been lost.", err)
		return err
	}
	defer f.Close()
	e := gob.NewEncoder(f)
	for k, v := range ws.weeks {
		err := e.Encode(weekRecord{Key: k, Value: v})
		if err != nil {
			log.Println("Failed to save a week record to the store. Some data may havve been lost.", err)
		}
	}
	return nil
}

func (ws *WeekStore) Get(year, week int) (*Week, error) {
	if w, present := ws.weeks[year*34+week]; present {
		return w, nil
	} else {
		return nil, errors.New("The requested week does not exist in the store.")
	}
}

func (ws *WeekStore) Put(week *Week) error {
	if _, present := ws.weeks[week.Year*34+week.WeekNum]; present {
		return errors.New("There is already a value in the week store with the provided key.")
	}
	ws.weeks[week.Year*34+week.WeekNum] = week
	return nil
}
