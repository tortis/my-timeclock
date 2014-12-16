package main

import (
	"errors"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Shift struct {
	Id     bson.ObjectId `bson:"_id,omitempty"`
	On     int64         `bson:"on,omitempty"`
	Off    int64         `bson:"off"`
	Active bool          `bson:"active"`
	User   bson.ObjectId `bson:user"user,omitempty"`
}

func NewShift() *Shift {
	return &Shift{
		On:     time.Now().Unix(),
		Active: true,
	}
}

func FullShift(on, off int64) (*Shift, error) {
	if off < on {
		return nil, errors.New("Off time is before on time.")
	}
	return &Shift{
		On:     on,
		Off:    off,
		Active: false,
	}, nil
}
