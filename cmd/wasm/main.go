//go:build js && wasm
// +build js,wasm

package main

import (
	"encoding/json"
	"log"
	"syscall/js"

	"github.com/dinnerdonebetter/backend/pkg/types"
)

func buildRenderRecipeToMermaid() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 1 {
			return ""
		}
		inputJSON := args[0].String()

		var recipe *types.Recipe
		if err := json.Unmarshal([]byte(inputJSON), &recipe); err != nil {
			log.Printf("error: %v\n", err)
			return ""
		}

		return mermaidDiagramFromRecipe(recipe)
	})
}

func main() {
	js.Global().Set("renderRecipeToMermaid", buildRenderRecipeToMermaid())
	<-make(chan bool)
}
