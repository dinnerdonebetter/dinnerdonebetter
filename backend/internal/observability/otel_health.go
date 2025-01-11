package observability

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

// EnsureOtelCollectorIsUp checks that a server is up and doesn't return until it's certain one way or the other.
func EnsureOtelCollectorIsUp(ctx context.Context, address string) {
	var (
		isDown           = true
		interval         = time.Second
		maxAttempts      = 50
		numberOfAttempts = 0
	)

	address = fmt.Sprintf("%s:13313/health", address)

	for isDown {
		if !otelIsUp(ctx, address) {
			log.Printf("waiting %s before pinging %s again", interval, address)
			time.Sleep(interval)

			numberOfAttempts++
			if numberOfAttempts >= maxAttempts {
				log.Fatal("Maximum number of attempts made, something's gone awry")
			}
		} else {
			isDown = false
		}
	}
}

// otelIsUp can check if an instance of our server is alive.
func otelIsUp(ctx context.Context, address string) bool {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, address, http.NoBody)
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}

	if err = res.Body.Close(); err != nil {
		log.Println("error closing body")
	}

	if res.StatusCode != http.StatusOK {
		log.Println("expected status code 200, got: ", res.Status)
	}

	return res.StatusCode == http.StatusOK
}
