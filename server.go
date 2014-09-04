package main

import "net/http"
import "fmt"
import "log"
import "strconv"
import "os"
import "os/signal"
import "encoding/json"
import "time"
import "flag"

var timeClock *TimeClock
var port *int
var store_file *string
var htdocs_dir *string

func init() {
	port = flag.Int("port", 713, "Port on which the web application will run.")
	store_file = flag.String("store", "weeks.gob", "File where time data will be stored.")
	htdocs_dir = flag.String("htdocs", "htdocs", "Path to the htdocs directory.")


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
	jbyte, err := json.Marshal(Status{OnClock: timeClock.onClock, TimeOn: timeClock.TimeOn().Hours()})
	println(string(jbyte))
	if err == nil {
		w.Write(jbyte)
	} else {
		fmt.Fprintf(w, "error")
	}
}

func handleWeek(w http.ResponseWriter, req *http.Request) {
	year, _ := strconv.Atoi(req.FormValue("year"))
	m, _ := strconv.Atoi(req.FormValue("month"))
	month := time.Month(m)
	day, _ := strconv.Atoi(req.FormValue("day"))
	fmt.Fprintf(w, timeClock.GetWeek(time.Date(year, month, day, 0, 0, 0, 0, time.Now().Location())).ToJSON())
}

func main() {
	flag.Parse()
	timeClock = NewTimeClock(*store_file)

	http.Handle("/", http.FileServer(http.Dir(*htdocs_dir)))
	http.HandleFunc("/in/", handleIn)
	http.HandleFunc("/out/", handleOut)
	http.HandleFunc("/week/", handleWeek)
	http.HandleFunc("/status/", handleStatus)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
}
