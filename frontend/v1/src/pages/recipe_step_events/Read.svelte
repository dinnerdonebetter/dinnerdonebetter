<script>
  import { Link } from "svelte-routing";

  let recipeStepEvent = {};

  const recipeStepEventID = window.location.pathname.replace("/recipe_step_events/", "");

  fetch(`/api/v1/recipe_step_events/${ recipeStepEventID }`)
    .then(response => response.json())
    .then(data => {
      recipeStepEvent = data;
    });
</script>

<!-- RecipeStepEvents.svelte -->

<div style="text-align: center; margin: 0px auto;">
  <p>
    name:
    {@html recipe_step_event.name}
  </p>
  <p>
    details:
    {@html recipe_step_event.details}
  </p>

  <Link to="/recipe_step_events">see all</Link>

</div>