<script>
  import { Link, navigate } from "svelte-routing";

  import Table from "../../components/Table.svelte";

  const columns = [
    {
      title: "ID",
      key: "id"
    },
    {
      title: "EventType",
      key: "event_type"
    },
    {
      title: "Done",
      key: "done"
    },
    {
      title: "RecipeIterationID",
      key: "recipe_iteration_id"
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
  let recipeStepEvents = [];

  function deleteRecipeStepEvent(row) {
    if (confirm("are you sure you want to delete this recipe step event?")) {
      fetch(`/api/v1/recipe_step_events/${row.id}`, {
        method: "DELETE"
      }).then(response => {
        if (response.status != 204) {
          console.error("something has gone awry");
        }
        recipeStepEvents = recipeStepEvents.filter(recipeStepEvent => {
          return recipeStepEvent.id != row.id;
        });
      });
    }
  }

  function goToRecipeStepEvent(row) {
    navigate(`/recipe_step_events/${row.id}`, { replace: true });
  }

  fetch("/api/v1/recipe_step_events/")
    .then(response => response.json())
    .then(data => {
      recipeStepEvents = data["recipe_step_events"];
    });
</script>

<!-- RecipeStepEvents.svelte -->

<Table
  {columns}
  tableStyle={'margin: 0px auto;'}
  rows={ recipeStepEvents }
  rowClickFunc={goToRecipeStepEvent}
  rowDeleteFunc={deleteRecipeStepEvent} />