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
      title: "Created On",
      key: "created_on"
    },
    {
      title: "Updated On",
      key: "updated_on"
    }
  ];
  let validIngredientTags = [];

  function deleteValidIngredientTag(row) {
    if (confirm("are you sure you want to delete this valid ingredient tag?")) {
      fetch(`/api/v1/valid_ingredient_tags/${row.id}`, {
        method: "DELETE"
      }).then(response => {
        if (response.status != 204) {
          console.error("something has gone awry");
        }
        validIngredientTags = validIngredientTags.filter(validIngredientTag => {
          return validIngredientTag.id != row.id;
        });
      });
    }
  }

  function goToValidIngredientTag(row) {
    navigate(`/valid_ingredient_tags/${row.id}`, { replace: true });
  }

  fetch("/api/v1/valid_ingredient_tags/")
    .then(response => response.json())
    .then(data => {
      validIngredientTags = data["valid_ingredient_tags"];
    });
</script>

<!-- ValidIngredientTags.svelte -->

<Table
  {columns}
  tableStyle={'margin: 0px auto;'}
  rows={ validIngredientTags }
  rowClickFunc={goToValidIngredientTag}
  rowDeleteFunc={deleteValidIngredientTag} />