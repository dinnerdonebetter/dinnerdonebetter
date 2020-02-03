<script>
  import { Link } from "svelte-routing";

  let invitation = {};

  const invitationID = window.location.pathname.replace("/invitations/", "");

  fetch(`/api/v1/invitations/${ invitationID }`)
    .then(response => response.json())
    .then(data => {
      invitation = data;
    });
</script>

<!-- Invitations.svelte -->

<div style="text-align: center; margin: 0px auto;">
  <p>
    name:
    {@html invitation.name}
  </p>
  <p>
    details:
    {@html invitation.details}
  </p>

  <Link to="/invitations">see all</Link>

</div>