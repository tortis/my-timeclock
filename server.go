package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

var (
	port       *int
	store_file *string
	htdocs_dir *string
	timeStore  *TimeStore
)

func init() {
	port = flag.Int("port", 7171, "Port on which the web application will run.")
	htdocs_dir = flag.String("htdocs", "htdocs", "Path to the htdocs directory.")
}

func main() {
	flag.Parse()
	var err error
	timeStore, err = OpenTimeStore("david", "asdf", "localhost", "27017", "timeclock", "shifts")
	if err != nil {
		log.Fatal("Failed to open time store: ", err)
	}
	defer timeStore.Close()

	http.Handle("/", http.FileServer(http.Dir(*htdocs_dir)))
	http.HandleFunc("/in", handleIn)
	http.HandleFunc("/out", handleOut)
	http.HandleFunc("/status", handleStatus)
	fmt.Println("Starting server on port " + strconv.Itoa(*port) + ".")
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
}

func handleIn(w http.ResponseWriter, req *http.Request) {
	err := timeStore.ClockIn()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, err.Error())
		return
	}
	fmt.Fprintf(w, "OK")
}

func handleOut(w http.ResponseWriter, req *http.Request) {
	err := timeStore.ClockOut()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, err.Error())
		return
	}
	fmt.Fprintf(w, "OK")
}

func handleStatus(w http.ResponseWriter, req *http.Request) {
	if timeStore.GetState() {
		fmt.Fprintf(w, "true")
	} else {
		fmt.Fprintf(w, "false")
	}
}

func handleShifts(w http.RequestWriter, req *http.Request) {

}
