<script>
  import { Link } from "svelte-routing";

  let recipeStep = {};

  const recipeStepID = window.location.pathname.replace("/recipe_steps/", "");

  fetch(`/api/v1/recipe_steps/${ recipeStepID }`)
    .then(response => response.json())
    .then(data => {
      recipeStep = data;
    });
</script>

<!-- RecipeSteps.svelte -->

<div style="text-align: center; margin: 0px auto;">
  <p>
    name:
    {@html recipe_step.name}
  </p>
  <p>
    details:
    {@html recipe_step.details}
  </p>

  <Link to="/recipe_steps">see all</Link>

</div>