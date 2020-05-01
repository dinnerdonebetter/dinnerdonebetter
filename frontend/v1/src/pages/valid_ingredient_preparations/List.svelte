<script>
  import { Link, navigate } from "svelte-routing";

  import Table from "../../components/Table.svelte";

  const columns = [
    {
      title: "ID",
      key: "id"
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
  let validIngredientPreparations = [];

  function deleteValidIngredientPreparation(row) {
    if (confirm("are you sure you want to delete this valid ingredient preparation?")) {
      fetch(`/api/v1/valid_ingredient_preparations/${row.id}`, {
        method: "DELETE"
      }).then(response => {
        if (response.status != 204) {
          console.error("something has gone awry");
        }
        validIngredientPreparations = validIngredientPreparations.filter(validIngredientPreparation => {
          return validIngredientPreparation.id != row.id;
        });
      });
    }
  }

  function goToValidIngredientPreparation(row) {
    navigate(`/valid_ingredient_preparations/${row.id}`, { replace: true });
  }

  fetch("/api/v1/valid_ingredient_preparations/")
    .then(response => response.json())
    .then(data => {
      validIngredientPreparations = data["valid_ingredient_preparations"];
    });
</script>

<!-- ValidIngredientPreparations.svelte -->

<Table
  {columns}
  tableStyle={'margin: 0px auto;'}
  rows={ validIngredientPreparations }
  rowClickFunc={goToValidIngredientPreparation}
  rowDeleteFunc={deleteValidIngredientPreparation} />