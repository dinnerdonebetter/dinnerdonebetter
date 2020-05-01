<script>
  import { Link, navigate } from "svelte-routing";

  import Table from "../../components/Table.svelte";

  const columns = [
    {
      title: "ID",
      key: "id"
    },
    {
      title: "StartedOn",
      key: "started_on"
    },
    {
      title: "EndedOn",
      key: "ended_on"
    },
    {
      title: "State",
      key: "state"
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
  let recipeIterationSteps = [];

  function deleteRecipeIterationStep(row) {
    if (confirm("are you sure you want to delete this recipe iteration step?")) {
      fetch(`/api/v1/recipe_iteration_steps/${row.id}`, {
        method: "DELETE"
      }).then(response => {
        if (response.status != 204) {
          console.error("something has gone awry");
        }
        recipeIterationSteps = recipeIterationSteps.filter(recipeIterationStep => {
          return recipeIterationStep.id != row.id;
        });
      });
    }
  }

  function goToRecipeIterationStep(row) {
    navigate(`/recipe_iteration_steps/${row.id}`, { replace: true });
  }

  fetch("/api/v1/recipe_iteration_steps/")
    .then(response => response.json())
    .then(data => {
      recipeIterationSteps = data["recipe_iteration_steps"];
    });
</script>

<!-- RecipeIterationSteps.svelte -->

<Table
  {columns}
  tableStyle={'margin: 0px auto;'}
  rows={ recipeIterationSteps }
  rowClickFunc={goToRecipeIterationStep}
  rowDeleteFunc={deleteRecipeIterationStep} />