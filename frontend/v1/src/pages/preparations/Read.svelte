<script>
  import { Link } from "svelte-routing";

  let preparation = {};

  const preparationID = window.location.pathname.replace("/preparations/", "");

  fetch(`/api/v1/preparations/${ preparationID }`)
    .then(response => response.json())
    .then(data => {
      preparation = data;
    });
</script>

<!-- Preparations.svelte -->

<div style="text-align: center; margin: 0px auto;">
  <p>
    name:
    {@html preparation.name}
  </p>
  <p>
    details:
    {@html preparation.details}
  </p>

  <Link to="/preparations">see all</Link>

</div>