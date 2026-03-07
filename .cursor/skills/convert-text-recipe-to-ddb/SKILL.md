---
name: convert-text-recipe-to-ddb
description: Converts text recipes into []*mealplanning.RecipeCreationRequestInput for the Dinner Done Better schema. Use when converting prose recipes, adding bootstrap recipes, or when the user asks about recipe step decomposition, vessel chaining, or the DDB recipe schema.
---

# Convert Text Recipe to RecipeCreationRequestInput

Converts prose recipe instructions into `[]*mealplanning.RecipeCreationRequestInput`. Conversion is **not** a 1:1 mapping—text steps must be decomposed into atomic, single-action steps with proper product chaining.

## Quick Start

1. Decompose each text step into discrete actions (e.g., one text step → many DDB steps)
2. Chain vessels through steps; each vessel product used in at most one subsequent step
3. Add completion conditions for state-based doneness (al dente, browned, at temperature)
4. Use bridge table entries from `Enumerations`; products from prior steps use `ProductOfRecipeStepIndex`/`ProductOfRecipeStepProductIndex`

## Reference

- Schema: [docs/recipes.md](docs/recipes.md)
- Examples: [backend/internal/domain/mealplanning/bootstrap/](backend/internal/domain/mealplanning/bootstrap/) - `recipe_0003_whole_roast_chicken.go`, `recipe_0018_stovetop_mac_and_cheese.go`, `recipe_0015_refried_beans.go`, `recipe_0006_simple_white_rice.go`
- For detailed schema, edge cases, and hangups: [reference.md](reference.md)

---

## Step Decomposition (Critical)

**Rule**: One text step often maps to many DDB steps. Break by distinct actions.

**Example**:
> "Bring a large pot of water to a boil. Add the spaghetti and cook according to package instructions. Reserve 1 cup of the cooking water. Drain the spaghetti and return to its pot."

**Decomposed into 6 steps**:
1. **Boil water** - preparation: boil; vessel: pot; yields: pot with boiling water
2. **Add spaghetti** - preparation: add; consumes: pot with boiling water; yields: pot with spaghetti in water
3. **Cook until al dente** - preparation: cook; completion condition: tender/al dente
4. **Reserve cooking liquid** - preparation: reserve; yields: reserved pasta water (1 cup)
5. **Drain spaghetti** - preparation: drain; uses colander; yields: drained spaghetti
6. **Return spaghetti to pot** - preparation: add; consumes: drained spaghetti + empty pot; yields: spaghetti in pot

**Heuristics**:
- Each step = one preparation method, one primary action
- "And" / "then" / "before" / "while" often indicate step boundaries
- Completion conditions get their own step with `CompletionConditions`
- Reserving liquid is its own step when explicitly mentioned
- Drain and return-to-vessel are separate steps

--

## Ingredient decomposition (critical)

If a recipe step calls for a minced clove of garlic, don't just put garlic as the ingredient with "minced" in the description, you must create a preceeding step that minces the garlic.

---

## Vessel and Instrument Chaining

**Vessels** (pots, pans, bowls, skillets):
- First use: `ValidPreparationVesselID` + output vessel product
- Subsequent uses: `ProductOfRecipeStepIndex` + `ProductOfRecipeStepProductIndex` from prior step
- **One vessel product → at most one subsequent consumer** (linear chain)

**Exception - Drain-and-return (temporary storage)**:
- When a step (e.g., drain) consumes a vessel but does not output it, a later step may reference the **last step that output that vessel**
- Example: Mac and cheese drain uses colander; "return to saucepan" uses saucepan from the rest step (before drain)
- See `recipe_0018_stovetop_mac_and_cheese.go` steps 9–10

**Instruments** (knife, spoon, tongs): Same chaining rules apply. **Pots are vessels, not instruments.**

---

## Aluminum Foil is an ingredient

Even when it's technically the cooking surface, aluminum foil is an ingredient, because a user following that recipe might have to purchase it.

---

## Reserve and Drain Patterns

**Reserve liquid**:
- Can be combined with drain in one step when natural ("drain, reserving liquid")
- Or separate step if reserve happens before drain
- Produces ingredient product: "reserved X liquid" with `MeasurementQuantity` and `MeasurementUnitID`
- See `recipe_0015_refried_beans.go` step 8

**Drain**:
- Uses colander (or similar) as vessel
- Consumes ingredient from prior step
- Produces: drained ingredient; optionally reserved liquid
- Original vessel (pot) may be "left behind" for a later return step—reference the last step that output it

---

## Completion Conditions

Use when a step is done based on state, not just time:
- `IngredientStateID` from `enums.IngredientStates`
- `Ingredients`: slice of ingredient indices in the step (0-based)
- `Notes`: human-readable condition (e.g., "pasta should be barely al dente")

Common states: `tender`, `browned`, `at temperature`, `shimmering`, `at desired consistency`, `pliable`, `dissolved`, `lightly charred`, `coated`, `combined`, `baked`.

---

## Bridge Table Requirements

- **Ingredients**: `ValidIngredientPreparationID` + `ValidIngredientMeasurementUnitID` (or `ProductOfRecipeStepIndex`/`ProductOfRecipeStepProductIndex` for products)
- **Instruments**: `ValidPreparationInstrumentID`
- **Vessels**: `ValidPreparationVesselID` (or product refs)
- Bridge entries must exist in `Enumerations`; use `enums.IngredientPreparations[prepID][ingredientID]`, etc.
- Products from prior steps skip bridge validation

---

## Product Types

- `ingredient`: food/liquid output
- `vessel`: container output (pot, skillet, etc.)
- `instrument`: tool output (rare)

Discrete vs continuous: set `ItemQuantity` for countable items (patties, cookies); omit for bulk (sauce, liquid).

-- 

## Prep tasks

A critical value proposition of the service is that we will save you time when cooking. Many components of recipes can be prepared ahead of time if stored properly. For instance, you can dice an onion up to 72 hours ahead of cook time if you put it in the fridge. Many recipes will fail to call out prep task opportunities, you should seek them out and document them proactively.

---

## Pre-Conversion Checklist

- [ ] Minimum 2 steps per recipe
- [ ] Each step: at least one instrument OR vessel
- [ ] All ingredients/instruments/vessels have valid bridge refs or product refs
- [ ] Vessel chain: each vessel product used in at most one subsequent step (except drain-and-return)
- [ ] Completion conditions reference valid `IngredientStateID` and correct ingredient indices

---

## Post-Conversion Checklist

- [ ] no step contains compound instructions
- [ ] all steps preparations match that of their instructions
- [ ] `make dry_run` succeeds from the backend folder
