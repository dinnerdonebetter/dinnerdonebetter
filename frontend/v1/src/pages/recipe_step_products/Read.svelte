<script>
  import { Link } from "svelte-routing";

  let recipeStepProduct = {};

  const recipeStepProductID = window.location.pathname.replace("/recipe_step_products/", "");

  fetch(`/api/v1/recipe_step_products/${ recipeStepProductID }`)
    .then(response => response.json())
    .then(data => {
      recipeStepProduct = data;
    });
</script>

<!-- RecipeStepProducts.svelte -->

<div style="text-align: center; margin: 0px auto;">
  <p>
    name:
    {@html recipe_step_product.name}
  </p>
  <p>
    details:
    {@html recipe_step_product.details}
  </p>

  <Link to="/recipe_step_products">see all</Link>

</div>