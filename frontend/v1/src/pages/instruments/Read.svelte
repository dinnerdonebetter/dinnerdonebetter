<script>
  import { Link } from "svelte-routing";

  let instrument = {};

  const instrumentID = window.location.pathname.replace("/instruments/", "");

  fetch(`/api/v1/instruments/${ instrumentID }`)
    .then(response => response.json())
    .then(data => {
      instrument = data;
    });
</script>

<!-- Instruments.svelte -->

<div style="text-align: center; margin: 0px auto;">
  <p>
    name:
    {@html instrument.name}
  </p>
  <p>
    details:
    {@html instrument.details}
  </p>

  <Link to="/instruments">see all</Link>

</div>