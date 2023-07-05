//go:build js && wasm
// +build js,wasm

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"syscall/js"

	"github.com/dustin/go-humanize/english"

	"github.com/dinnerdonebetter/backend/pkg/types"
)

type provisionCount struct {
	ingredients, instruments, vessels uint
}

func stepProvidesWhatToOtherStep(recipe *types.Recipe, fromStepIndex, toStepIndex uint) string {
	from, to := recipe.Steps[fromStepIndex], recipe.Steps[toStepIndex]
	provides := []string{}

	count := provisionCount{}
	for _, product := range from.Products {
		for _, step := range recipe.Steps {
			if step.ID != to.ID {
				continue
			}

			for _, ingredient := range step.Ingredients {
				if ingredient.RecipeStepProductID != nil && *ingredient.RecipeStepProductID == product.ID {
					count.ingredients++
				}
			}

			for _, instrument := range step.Instruments {
				if instrument.RecipeStepProductID != nil && *instrument.RecipeStepProductID == product.ID {
					count.instruments++
				}
			}

			for _, vessel := range step.Vessels {
				if vessel.RecipeStepProductID != nil && *vessel.RecipeStepProductID == product.ID {
					count.vessels++
				}
			}
		}
	}

	renderCount := func(x uint, typ string) string {
		/*
			unnecessary Sprintf, but I might do something like this later:

			var prefix string
			if x == 1 {
				prefix = "an"
			}
		*/

		return strings.TrimSpace(fmt.Sprintf(" %s ", english.PluralWord(int(x), typ, fmt.Sprintf("%ss", typ))))
	}

	if count.ingredients > 0 {
		provides = append(provides, renderCount(count.ingredients, "ingredient"))
	}

	if count.instruments > 0 {
		provides = append(provides, renderCount(count.ingredients, "instrument"))
	}

	if count.vessels > 0 {
		provides = append(provides, renderCount(count.vessels, "vessel"))
	}

	return english.OxfordWordSeries(provides, "and")
}

func graphIDForStep(step *types.RecipeStep) int64 {
	return int64(step.Index + 1)
}

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
		var mermaid strings.Builder
		mermaid.WriteString("flowchart TD;\n")

		for _, step := range recipe.Steps {
			mermaid.WriteString(fmt.Sprintf("	Step%d[\"Step #%d (%s)\"];\n", graphIDForStep(step), graphIDForStep(step), step.Preparation.Name))
		}

		for i := range recipe.Steps {
			for j := range recipe.Steps {
				if i == j {
					continue
				}

				if provides := stepProvidesWhatToOtherStep(recipe, uint(i), uint(j)); provides != "" {
					mermaid.WriteString(fmt.Sprintf("\tStep%d -->|%s| Step%d;\n", graphIDForStep(recipe.Steps[i]), provides, graphIDForStep(recipe.Steps[j])))
				}
			}
		}

		for i := range recipe.PrepTasks {
			prepTask := recipe.PrepTasks[i]

			mermaid.WriteString(fmt.Sprintf("subgraph %d [\"%s (prep task #%d)\"]\n", i, prepTask.Name, i+1))
			for j := range prepTask.TaskSteps {
				mermaid.WriteString(fmt.Sprintf("Step%d", recipe.FindStepIndexByID(prepTask.TaskSteps[j].BelongsToRecipeStep)))
			}
			mermaid.WriteString("end\n")
		}

		return mermaid.String()
	})
}

func main() {
	js.Global().Set("renderRecipeToMermaid", buildRenderRecipeToMermaid())
	<-make(chan bool)
}
