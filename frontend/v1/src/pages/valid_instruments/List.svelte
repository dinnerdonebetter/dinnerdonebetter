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
  let validInstruments = [];

  function deleteValidInstrument(row) {
    if (confirm("are you sure you want to delete this valid instrument?")) {
      fetch(`/api/v1/valid_instruments/${row.id}`, {
        method: "DELETE"
      }).then(response => {
        if (response.status != 204) {
          console.error("something has gone awry");
        }
        validInstruments = validInstruments.filter(validInstrument => {
          return validInstrument.id != row.id;
        });
      });
    }
  }

  function goToValidInstrument(row) {
    navigate(`/valid_instruments/${row.id}`, { replace: true });
  }

  fetch("/api/v1/valid_instruments/")
    .then(response => response.json())
    .then(data => {
      validInstruments = data["valid_instruments"];
    });
</script>

<!-- ValidInstruments.svelte -->

<Table
  {columns}
  tableStyle={'margin: 0px auto;'}
  rows={ validInstruments }
  rowClickFunc={goToValidInstrument}
  rowDeleteFunc={deleteValidInstrument} />