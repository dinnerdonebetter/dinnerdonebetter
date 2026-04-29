package main

import (
	"dagger/ci/internal/dagger"
)

// Image versions are pinned here so changes are visible in one place.
const (
	goImage         = "golang:1.26-alpine"
	nodeImage       = "node:20-bookworm-slim"
	npmVersion      = "11"
	golangciImage   = "golangci/golangci-lint:v2.10.1"
	shellcheckImage = "koalaman/shellcheck-alpine:stable"
)

// Cache volume names. One per concern, shared across functions.
const (
	cacheGoMod      = "ddb-go-mod"
	cacheGoBuild    = "ddb-go-build"
	cacheGolangci   = "ddb-golangci-lint"
	cacheNpm        = "ddb-npm"
)

// baseGo returns a golang container with the toolchain CI needs (make, bash,
// git, gcc/musl-dev for `-race`) and the Go module + build caches mounted.
//
// All backend Go functions start from this; per-function adjustments (like
// docker.sock for testcontainers) layer on top.
func (m *Ci) baseGo() *dagger.Container {
	return dag.Container().
		From(goImage).
		WithExec([]string{"apk", "add", "--no-cache", "make", "bash", "git", "gcc", "musl-dev"}).
		WithEnvVariable("CGO_ENABLED", "1").
		WithEnvVariable("GOPATH", "/root/go").
		WithEnvVariable("GOCACHE", "/root/.cache/go-build").
		WithEnvVariable("GOFLAGS", "-p=4").
		WithMountedCache("/root/go/pkg/mod", dag.CacheVolume(cacheGoMod)).
		WithMountedCache("/root/.cache/go-build", dag.CacheVolume(cacheGoBuild))
}

// baseNode returns a node:20 container with npm@11 and git installed,
// and the npm cache mounted. Suitable for any frontend function.
//
// git is required because `npm run lint` invokes `pin:check`, which does
// `git diff --exit-code` to enforce that pinned dependencies haven't drifted.
func (m *Ci) baseNode() *dagger.Container {
	return dag.Container().
		From(nodeImage).
		WithExec([]string{"apt-get", "update"}).
		WithExec([]string{"apt-get", "install", "-y", "--no-install-recommends", "git", "ca-certificates"}).
		WithExec([]string{"rm", "-rf", "/var/lib/apt/lists/*"}).
		WithExec([]string{"npm", "install", "-g", "npm@" + npmVersion}).
		WithMountedCache("/root/.npm", dag.CacheVolume(cacheNpm))
}
