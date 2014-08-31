package main

import "os"
import "encoding/gob"
import "log"
import "errors"
import "strconv"

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
