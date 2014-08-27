package main

import "net/http"
import "fmt"
import "time"
import "log"
import "strconv"

func handleIn(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "In!")
}

func handleOut(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Out!")
}

func main() {
	b1 := NewTimeBlock()
	b2 := NewTimeBlock()
	time.Sleep(time.Second * 5)
	b1.End()
	b2.End()
	d := NewDay()
	d.AddTimeBlock(b1)
	d.AddTimeBlock(b2)
	println(d.TimeWorked())

	weekStore := NewWeekStore("weeks.gob")
	println("Weeks loaded from store: "+strconv.Itoa(len(weekStore.weeks)))
	weekNow := NewWeek()
	weekNow.Set(time.Monday, d)
	weekStore.Put(weekNow)
	weekStore.SaveWeeks()

	http.Handle("/", http.FileServer(http.Dir("htdocs")))
	http.HandleFunc("/in/", handleIn)
	http.HandleFunc("/out/", handleOut)
	log.Fatal(http.ListenAndServe(":88", nil))
}
