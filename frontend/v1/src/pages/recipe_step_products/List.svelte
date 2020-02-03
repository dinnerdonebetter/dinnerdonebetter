<script>
  import { Link, navigate } from "svelte-routing";

  import Table from "../../components/Table.svelte";

  const columns = [
    {
      title: "ID",
      key: "id"
    },
    {
      title: "Name",
      key: "name"
    },
    {
      title: "RecipeStepID",
      key: "recipe_step_id"
    },
    {
      title: "Created On",
      key: "created_on"
    },
    {
      title: "Updated On",
      key: "updated_on"
    }
  ];
  let recipeStepProducts = [];

  function deleteRecipeStepProduct(row) {
    if (confirm("are you sure you want to delete this recipe step product?")) {
      fetch(`/api/v1/recipe_step_products/${row.id}`, {
        method: "DELETE"
      }).then(response => {
        if (response.status != 204) {
          console.error("something has gone awry");
        }
        recipeStepProducts = recipeStepProducts.filter(recipeStepProduct => {
          return recipeStepProduct.id != row.id;
        });
      });
    }
  }

  function goToRecipeStepProduct(row) {
    navigate(`/recipe_step_products/${row.id}`, { replace: true });
  }

  fetch("/api/v1/recipe_step_products/")
    .then(response => response.json())
    .then(data => {
      recipeStepProducts = data["recipe_step_products"];
    });
</script>

<!-- RecipeStepProducts.svelte -->

<Table
  {columns}
  tableStyle={'margin: 0px auto;'}
  rows={ recipeStepProducts }
  rowClickFunc={goToRecipeStepProduct}
  rowDeleteFunc={deleteRecipeStepProduct} />