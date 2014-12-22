package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

var (
	port       *int
	mongo_url  *string
	mongo_port *string
	mongo_db   *string
	mongo_col  *string
	mongo_user *string
	mongo_pwd  *string
	htdocs_dir *string
	timeStore  *TimeStore
)

func init() {
	port = flag.Int("port", 8080, "Port on which the web application will run.")
	mongo_url = flag.String("dburl", "localhost", "URL of mongodb.")
	mongo_port = flag.String("dbport", "27017", "Port to connect to mongo.")
	mongo_db = flag.String("dbname", "timeclock", "Mongo database the app will use.")
	mongo_col = flag.String("dbcol", "shifts", "Collection to store shifts in.")
	mongo_user = flag.String("dbuser", "", "Database username.")
	mongo_pwd = flag.String("dbpwd", "", "Database password.")
	htdocs_dir = flag.String("htdocs", "htdocs", "Path to the htdocs directory.")
}

func main() {
	flag.Parse()
	var err error
	timeStore, err = OpenTimeStore(*mongo_user, *mongo_pwd, *mongo_url, *mongo_port, *mongo_db, *mongo_col)
	if err != nil {
		log.Fatal("Failed to open time store: ", err)
	}
	defer timeStore.Close()

	http.Handle("/", http.FileServer(http.Dir(*htdocs_dir)))
	http.HandleFunc("/in", handleIn)
	http.HandleFunc("/out", handleOut)
	http.HandleFunc("/toggle", handleToggle)
	http.HandleFunc("/status", handleStatus)
	http.HandleFunc("/shifts", handleShifts)
	http.HandleFunc("/week", handleWeek)
	http.HandleFunc("/createshift", handleCreateShift)
	http.HandleFunc("/editshift", handleEditShift)
	http.HandleFunc("/deleteshift", handleDeleteShift)
	log.Println("Starting server on port " + strconv.Itoa(*port) + ".")
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

func handleToggle(w http.ResponseWriter, req *http.Request) {
	if timeStore.GetState() {
		handleOut(w, req)
	} else {
		handleIn(w, req)
	}
}

func handleStatus(w http.ResponseWriter, req *http.Request) {
	if timeStore.GetState() {
		fmt.Fprintf(w, "true")
	} else {
		fmt.Fprintf(w, "false")
	}
}

func handleShifts(w http.ResponseWriter, req *http.Request) {
	fromString := req.FormValue("from")
	toString := req.FormValue("to")
	from, err := strconv.ParseInt(fromString, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid 'from' date.")
		return
	}
	to, err := strconv.ParseInt(toString, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid 'to' date.")
		return
	}
	shifts, err := timeStore.GetShifts(from, to)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, err.Error())
		return
	}
	jss, _ := json.Marshal(&shifts)
	fmt.Fprint(w, string(jss))
}

func handleWeek(w http.ResponseWriter, req *http.Request) {
	sundayString := req.FormValue("sunday")
	sunday, err := strconv.ParseInt(sundayString, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid sunday date.")
		return
	}
	week, err := timeStore.GetWeek(sunday)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Failed to get week: ", err)
		return
	}
	jss, _ := json.Marshal(&week)
	fmt.Fprint(w, string(jss))
}

func handleCreateShift(w http.ResponseWriter, req *http.Request) {
	onString := req.FormValue("on")
	offString := req.FormValue("off")
	on, err := strconv.ParseInt(onString, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid on time")
		return
	}
	off, err := strconv.ParseInt(offString, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid off time")
		return
	}
	err = timeStore.CreateShift(on, off)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Server failed to create the shift.")
		return
	}
	fmt.Fprint(w, "OK")
}

func handleEditShift(w http.ResponseWriter, req *http.Request) {
	hex_id := req.FormValue("id")
	onString := req.FormValue("on")
	offString := req.FormValue("off")
	var on, off int64
	var err error
	// If on was not provided, set to -1 to indicate no change
	if onString == "" {
		on = -1
	} else {
		on, err = strconv.ParseInt(onString, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Invalid on time")
			return
		}
	}

	// If off was not provided, set to -1 to indicate no change
	if offString == "" {
		off = -1
	} else {
		off, err = strconv.ParseInt(offString, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Invalid off time")
			return
		}
	}
	err = timeStore.ModifyShift(hex_id, on, off)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Failed to modify shift: ", err)
		return
	}
	fmt.Fprint(w, "OK")
}

func handleDeleteShift(w http.ResponseWriter, req *http.Request) {
	hex_id := req.FormValue("id")
	err := timeStore.DeleteShift(hex_id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Failed to delete shift: ", err)
		return
	}
	fmt.Fprint(w, "OK")
}
