<script>
  import { Link } from "svelte-routing";

  let validPreparation = {};

  const validPreparationID = window.location.pathname.replace("/valid_preparations/", "");

  fetch(`/api/v1/valid_preparations/${ validPreparationID }`)
    .then(response => response.json())
    .then(data => {
      validPreparation = data;
    });
</script>

<!-- ValidPreparations.svelte -->

<div style="text-align: center; margin: 0px auto;">
  <p>
    name:
    {@html valid_preparation.name}
  </p>
  <p>
    details:
    {@html valid_preparation.details}
  </p>

  <Link to="/valid_preparations">see all</Link>

</div>