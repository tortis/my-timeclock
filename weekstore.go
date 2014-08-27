package main

type WeekStore struct {
	Weeks map[int]Week
	// The integer key to get a Week is (year*54+ISOWeek)
	fname string
}

type weekRecord struct {
	Key   int
	Value Week
}

func NewWeekStore(fname string) *WeekStore {
	return &WeekStore{
		Weeks: make(map[int]Week),
		fname: fname,
	}
}
