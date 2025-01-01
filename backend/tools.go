//go:build tools
// +build tools

package tools

import (
	_ "github.com/4meepo/tagalign/cmd/tagalign"
	_ "github.com/boyter/scc"
	_ "github.com/daixiang0/gci"
	_ "github.com/google/wire/cmd/wire"
	_ "github.com/sqlc-dev/sqlc/cmd/sqlc"
	_ "golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment"
)
