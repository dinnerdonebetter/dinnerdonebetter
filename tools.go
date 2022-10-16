//go:build tools
// +build tools

package tools

import (
	_ "github.com/kyleconroy/sqlc"
	_ "golang.org/x/tools/go/analysis/passes/fieldalignment"
)
