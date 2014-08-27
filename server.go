package main

import "net/http"
import "fmt"
import "log"
import "time"

func handleIn(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "In!")
}

func handleOut(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Out!")
}

func main() {
	timeClock := NewTimeClock()
	timeClock.ClockIn()
	time.Sleep(time.Second * 10)
	timeClock.ClockOut()
	println(timeClock.TimeToday())
	println(timeClock.TimeThisWeek())

	http.Handle("/", http.FileServer(http.Dir("htdocs")))
	http.HandleFunc("/in/", handleIn)
	http.HandleFunc("/out/", handleOut)
	log.Fatal(http.ListenAndServe(":88", nil))
}
