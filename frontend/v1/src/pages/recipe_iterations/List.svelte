<script>
  import { Link, navigate } from "svelte-routing";

  import Table from "../../components/Table.svelte";

  const columns = [
    {
      title: "ID",
      key: "id"
    },
    {
      title: "RecipeID",
      key: "recipe_id"
    },
    {
      title: "EndDifficultyRating",
      key: "end_difficulty_rating"
    },
    {
      title: "EndComplexityRating",
      key: "end_complexity_rating"
    },
    {
      title: "EndTasteRating",
      key: "end_taste_rating"
    },
    {
      title: "EndOverallRating",
      key: "end_overall_rating"
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
  let recipeIterations = [];

  function deleteRecipeIteration(row) {
    if (confirm("are you sure you want to delete this recipe iteration?")) {
      fetch(`/api/v1/recipe_iterations/${row.id}`, {
        method: "DELETE"
      }).then(response => {
        if (response.status != 204) {
          console.error("something has gone awry");
        }
        recipeIterations = recipeIterations.filter(recipeIteration => {
          return recipeIteration.id != row.id;
        });
      });
    }
  }

  function goToRecipeIteration(row) {
    navigate(`/recipe_iterations/${row.id}`, { replace: true });
  }

  fetch("/api/v1/recipe_iterations/")
    .then(response => response.json())
    .then(data => {
      recipeIterations = data["recipe_iterations"];
    });
</script>

<!-- RecipeIterations.svelte -->

<Table
  {columns}
  tableStyle={'margin: 0px auto;'}
  rows={ recipeIterations }
  rowClickFunc={goToRecipeIteration}
  rowDeleteFunc={deleteRecipeIteration} />