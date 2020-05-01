<script>
  import { Link } from "svelte-routing";

  let recipeStepPreparation = {};

  const recipeStepPreparationID = window.location.pathname.replace("/recipe_step_preparations/", "");

  fetch(`/api/v1/recipe_step_preparations/${ recipeStepPreparationID }`)
    .then(response => response.json())
    .then(data => {
      recipeStepPreparation = data;
    });
</script>

<!-- RecipeStepPreparations.svelte -->

<div style="text-align: center; margin: 0px auto;">
  <p>
    name:
    {@html recipe_step_preparation.name}
  </p>
  <p>
    details:
    {@html recipe_step_preparation.details}
  </p>

  <Link to="/recipe_step_preparations">see all</Link>

</div>