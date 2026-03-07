# Recipe Conversion Reference

Detailed schema reference, clarifications, edge cases, and hangups for converting text recipes to `RecipeCreationRequestInput`.

## Clarifications and Edge Cases

### Pot vs Vessel

In DDB, **pots are vessels**, not instruments. Use `ValidPreparationVesselID` or vessel product references for pots, pans, skillets, and bowls.

### "Cook According to Package"

Map to a completion condition (e.g., al dente/tender) plus optional `EstimatedTimeInSeconds`. Avoid hardcoding package-specific times; use `IngredientStateID` for doneness.

### Parallel Steps ("Meanwhile, mix...")

These are separate steps. Order by logical dependencies—place the parallel step after the step it supports or can run alongside.

### Optional Steps

Set `Optional: true` when a step may be skipped (e.g., "if you have more beans, measure and reserve the rest").

### Multiple Outputs from One Step

A step can produce multiple products. Example: drain step yields both "drained cooked beans" and "reserved bean-cooking liquid". Use distinct `Index` values for each product.

### Vessel State Names

Use descriptive names to clarify state in the chain:
- "pot with boiling water"
- "empty pot"
- "pot with cooked spaghetti"
- "large saucepan"

### Index vs OptionIndex

- **Index**: Position within the step. Items with the same `Index` belong to the same option group.
- **OptionIndex**: Position within the option group. `0` = primary, `1` = first alternative.
- Use for option groups (butter OR margarine); each regular item gets a unique `Index`.

### ProductOfRecipeStepProductIndex

When a step has multiple products, reference by `ProductOfRecipeStepProductIndex` (0-based). Product at index 0 is the first product, index 1 is the second, etc.

---

## Hangups to Anticipate

### Bridge Table Gaps

New preparations, ingredients, vessels, or instruments may need to be added to `enumerations.go` before conversion. Check `Enumerations` for existing entries; add via the appropriate `createVIP`, `createVPI`, `createVPV` helpers if missing.

### Ambiguous Ordering

When the text says "meanwhile" or "while X cooks", infer dependency order. The parallel step often goes after the step it supports (e.g., mix cheese mixture after pasta starts cooking).

### Implicit Vessels

"In a large pot" implies the pot is a vessel from the start. "Return to pot" implies the same pot—track it through the chain. Do not introduce a new vessel unless the recipe explicitly requires it.

### Reserve Before vs During Drain

- **"Reserve 1 cup before draining"**: Separate step—reserve first, then drain.
- **"Drain, reserving 1 cup"**: Combine in one drain step with two products (drained ingredient + reserved liquid).

### Scaling

Ensure `MeasurementQuantity` and `ItemQuantity` are set correctly:
- **Discrete products** (patties, cookies): Set `ItemQuantity` and `MeasurementQuantity` (per-item). Count scales; per-item size stays constant.
- **Continuous products** (sauce, liquid): Omit `ItemQuantity`. Total quantity scales proportionally.

---

## Prep Tasks

Identify steps that can be done ahead of time. Create `RecipePrepTaskWithinRecipeCreationRequestInput` with:
- `RecipeSteps`: `BelongsToRecipeStepIndex` and `SatisfiesRecipeStep` for each step
- `SatisfiesRecipeStep: true` on the final step of the prep task (the step that "completes" the prep)
- `TimeBufferBeforeRecipeInSeconds`, `StorageType`, `StorageTemperatureInCelsius` as appropriate

---

## Domain Types

- `RecipeCreationRequestInput`: [backend/internal/domain/mealplanning/recipe.go](backend/internal/domain/mealplanning/recipe.go)
- `RecipeStepCreationRequestInput`: [backend/internal/domain/mealplanning/recipe_step.go](backend/internal/domain/mealplanning/recipe_step.go)
- `RecipeStepProductCreationRequestInput`: [backend/internal/domain/mealplanning/recipe_step_product.go](backend/internal/domain/mealplanning/recipe_step_product.go)
- `RecipeStepCompletionConditionCreationRequestInput`: [backend/internal/domain/mealplanning/recipe_step_completion_condition.go](backend/internal/domain/mealplanning/recipe_step_completion_condition.go)

---

## Example Flow: Spaghetti Step

```
Text: "Bring a large pot of water to a boil. Add the spaghetti and cook according to package instructions. Reserve 1 cup of the cooking water. Drain the spaghetti and return to its pot."

Step 0 (boil): Pot + water → pot with boiling water
Step 1 (add): Pot with boiling water + spaghetti → pot with spaghetti in water
Step 2 (cook): Pot with spaghetti → pot with cooked spaghetti (CompletionCondition: tender)
Step 3 (reserve): Ladle 1 cup from pot → reserved pasta water
Step 4 (drain): Colander + cooked spaghetti → drained spaghetti (pot "left behind")
Step 5 (add): Drained spaghetti + pot from step 2 → spaghetti in pot
```

Vessel chain: Step 0 → 1 → 2 (pot). Step 4 uses colander. Step 5 references pot from step 2 (temporary storage exception).
