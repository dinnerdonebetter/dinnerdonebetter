<script>
  import { Link } from "svelte-routing";

  let validIngredientPreparation = {};

  const validIngredientPreparationID = window.location.pathname.replace("/valid_ingredient_preparations/", "");

  fetch(`/api/v1/valid_ingredient_preparations/${ validIngredientPreparationID }`)
    .then(response => response.json())
    .then(data => {
      validIngredientPreparation = data;
    });
</script>

<!-- ValidIngredientPreparations.svelte -->

<div style="text-align: center; margin: 0px auto;">
  <p>
    name:
    {@html valid_ingredient_preparation.name}
  </p>
  <p>
    details:
    {@html valid_ingredient_preparation.details}
  </p>

  <Link to="/valid_ingredient_preparations">see all</Link>

</div>