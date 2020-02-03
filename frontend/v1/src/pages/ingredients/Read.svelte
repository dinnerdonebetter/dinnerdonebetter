<script>
  import { Link } from "svelte-routing";

  let ingredient = {};

  const ingredientID = window.location.pathname.replace("/ingredients/", "");

  fetch(`/api/v1/ingredients/${ ingredientID }`)
    .then(response => response.json())
    .then(data => {
      ingredient = data;
    });
</script>

<!-- Ingredients.svelte -->

<div style="text-align: center; margin: 0px auto;">
  <p>
    name:
    {@html ingredient.name}
  </p>
  <p>
    details:
    {@html ingredient.details}
  </p>

  <Link to="/ingredients">see all</Link>

</div>