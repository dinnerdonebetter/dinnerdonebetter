<script>
  import { Link } from "svelte-routing";

  let recipeStepIngredient = {};

  const recipeStepIngredientID = window.location.pathname.replace("/recipe_step_ingredients/", "");

  fetch(`/api/v1/recipe_step_ingredients/${ recipeStepIngredientID }`)
    .then(response => response.json())
    .then(data => {
      recipeStepIngredient = data;
    });
</script>

<!-- RecipeStepIngredients.svelte -->

<div style="text-align: center; margin: 0px auto;">
  <p>
    name:
    {@html recipe_step_ingredient.name}
  </p>
  <p>
    details:
    {@html recipe_step_ingredient.details}
  </p>

  <Link to="/recipe_step_ingredients">see all</Link>

</div>