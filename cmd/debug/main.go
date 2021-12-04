package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

const (
	dbURL = "postgres://prixfixe_api:password@api-database.cluster-ctj4wxgujo7g.us-east-1.rds.amazonaws.com:5432/prixfixe"
)

var (
	db *sql.DB
)

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Server started at port 80")

	var err error
	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	log.Fatal(http.ListenAndServe(":80", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handler invoked for request: %s %s", r.Method, r.URL.String())

	var version string
	row := db.QueryRow("SELECT VERSION()")
	if err := row.Scan(&version); err != nil {
		log.Printf("error querying database: %v", err)
		fmt.Fprintf(w, "Error: %s", err)
	}

	fmt.Fprintf(w, "Hello, there from version %s!", version)
}
