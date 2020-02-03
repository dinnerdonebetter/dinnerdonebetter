<script>
  import { Link, navigate } from "svelte-routing";

  import Table from "../../components/Table.svelte";

  const columns = [
    {
      title: "ID",
      key: "id"
    },
    {
      title: "Code",
      key: "code"
    },
    {
      title: "Consumed",
      key: "consumed"
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
  let invitations = [];

  function deleteInvitation(row) {
    if (confirm("are you sure you want to delete this invitation?")) {
      fetch(`/api/v1/invitations/${row.id}`, {
        method: "DELETE"
      }).then(response => {
        if (response.status != 204) {
          console.error("something has gone awry");
        }
        invitations = invitations.filter(invitation => {
          return invitation.id != row.id;
        });
      });
    }
  }

  function goToInvitation(row) {
    navigate(`/invitations/${row.id}`, { replace: true });
  }

  fetch("/api/v1/invitations/")
    .then(response => response.json())
    .then(data => {
      invitations = data["invitations"];
    });
</script>

<!-- Invitations.svelte -->

<Table
  {columns}
  tableStyle={'margin: 0px auto;'}
  rows={ invitations }
  rowClickFunc={goToInvitation}
  rowDeleteFunc={deleteInvitation} />