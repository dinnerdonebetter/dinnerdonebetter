package main

import (
	"context"

	"dagger/ci/internal/dagger"
)

// frontendBase mounts the source at /src, runs `npm ci` in /src/frontend, and
// returns a container ready to run any frontend npm script.
func (m *Ci) frontendBase(source *dagger.Directory) *dagger.Container {
	return m.baseNode().
		WithMountedDirectory("/src", source).
		WithWorkdir("/src/frontend").
		WithExec([]string{"npm", "ci"})
}

// FrontendBuild runs `npm run build` (consumer + admin workspaces).
func (m *Ci) FrontendBuild(
	ctx context.Context,
	// +ignore=["frontend/node_modules", "backend/vendor", "ios/build", "**/node_modules"]
	source *dagger.Directory,
) (string, error) {
	return m.frontendBase(source).
		WithExec([]string{"npm", "run", "build"}).
		Stdout(ctx)
}

// FrontendLint runs `npm run lint` (which includes pin:check via git diff)
// followed by `npm run check` (svelte-check).
func (m *Ci) FrontendLint(
	ctx context.Context,
	// +ignore=["frontend/node_modules", "backend/vendor", "ios/build", "**/node_modules"]
	source *dagger.Directory,
) (string, error) {
	return m.frontendBase(source).
		WithExec([]string{"npm", "run", "lint"}).
		WithExec([]string{"npm", "run", "check"}).
		Stdout(ctx)
}

// FrontendFormatCheck runs `npm run format:check` (Prettier).
func (m *Ci) FrontendFormatCheck(
	ctx context.Context,
	// +ignore=["frontend/node_modules", "backend/vendor", "ios/build", "**/node_modules"]
	source *dagger.Directory,
) (string, error) {
	return m.frontendBase(source).
		WithExec([]string{"npm", "run", "format:check"}).
		Stdout(ctx)
}

// FrontendUnitTests runs `npm run test` across both workspaces.
func (m *Ci) FrontendUnitTests(
	ctx context.Context,
	// +ignore=["frontend/node_modules", "backend/vendor", "ios/build", "**/node_modules"]
	source *dagger.Directory,
) (string, error) {
	return m.frontendBase(source).
		WithExec([]string{"npm", "run", "test"}).
		Stdout(ctx)
}
