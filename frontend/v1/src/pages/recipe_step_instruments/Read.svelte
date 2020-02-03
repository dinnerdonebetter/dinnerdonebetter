<script>
  import { Link } from "svelte-routing";

  let recipeStepInstrument = {};

  const recipeStepInstrumentID = window.location.pathname.replace("/recipe_step_instruments/", "");

  fetch(`/api/v1/recipe_step_instruments/${ recipeStepInstrumentID }`)
    .then(response => response.json())
    .then(data => {
      recipeStepInstrument = data;
    });
</script>

<!-- RecipeStepInstruments.svelte -->

<div style="text-align: center; margin: 0px auto;">
  <p>
    name:
    {@html recipe_step_instrument.name}
  </p>
  <p>
    details:
    {@html recipe_step_instrument.details}
  </p>

  <Link to="/recipe_step_instruments">see all</Link>

</div>