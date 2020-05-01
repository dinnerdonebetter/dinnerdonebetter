<script>
  import { Link } from "svelte-routing";

  let recipeTag = {};

  const recipeTagID = window.location.pathname.replace("/recipe_tags/", "");

  fetch(`/api/v1/recipe_tags/${ recipeTagID }`)
    .then(response => response.json())
    .then(data => {
      recipeTag = data;
    });
</script>

<!-- RecipeTags.svelte -->

<div style="text-align: center; margin: 0px auto;">
  <p>
    name:
    {@html recipe_tag.name}
  </p>
  <p>
    details:
    {@html recipe_tag.details}
  </p>

  <Link to="/recipe_tags">see all</Link>

</div>