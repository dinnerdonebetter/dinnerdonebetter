package recipeanalysis

import (
	"context"
	"errors"
	"fmt"
	"image"
	"strings"
	"time"

	"github.com/goccy/go-graphviz"
	"github.com/heimdalr/dag"
	"github.com/segmentio/ksuid"
	"gonum.org/v1/gonum/graph"
	dotencoding "gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"

	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

var errNotAcyclic = errors.New("recipe is not acyclic")

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
var errRecipeStepNotFound = errors.New("recipe step not found")

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

func findStepIDForRecipeStepProductID(recipe *types.Recipe, recipeStepProductID string) (string, error) {
	for _, step := range recipe.Steps {
		for _, product := range step.Products {
			if product.ID == recipeStepProductID {
				return fmt.Sprintf("%d", step.Index+1), nil
			}
		}
	}

	return "", errRecipeStepIDNotFound
}

func findStepIndexForRecipeStepID(recipe *types.Recipe, recipeStepID string) (int64, error) {
	for _, step := range recipe.Steps {
		if step.ID == recipeStepID {
			return graphIDForStep(step), nil
		}
	}

	return -1, errRecipeStepIDNotFound
}

func findStepForIndex(recipe *types.Recipe, index int64) (*types.RecipeStep, error) {
	for _, step := range recipe.Steps {
		if int64(step.Index+1) == index {
			return step, nil
		}
	}

	return nil, errRecipeStepNotFound
}

func graphIDForStep(step *types.RecipeStep) int64 {
	return int64(step.Index + 1)
}

// RecipeAnalyzer generates DAG diagrams from recipes.
type RecipeAnalyzer interface {
	GenerateDAGDiagramForRecipe(ctx context.Context, recipe *types.Recipe) (image.Image, error)
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

	recipeGraph, err := g.makeGraphForRecipe(ctx, recipe)
	if err != nil {
		return nil, err
	}

	img, err := g.renderGraph(ctx, recipeGraph)
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

	directedCycles := topo.DirectedCyclesIn(recipeGraph)
	if len(directedCycles) > 0 {
		return nil, errNotAcyclic
	}

	if _, err := g.makeDAGForRecipe(ctx, recipe); err != nil {
		return nil, fmt.Errorf("parsing recipe as DAG: %w", err)
	}

	return recipeGraph, nil
}

type RecipeStepIdentifier struct {
	recipeStep *types.RecipeStep
}

func (i *RecipeStepIdentifier) ID() string {
	return fmt.Sprintf("%d", i.recipeStep.Index+1)
}

// makeDAGForRecipe makes a proper DAG for the provided Recipe. Gonum has the notion of a directed graph, but
// doesn't seem to really give a rat's ass if it's acyclic, the way that the DAG library does.
func (g *recipeAnalyzer) makeDAGForRecipe(ctx context.Context, recipe *types.Recipe) (*dag.DAG, error) {
	_, span := g.tracer.StartSpan(ctx)
	defer span.End()

	recipeGraph := dag.NewDAG()

	for _, step := range recipe.Steps {
		if _, err := recipeGraph.AddVertex(&RecipeStepIdentifier{recipeStep: step}); err != nil {
			return nil, fmt.Errorf("adding step %d to graph: %w", step.Index, err)
		}
	}

	for _, step := range recipe.Steps {
		for _, ingredient := range step.Ingredients {
			if !ingredient.ProductOfRecipeStep {
				continue
			}

			toStepID, err := findStepIDForRecipeStepProductID(recipe, *ingredient.RecipeStepProductID)
			if err != nil {
				return nil, fmt.Errorf("finding step ID for step %d for recipe step product ID: %w", step.Index, err)
			}

			if err = recipeGraph.AddEdge(fmt.Sprintf("%d", step.Index+1), toStepID); err != nil {
				return nil, fmt.Errorf("adding step %d to graph: %w", step.Index, err)
			}
		}

		for _, instrument := range step.Instruments {
			if !instrument.ProductOfRecipeStep || instrument.RecipeStepProductID == nil {
				continue
			}

			toStep, err := findStepIDForRecipeStepProductID(recipe, *instrument.RecipeStepProductID)
			if err != nil {
				return nil, fmt.Errorf("finding step ID for step %d for recipe step ingredient ID: %w", step.Index, err)
			}

			if err = recipeGraph.AddEdge(step.ID, toStep); err != nil {
				return nil, fmt.Errorf("adding step %d to graph: %w", step.Index, err)
			}
		}
	}

	recipeGraph.ReduceTransitively()

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

// findStepsEligibleForMealPlanTasks finds steps eligible for meal plan task creation.
func (g *recipeAnalyzer) findStepsEligibleForMealPlanTasks(ctx context.Context, recipe *types.Recipe) ([]*types.RecipeStep, error) {
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
					if product.MaximumStorageDurationInSeconds == nil || strings.TrimSpace(product.StorageInstructions) == "" {
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

func findAllAncestorsForNode(recipeGraph graph.Directed, node graph.Node) []graph.Node {
	out := []graph.Node{}

	ancestors := recipeGraph.To(node.ID())

	for ancestors.Next() {
		n := ancestors.Node()

		if !nodeHasMultipleInboundVertices(recipeGraph, n) {
			out = append(out, n)
			out = append(out, findAllAncestorsForNode(recipeGraph, n)...)
		}
	}

	return out
}

func (g *recipeAnalyzer) findStepsWhichCulminateInStorablePreparedIngredients(ctx context.Context, recipe *types.Recipe) ([][]*types.RecipeStep, error) {
	ctx, span := g.tracer.StartSpan(ctx)
	defer span.End()

	recipeGraph, err := g.makeGraphForRecipe(ctx, recipe)
	if err != nil {
		return nil, err
	}

	junctionNodes := []graph.Node{}

	nodes := recipeGraph.Nodes()
	for nodes.Next() {
		n := nodes.Node()

		if nodeHasMultipleInboundVertices(recipeGraph, n) {
			junctionNodes = append(junctionNodes, n)
		}
	}

	stepSets := [][]graph.Node{}
	for _, junctionNode := range junctionNodes {
		ancestors := findAllAncestorsForNode(recipeGraph, junctionNode)
		stepSets = append(stepSets, ancestors)
	}

	out := [][]*types.RecipeStep{}
	for _, set := range stepSets {
		steps := []*types.RecipeStep{}
		for _, step := range set {
			s, findStepErr := findStepForIndex(recipe, step.ID())
			if findStepErr != nil {
				return nil, findStepErr
			}
			steps = append(steps, s)
		}
		out = append(out, steps)
	}

	return out, nil
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

func determineCreationMinAndMaxTimesForRecipeStep(steps []*types.RecipeStep, mealPlanEvent *types.MealPlanEvent) (cannotCompleteBefore, cannotCompleteAfter time.Time) {
	var (
		shortestDuration,
		longestDuration uint32
	)

	for _, step := range steps {
		for _, product := range step.Products {
			if product.MaximumStorageDurationInSeconds != nil {
				if *product.MaximumStorageDurationInSeconds < shortestDuration || shortestDuration == 0 {
					shortestDuration = *product.MaximumStorageDurationInSeconds
				}

				if *product.MaximumStorageDurationInSeconds > longestDuration || longestDuration == 0 {
					longestDuration = *product.MaximumStorageDurationInSeconds
				}
			}
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
	logger.WithValue("frozen_steps_qty", len(frozenIngredientSteps)).Info("creating frozen stepSet inputs")

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

		taskID := ksuid.New().String()
		inputs = append(inputs, &types.MealPlanTaskDatabaseCreationInput{
			ID:                   taskID,
			CannotCompleteBefore: mealPlanEvent.StartsAt.Add(2 * -time.Hour * 24),
			CannotCompleteAfter:  mealPlanEvent.StartsAt.Add(-time.Hour * 24),
			CreationExplanation:  explanation,
			MealPlanOptionID:     mealPlanOptionID,
			RecipeSteps: []*types.MealPlanTaskRecipeStepDatabaseCreationInput{
				{
					ID:                    ksuid.New().String(),
					BelongsToMealPlanTask: taskID,
					SatisfiesRecipeStep:   stepID,
				},
			},
		})
	}

	_, graphErr := g.findStepsEligibleForMealPlanTasks(ctx, recipe)
	if graphErr != nil {
		return nil, observability.PrepareAndLogError(graphErr, logger, nil, "finding steps eligible for advanced preparation")
	}

	stepSets, graphErr := g.findStepsWhichCulminateInStorablePreparedIngredients(ctx, recipe)
	if graphErr != nil {
		return nil, observability.PrepareAndLogError(graphErr, logger, nil, "finding steps eligible for advanced preparation")
	}

	logger.WithValue("steps_qty", len(stepSets)).Info("creating inputs from stepSets")

	for _, stepSet := range stepSets {
		cannotCompleteBefore, cannotCompleteAfter := determineCreationMinAndMaxTimesForRecipeStep(stepSet, mealPlanEvent)

		taskID := ksuid.New().String()
		input := &types.MealPlanTaskDatabaseCreationInput{
			ID:                   taskID,
			CannotCompleteBefore: cannotCompleteBefore,
			CannotCompleteAfter:  cannotCompleteAfter,
			CreationExplanation:  storagePrepCreationExplanation,
			MealPlanOptionID:     mealPlanOptionID,
			RecipeSteps:          []*types.MealPlanTaskRecipeStepDatabaseCreationInput{},
		}

		for _, step := range stepSet {
			input.RecipeSteps = append(input.RecipeSteps,
				&types.MealPlanTaskRecipeStepDatabaseCreationInput{
					ID:                    ksuid.New().String(),
					BelongsToMealPlanTask: taskID,
					SatisfiesRecipeStep:   step.ID,
				},
			)
		}

		inputs = append(inputs, input)
	}

	return inputs, nil
}
