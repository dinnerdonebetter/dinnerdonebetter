<script>
  import { Link, navigate } from "svelte-routing";

  import Table from "../../components/Table.svelte";

  const columns = [
    {
      title: "ID",
      key: "id"
    },
    {
      title: "Source",
      key: "source"
    },
    {
      title: "Mimetype",
      key: "mimetype"
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
  let iterationMedias = [];

  function deleteIterationMedia(row) {
    if (confirm("are you sure you want to delete this iteration media?")) {
      fetch(`/api/v1/iteration_medias/${row.id}`, {
        method: "DELETE"
      }).then(response => {
        if (response.status != 204) {
          console.error("something has gone awry");
        }
        iterationMedias = iterationMedias.filter(iterationMedia => {
          return iterationMedia.id != row.id;
        });
      });
    }
  }

  function goToIterationMedia(row) {
    navigate(`/iteration_medias/${row.id}`, { replace: true });
  }

  fetch("/api/v1/iteration_medias/")
    .then(response => response.json())
    .then(data => {
      iterationMedias = data["iteration_medias"];
    });
</script>

<!-- IterationMedias.svelte -->

<Table
  {columns}
  tableStyle={'margin: 0px auto;'}
  rows={ iterationMedias }
  rowClickFunc={goToIterationMedia}
  rowDeleteFunc={deleteIterationMedia} />