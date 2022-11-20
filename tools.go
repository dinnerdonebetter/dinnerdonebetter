//go:build tools
// +build tools

package tools

import (
	_ "github.com/daixiang0/gci"
	_ "github.com/google/wire/cmd/wire"
	_ "github.com/kyleconroy/sqlc"
	_ "golang.org/x/tools/cmd/goimports"
	_ "golang.org/x/tools/go/analysis/passes/fieldalignment"
)
