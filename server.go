package main

import "net/http"
import "fmt"
import "time"
import "log"

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

	http.Handle("/", http.FileServer(http.Dir("htdocs")))
	http.HandleFunc("/in/", handleIn)
	http.HandleFunc("/out/", handleOut)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
