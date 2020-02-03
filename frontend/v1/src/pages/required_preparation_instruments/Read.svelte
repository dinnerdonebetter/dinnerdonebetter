<script>
  import { Link } from "svelte-routing";

  let requiredPreparationInstrument = {};

  const requiredPreparationInstrumentID = window.location.pathname.replace("/required_preparation_instruments/", "");

  fetch(`/api/v1/required_preparation_instruments/${ requiredPreparationInstrumentID }`)
    .then(response => response.json())
    .then(data => {
      requiredPreparationInstrument = data;
    });
</script>

<!-- RequiredPreparationInstruments.svelte -->

<div style="text-align: center; margin: 0px auto;">
  <p>
    name:
    {@html required_preparation_instrument.name}
  </p>
  <p>
    details:
    {@html required_preparation_instrument.details}
  </p>

  <Link to="/required_preparation_instruments">see all</Link>

</div>