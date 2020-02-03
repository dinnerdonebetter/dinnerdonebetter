<script>
  import { Link } from "svelte-routing";

  let recipeIteration = {};

  const recipeIterationID = window.location.pathname.replace("/recipe_iterations/", "");

  fetch(`/api/v1/recipe_iterations/${ recipeIterationID }`)
    .then(response => response.json())
    .then(data => {
      recipeIteration = data;
    });
</script>

<!-- RecipeIterations.svelte -->

<div style="text-align: center; margin: 0px auto;">
  <p>
    name:
    {@html recipe_iteration.name}
  </p>
  <p>
    details:
    {@html recipe_iteration.details}
  </p>

  <Link to="/recipe_iterations">see all</Link>

</div>