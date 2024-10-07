package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

const (
	hitThreshold = 1
)

var hitCount uint8

func main() {
	shutdownChan := make(chan bool)
	go func() {
		<-shutdownChan
		os.Exit(0)
	}()

	http.HandleFunc("/completed/{test_name}", func(res http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodPost {
			hitCount++
			log.Printf("%s completed\n", req.PathValue("test_name"))

			if hitCount >= hitThreshold {
				go func() {
					<-time.After(time.Second)
					log.Println("shutting down")
					shutdownChan <- true
				}()
			}
		}

		res.WriteHeader(http.StatusOK)
	})

	server := &http.Server{
		Addr:              ":9999",
		ReadHeaderTimeout: 3 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
