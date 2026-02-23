package recipeanalysis

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	mealplanningkeys "github.com/dinnerdonebetter/backend/internal/domain/mealplanning/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"

	"github.com/dustin/go-humanize/english"
	"github.com/hako/durafmt"
	"github.com/heimdalr/dag"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"
)

var (
	_ graph.Node = (*recipeStepGraphNode)(nil)

	errNotAcyclic = errors.New("recipe is not acyclic")
)

// recipeStepGraphNode is a node in an implicit graph.
type recipeStepGraphNode struct {
	id int64
}

// newGraphNode returns a new recipeStepGraphNode.
func newGraphNode(id int64) graph.Node {
	return &recipeStepGraphNode{id: id}
}

func (g recipeStepGraphNode) ID() int64 {
	return g.id
}

var errRecipeStepIDNotFound = errors.New("recipe step ID not found")

func findStepIndexForRecipeStepProductID(recipe *mealplanning.Recipe, recipeStepProductID string) (int64, error) {
	if step := recipe.FindStepForRecipeStepProductID(recipeStepProductID); step != nil {
		return graphIDForStep(step), nil
	}

	return -1, errRecipeStepIDNotFound
}

func findStepIDForRecipeStepProductID(recipe *mealplanning.Recipe, recipeStepProductID string) (string, error) {
	if step := recipe.FindStepForRecipeStepProductID(recipeStepProductID); step != nil {
		return fmt.Sprintf("%d", graphIDForStep(step)), nil
	}

	return "", errRecipeStepIDNotFound
}

func findStepIndexForRecipeStepID(recipe *mealplanning.Recipe, recipeStepID string) (int64, error) {
	for _, step := range recipe.Steps {
		if step.ID == recipeStepID {
			return graphIDForStep(step), nil
		}
	}

	return -1, errRecipeStepIDNotFound
}

func graphIDForStep(step *mealplanning.RecipeStep) int64 {
	return int64(step.Index + 1)
}

// RecipeAnalyzer analyzes recipes for insights (ugh).
type RecipeAnalyzer interface {
	MakeGraphForRecipe(ctx context.Context, recipe *mealplanning.Recipe) (*simple.DirectedGraph, error)
	ValidateRecipeCreationRequestInputIsDAG(ctx context.Context, input *mealplanning.RecipeCreationRequestInput) error
	GenerateMealPlanTasksForRecipe(ctx context.Context, mealPlanOptionID string, recipe *mealplanning.Recipe) ([]*mealplanning.MealPlanTaskDatabaseCreationInput, error)
	RenderMermaidDiagramForRecipe(ctx context.Context, recipe *mealplanning.Recipe) string
	RenderGraphvizDiagramForRecipe(ctx context.Context, recipe *mealplanning.Recipe) string
}

var _ RecipeAnalyzer = (*recipeAnalyzer)(nil)

// recipeAnalyzer creates graphs from recipes.
type recipeAnalyzer struct {
	tracer tracing.Tracer
	logger logging.Logger
}

// NewRecipeAnalyzer creates a recipeAnalyzer.
func NewRecipeAnalyzer(logger logging.Logger, tracerProvider tracing.TracerProvider) RecipeAnalyzer {
	return &recipeAnalyzer{
		logger: logging.EnsureLogger(logger).WithName("recipe_analyzer"),
		tracer: tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer("recipe_grapher")),
	}
}

func (g *recipeAnalyzer) MakeGraphForRecipe(ctx context.Context, recipe *mealplanning.Recipe) (*simple.DirectedGraph, error) {
	_, span := g.tracer.StartSpan(ctx)
	defer span.End()

	recipeGraph := simple.NewDirectedGraph()

	for _, step := range recipe.Steps {
		recipeGraph.AddNode(newGraphNode(graphIDForStep(step)))
	}

	for _, step := range recipe.Steps {
		for _, ingredient := range step.Ingredients {
			if ingredient.RecipeStepProductID == nil {
				continue
			}

			fromStep, err := findStepIndexForRecipeStepProductID(recipe, *ingredient.RecipeStepProductID)
			if err != nil {
				return nil, err
			}

			from := recipeGraph.Node(fromStep)
			to := recipeGraph.Node(graphIDForStep(step))
			recipeGraph.SetEdge(simple.Edge{F: from, T: to})
		}

		for _, instrument := range step.Instruments {
			if instrument.RecipeStepProductID == nil {
				continue
			}

			fromStep, err := findStepIndexForRecipeStepProductID(recipe, *instrument.RecipeStepProductID)
			if err != nil {
				return nil, err
			}

			from := recipeGraph.Node(fromStep)
			to := recipeGraph.Node(graphIDForStep(step))
			recipeGraph.SetEdge(simple.Edge{F: from, T: to})
		}

		for _, vessel := range step.Vessels {
			if vessel.RecipeStepProductID == nil {
				continue
			}

			fromStep, err := findStepIndexForRecipeStepProductID(recipe, *vessel.RecipeStepProductID)
			if err != nil {
				return nil, err
			}

			from := recipeGraph.Node(fromStep)
			to := recipeGraph.Node(graphIDForStep(step))
			recipeGraph.SetEdge(simple.Edge{F: from, T: to})
		}
	}

	directedCycles := topo.DirectedCyclesIn(recipeGraph)
	if len(directedCycles) > 0 {
		return nil, errNotAcyclic
	}

	if _, err := g.makeDAGForRecipe(ctx, recipe); err != nil {
		return nil, fmt.Errorf("parsing recipe as DAG: %w", err)
	}

	return recipeGraph, nil
}

// ValidateRecipeCreationRequestInputIsDAG validates that a RecipeCreationRequestInput represents a valid DAG.
// It builds a graph from the request input and checks for cycles.
// ProductOfRecipeStepIndex is treated as an array index into the Steps slice.
func (g *recipeAnalyzer) ValidateRecipeCreationRequestInputIsDAG(ctx context.Context, input *mealplanning.RecipeCreationRequestInput) error {
	_, span := g.tracer.StartSpan(ctx)
	defer span.End()

	recipeGraph := simple.NewDirectedGraph()

	// Add nodes for each step (using array index+1 as the graph ID, similar to graphIDForStep)
	for stepIdx := range input.Steps {
		graphID := int64(stepIdx + 1)
		recipeGraph.AddNode(newGraphNode(graphID))
	}

	// Build edges based on ProductOfRecipeStepIndex references
	for stepIdx := range input.Steps {
		currentStepGraphID := int64(stepIdx + 1)
		step := input.Steps[stepIdx]

		// Check ingredients for references to other steps
		for _, ingredient := range step.Ingredients {
			if ingredient.ProductOfRecipeStepIndex == nil {
				continue
			}

			fromStepArrayIndex := *ingredient.ProductOfRecipeStepIndex

			// Skip cross-recipe references (they don't affect this recipe's DAG)
			// Check this first before validating index bounds
			// If RecipeStepProductRecipeID is set (even if empty), it indicates a cross-recipe reference
			// The empty string case can happen when getRecipeIDBySlug returns nil but we still want to indicate
			// this is a cross-recipe reference that will be resolved later
			if ingredient.RecipeStepProductRecipeID != nil {
				continue
			}

			// Validate that the referenced step array index exists
			if fromStepArrayIndex >= uint64(len(input.Steps)) {
				return fmt.Errorf("ingredient in step at array index %d references invalid step array index %d", stepIdx, fromStepArrayIndex)
			}

			// Prevent self-references (which would create a cycle)
			if fromStepArrayIndex == uint64(stepIdx) {
				return fmt.Errorf("%w: step at array index %d references itself", errNotAcyclic, stepIdx)
			}

			fromStepGraphID := int64(fromStepArrayIndex + 1)
			from := recipeGraph.Node(fromStepGraphID)
			to := recipeGraph.Node(currentStepGraphID)
			recipeGraph.SetEdge(simple.Edge{F: from, T: to})
		}

		// Check instruments for references to other steps
		for _, instrument := range step.Instruments {
			if instrument.ProductOfRecipeStepIndex == nil {
				continue
			}

			fromStepArrayIndex := *instrument.ProductOfRecipeStepIndex
			// Validate that the referenced step array index exists
			if fromStepArrayIndex >= uint64(len(input.Steps)) {
				return fmt.Errorf("instrument in step at array index %d references invalid step array index %d", stepIdx, fromStepArrayIndex)
			}

			// Prevent self-references (which would create a cycle)
			if fromStepArrayIndex == uint64(stepIdx) {
				return fmt.Errorf("%w: step at array index %d references itself", errNotAcyclic, stepIdx)
			}

			fromStepGraphID := int64(fromStepArrayIndex + 1)
			from := recipeGraph.Node(fromStepGraphID)
			to := recipeGraph.Node(currentStepGraphID)
			recipeGraph.SetEdge(simple.Edge{F: from, T: to})
		}

		// Check vessels for references to other steps
		for _, vessel := range step.Vessels {
			if vessel.ProductOfRecipeStepIndex == nil {
				continue
			}

			fromStepArrayIndex := *vessel.ProductOfRecipeStepIndex
			// Validate that the referenced step array index exists
			if fromStepArrayIndex >= uint64(len(input.Steps)) {
				return fmt.Errorf("vessel in step at array index %d references invalid step array index %d", stepIdx, fromStepArrayIndex)
			}

			// Prevent self-references (which would create a cycle)
			if fromStepArrayIndex == uint64(stepIdx) {
				return fmt.Errorf("%w: step at array index %d references itself", errNotAcyclic, stepIdx)
			}

			fromStepGraphID := int64(fromStepArrayIndex + 1)
			from := recipeGraph.Node(fromStepGraphID)
			to := recipeGraph.Node(currentStepGraphID)
			recipeGraph.SetEdge(simple.Edge{F: from, T: to})
		}
	}

	// Check for cycles
	directedCycles := topo.DirectedCyclesIn(recipeGraph)
	if len(directedCycles) > 0 {
		return fmt.Errorf("%w: recipe contains %d cycle(s)", errNotAcyclic, len(directedCycles))
	}

	return nil
}

type RecipeStepIdentifier struct {
	recipeStep *mealplanning.RecipeStep
}

func (i *RecipeStepIdentifier) ID() string {
	return fmt.Sprintf("%d", graphIDForStep(i.recipeStep))
}

// makeDAGForRecipe makes a proper DAG for the provided Recipe.
func (g *recipeAnalyzer) makeDAGForRecipe(ctx context.Context, recipe *mealplanning.Recipe) (*dag.DAG, error) {
	_, span := g.tracer.StartSpan(ctx)
	defer span.End()

	recipeGraph := dag.NewDAG()

	for _, step := range recipe.Steps {
		if _, err := recipeGraph.AddVertex(&RecipeStepIdentifier{recipeStep: step}); err != nil {
			return nil, fmt.Errorf("adding initial step %d to graph: %w", step.Index, err)
		}
	}

	for _, step := range recipe.Steps {
		for _, ingredient := range step.Ingredients {
			if ingredient.RecipeStepProductID == nil {
				continue
			}

			toStepID, err := findStepIDForRecipeStepProductID(recipe, *ingredient.RecipeStepProductID)
			if err != nil {
				return nil, fmt.Errorf("finding step ID for step %d for recipe step product ID: %w", step.Index, err)
			}

			if err = recipeGraph.AddEdge(fmt.Sprintf("%d", graphIDForStep(step)), toStepID); err != nil {
				return nil, fmt.Errorf("adding recipe step %d to graph: %w", step.Index, err)
			}
		}

		for _, instrument := range step.Instruments {
			if instrument.RecipeStepProductID == nil {
				continue
			}

			toStep, err := findStepIDForRecipeStepProductID(recipe, *instrument.RecipeStepProductID)
			if err != nil {
				return nil, fmt.Errorf("finding step ID for step %d for recipe step ingredient ID: %w", step.Index, err)
			}

			if err = recipeGraph.AddEdge(fmt.Sprintf("%d", graphIDForStep(step)), toStep); err != nil {
				var dupeErr dag.EdgeDuplicateError
				if errors.As(err, &dupeErr) {
					continue
				}

				return nil, fmt.Errorf("adding ingredient step %d to graph: %w", step.Index, err)
			}
		}
	}

	recipeGraph.ReduceTransitively()

	return recipeGraph, nil
}

// frozenIngredientDefrostStepsFilter iterates through a recipe and returns
// the list of ingredients within that are indicated as kept frozen.
func frozenIngredientDefrostStepsFilter(recipe *mealplanning.Recipe) map[string][]int {
	out := map[string][]int{}

	for _, recipeStep := range recipe.Steps {
		ingredientIndices := []int{}
		for i, ingredient := range recipeStep.Ingredients {
			// if it's a valid ingredient
			if ingredient.Ingredient != nil &&
				// if the ingredient has storage temperature set
				ingredient.Ingredient.StorageTemperatureInCelsius.Min != nil &&
				// the ingredient's storage temperature is set to something about freezing temperature.
				*ingredient.Ingredient.StorageTemperatureInCelsius.Min <= 3 {
				ingredientIndices = append(ingredientIndices, i)
			}
		}

		if len(ingredientIndices) > 0 {
			out[recipeStep.ID] = ingredientIndices
		}
	}

	return out
}

func buildThawStepCreationExplanation(recipeStepIndex int64, ingredientIndices ...int) string {
	if len(ingredientIndices) == 0 {
		return ""
	}

	stringIndices := []string{}
	for _, i := range ingredientIndices {
		stringIndices = append(stringIndices, fmt.Sprintf("#%d", i+1))
	}

	var d string
	if len(ingredientIndices) > 1 {
		d = "ingredients"
	} else if len(ingredientIndices) == 1 {
		d = "ingredient"
	}

	return fmt.Sprintf("frozen %s (%s) for step #%d might need to be thawed ahead of time", d, strings.Join(stringIndices, ", "), recipeStepIndex)
}

func (g *recipeAnalyzer) generateMealPlanTasksForFrozenIngredients(ctx context.Context, mealPlanOptionID string, recipe *mealplanning.Recipe) []*mealplanning.MealPlanTaskDatabaseCreationInput {
	_, span := g.tracer.StartSpan(ctx)
	defer span.End()

	logger := g.logger.Clone().WithValue(mealplanningkeys.RecipeIDKey, recipe.ID)

	frozenIngredientSteps := frozenIngredientDefrostStepsFilter(recipe)
	logger.WithValue("frozen_steps_qty", len(frozenIngredientSteps)).Info("creating frozen stepSet inputs")

	outputs := []*mealplanning.MealPlanTaskDatabaseCreationInput{}
	for stepID, ingredientIndices := range frozenIngredientSteps {
		stepIndex, err := findStepIndexForRecipeStepID(recipe, stepID)
		if err != nil {
			observability.AcknowledgeError(err, logger, span, "determining recipe stepSet index for stepSet ID")
			continue
		}

		explanation := buildThawStepCreationExplanation(stepIndex, ingredientIndices...)
		if explanation == "" {
			continue
		}

		outputs = append(outputs, &mealplanning.MealPlanTaskDatabaseCreationInput{
			ID:                  identifiers.New(),
			CreationExplanation: explanation,
			MealPlanOptionID:    mealPlanOptionID,
		})
	}

	return outputs
}

const recipeTaskStepCreationExplanation = "recipe prep task exists for steps"

func (g *recipeAnalyzer) GenerateMealPlanTasksForRecipe(ctx context.Context, mealPlanOptionID string, recipe *mealplanning.Recipe) ([]*mealplanning.MealPlanTaskDatabaseCreationInput, error) {
	ctx, span := g.tracer.StartSpan(ctx)
	defer span.End()

	inputs := g.generateMealPlanTasksForFrozenIngredients(ctx, mealPlanOptionID, recipe)

	for _, prepTask := range recipe.PrepTasks {
		inputs = append(inputs, &mealplanning.MealPlanTaskDatabaseCreationInput{
			AssignedToUser:      nil,
			CreationExplanation: recipeTaskStepCreationExplanation,
			StatusExplanation:   "",
			MealPlanOptionID:    mealPlanOptionID,

			RecipePrepTaskID: prepTask.ID,
			ID:               identifiers.New(),
		})
	}

	return inputs, nil
}

type provisionCount struct {
	ingredients, instruments, vessels uint
}

func stepProvidesWhatToOtherStep(recipe *mealplanning.Recipe, fromStepIndex, toStepIndex uint) string {
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

func (g *recipeAnalyzer) RenderMermaidDiagramForRecipe(ctx context.Context, recipe *mealplanning.Recipe) string {
	_, span := g.tracer.StartSpan(ctx)
	defer span.End()

	var mermaid strings.Builder
	mermaid.WriteString("flowchart TD;\n")

	for _, step := range recipe.Steps {
		if _, err := fmt.Fprintf(&mermaid, "	Step%d[\"Step #%d (%s)\"];\n", graphIDForStep(step), graphIDForStep(step), step.Preparation.Name); err != nil {
			observability.AcknowledgeError(err, g.logger, span, "writing mermaid step node")
		}
	}

	for i := range recipe.Steps {
		for j := range recipe.Steps {
			if i == j {
				continue
			}

			if provides := stepProvidesWhatToOtherStep(recipe, uint(i), uint(j)); provides != "" {
				if _, err := fmt.Fprintf(&mermaid, "\tStep%d -->|%s| Step%d;\n", graphIDForStep(recipe.Steps[i]), provides, graphIDForStep(recipe.Steps[j])); err != nil {
					observability.AcknowledgeError(err, g.logger, span, "writing mermaid step edge")
				}
			}
		}
	}

	for i := range recipe.PrepTasks {
		prepTask := recipe.PrepTasks[i]

		if _, err := fmt.Fprintf(&mermaid, "subgraph %d [\"%s (prep task #%d)\"]\n", i, prepTask.Name, i+1); err != nil {
			observability.AcknowledgeError(err, g.logger, span, "writing mermaid subgraph header")
		}
		for j := range prepTask.TaskSteps {
			if _, err := fmt.Fprintf(&mermaid, "Step%d\n", recipe.FindStepIndexByID(prepTask.TaskSteps[j].BelongsToRecipeStep)); err != nil {
				observability.AcknowledgeError(err, g.logger, span, "writing mermaid subgraph step")
			}
		}
		mermaid.WriteString("end\n")
	}

	return mermaid.String()
}

var colorNames = []string{
	"CornflowerBlue",
	"purple",
	"Salmon",
	"PeachPuff",
	"Tomato",
	"Orange",
	"Chocolate",
	"LemonChiffon",
}

const graphvizPreamble = `strict digraph {
	fontname="Outfit,Helvetica,Arial,sans-serif";
	node [fontname="Helvetica,Arial,sans-serif" shape=egg];
	edge [fontname="Helvetica,Arial,sans-serif"];
    concentrate=true;
    compound=true;
	rankdir=TB;
	start=1;
`

func (g *recipeAnalyzer) RenderGraphvizDiagramForRecipe(ctx context.Context, recipe *mealplanning.Recipe) string {
	_, span := g.tracer.StartSpan(ctx)
	defer span.End()

	var graphViz strings.Builder
	graphViz.WriteString(graphvizPreamble + "\n")

	for _, step := range recipe.Steps {
		if _, err := fmt.Fprintf(&graphViz, "\tStep%d [label=\"Step #%d (%s)\"];\n", graphIDForStep(step), graphIDForStep(step), step.Preparation.Name); err != nil {
			observability.AcknowledgeError(err, g.logger, span, "writing graphviz step node")
		}
	}

	for i := range recipe.Steps {
		for j := range recipe.Steps {
			if i == j {
				continue
			}

			if provides := stepProvidesWhatToOtherStep(recipe, uint(i), uint(j)); provides != "" {
				stepLabel := ""
				if recipe.Steps[i].EstimatedTimeInSeconds.Min != nil && *recipe.Steps[i].EstimatedTimeInSeconds.Min > 0 {
					stepLabel = durafmt.Parse(time.Duration(*recipe.Steps[i].EstimatedTimeInSeconds.Min) * time.Second).String()
				}

				if _, err := fmt.Fprintf(&graphViz, "\tStep%d -> Step%d [color=\"black\" label=%q];\n", graphIDForStep(recipe.Steps[i]), graphIDForStep(recipe.Steps[j]), stepLabel); err != nil {
					observability.AcknowledgeError(err, g.logger, span, "writing graphviz step edge")
				}
			}
		}
	}

	for i := range recipe.PrepTasks {
		prepTask := recipe.PrepTasks[i]

		if _, err := fmt.Fprintf(&graphViz, "\n\tsubgraph cluster_%d {\n\t\tnode [style=filled];\n\t\t", i); err != nil {
			observability.AcknowledgeError(err, g.logger, span, "writing graphviz subgraph header")
		}
		for j := range prepTask.TaskSteps {
			var arrow string
			if j != len(prepTask.TaskSteps)-1 {
				arrow = " -> "
			}
			if _, err := fmt.Fprintf(&graphViz, "Step%d%s", graphIDForStep(recipe.Steps[recipe.FindStepIndexByID(prepTask.TaskSteps[j].BelongsToRecipeStep)]), arrow); err != nil {
				observability.AcknowledgeError(err, g.logger, span, "writing graphviz subgraph step")
			}
		}
		colorToUse := colorNames[i%len(colorNames)]

		if _, err := fmt.Fprintf(&graphViz, ";\n\t\tlabel = \"(prep task #%d)\";\n\t\tcolor=%s;\n\t}\n", i+1, colorToUse); err != nil {
			observability.AcknowledgeError(err, g.logger, span, "writing graphviz subgraph footer")
		}
	}

	graphViz.WriteString("}\n")

	return graphViz.String()
}
