package workers

import (
	"bytes"
	"context"
	"log"

	"github.com/goccy/go-graphviz"
	"github.com/prixfixeco/api_server/pkg/types"
	"gonum.org/v1/gonum/graph"
	dotencoding "gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"

	"github.com/prixfixeco/api_server/internal/customerdata"
	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/encoding"
	"github.com/prixfixeco/api_server/internal/messagequeue"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
)

const (
	advancedPrepStepCreationEnsurerWorkerName = "advanced_prep_step_creation_ensurer"
)

// AdvancedPrepStepCreationEnsurerWorker ensurers advanced prep steps are created.
type AdvancedPrepStepCreationEnsurerWorker struct {
	logger                logging.Logger
	tracer                tracing.Tracer
	encoder               encoding.ClientEncoder
	dataManager           database.DataManager
	postUpdatesPublisher  messagequeue.Publisher
	customerDataCollector customerdata.Collector
}

// ProvideAdvancedPrepStepCreationEnsurerWorker provides a AdvancedPrepStepCreationEnsurerWorker.
func ProvideAdvancedPrepStepCreationEnsurerWorker(
	logger logging.Logger,
	dataManager database.DataManager,
	postUpdatesPublisher messagequeue.Publisher,
	customerDataCollector customerdata.Collector,
	tracerProvider tracing.TracerProvider,
) *AdvancedPrepStepCreationEnsurerWorker {
	return &AdvancedPrepStepCreationEnsurerWorker{
		logger:                logging.EnsureLogger(logger).WithName(advancedPrepStepCreationEnsurerWorkerName),
		tracer:                tracing.NewTracer(tracerProvider.Tracer(advancedPrepStepCreationEnsurerWorkerName)),
		encoder:               encoding.ProvideClientEncoder(logger, tracerProvider, encoding.ContentTypeJSON),
		dataManager:           dataManager,
		postUpdatesPublisher:  postUpdatesPublisher,
		customerDataCollector: customerDataCollector,
	}
}

// HandleMessage handles a pending write.
func (w *AdvancedPrepStepCreationEnsurerWorker) HandleMessage(ctx context.Context, _ []byte) error {
	ctx, span := w.tracer.StartSpan(ctx)
	defer span.End()

	return w.ensureAdvancedPrepStepsAreCreated(ctx)
}

func (w *AdvancedPrepStepCreationEnsurerWorker) ensureAdvancedPrepStepsAreCreated(ctx context.Context) error {
	_, span := w.tracer.StartSpan(ctx)
	defer span.End()

	logger := w.logger.Clone()

	// get all the options that have been chosen for finalized meal plans
	// that are coming up in the week and haven't had advanced prep steps
	// created for them yet
	// w.dataManager.GetFinalizedMealPlanIDsForTheNextWeek(ctx)

	// iterate through the recipes

	logger.Info("ensureAdvancedPrepStepsAreCreated invoked")

	return nil
}

var _ graph.Node = (*graphNode)(nil)

// graphNode is a node in an implicit graph.
type graphNode struct {
	id        int64
	neighbors []graph.Node
	roots     []*graphNode
}

// NewGraphNode returns a new graphNode.
func NewGraphNode(id int64) *graphNode {
	return &graphNode{id: id}
}

func (g graphNode) ID() int64 {
	return g.id
}

func findStepIndexForRecipeStepProductID(recipe *types.Recipe, recipeStepProductID string) int64 {
	for _, step := range recipe.Steps {
		for _, product := range step.Products {
			if product.ID == recipeStepProductID {
				return int64(step.Index)
			}
		}
	}

	log.Fatal("invalid recipe step ID")

	return -1
}

func makeGraphForRecipe(ctx context.Context, recipe *types.Recipe) (graph.Graph, error) {
	recipeGraph := simple.NewDirectedGraph()

	for _, step := range recipe.Steps {
		recipeGraph.AddNode(NewGraphNode(int64(step.Index)))

		for _, ingredient := range step.Ingredients {
			if ingredient.ProductOfRecipeStep {
				toStep := findStepIndexForRecipeStepProductID(recipe, *ingredient.RecipeStepProductID)
				from := NewGraphNode(toStep)
				to := recipeGraph.Node(int64(step.Index))
				recipeGraph.SetEdge(simple.Edge{F: from, T: to})

			}
		}

		for _, instrument := range step.Instruments {
			if instrument.ProductOfRecipeStep {
				from := NewGraphNode(findStepIndexForRecipeStepProductID(recipe, *instrument.RecipeStepProductID))
				to := recipeGraph.Node(int64(step.Index))
				recipeGraph.SetEdge(simple.Edge{F: from, T: to})
			}
		}
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
	for id, _ := range nodesWithOnlyDependencies {
		if _, ok := assholeNodes[id]; !ok {
			for _, step := range recipe.Steps {
				if int64(step.Index) == id {
					stepsWorthNotifyingAbout = append(stepsWorthNotifyingAbout, step)
				}
			}
		}
	}

	return recipeGraph, nil
}

func renderGraph(g graph.Graph, path string) error {
	dotBytes, err := dotencoding.Marshal(g, "cheese buldak", "", "")
	if err != nil {
		return err
	}

	gv := graphviz.New()
	pg, err := graphviz.ParseBytes(dotBytes)

	var buf bytes.Buffer
	if err = gv.Render(pg, graphviz.PNG, &buf); err != nil {
		return err
	}

	if err = gv.RenderFilename(pg, graphviz.PNG, path); err != nil {
		return err
	}

	return nil
}
