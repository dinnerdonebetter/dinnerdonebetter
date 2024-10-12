package golang

import (
	"embed"
	"fmt"
)

var (
	//go:embed statics/*.go
	staticFiles embed.FS
)

func fetchStaticFile(name string) string {
	file, err := staticFiles.ReadFile(fmt.Sprintf("statics/%s.go", name))
	if err != nil {
		panic(err)
	}

	return string(file)
}

var baseFiles = map[string]string{
	"client":                fetchStaticFile("client"),
	"client_test":           fetchStaticFile("client_test"),
	"client_options":        fetchStaticFile("client_options"),
	"client_options_test":   fetchStaticFile("client_options_test"),
	"doc":                   fetchStaticFile("doc"),
	"errors":                fetchStaticFile("errors"),
	"helpers":               fetchStaticFile("helpers"),
	"helpers_test":          fetchStaticFile("helpers_test"),
	"image_uploading":       fetchStaticFile("image_uploading"),
	"roundtripper":          fetchStaticFile("roundtripper"),
	"roundtripper_test":     fetchStaticFile("roundtripper_test"),
	"test_helpers_test":     fetchStaticFile("test_helpers_test"),
	"mock_read_closer_test": fetchStaticFile("mock_read_closer_test"),
}
