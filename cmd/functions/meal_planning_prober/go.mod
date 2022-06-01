module github.com/prixfixeco/api_server/cmd/functions/meal_planning_prober

// these have to be at 1.16 per Cloud Functions requirements: https://cloud.google.com/functions/docs/concepts/exec#runtimes
go 1.16

replace github.com/prixfixeco/api_server => ../../../

require (
	github.com/GoogleCloudPlatform/functions-framework-go v1.5.3
	github.com/prixfixeco/api_server v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/otel v1.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0 // indirect
)
