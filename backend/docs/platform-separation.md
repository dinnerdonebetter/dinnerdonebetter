# Platform Separation Roadmap

The goal of this repo is to be a reusable service template: someone should be able to fork it, remove the meal planning domain, and replace it with their own without touching core infrastructure. This document tracks the work needed to make that as painless as possible.

See the [template philosophy](../../README.md) for background.

---

## Open Tasks

### High Priority

- [ ] **Make authorization permissions domain-pluggable**
  Refactor `internal/authorization/permissions.go` so the generic role sets (`ServiceAdmin`, `ServiceDataAdmin`, `AccountAdmin`, `AccountMember`) contain only platform-level permissions. Introduce a registration pattern (mirroring the DI registration module) so a domain can inject its own permission constants into the appropriate roles at startup. Currently almost the entire content of those role sets is meal-planning-specific constants with no injection mechanism.

- [ ] **Make the top-level gRPC server struct domain-agnostic**
  Refactor `internal/build/services/api/grpc/grpc_service.go` and `extras.go` so `MealPlanningServiceServer` is not embedded in the struct definition or hardcoded in `BuildRegistrationFuncs`/`AggregateMethodPermissions` signatures. The domain registration module should handle its own server registration. Currently the project will not compile without the meal planning domain present because it is embedded as a named field in the server struct.

### Medium Priority

- [ ] **Make `ServicesConfig` and the `configurations` union type domain-agnostic**
  Remove the `MealPlanning` field from `ServicesConfig` in `internal/config/services_config.go` and remove the three meal-planning config types from the `configurations` union in `internal/config/configs.go`. Domain config should be registered and validated through the domain registration module rather than hardcoded into the shared config struct.

- [ ] **Make the search data index scheduler pluggable**
  Refactor `internal/build/jobs/search_data_index_scheduler/build.go` and `indexers.go` so `ProvideIndexFunctions` accepts domain indexers via the registration pattern rather than hardcoding the merge of generic + meal-planning indexers.

- [ ] **Make data change message handler searchers pluggable**
  Refactor `internal/build/functions/data_change_message_handler/searchers.go` so the 8 meal-planning-specific search index providers are registered by the domain registration module rather than hardcoded in an infrastructure file.

- [ ] **Write a domain-swap script**
  Write a bash script (e.g. `scripts/swap-domain.sh`) that removes all files/directories labeled `// Domain: mealplanning`, removes the three meal-plan worker `cmd/` directories, and scaffolds a minimal placeholder domain (registration module, empty proto, stub repository and service) so the project compiles clean out of the box after running it.

### Low Priority

- [ ] **Fix MCP server build to use the registration pattern**
  Refactor `internal/build/services/mcp/build.go` to go through `mealplanningregistration` like the gRPC API builder does, rather than directly importing the meal planning repo. Currently inconsistent with the registration pattern used everywhere else.

- [ ] **Make the public client interface domain-agnostic**
  Remove `MealPlanningServiceClient` from the `Client` interface, struct, and `BuildClient` constructor in `pkg/client/client.go`. Consider a pattern where domain clients are accessed via a separate domain-specific client builder or type assertion so the base client compiles without the domain.

- [ ] **Write a fork/copy guide**
  Document the intended swap process: what to delete, what to replace, how the registration module pattern works, and how to run `make querier` / `make proto` after wiring in a new domain.

---

## Reference: Known Coupling Points

The following infrastructure-adjacent files contain hardcoded meal-planning identifiers that a fork must edit even though they are intended to be generic platform code.

| File | Lines | What's coupled |
|---|---|---|
| `internal/build/services/api/grpc/grpc_service.go` | 10, 29, 47, 65 | `MealPlanningServiceServer` in struct definition and constructor |
| `internal/build/services/api/grpc/extras.go` | 16, 33, 79, 110–170, 226–244 | Meal-planning in `BuildRegistrationFuncs`, `AggregateMethodPermissions`, and DI provider closures |
| `internal/build/services/api/grpc/config.go` | 9, 122–125 | `mealplanningcfg.Config` wired into DI |
| `internal/build/services/api/grpc/build.go` | 15, 160 | `mealplanningregistration` import + call (acceptable pattern, single edit) |
| `internal/build/services/mcp/build.go` | 11, 49 | `mealplanningrepo` directly imported, bypasses registration pattern |
| `internal/build/functions/data_change_message_handler/build.go` | 7, 68 | `mealplanningregistration.RegisterForDataChangeHandler` (acceptable pattern, single edit) |
| `internal/build/functions/data_change_message_handler/searchers.go` | 7, 27–235 | 8 meal-planning searcher providers hardcoded in an infrastructure file |
| `internal/build/jobs/search_data_index_scheduler/build.go` | 8, 11, 46–52 | `mealplanning` domain + repo imported in a generic job |
| `internal/build/jobs/search_data_index_scheduler/indexers.go` | 7–22 | `ProvideIndexFunctions` hardcodes merge of generic + meal-planning indexers |
| `internal/authorization/permissions.go` | 25–268 | All 4 role permission sets contain meal-planning permissions inline |
| `internal/config/configs.go` | 69–71 | `configurations` union type lists 3 meal-planning config types |
| `internal/config/services_config.go` | 10, 27, 40 | `MealPlanning mealplanningcfg.Config` in `ServicesConfig` struct |
| `internal/config/environment.go` | 24–26, 134, 241, 265 | 3 meal-planning path fields + `renderMealPlanningConfigs` call |
| `pkg/client/client.go` | 19, 52, 74, 101 | `MealPlanningServiceClient` in the public `Client` interface and `BuildClient` |
