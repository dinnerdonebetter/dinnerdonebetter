<script>
  import { Link } from "svelte-routing";

  let recipe = {};

  const recipeID = window.location.pathname.replace("/recipes/", "");

  fetch(`/api/v1/recipes/${ recipeID }`)
    .then(response => response.json())
    .then(data => {
      recipe = data;
    });
</script>

<!-- Recipes.svelte -->

<div style="text-align: center; margin: 0px auto;">
  <p>
    name:
    {@html recipe.name}
  </p>
  <p>
    details:
    {@html recipe.details}
  </p>

  <Link to="/recipes">see all</Link>

</div>