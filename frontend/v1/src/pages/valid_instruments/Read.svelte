<script>
  import { Link } from "svelte-routing";

  let validInstrument = {};

  const validInstrumentID = window.location.pathname.replace("/valid_instruments/", "");

  fetch(`/api/v1/valid_instruments/${ validInstrumentID }`)
    .then(response => response.json())
    .then(data => {
      validInstrument = data;
    });
</script>

<!-- ValidInstruments.svelte -->

<div style="text-align: center; margin: 0px auto;">
  <p>
    name:
    {@html valid_instrument.name}
  </p>
  <p>
    details:
    {@html valid_instrument.details}
  </p>

  <Link to="/valid_instruments">see all</Link>

</div>