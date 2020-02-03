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
  let instruments = [];

  function deleteInstrument(row) {
    if (confirm("are you sure you want to delete this instrument?")) {
      fetch(`/api/v1/instruments/${row.id}`, {
        method: "DELETE"
      }).then(response => {
        if (response.status != 204) {
          console.error("something has gone awry");
        }
        instruments = instruments.filter(instrument => {
          return instrument.id != row.id;
        });
      });
    }
  }

  function goToInstrument(row) {
    navigate(`/instruments/${row.id}`, { replace: true });
  }

  fetch("/api/v1/instruments/")
    .then(response => response.json())
    .then(data => {
      instruments = data["instruments"];
    });
</script>

<!-- Instruments.svelte -->

<Table
  {columns}
  tableStyle={'margin: 0px auto;'}
  rows={ instruments }
  rowClickFunc={goToInstrument}
  rowDeleteFunc={deleteInstrument} />