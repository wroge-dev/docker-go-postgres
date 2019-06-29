package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func connect() (db *sql.DB, err error) {
	database := os.Getenv("POSTGRES_DB")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	connStr := fmt.Sprintf("host=db port=5432 user=%s dbname=%s password=%s sslmode=disable", user, database, password)
	return sql.Open("postgres", connStr)
}

func query(q string) (rows *sql.Rows, err error) {
	db, err := connect()
	if err != nil {
		return
	}
	return db.Query(q)
}

func send(w http.ResponseWriter, o interface{}) {
	json, err := json.Marshal(o)
	if err == nil {
		_, err = w.Write(json)
	}
}

func sendError(w http.ResponseWriter, errMsg string, status int) {
	w.WriteHeader(status)
	send(w, errMsg)
}

// Name struct for JSON response
type Name struct {
	ID   int
	Name string
}

func getName(w http.ResponseWriter, r *http.Request) {
	var name Name
	rows, err := query("SELECT id, nm FROM names WHERE id=" + mux.Vars(r)["id"] + ";")
	if err != nil {
		sendError(w, "500 - Internal Server Error", http.StatusInternalServerError)
		return
	}
	if rows.Next() {
		var id int
		var nm string
		if err := rows.Scan(&id, &nm); err != nil {
			sendError(w, "500 - Internal Server Error", http.StatusInternalServerError)
			return
		}
		name = Name{id, nm}
	} else {
		sendError(w, "400 - Bad Request", http.StatusBadRequest)
		return
	}
	send(w, name)
}

func getNames(w http.ResponseWriter, r *http.Request) {
	var names []Name
	var count int
	rows, err := query("SELECT id, nm FROM names;")
	if err != nil {
		sendError(w, "500 - Internal Server Error", http.StatusInternalServerError)
	}
	for rows.Next() {
		var id int
		var nm string
		if err := rows.Scan(&id, &nm); err != nil {
			sendError(w, "500 - Internal Server Error", http.StatusInternalServerError)
		}
		names = append(names, Name{id, nm})
		count++
	}
	if len(names) == 0 {
		sendError(w, "400 - Bad Request", http.StatusInternalServerError)
	}
	o := struct {
		Count int
		Names []Name
	}{
		count,
		names,
	}
	send(w, o)
}

func redirect(w http.ResponseWriter, req *http.Request) {
	target := "https://" + req.Host + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}
	log.Printf("redirect to: %s", target)
	http.Redirect(w, req, target, http.StatusTemporaryRedirect)
}

func main() {
	go http.ListenAndServe(":80", http.HandlerFunc(redirect))
	r := mux.NewRouter()
	r.Handle("/", http.FileServer(http.Dir("./html")))
	r.HandleFunc("/name/{id}", getName)
	r.HandleFunc("/names", getNames)
	log.Fatal(http.ListenAndServeTLS(":443", "./cert/server.crt", "./cert/server.key", r))
}
