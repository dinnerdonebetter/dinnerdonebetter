<script>
  import { Link } from "svelte-routing";

  let validIngredientTag = {};

  const validIngredientTagID = window.location.pathname.replace("/valid_ingredient_tags/", "");

  fetch(`/api/v1/valid_ingredient_tags/${ validIngredientTagID }`)
    .then(response => response.json())
    .then(data => {
      validIngredientTag = data;
    });
</script>

<!-- ValidIngredientTags.svelte -->

<div style="text-align: center; margin: 0px auto;">
  <p>
    name:
    {@html valid_ingredient_tag.name}
  </p>
  <p>
    details:
    {@html valid_ingredient_tag.details}
  </p>

  <Link to="/valid_ingredient_tags">see all</Link>

</div>