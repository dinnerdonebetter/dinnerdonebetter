<script>
  import { Link, navigate } from "svelte-routing";

  import Table from "../../components/Table.svelte";

  const columns = [
    {
      title: "ID",
      key: "id"
    },
    {
      title: "ValidPreparationID",
      key: "valid_preparation_id"
    },
    {
      title: "Notes",
      key: "notes"
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
  let recipeStepPreparations = [];

  function deleteRecipeStepPreparation(row) {
    if (confirm("are you sure you want to delete this recipe step preparation?")) {
      fetch(`/api/v1/recipe_step_preparations/${row.id}`, {
        method: "DELETE"
      }).then(response => {
        if (response.status != 204) {
          console.error("something has gone awry");
        }
        recipeStepPreparations = recipeStepPreparations.filter(recipeStepPreparation => {
          return recipeStepPreparation.id != row.id;
        });
      });
    }
  }

  function goToRecipeStepPreparation(row) {
    navigate(`/recipe_step_preparations/${row.id}`, { replace: true });
  }

  fetch("/api/v1/recipe_step_preparations/")
    .then(response => response.json())
    .then(data => {
      recipeStepPreparations = data["recipe_step_preparations"];
    });
</script>

<!-- RecipeStepPreparations.svelte -->

<Table
  {columns}
  tableStyle={'margin: 0px auto;'}
  rows={ recipeStepPreparations }
  rowClickFunc={goToRecipeStepPreparation}
  rowDeleteFunc={deleteRecipeStepPreparation} />