package main

import (
	"fmt"
	"log"
	"net/http"
)

/*
CREATE TABLE IF NOT EXISTS stuff (
	"name" TEXT NOT NULL,
	"description" TEXT NOT NULL
);
*/

const (
// dbURL = "postgres://prixfixe_api"
)

// var (
// 	db *sql.DB
// )

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Server started at port 80")

	// var err error
	// db, err = sql.Open("postgres", dbURL)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	log.Fatal(http.ListenAndServe(":80", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handler invoked for request: %s %s", r.Method, r.URL.String())

	// age := 21
	// rows, err := db.Query("SELECT  FROM users WHERE age = $1", age)

	fmt.Fprintf(w, "Hello, there!")
}
