<script>
  import { Link, navigate } from "svelte-routing";

  import Table from "../../components/Table.svelte";

  const columns = [
    {
      title: "ID",
      key: "id"
    },
    {
      title: "ValidIngredientTagID",
      key: "valid_ingredient_tag_id"
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
  let ingredientTagMappings = [];

  function deleteIngredientTagMapping(row) {
    if (confirm("are you sure you want to delete this ingredient tag mapping?")) {
      fetch(`/api/v1/ingredient_tag_mappings/${row.id}`, {
        method: "DELETE"
      }).then(response => {
        if (response.status != 204) {
          console.error("something has gone awry");
        }
        ingredientTagMappings = ingredientTagMappings.filter(ingredientTagMapping => {
          return ingredientTagMapping.id != row.id;
        });
      });
    }
  }

  function goToIngredientTagMapping(row) {
    navigate(`/ingredient_tag_mappings/${row.id}`, { replace: true });
  }

  fetch("/api/v1/ingredient_tag_mappings/")
    .then(response => response.json())
    .then(data => {
      ingredientTagMappings = data["ingredient_tag_mappings"];
    });
</script>

<!-- IngredientTagMappings.svelte -->

<Table
  {columns}
  tableStyle={'margin: 0px auto;'}
  rows={ ingredientTagMappings }
  rowClickFunc={goToIngredientTagMapping}
  rowDeleteFunc={deleteIngredientTagMapping} />