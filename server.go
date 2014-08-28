package main

import "net/http"
import "fmt"
import "log"
import "time"
import "strconv"
import "os"
import "os/signal"

var timeClock *TimeClock

func init() {
	timeClock = NewTimeClock()

	// Register to receive interrupt signal.
	// Write to the weekstore when signal is
	// received.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		timeClock.store.SaveWeeks()
		os.Exit(0)
	}()

}

func handleIn(w http.ResponseWriter, req *http.Request) {
	if timeClock.ClockIn() {
		fmt.Fprintf(w, "true")
	} else {
		fmt.Fprintf(w, "false")
	}
}

func handleOut(w http.ResponseWriter, req *http.Request) {
	if timeClock.ClockOut() {
		fmt.Fprintf(w, "true")
	} else {
		fmt.Fprintf(w, "false")
	}
}

func handleStatus(w http.ResponseWriter, req *http.Request) {
	if timeClock.onClock {
		fmt.Fprintf(w, "true")
	} else {
		fmt.Fprintf(w, "false")
	}
}

func handleWeek(w http.ResponseWriter, req *http.Request) {
	year, _ := strconv.Atoi(req.FormValue("year"))
	week, _ := strconv.Atoi(req.FormValue("week"))
	fmt.Fprintf(w, timeClock.GetWeek(year, week).ToJSON())
}

func main() {
	timeClock.ClockIn()
	time.Sleep(time.Second * 10)
	timeClock.ClockOut()
	println(timeClock.TimeToday())
	println(timeClock.TimeThisWeek())

	println(timeClock.GetWeek(time.Now().ISOWeek()).ToJSON())

	http.Handle("/", http.FileServer(http.Dir("htdocs")))
	http.HandleFunc("/in/", handleIn)
	http.HandleFunc("/out/", handleOut)
	http.HandleFunc("/week/", handleWeek)
	http.HandleFunc("/status/", handleStatus)
	log.Fatal(http.ListenAndServe(":88", nil))
}
