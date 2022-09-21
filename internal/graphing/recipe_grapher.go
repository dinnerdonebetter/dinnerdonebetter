package graphing

import (
	"context"
	"errors"
	"image"

	"github.com/goccy/go-graphviz"
	"gonum.org/v1/gonum/graph"
	dotencoding "gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"

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

func graphIDForStep(step *types.RecipeStep) int64 {
	return int64(step.Index + 1)
}

// RecipeDAGDiagramGenerator generates DAG diagrams from recipes.
type RecipeDAGDiagramGenerator interface {
	GenerateDAGDiagramForRecipe(ctx context.Context, recipe *types.Recipe) (image.Image, error)
}

var _ RecipeDAGDiagramGenerator = (*RecipeGrapher)(nil)

// RecipeGrapher creates graphs from recipes.
type RecipeGrapher struct {
	tracer tracing.Tracer
}

// NewRecipeGrapher creates a RecipeGrapher.
func NewRecipeGrapher(tracerProvider tracing.TracerProvider) RecipeDAGDiagramGenerator {
	return &RecipeGrapher{
		tracer: tracing.NewTracer(tracerProvider.Tracer("recipe_grapher")),
	}
}

// GenerateDAGDiagramForRecipe generates DAG diagrams for a given recipe.
func (g *RecipeGrapher) GenerateDAGDiagramForRecipe(ctx context.Context, recipe *types.Recipe) (image.Image, error) {
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

	_, err = g.findStepsEligibleFor(ctx, recipe)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func (g *RecipeGrapher) makeGraphForRecipe(ctx context.Context, recipe *types.Recipe) (*simple.DirectedGraph, error) {
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
			if !instrument.ProductOfRecipeStep {
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

func (g *RecipeGrapher) renderGraph(ctx context.Context, recipeGraph graph.Graph) (image.Image, error) {
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

func (g *RecipeGrapher) findStepsEligibleFor(ctx context.Context, recipe *types.Recipe) ([]*types.RecipeStep, error) {
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
			toLen := to.Len()
			_ = toLen

			for to.Next() {
				f := to.Node()
				fID := f.ID()

				toParent := recipeGraph.To(fID)
				toParentLen := to.Len()
				_ = toParentLen

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
				if graphIDForStep(step) == id {
					stepsWorthNotifyingAbout = append(stepsWorthNotifyingAbout, step)
				}
			}
		}
	}

	return stepsWorthNotifyingAbout, nil
}
