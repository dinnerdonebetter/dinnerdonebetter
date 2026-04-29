package main

import (
	"context"

	"dagger/ci/internal/dagger"
)

// Common +ignore list for source uploads; CI runners are clean but local dev
// machines have heavy gitignored directories that would otherwise be uploaded
// on every call.
//
// Used as the +ignore directive value on every public function's source param.

// backendBase mounts the source at /src and chdirs to /src/backend.
func (m *Ci) backendBase(source *dagger.Directory) *dagger.Container {
	return m.baseGo().
		WithMountedDirectory("/src", source).
		WithWorkdir("/src/backend")
}

// BackendShellcheck runs shellcheck against every backend/scripts/*.sh file,
// mirroring backend's `make shellcheck` target without nested Docker.
func (m *Ci) BackendShellcheck(
	ctx context.Context,
	// +ignore=["frontend/node_modules", "backend/vendor", "ios/build", "**/node_modules"]
	source *dagger.Directory,
) (string, error) {
	return dag.Container().
		From(shellcheckImage).
		WithMountedDirectory("/workdir", source).
		WithWorkdir("/workdir").
		WithExec([]string{"sh", "-c", `
set -eu
found=0
for f in backend/scripts/*.sh; do
  [ -e "$f" ] || continue
  found=1
  echo "Checking $f..."
  shellcheck --source-path=SCRIPTDIR -x "$f"
done
if [ "$found" = "0" ]; then
  echo "No shell scripts found in backend/scripts"
fi
echo "All shell scripts passed shellcheck!"
`}).
		Stdout(ctx)
}

// BackendBuild runs `make build` (which builds every backend Go package).
func (m *Ci) BackendBuild(
	ctx context.Context,
	// +ignore=["frontend/node_modules", "backend/vendor", "ios/build", "**/node_modules"]
	source *dagger.Directory,
) (string, error) {
	return m.backendBase(source).
		WithExec([]string{"make", "vendor"}).
		WithExec([]string{"make", "build"}).
		Stdout(ctx)
}

// BackendUnitTests runs `make test` (race + shuffle, excludes integration/cmd/generated).
func (m *Ci) BackendUnitTests(
	ctx context.Context,
	// +ignore=["frontend/node_modules", "backend/vendor", "ios/build", "**/node_modules"]
	source *dagger.Directory,
) (string, error) {
	return m.backendBase(source).
		WithExec([]string{"make", "test"}).
		Stdout(ctx)
}

// BackendFormatCheck mirrors backend_formatting.yaml: fail if `gofmt -l` finds
// any file outside of vendor/ that needs formatting.
func (m *Ci) BackendFormatCheck(
	ctx context.Context,
	// +ignore=["frontend/node_modules", "backend/vendor", "ios/build", "**/node_modules"]
	source *dagger.Directory,
) (string, error) {
	return m.backendBase(source).
		WithExec([]string{"sh", "-c", `if [ $(gofmt -l . | grep -Ev '^vendor/' | head -c1 | wc -c) -ne 0 ]; then gofmt -l . | grep -Ev '^vendor/'; exit 1; fi; echo "all backend Go files are gofmt-clean"`}).
		Stdout(ctx)
}

// BackendLint runs golangci-lint against the backend module.
func (m *Ci) BackendLint(
	ctx context.Context,
	// +ignore=["frontend/node_modules", "backend/vendor", "ios/build", "**/node_modules"]
	source *dagger.Directory,
) (string, error) {
	return dag.Container().
		From(golangciImage).
		WithMountedDirectory("/src", source).
		WithWorkdir("/src/backend").
		WithMountedCache("/root/go/pkg/mod", dag.CacheVolume(cacheGoMod)).
		WithMountedCache("/root/.cache/go-build", dag.CacheVolume(cacheGoBuild)).
		WithMountedCache("/root/.cache/golangci-lint", dag.CacheVolume(cacheGolangci)).
		WithEnvVariable("GOPATH", "/root/go").
		WithEnvVariable("GOCACHE", "/root/.cache/go-build").
		WithEnvVariable("GOLANGCI_LINT_CACHE", "/root/.cache/golangci-lint").
		WithExec([]string{"golangci-lint", "run", "--timeout=15m"}).
		Stdout(ctx)
}

// generated-files diff helper: run a make target inside backendBase, then
// `git diff --exit-code -- backend` from the repo root. Vendor changes are
// gitignored and won't trigger the check. Scoping to backend/ keeps the
// check focused on what the generator actually writes.
func (m *Ci) backendGen(ctx context.Context, source *dagger.Directory, makeTargets ...string) (string, error) {
	c := m.backendBase(source).
		WithExec([]string{"make", "vendor"})
	c = c.WithExec(append([]string{"make"}, makeTargets...))
	return c.
		WithWorkdir("/src").
		WithExec([]string{"git", "diff", "--exit-code", "--", "backend"}).
		Stdout(ctx)
}

// BackendGenConfigs regenerates config structs and fails if anything drifted.
func (m *Ci) BackendGenConfigs(
	ctx context.Context,
	// +ignore=["frontend/node_modules", "backend/vendor", "ios/build", "**/node_modules"]
	source *dagger.Directory,
) (string, error) {
	return m.backendGen(ctx, source, "configs")
}

// BackendGenQueries regenerates queries + struct usage check + format_golang
// and fails if anything drifted.
func (m *Ci) BackendGenQueries(
	ctx context.Context,
	// +ignore=["frontend/node_modules", "backend/vendor", "ios/build", "**/node_modules"]
	source *dagger.Directory,
) (string, error) {
	return m.backendGen(ctx, source, "queries", "struct_usage_check", "format_golang")
}

// BackendGenEnvVars regenerates env-var constants and fails if anything drifted.
func (m *Ci) BackendGenEnvVars(
	ctx context.Context,
	// +ignore=["frontend/node_modules", "backend/vendor", "ios/build", "**/node_modules"]
	source *dagger.Directory,
) (string, error) {
	return m.backendGen(ctx, source, "env_vars")
}

// BackendGoModTidy ensures go.mod / go.sum are tidy. Skips `make vendor` to
// match the existing GHA workflow exactly.
func (m *Ci) BackendGoModTidy(
	ctx context.Context,
	// +ignore=["frontend/node_modules", "backend/vendor", "ios/build", "**/node_modules"]
	source *dagger.Directory,
) (string, error) {
	return m.backendBase(source).
		WithExec([]string{"go", "mod", "tidy"}).
		WithWorkdir("/src").
		WithExec([]string{"git", "diff", "--exit-code", "--", "backend/go.mod", "backend/go.sum"}).
		Stdout(ctx)
}

// BackendIntegrationTests runs `make integration_tests`. Tests use
// testcontainers-go internally to spin up Postgres, so the function only
// needs Go + a Docker daemon. The host docker socket is forwarded in.
func (m *Ci) BackendIntegrationTests(
	ctx context.Context,
	// +ignore=["frontend/node_modules", "backend/vendor", "ios/build", "**/node_modules"]
	source *dagger.Directory,
	// dockerSocket is the host's Docker socket. Pass `--docker-socket=/var/run/docker.sock`.
	dockerSocket *dagger.Socket,
) (string, error) {
	return m.backendBase(source).
		WithUnixSocket("/var/run/docker.sock", dockerSocket).
		WithExec([]string{"apk", "add", "--no-cache", "docker-cli"}).
		// testcontainers-go spawns sibling containers via the host docker
		// daemon; their ports are published on the host's localhost, which the
		// dagger container can't reach. Point testcontainers at the host's
		// gateway address so it connects via the bridge network instead.
		WithEnvVariable("TESTCONTAINERS_HOST_OVERRIDE", "host.docker.internal").
		WithExec([]string{"make", "vendor"}).
		WithExec([]string{"make", "integration_tests"}).
		Stdout(ctx)
}
