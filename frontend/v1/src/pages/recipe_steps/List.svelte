<script>
  import { Link, navigate } from "svelte-routing";

  import Table from "../../components/Table.svelte";

  const columns = [
    {
      title: "ID",
      key: "id"
    },
    {
      title: "Index",
      key: "index"
    },
    {
      title: "ValidPreparationID",
      key: "valid_preparation_id"
    },
    {
      title: "PrerequisiteStepID",
      key: "prerequisite_step_id"
    },
    {
      title: "MinEstimatedTimeInSeconds",
      key: "min_estimated_time_in_seconds"
    },
    {
      title: "MaxEstimatedTimeInSeconds",
      key: "max_estimated_time_in_seconds"
    },
    {
      title: "YieldsProductName",
      key: "yields_product_name"
    },
    {
      title: "YieldsQuantity",
      key: "yields_quantity"
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
  let recipeSteps = [];

  function deleteRecipeStep(row) {
    if (confirm("are you sure you want to delete this recipe step?")) {
      fetch(`/api/v1/recipe_steps/${row.id}`, {
        method: "DELETE"
      }).then(response => {
        if (response.status != 204) {
          console.error("something has gone awry");
        }
        recipeSteps = recipeSteps.filter(recipeStep => {
          return recipeStep.id != row.id;
        });
      });
    }
  }

  function goToRecipeStep(row) {
    navigate(`/recipe_steps/${row.id}`, { replace: true });
  }

  fetch("/api/v1/recipe_steps/")
    .then(response => response.json())
    .then(data => {
      recipeSteps = data["recipe_steps"];
    });
</script>

<!-- RecipeSteps.svelte -->

<Table
  {columns}
  tableStyle={'margin: 0px auto;'}
  rows={ recipeSteps }
  rowClickFunc={goToRecipeStep}
  rowDeleteFunc={deleteRecipeStep} />