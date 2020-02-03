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
      title: "PreparationID",
      key: "preparation_id"
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
  let requiredPreparationInstruments = [];

  function deleteRequiredPreparationInstrument(row) {
    if (confirm("are you sure you want to delete this required preparation instrument?")) {
      fetch(`/api/v1/required_preparation_instruments/${row.id}`, {
        method: "DELETE"
      }).then(response => {
        if (response.status != 204) {
          console.error("something has gone awry");
        }
        requiredPreparationInstruments = requiredPreparationInstruments.filter(requiredPreparationInstrument => {
          return requiredPreparationInstrument.id != row.id;
        });
      });
    }
  }

  function goToRequiredPreparationInstrument(row) {
    navigate(`/required_preparation_instruments/${row.id}`, { replace: true });
  }

  fetch("/api/v1/required_preparation_instruments/")
    .then(response => response.json())
    .then(data => {
      requiredPreparationInstruments = data["required_preparation_instruments"];
    });
</script>

<!-- RequiredPreparationInstruments.svelte -->

<Table
  {columns}
  tableStyle={'margin: 0px auto;'}
  rows={ requiredPreparationInstruments }
  rowClickFunc={goToRequiredPreparationInstrument}
  rowDeleteFunc={deleteRequiredPreparationInstrument} />