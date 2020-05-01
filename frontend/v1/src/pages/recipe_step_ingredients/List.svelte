<script>
  import { Link, navigate } from "svelte-routing";

  import Table from "../../components/Table.svelte";

  const columns = [
    {
      title: "ID",
      key: "id"
    },
    {
      title: "ValidIngredientID",
      key: "valid_ingredient_id"
    },
    {
      title: "IngredientNotes",
      key: "ingredient_notes"
    },
    {
      title: "QuantityType",
      key: "quantity_type"
    },
    {
      title: "QuantityValue",
      key: "quantity_value"
    },
    {
      title: "QuantityNotes",
      key: "quantity_notes"
    },
    {
      title: "ProductOfRecipeStepID",
      key: "product_of_recipe_step_id"
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
  let recipeStepIngredients = [];

  function deleteRecipeStepIngredient(row) {
    if (confirm("are you sure you want to delete this recipe step ingredient?")) {
      fetch(`/api/v1/recipe_step_ingredients/${row.id}`, {
        method: "DELETE"
      }).then(response => {
        if (response.status != 204) {
          console.error("something has gone awry");
        }
        recipeStepIngredients = recipeStepIngredients.filter(recipeStepIngredient => {
          return recipeStepIngredient.id != row.id;
        });
      });
    }
  }

  function goToRecipeStepIngredient(row) {
    navigate(`/recipe_step_ingredients/${row.id}`, { replace: true });
  }

  fetch("/api/v1/recipe_step_ingredients/")
    .then(response => response.json())
    .then(data => {
      recipeStepIngredients = data["recipe_step_ingredients"];
    });
</script>

<!-- RecipeStepIngredients.svelte -->

<Table
  {columns}
  tableStyle={'margin: 0px auto;'}
  rows={ recipeStepIngredients }
  rowClickFunc={goToRecipeStepIngredient}
  rowDeleteFunc={deleteRecipeStepIngredient} />