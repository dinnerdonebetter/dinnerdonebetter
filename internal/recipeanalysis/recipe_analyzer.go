package recipeanalysis

import (
	"context"
	"errors"
	"fmt"
	"image"
	"strings"
	"time"

	"github.com/goccy/go-graphviz"
	"github.com/segmentio/ksuid"
	"gonum.org/v1/gonum/graph"
	dotencoding "gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"

	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

var _ graph.Node = (*recipeStepGraphNode)(nil)

// recipeStepGraphNode is a node in an implicit graph.
type recipeStepGraphNode struct {
	id int64
}

// NewGraphNode returns a new recipeStepGraphNode.
func NewGraphNode(id int64) graph.Node {
	return &recipeStepGraphNode{id: id}
}

func (g recipeStepGraphNode) ID() int64 {
	return g.id
}

var errRecipeStepIDNotFound = errors.New("recipe step ID not found")

func findStepIndexForRecipeStepProductID(recipe *types.Recipe, recipeStepProductID string) (int64, error) {
	for _, step := range recipe.Steps {
		for _, product := range step.Products {
			if product.ID == recipeStepProductID {
				return graphIDForStep(step), nil
			}
		}
	}

	return -1, errRecipeStepIDNotFound
}

func findStepIndexForRecipeStepID(recipe *types.Recipe, recipeStepID string) (int64, error) {
	for _, step := range recipe.Steps {
		if step.ID == recipeStepID {
			return graphIDForStep(step), nil
		}
	}

	return -1, errRecipeStepIDNotFound
}

func graphIDForStep(step *types.RecipeStep) int64 {
	return int64(step.Index + 1)
}

// RecipeAnalyzer generates DAG diagrams from recipes.
type RecipeAnalyzer interface {
	GenerateDAGDiagramForRecipe(ctx context.Context, recipe *types.Recipe) (image.Image, error)
	FindStepsEligibleForMealPlanTasks(ctx context.Context, recipe *types.Recipe) ([]*types.RecipeStep, error)
	GenerateMealPlanTasksForRecipe(ctx context.Context, mealPlanEvent *types.MealPlanEvent, mealPlanOptionID string, recipe *types.Recipe) ([]*types.MealPlanTaskDatabaseCreationInput, error)
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
		tracer: tracing.NewTracer(tracerProvider.Tracer("recipe_grapher")),
	}
}

// GenerateDAGDiagramForRecipe generates DAG diagrams for a given recipe.
func (g *recipeAnalyzer) GenerateDAGDiagramForRecipe(ctx context.Context, recipe *types.Recipe) (image.Image, error) {
	ctx, span := g.tracer.StartSpan(ctx)
	defer span.End()

	dag, err := g.makeGraphForRecipe(ctx, recipe)
	if err != nil {
		return nil, err
	}

	img, err := g.renderGraph(ctx, dag)
	if err != nil {
		return nil, err
	}

	_, err = g.FindStepsEligibleForMealPlanTasks(ctx, recipe)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func (g *recipeAnalyzer) makeGraphForRecipe(ctx context.Context, recipe *types.Recipe) (*simple.DirectedGraph, error) {
	_, span := g.tracer.StartSpan(ctx)
	defer span.End()

	recipeGraph := simple.NewDirectedGraph()

	for _, step := range recipe.Steps {
		recipeGraph.AddNode(NewGraphNode(graphIDForStep(step)))
	}

	for _, step := range recipe.Steps {
		for _, ingredient := range step.Ingredients {
			if !ingredient.ProductOfRecipeStep {
				continue
			}

			toStep, err := findStepIndexForRecipeStepProductID(recipe, *ingredient.RecipeStepProductID)
			if err != nil {
				return nil, err
			}

			from := recipeGraph.Node(toStep)
			to := recipeGraph.Node(graphIDForStep(step))
			recipeGraph.SetEdge(simple.Edge{F: from, T: to})
		}

		for _, instrument := range step.Instruments {
			if !instrument.ProductOfRecipeStep || instrument.RecipeStepProductID == nil {
				continue
			}

			toStep, err := findStepIndexForRecipeStepProductID(recipe, *instrument.RecipeStepProductID)
			if err != nil {
				return nil, err
			}

			from := recipeGraph.Node(toStep)
			to := recipeGraph.Node(graphIDForStep(step))
			recipeGraph.SetEdge(simple.Edge{F: from, T: to})
		}
	}

	return recipeGraph, nil
}

func (g *recipeAnalyzer) renderGraph(ctx context.Context, recipeGraph graph.Graph) (image.Image, error) {
	_, span := g.tracer.StartSpan(ctx)
	defer span.End()

	dotBytes, err := dotencoding.Marshal(recipeGraph, "", "", "")
	if err != nil {
		return nil, err
	}

	gv := graphviz.New()
	pg, err := graphviz.ParseBytes(dotBytes)
	if err != nil {
		return nil, err
	}

	img, err := gv.RenderImage(pg)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// FindStepsEligibleForMealPlanTasks finds steps eligible for meal plan task creation.
func (g *recipeAnalyzer) FindStepsEligibleForMealPlanTasks(ctx context.Context, recipe *types.Recipe) ([]*types.RecipeStep, error) {
	ctx, span := g.tracer.StartSpan(ctx)
	defer span.End()

	recipeGraph, err := g.makeGraphForRecipe(ctx, recipe)
	if err != nil {
		return nil, err
	}

	nodesWithOnlyDependencies := map[int64]graph.Node{}
	nodes := recipeGraph.Nodes()
	for nodes.Next() {
		n := nodes.Node()

		to := recipeGraph.To(n.ID())

		if to.Len() == 0 {
			nodesWithOnlyDependencies[n.ID()] = n
		}
	}

	assholeNodes := map[int64]graph.Node{}
	for _, node := range nodesWithOnlyDependencies {
		nodes = recipeGraph.From(node.ID())
		for nodes.Next() {
			n := nodes.Node()
			nID := n.ID()

			to := recipeGraph.To(nID)
			for to.Next() {
				f := to.Node()
				fID := f.ID()

				toParent := recipeGraph.To(fID)
				if toParent.Len() > 0 {
					assholeNodes[node.ID()] = node
				}
			}
		}
	}

	stepsWorthNotifyingAbout := []*types.RecipeStep{}
	for id := range nodesWithOnlyDependencies {
		if _, ok := assholeNodes[id]; !ok {
			for _, step := range recipe.Steps {
				allProductsHaveStorageInstructionsAndDurations := true
				for _, product := range step.Products {
					if product.MaximumStorageDurationInSeconds == 0 || strings.TrimSpace(product.StorageInstructions) == "" {
						allProductsHaveStorageInstructionsAndDurations = false
					}
				}

				if graphIDForStep(step) == id && allProductsHaveStorageInstructionsAndDurations {
					stepsWorthNotifyingAbout = append(stepsWorthNotifyingAbout, step)
				}
			}
		}
	}

	return stepsWorthNotifyingAbout, nil
}

// frozenIngredientDefrostStepsFilter iterates through a recipe and returns
// the list of ingredients within that are indicated as kept frozen.
func frozenIngredientDefrostStepsFilter(recipe *types.Recipe) map[string][]int {
	out := map[string][]int{}

	for _, recipeStep := range recipe.Steps {
		ingredientIndices := []int{}
		for i, ingredient := range recipeStep.Ingredients {
			// if it's a valid ingredient
			if ingredient.Ingredient != nil &&
				// if the ingredient has storage temperature set
				ingredient.Ingredient.MinimumIdealStorageTemperatureInCelsius != nil &&
				// the ingredient's storage temperature is set to something about freezing temperature.
				*ingredient.Ingredient.MinimumIdealStorageTemperatureInCelsius <= 3 {
				ingredientIndices = append(ingredientIndices, i)
			}
		}

		if len(ingredientIndices) > 0 {
			out[recipeStep.ID] = ingredientIndices
		}
	}

	return out
}

func whicheverIsLater(t1, t2 time.Time) time.Time {
	if t2.After(t1) {
		return t2
	}
	return t1
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

func determineCreationMinAndMaxTimesForRecipeStep(step *types.RecipeStep, mealPlanEvent *types.MealPlanEvent) (cannotCompleteBefore, cannotCompleteAfter time.Time) {
	var (
		shortestDuration,
		longestDuration uint32
	)

	for _, product := range step.Products {
		if product.MaximumStorageDurationInSeconds < shortestDuration || shortestDuration == 0 {
			shortestDuration = product.MaximumStorageDurationInSeconds
		}

		if product.MaximumStorageDurationInSeconds > longestDuration || longestDuration == 0 {
			longestDuration = product.MaximumStorageDurationInSeconds
		}
	}

	cannotCompleteBefore = whicheverIsLater(time.Now(), mealPlanEvent.StartsAt.Add(time.Duration(shortestDuration)*-time.Second))
	cannotCompleteAfter = whicheverIsLater(mealPlanEvent.StartsAt, mealPlanEvent.StartsAt.Add(time.Duration(longestDuration)*-time.Second))

	return cannotCompleteBefore, cannotCompleteAfter
}

const storagePrepCreationExplanation = "adequate storage instructions for early step"

func (g *recipeAnalyzer) GenerateMealPlanTasksForRecipe(ctx context.Context, mealPlanEvent *types.MealPlanEvent, mealPlanOptionID string, recipe *types.Recipe) ([]*types.MealPlanTaskDatabaseCreationInput, error) {
	ctx, span := g.tracer.StartSpan(ctx)
	defer span.End()

	logger := g.logger.Clone().WithValues(map[string]interface{}{
		keys.MealPlanEventIDKey:  mealPlanEvent.ID,
		keys.MealPlanOptionIDKey: mealPlanOptionID,
		keys.RecipeIDKey:         recipe.ID,
	})

	inputs := []*types.MealPlanTaskDatabaseCreationInput{}

	frozenIngredientSteps := frozenIngredientDefrostStepsFilter(recipe)
	logger.WithValue("frozen_steps_qty", len(frozenIngredientSteps)).Info("creating frozen step inputs")

	for stepID, ingredientIndices := range frozenIngredientSteps {
		stepIndex, err := findStepIndexForRecipeStepID(recipe, stepID)
		if err != nil {
			observability.AcknowledgeError(err, logger, span, "determining recipe step index for step ID")
			continue
		}

		explanation := buildThawStepCreationExplanation(stepIndex, ingredientIndices...)
		if explanation == "" {
			continue
		}

		inputs = append(inputs, &types.MealPlanTaskDatabaseCreationInput{
			ID:                   ksuid.New().String(),
			CannotCompleteBefore: mealPlanEvent.StartsAt.Add(2 * -time.Hour * 24),
			CannotCompleteAfter:  mealPlanEvent.StartsAt.Add(-time.Hour * 24),
			Status:               types.MealPlanTaskStatusUnfinished,
			CreationExplanation:  explanation,
			MealPlanOptionID:     mealPlanOptionID,
			RecipeStepID:         stepID,
		})
	}

	steps, graphErr := g.FindStepsEligibleForMealPlanTasks(ctx, recipe)
	if graphErr != nil {
		return nil, observability.PrepareAndLogError(graphErr, logger, nil, "generating graph for recipe")
	}

	logger.WithValue("steps_qty", len(steps)).Info("creating inputs from steps")

	for _, step := range steps {
		cannotCompleteBefore, cannotCompleteAfter := determineCreationMinAndMaxTimesForRecipeStep(step, mealPlanEvent)

		inputs = append(inputs, &types.MealPlanTaskDatabaseCreationInput{
			ID:                   ksuid.New().String(),
			CannotCompleteBefore: cannotCompleteBefore,
			CannotCompleteAfter:  cannotCompleteAfter,
			Status:               types.MealPlanTaskStatusUnfinished,
			CreationExplanation:  storagePrepCreationExplanation,
			MealPlanOptionID:     mealPlanOptionID,
			RecipeStepID:         step.ID,
		})
	}

	return inputs, nil
}
