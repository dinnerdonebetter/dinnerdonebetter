<script>
  import { Link } from "svelte-routing";

  let ingredientTagMapping = {};

  const ingredientTagMappingID = window.location.pathname.replace("/ingredient_tag_mappings/", "");

  fetch(`/api/v1/ingredient_tag_mappings/${ ingredientTagMappingID }`)
    .then(response => response.json())
    .then(data => {
      ingredientTagMapping = data;
    });
</script>

<!-- IngredientTagMappings.svelte -->

<div style="text-align: center; margin: 0px auto;">
  <p>
    name:
    {@html ingredient_tag_mapping.name}
  </p>
  <p>
    details:
    {@html ingredient_tag_mapping.details}
  </p>

  <Link to="/ingredient_tag_mappings">see all</Link>

</div>