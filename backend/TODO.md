# Recipe Creation Strictness: Bridge Table ID Implementation

## Overview

This document tracks the implementation of stricter recipe creation that requires bridge table IDs 
(ValidIngredientPreparation, ValidPreparationInstrument, ValidPreparationVessel, ValidIngredientMeasurementUnit) 
instead of raw component IDs.

## Phase 1: Add New Fields (Additive, Non-Breaking)

### 1.1 Update Domain Types - Add Bridge Table ID Fields

- [x] `recipe_step_ingredient.go`
  - [x] Add `ValidIngredientPreparationID *string` to `RecipeStepIngredientCreationRequestInput`
  - [x] Add `ValidIngredientMeasurementUnitID *string` to `RecipeStepIngredientCreationRequestInput`
  - [x] Add `ValidIngredientPreparationID *string` to `RecipeStepIngredientDatabaseCreationInput`
  - [x] Add `ValidIngredientMeasurementUnitID *string` to `RecipeStepIngredientDatabaseCreationInput`

- [x] `recipe_step_instrument.go`
  - [x] Add `ValidPreparationInstrumentID *string` to `RecipeStepInstrumentCreationRequestInput`
  - [x] Add `ValidPreparationInstrumentID *string` to `RecipeStepInstrumentDatabaseCreationInput`

- [x] `recipe_step_vessel.go`
  - [x] Add `ValidPreparationVesselID *string` to `RecipeStepVesselCreationRequestInput`
  - [x] Add `ValidPreparationVesselID *string` to `RecipeStepVesselDatabaseCreationInput`

### 1.2 Add Helper Methods on Recipe Input Types

- [x] `recipe.go` (or `recipe_step.go`)
  - [x] Add `GetAllValidIngredientPreparationIDs() []string` to `RecipeDatabaseCreationInput`
  - [x] Add `GetAllValidIngredientMeasurementUnitIDs() []string` to `RecipeDatabaseCreationInput`
  - [x] Add `GetAllValidPreparationInstrumentIDs() []string` to `RecipeDatabaseCreationInput`
  - [x] Add `GetAllValidPreparationVesselIDs() []string` to `RecipeDatabaseCreationInput`
  - [x] Write unit tests for these helper methods

### 1.3 Update Converters

- [x] `converters/recipe_step_ingredients.go`
  - [x] Pass through `ValidIngredientPreparationID` from request to database input
  - [x] Pass through `ValidIngredientMeasurementUnitID` from request to database input

- [x] `converters/recipe_step_instruments.go`
  - [x] Pass through `ValidPreparationInstrumentID` from request to database input

- [x] `converters/recipe_step_vessels.go`
  - [x] Pass through `ValidPreparationVesselID` from request to database input

### 1.4 Update Fakes

- [x] `fakes/recipe_step_ingredient.go`
  - [x] Update fake builders to optionally include new bridge table ID fields

- [x] `fakes/recipe_step_instrument.go`
  - [x] Update fake builders to optionally include new bridge table ID fields

- [x] `fakes/recipe_step_vessel.go`
  - [x] Update fake builders to optionally include new bridge table ID fields

### 1.5 Verify Phase 1

- [ ] Run `make format lint` - should pass
- [ ] Run unit tests - should pass (new fields are optional)
- [ ] Run integration tests - should pass (new fields are optional)

---

## Phase 2: Add Batch Query Methods for Bridge Tables

### 2.1 Update Data Manager Interfaces

- [ ] `valid_ingredient_preparation.go`
  - [ ] Add `GetValidIngredientPreparationsByIDs(ctx context.Context, ids []string) (map[string]*ValidIngredientPreparation, error)` to interface

- [ ] `valid_ingredient_measurement_unit.go`
  - [ ] Add `GetValidIngredientMeasurementUnitsByIDs(ctx context.Context, ids []string) (map[string]*ValidIngredientMeasurementUnit, error)` to interface

- [ ] `valid_preparation_instrument.go`
  - [ ] Add `GetValidPreparationInstrumentsByIDs(ctx context.Context, ids []string) (map[string]*ValidPreparationInstrument, error)` to interface

- [ ] `valid_preparation_vessel.go`
  - [ ] Add `GetValidPreparationVesselsByIDs(ctx context.Context, ids []string) (map[string]*ValidPreparationVessel, error)` to interface

### 2.2 Add SQL Queries

- [ ] `codegen/queries/valid_ingredient_preparations.go`
  - [ ] Add query for `GetValidIngredientPreparationsByIDs` (using `ANY($1::text[])`)

- [ ] `codegen/queries/valid_ingredient_measurement_units.go`
  - [ ] Add query for `GetValidIngredientMeasurementUnitsByIDs`

- [ ] `codegen/queries/valid_preparation_instruments.go`
  - [ ] Add query for `GetValidPreparationInstrumentsByIDs`

- [ ] `codegen/queries/valid_preparation_vessels.go`
  - [ ] Add query for `GetValidPreparationVesselsByIDs`

- [ ] Run codegen to generate query code

### 2.3 Implement Repository Methods

- [ ] `postgres/mealplanning/valid_ingredient_preparations.go`
  - [ ] Implement `GetValidIngredientPreparationsByIDs`
  - [ ] Return `map[string]*ValidIngredientPreparation` keyed by ID

- [ ] `postgres/mealplanning/valid_ingredient_measurement_units.go`
  - [ ] Implement `GetValidIngredientMeasurementUnitsByIDs`

- [ ] `postgres/mealplanning/valid_preparation_instruments.go`
  - [ ] Implement `GetValidPreparationInstrumentsByIDs`

- [ ] `postgres/mealplanning/valid_preparation_vessels.go`
  - [ ] Implement `GetValidPreparationVesselsByIDs`

### 2.4 Add Unit Tests for Batch Query Methods

- [ ] `postgres/mealplanning/valid_ingredient_preparations_test.go`
  - [ ] Test `GetValidIngredientPreparationsByIDs` with valid IDs
  - [ ] Test with empty list
  - [ ] Test with non-existent IDs (should return partial results)

- [ ] Same for other three bridge tables

### 2.5 Verify Phase 2

- [ ] Run `make format lint` - should pass
- [ ] Run unit tests - should pass
- [ ] Run integration tests - should pass

---

## Phase 3: Implement RecipeValidator

### 3.1 Create Validator

- [ ] Create `recipe_validator.go` in `internal/domain/mealplanning/`
  - [ ] Define `RecipeValidator` struct with map fields for each bridge table type
  - [ ] Implement `NewRecipeValidator(...)` constructor
  - [ ] Implement `ValidateAndPopulate(input *RecipeDatabaseCreationInput) error`
  - [ ] Implement `validateStep(...)` 
  - [ ] Implement `validateAndPopulateIngredient(...)` - validates VIP + VIMU, populates IngredientID + MeasurementUnitID
  - [ ] Implement `validateAndPopulateInstrument(...)` - validates VPI, populates InstrumentID
  - [ ] Implement `validateAndPopulateVessel(...)` - validates VPV, populates VesselID
  - [ ] Handle skip logic for recipe step products (outputs from previous steps)

### 3.2 Write Validator Unit Tests

- [ ] Create `recipe_validator_test.go`
  - [ ] Test valid recipe with all bridge IDs present and matching
  - [ ] Test missing ValidIngredientPreparationID
  - [ ] Test ValidIngredientPreparation with wrong preparation
  - [ ] Test ValidIngredientMeasurementUnit for wrong ingredient
  - [ ] Test missing ValidPreparationInstrumentID
  - [ ] Test ValidPreparationInstrument with wrong preparation
  - [ ] Test missing ValidPreparationVesselID
  - [ ] Test ValidPreparationVessel with wrong preparation
  - [ ] Test recipe step products are skipped (no validation required)
  - [ ] Test that derived IDs are correctly populated after validation

### 3.3 Verify Phase 3

- [ ] Run `make format lint` - should pass
- [ ] Run unit tests - should pass
- [ ] Run integration tests - should pass (validator not yet wired up)

---

## Phase 4: Wire Up Validation in Repository (Validation Only)

### 4.1 Update Recipe Repository

- [ ] `postgres/mealplanning/recipes.go`
  - [ ] In `CreateRecipe`, after receiving input:
    - [ ] Collect bridge table IDs using helper methods
    - [ ] Only proceed with validation if any bridge table IDs are present
    - [ ] Batch fetch bridge table records
    - [ ] Create `RecipeValidator` with fetched maps
    - [ ] Call `validator.ValidateAndPopulate(input)`
    - [ ] If validation fails, return error
    - [ ] If validation passes, derived IDs are now populated, continue with existing insert logic

### 4.2 Update Repository Dependencies

- [ ] Ensure recipe repository has access to bridge table query methods
  - [ ] May need to inject additional repository dependencies or use a shared querier

### 4.3 Verify Phase 4 (Backward Compatible)

- [ ] Run `make format lint` - should pass
- [ ] Run unit tests - should pass
- [ ] Run integration tests - should pass (still using old fields, new fields optional)

---

## Phase 5: Update Integration Tests to Use Bridge Table IDs

### 5.1 Audit Existing Integration Tests

- [ ] Identify all integration tests that create recipes
  - [ ] `tests_integration/apiserver/recipes_test.go`
  - [ ] Any other files creating recipes

### 5.2 Update Test Helpers/Fixtures

- [ ] Ensure test setup creates necessary bridge table entries:
  - [ ] ValidIngredientPreparation entries for all ingredient+preparation combos used
  - [ ] ValidIngredientMeasurementUnit entries for all ingredient+unit combos used
  - [ ] ValidPreparationInstrument entries for all preparation+instrument combos used
  - [ ] ValidPreparationVessel entries for all preparation+vessel combos used

- [ ] Create helper functions to look up bridge table IDs in tests:
  - [ ] `getValidIngredientPreparationID(preparationID, ingredientID) string`
  - [ ] `getValidIngredientMeasurementUnitID(ingredientID, unitID) string`
  - [ ] `getValidPreparationInstrumentID(preparationID, instrumentID) string`
  - [ ] `getValidPreparationVesselID(preparationID, vesselID) string`

### 5.3 Migrate Integration Tests

- [ ] Update recipe creation tests to:
  - [ ] First create bridge table entries (or use pre-existing seed data)
  - [ ] Use bridge table IDs in recipe creation requests
  - [ ] Remove raw IngredientID, InstrumentID, VesselID, MeasurementUnitID from requests

### 5.4 Verify Phase 5

- [ ] Run integration tests - should pass with new bridge table IDs
- [ ] Verify validation errors are returned when bridge table IDs are invalid

---

## Phase 6: Update Bootstrap Code

### 6.1 Update Enumerations

- [ ] `bootstrapping/enumerations.go`
  - [ ] Add maps for bridge table lookups:
    - [ ] `IngredientPreparations map[string]map[string]*ValidIngredientPreparation` (keyed by [preparation][ingredient])
    - [ ] `IngredientMeasurementUnits map[string]map[string]*ValidIngredientMeasurementUnit` (keyed by [ingredient][unit])
    - [ ] `PreparationInstruments map[string]map[string]*ValidPreparationInstrument` (keyed by [preparation][instrument])
    - [ ] `PreparationVessels map[string]map[string]*ValidPreparationVessel` (keyed by [preparation][vessel])
  - [ ] Populate these maps during enumeration loading

### 6.2 Update Bootstrap Recipes

- [ ] `bootstrapping/recipe_refried_beans.go`
  - [ ] Replace `IngredientID: pointer.To(ingredientMap["garlic"].ID)` with bridge table ID lookups
  - [ ] Replace `MeasurementUnitID: unitMeasurementUnit.ID` with bridge table ID lookups
  - [ ] Replace `InstrumentID: pointer.To(instruments["knife"].ID)` with bridge table ID lookups
  - [ ] Replace `VesselID: pointer.To(vessels["pot"].ID)` with bridge table ID lookups

- [ ] `bootstrapping/recipe_pay_de_elote.go`
  - [ ] Same updates

- [ ] Any other bootstrap recipes

### 6.3 Ensure Bridge Table Seed Data Exists

- [ ] Audit what ingredient+preparation combinations are used in bootstrap recipes
- [ ] Audit what ingredient+unit combinations are used
- [ ] Audit what preparation+instrument combinations are used
- [ ] Audit what preparation+vessel combinations are used
- [ ] Add any missing bridge table entries to seed data

### 6.4 Verify Phase 6

- [ ] Run bootstrap code locally
- [ ] Verify recipes are created successfully
- [ ] Run full test suite

---

## Phase 7: Make Bridge Table IDs Required, Remove Old Fields

### 7.1 Update Validation to Require Bridge Table IDs

- [ ] `recipe_step_ingredient.go`
  - [ ] Update `ValidateWithContext` to require `ValidIngredientPreparationID` (when not a recipe step product)
  - [ ] Update `ValidateWithContext` to require `ValidIngredientMeasurementUnitID` (when not a recipe step product)
  - [ ] Remove `IngredientID` requirement from validation
  - [ ] Remove `MeasurementUnitID` requirement from validation

- [ ] `recipe_step_instrument.go`
  - [ ] Update `ValidateWithContext` to require `ValidPreparationInstrumentID` (when not a recipe step product)
  - [ ] Remove `InstrumentID` requirement from validation

- [ ] `recipe_step_vessel.go`
  - [ ] Update `ValidateWithContext` to require `ValidPreparationVesselID` (when not a recipe step product)
  - [ ] Remove `VesselID` requirement from validation

### 7.2 Remove Old Fields from Request Inputs

- [ ] `recipe_step_ingredient.go`
  - [ ] Remove `IngredientID *string` from `RecipeStepIngredientCreationRequestInput`
  - [ ] Remove `MeasurementUnitID string` from `RecipeStepIngredientCreationRequestInput`

- [ ] `recipe_step_instrument.go`
  - [ ] Remove `InstrumentID *string` from `RecipeStepInstrumentCreationRequestInput`

- [ ] `recipe_step_vessel.go`
  - [ ] Remove `VesselID *string` from `RecipeStepVesselCreationRequestInput`

### 7.3 Update Converters

- [ ] Remove code that copies old fields
- [ ] Ensure only bridge table IDs are passed through

### 7.4 Update All Tests

- [ ] Fix any unit tests still using old fields
- [ ] Fix any integration tests still using old fields
- [ ] Update fakes to not generate old fields

### 7.5 Final Verification

- [ ] Run `make format lint` - should pass
- [ ] Run all unit tests - should pass
- [ ] Run all integration tests - should pass
- [ ] Manual testing of recipe creation flow

---

## Phase 8: Cleanup and Documentation

### 8.1 Update API Documentation

- [ ] Document new required fields in API docs
- [ ] Document that bridge table IDs are now required for recipe creation
- [ ] Add examples showing the new request format

### 8.2 Update recipes.md

- [ ] Document the validation behavior
- [ ] Document error messages for invalid bridge table IDs
- [ ] Update "Creating a New Recipe" section

### 8.3 Code Cleanup

- [ ] Remove any TODO comments added during migration
- [ ] Review for dead code
- [ ] Ensure consistent error messages

### 8.4 Final Review

- [ ] Code review of all changes
- [ ] Run full test suite one more time
- [ ] Tag/release if appropriate

---

## Notes

### Handling Recipe Step Products

Recipe step products (outputs from previous steps used as inputs to later steps) do NOT need bridge table IDs because:
- They're not ValidIngredients - they're ephemeral products of the recipe
- The preparation compatibility was already validated when the product was created
- They're identified by `RecipeStepProductID` or `ProductOfRecipeStepIndex`/`ProductOfRecipeStepProductIndex`

The validator should skip validation for any ingredient/instrument/vessel that has a `RecipeStepProductID` or product index set.

### Error Message Format

Suggest consistent error format:
```
step {stepIndex} ingredient {ingredientIndex}: {specific error}
step {stepIndex} instrument {instrumentIndex}: {specific error}
step {stepIndex} vessel {vesselIndex}: {specific error}
```

### Bridge Table Lookup Maps in Bootstrap

For bootstrap code, consider a helper like:
```go
func (e *Enumerations) GetVIP(preparation, ingredient string) string {
    if e.IngredientPreparations[preparation] == nil {
        panic(fmt.Sprintf("no preparations found for %q", preparation))
    }
    vip := e.IngredientPreparations[preparation][ingredient]
    if vip == nil {
        panic(fmt.Sprintf("no ValidIngredientPreparation for %q + %q", preparation, ingredient))
    }
    return vip.ID
}
```

This fails fast during development if bridge table data is missing.
