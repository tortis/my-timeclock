package main

import "net/http"
import "fmt"

func handleIn(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "In!")
}

func handleOut(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Out!")
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("htdocs")))
	http.HandleFunc("/in/", handleIn)
	http.HandleFunc("/out/", handleOut)
	http.ListenAndServe(":8080", nil)
}
