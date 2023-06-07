# Dinner Done Better API

## dev dependencies

The following tools are prerequisites for development work:

- [go](https://golang.org/) 1.20
- [docker](https://docs.docker.com/get-docker/) &&  [docker-compose](https://docs.docker.com/compose/install/)
- [wire](https://github.com/google/wire) for dependency management
- [make](https://www.gnu.org/software/make/) for task running
- [terraform](https://learn.hashicorp.com/tutorials/terraform/install-cli) for deploying/formatting
- [fieldalignment](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/fieldalignment) (`go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest`)

## dev setup

It's a good idea to run `make quicktest lint integration_tests` before commits.

## running the server

1. clone this repository
2. run `make dev`
3. [http://localhost:8000/](http://localhost:8000/)

## Infrastructure

```mermaid
flowchart TD
    APIServer
    Cron(GCP Cloud Scheduler)
    DataChangesWorker(Data Changes Worker)
    MealPlanFinalizerWorker(Meal Plan Finalizer)
    MealPlanGroceryListInitializerWorker(Meal Plan Grocery List Initializer)
    MealPlanTaskCreatorWorker(Meal Plan Task Creator)
    OutboundEmailerWorker(Outbound Emailer)
    SearchDataIndexSchedulerWorker(Search Data Index Scheduler)
    SearchDataIndexerWorker(Search Data Indexer)
    Cron-->MealPlanGroceryListInitializerWorker
    Cron-->SearchDataIndexSchedulerWorker
    Cron-->MealPlanTaskCreatorWorker
    Cron-->MealPlanFinalizerWorker
    DataChangesWorker-->OutboundEmailerWorker
    DataChangesWorker-->SearchDataIndexerWorker
    OutboundEmailerWorker-->DataChangesWorker
    MealPlanGroceryListInitializerWorker-->DataChangesWorker
    MealPlanTaskCreatorWorker-->DataChangesWorker
    MealPlanFinalizerWorker-->DataChangesWorker
    SearchDataIndexSchedulerWorker-->SearchDataIndexerWorker
    APIServer-->DataChangesWorker
    APIServer-->SearchDataIndexerWorker
```
