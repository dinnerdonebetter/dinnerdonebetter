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
      title: "AllergyWarning",
      key: "allergy_warning"
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
  let preparations = [];

  function deletePreparation(row) {
    if (confirm("are you sure you want to delete this preparation?")) {
      fetch(`/api/v1/preparations/${row.id}`, {
        method: "DELETE"
      }).then(response => {
        if (response.status != 204) {
          console.error("something has gone awry");
        }
        preparations = preparations.filter(preparation => {
          return preparation.id != row.id;
        });
      });
    }
  }

  function goToPreparation(row) {
    navigate(`/preparations/${row.id}`, { replace: true });
  }

  fetch("/api/v1/preparations/")
    .then(response => response.json())
    .then(data => {
      preparations = data["preparations"];
    });
</script>

<!-- Preparations.svelte -->

<Table
  {columns}
  tableStyle={'margin: 0px auto;'}
  rows={ preparations }
  rowClickFunc={goToPreparation}
  rowDeleteFunc={deletePreparation} />