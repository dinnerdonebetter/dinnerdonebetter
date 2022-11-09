module github.com/prixfixeco/backend/cmd/functions/meal_plan_finalizer

// these have to be at 1.16 per Cloud Functions requirements: https://cloud.google.com/functions/docs/concepts/exec#runtimes
go 1.16

replace github.com/prixfixeco/backend => ../../../

require (
	github.com/GoogleCloudPlatform/functions-framework-go v1.5.2
	github.com/goccy/go-graphviz v0.0.9
	github.com/prixfixeco/backend v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/otel v1.3.0
	gonum.org/v1/gonum v0.12.0
)
