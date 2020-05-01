<script>
  import { Link } from "svelte-routing";

  let recipeIterationStep = {};

  const recipeIterationStepID = window.location.pathname.replace("/recipe_iteration_steps/", "");

  fetch(`/api/v1/recipe_iteration_steps/${ recipeIterationStepID }`)
    .then(response => response.json())
    .then(data => {
      recipeIterationStep = data;
    });
</script>

<!-- RecipeIterationSteps.svelte -->

<div style="text-align: center; margin: 0px auto;">
  <p>
    name:
    {@html recipe_iteration_step.name}
  </p>
  <p>
    details:
    {@html recipe_iteration_step.details}
  </p>

  <Link to="/recipe_iteration_steps">see all</Link>

</div>