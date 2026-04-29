package main

import (
	"context"

	"dagger/ci/internal/dagger"
)

// MarkdownLint runs markdownlint-cli over every .md file in the repo,
// mirroring scripts/lint_markdown.sh. Fails on unfixable lint issues.
func (m *Ci) MarkdownLint(
	ctx context.Context,
	// +ignore=["frontend/node_modules", "backend/vendor", "ios/build", "**/node_modules"]
	source *dagger.Directory,
) (string, error) {
	return dag.Container().
		From("ghcr.io/igorshubovych/markdownlint-cli:latest").
		WithMountedDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{
			"markdownlint",
			"--ignore", "**/vendor/**",
			"--ignore", "**/node_modules/**",
			"--ignore", "**/ios/build/**",
			"--fix",
			"--disable=MD013",
			"**/*.md",
		}).
		Stdout(ctx)
}
