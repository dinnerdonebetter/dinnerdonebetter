<script>
  import { Link, navigate } from "svelte-routing";

  import Table from "../../components/Table.svelte";

  const columns = [
    {
      title: "ID",
      key: "id"
    },
    {
      title: "InstrumentID",
      key: "instrument_id"
    },
    {
      title: "RecipeStepID",
      key: "recipe_step_id"
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
  let recipeStepInstruments = [];

  function deleteRecipeStepInstrument(row) {
    if (confirm("are you sure you want to delete this recipe step instrument?")) {
      fetch(`/api/v1/recipe_step_instruments/${row.id}`, {
        method: "DELETE"
      }).then(response => {
        if (response.status != 204) {
          console.error("something has gone awry");
        }
        recipeStepInstruments = recipeStepInstruments.filter(recipeStepInstrument => {
          return recipeStepInstrument.id != row.id;
        });
      });
    }
  }

  function goToRecipeStepInstrument(row) {
    navigate(`/recipe_step_instruments/${row.id}`, { replace: true });
  }

  fetch("/api/v1/recipe_step_instruments/")
    .then(response => response.json())
    .then(data => {
      recipeStepInstruments = data["recipe_step_instruments"];
    });
</script>

<!-- RecipeStepInstruments.svelte -->

<Table
  {columns}
  tableStyle={'margin: 0px auto;'}
  rows={ recipeStepInstruments }
  rowClickFunc={goToRecipeStepInstrument}
  rowDeleteFunc={deleteRecipeStepInstrument} />