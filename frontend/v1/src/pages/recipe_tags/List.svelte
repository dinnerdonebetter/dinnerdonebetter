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
  let recipeTags = [];

  function deleteRecipeTag(row) {
    if (confirm("are you sure you want to delete this recipe tag?")) {
      fetch(`/api/v1/recipe_tags/${row.id}`, {
        method: "DELETE"
      }).then(response => {
        if (response.status != 204) {
          console.error("something has gone awry");
        }
        recipeTags = recipeTags.filter(recipeTag => {
          return recipeTag.id != row.id;
        });
      });
    }
  }

  function goToRecipeTag(row) {
    navigate(`/recipe_tags/${row.id}`, { replace: true });
  }

  fetch("/api/v1/recipe_tags/")
    .then(response => response.json())
    .then(data => {
      recipeTags = data["recipe_tags"];
    });
</script>

<!-- RecipeTags.svelte -->

<Table
  {columns}
  tableStyle={'margin: 0px auto;'}
  rows={ recipeTags }
  rowClickFunc={goToRecipeTag}
  rowDeleteFunc={deleteRecipeTag} />