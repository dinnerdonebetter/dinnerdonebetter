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
      title: "Variant",
      key: "variant"
    },
    {
      title: "Description",
      key: "description"
    },
    {
      title: "Warning",
      key: "warning"
    },
    {
      title: "ContainsEgg",
      key: "contains_egg"
    },
    {
      title: "ContainsDairy",
      key: "contains_dairy"
    },
    {
      title: "ContainsPeanut",
      key: "contains_peanut"
    },
    {
      title: "ContainsTreeNut",
      key: "contains_tree_nut"
    },
    {
      title: "ContainsSoy",
      key: "contains_soy"
    },
    {
      title: "ContainsWheat",
      key: "contains_wheat"
    },
    {
      title: "ContainsShellfish",
      key: "contains_shellfish"
    },
    {
      title: "ContainsSesame",
      key: "contains_sesame"
    },
    {
      title: "ContainsFish",
      key: "contains_fish"
    },
    {
      title: "ContainsGluten",
      key: "contains_gluten"
    },
    {
      title: "AnimalFlesh",
      key: "animal_flesh"
    },
    {
      title: "AnimalDerived",
      key: "animal_derived"
    },
    {
      title: "ConsideredStaple",
      key: "considered_staple"
    },
    {
      title: "Icon",
      key: "icon"
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
  let validIngredients = [];

  function deleteValidIngredient(row) {
    if (confirm("are you sure you want to delete this valid ingredient?")) {
      fetch(`/api/v1/valid_ingredients/${row.id}`, {
        method: "DELETE"
      }).then(response => {
        if (response.status != 204) {
          console.error("something has gone awry");
        }
        validIngredients = validIngredients.filter(validIngredient => {
          return validIngredient.id != row.id;
        });
      });
    }
  }

  function goToValidIngredient(row) {
    navigate(`/valid_ingredients/${row.id}`, { replace: true });
  }

  fetch("/api/v1/valid_ingredients/")
    .then(response => response.json())
    .then(data => {
      validIngredients = data["valid_ingredients"];
    });
</script>

<!-- ValidIngredients.svelte -->

<Table
  {columns}
  tableStyle={'margin: 0px auto;'}
  rows={ validIngredients }
  rowClickFunc={goToValidIngredient}
  rowDeleteFunc={deleteValidIngredient} />