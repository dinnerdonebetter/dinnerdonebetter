package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/server/http/build"
)

const (
	/* #nosec G101 */
	debugCookieHashKey = "HEREISA32CHARSECRETWHICHISMADEUP"
	/* #nosec G101 */
	debugCookieBlockKey = "DIFFERENT32CHARSECRETTHATIMADEUP"
)

func main() {
	ctx := context.Background()

	rawCfg, err := os.ReadFile("environments/dev/config_files/service-config.json")
	if err != nil {
		log.Fatal(err)
	}

	var cfg *config.InstanceConfig
	if err = json.Unmarshal(rawCfg, &cfg); err != nil {
		log.Fatal(err)
	}

	cfg.Neutralize()

	// build our server struct.
	srv, err := build.Build(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}

	for _, route := range srv.Router().Routes() {
		fmt.Printf("%s %s\n", route.Method, route.Path)
	}
}
