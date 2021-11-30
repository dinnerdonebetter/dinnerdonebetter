package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Server started at port 80")
	log.Fatal(http.ListenAndServe(":80", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handler invoked for request: %q", r.URL.String())
	fmt.Fprintf(w, "Hello, there\n")
}
