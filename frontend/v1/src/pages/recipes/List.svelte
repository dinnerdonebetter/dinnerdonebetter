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
      title: "Source",
      key: "source"
    },
    {
      title: "Description",
      key: "description"
    },
    {
      title: "InspiredByRecipeID",
      key: "inspired_by_recipe_id"
    },
    {
      title: "Private",
      key: "private"
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
  let recipes = [];

  function deleteRecipe(row) {
    if (confirm("are you sure you want to delete this recipe?")) {
      fetch(`/api/v1/recipes/${row.id}`, {
        method: "DELETE"
      }).then(response => {
        if (response.status != 204) {
          console.error("something has gone awry");
        }
        recipes = recipes.filter(recipe => {
          return recipe.id != row.id;
        });
      });
    }
  }

  function goToRecipe(row) {
    navigate(`/recipes/${row.id}`, { replace: true });
  }

  fetch("/api/v1/recipes/")
    .then(response => response.json())
    .then(data => {
      recipes = data["recipes"];
    });
</script>

<!-- Recipes.svelte -->

<Table
  {columns}
  tableStyle={'margin: 0px auto;'}
  rows={ recipes }
  rowClickFunc={goToRecipe}
  rowDeleteFunc={deleteRecipe} />