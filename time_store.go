package main

import (
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

type TimeStore struct {
	db_session *mgo.Session
	db_name    string
	col_name   string
	state      bool
}

func OpenTimeStore(mgo_user, mgo_pwd, mgo_url, mgo_port, mgo_db, col string) (*TimeStore, error) {
	ts := &TimeStore{
		db_name:  mgo_db,
		col_name: col,
	}
	var err error

	// Dial the mongo session
	ts.db_session, err = mgo.Dial(mgo_user + ":" + mgo_pwd + "@" + mgo_url + ":" + mgo_port + "/" + mgo_db)
	if err != nil {
		log.Println("Failed to Dial")
		return nil, err
	}

	// Verify the database and collections exist
	err = ts.verify_db()
	if err != nil {
		return nil, err
	}
	ts.state = ts.get_state()
	return ts, nil

}

func (ts *TimeStore) ClockIn() error {
	if ts.state == true {
		return errors.New("Already clocked in.")
	}

	s := NewShift()
	shifts := ts.db_session.DB(ts.db_name).C(ts.col_name)
	err := shifts.Insert(s)
	if err != nil {
		return err
	}

	ts.state = true
	return nil
}

func (ts *TimeStore) Close() {
	ts.db_session.Close()
}

func (ts *TimeStore) ClockOut() error {
	if ts.state == false {
		return errors.New("Not clocked on.")
	}

	shifts := ts.db_session.DB(ts.db_name).C(ts.col_name)
	err := shifts.Update(bson.M{"active": true}, bson.M{"$set":bson.M{"off": time.Now().Unix(), "active": false}})
	if err != nil {
		return err
	}
	ts.state = false
	return nil
}

func (ts *TimeStore) GetState() bool {
	return ts.state
}

func (ts *TimeStore) GetShifts(from, to int64) ([]*Shift, error) {
	shifts := ts.db_session.DB(ts.db_name).C(ts.col_name)
	q := shifts.Find(bson.M{"$or": []bson.M{{"on": bson.M{"$gte": from}}, bson.M{"off": bson.M{"$lt": to}}}}) // da faq
	num, err := q.Count()
	if err != nil {
		return nil, err
	}
	ss := make([]*Shift, num)
	err = q.All(&ss)
	if err != nil {
		return nil, err
	}
	return ss, nil
}

func (ts *TimeStore) GetShift(id bson.ObjectId) (*Shift, error) {
	s := &Shift{}
	shifts := ts.db_session.DB(ts.db_name).C(ts.col_name)
	err := shifts.FindId(id).One(s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (ts *TimeStore) DeleteShift(id bson.ObjectId) error {
	shifts := ts.db_session.DB(ts.db_name).C(ts.col_name)
	err := shifts.RemoveId(id)
	if err != nil {
		return err
	}
	return nil
}

func (ts *TimeStore) ModifyShift(id bson.ObjectId, on, off int64) error {
	if off < on {
		return errors.New("Off time is before on time.")
	}
	shifts := ts.db_session.DB(ts.db_name).C(ts.col_name)
	err := shifts.UpdateId(id, bson.M{"$set":bson.M{"on": on, "off": off, "active": false}})
	if err != nil {
		return err
	}
	return nil
}

func (ts *TimeStore) CreateShift(on, off int64) error {
	s, err := FullShift(on, off)
	if err != nil {
		return err
	}
	shifts := ts.db_session.DB(ts.db_name).C(ts.col_name)
	err = shifts.Insert(s)
	if err != nil {
		return err
	}
	return nil
}

func (ts *TimeStore) verify_db() error {
	db := ts.db_session.DB(ts.db_name)
	if db.Name != ts.db_name {
		return errors.New("The database: " + ts.db_name + " does not exist.")
	}
	cnames, err := db.CollectionNames()
	if err != nil {
		return err
	}
	for _, name := range cnames {
		if name == ts.col_name {
			return nil
		}
	}
	return errors.New("The collection: " + ts.col_name + " does not exist.")
}

func (ts *TimeStore) get_state() bool {
	shifts := ts.db_session.DB(ts.db_name).C(ts.col_name)
	count, err := shifts.Find(bson.M{"active": true}).Count()
	if err != nil {
		log.Println("get_state() active Count failed:", err)
		return false
	}
	if count > 0 {
		return true
	} else {
		return false
	}
}
