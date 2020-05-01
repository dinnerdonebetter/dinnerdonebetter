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
      title: "Description",
      key: "description"
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
  let validPreparations = [];

  function deleteValidPreparation(row) {
    if (confirm("are you sure you want to delete this valid preparation?")) {
      fetch(`/api/v1/valid_preparations/${row.id}`, {
        method: "DELETE"
      }).then(response => {
        if (response.status != 204) {
          console.error("something has gone awry");
        }
        validPreparations = validPreparations.filter(validPreparation => {
          return validPreparation.id != row.id;
        });
      });
    }
  }

  function goToValidPreparation(row) {
    navigate(`/valid_preparations/${row.id}`, { replace: true });
  }

  fetch("/api/v1/valid_preparations/")
    .then(response => response.json())
    .then(data => {
      validPreparations = data["valid_preparations"];
    });
</script>

<!-- ValidPreparations.svelte -->

<Table
  {columns}
  tableStyle={'margin: 0px auto;'}
  rows={ validPreparations }
  rowClickFunc={goToValidPreparation}
  rowDeleteFunc={deleteValidPreparation} />