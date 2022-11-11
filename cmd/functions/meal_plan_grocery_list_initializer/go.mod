module github.com/prixfixeco/backend/cmd/functions/meal_plan_finalizer

go 1.19

replace github.com/prixfixeco/backend => ../../../

require (
	github.com/GoogleCloudPlatform/functions-framework-go v1.5.2
	github.com/prixfixeco/backend v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/otel v1.3.0
)
