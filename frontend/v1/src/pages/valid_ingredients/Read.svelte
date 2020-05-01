<script>
  import { Link } from "svelte-routing";

  let validIngredient = {};

  const validIngredientID = window.location.pathname.replace("/valid_ingredients/", "");

  fetch(`/api/v1/valid_ingredients/${ validIngredientID }`)
    .then(response => response.json())
    .then(data => {
      validIngredient = data;
    });
</script>

<!-- ValidIngredients.svelte -->

<div style="text-align: center; margin: 0px auto;">
  <p>
    name:
    {@html valid_ingredient.name}
  </p>
  <p>
    details:
    {@html valid_ingredient.details}
  </p>

  <Link to="/valid_ingredients">see all</Link>

</div>