<script>
  import { Link } from "svelte-routing";

  let name = "";
  let details = "";
  let canSubmit = false;

  function evaluateSubmission() {
    canSubmit = name != "" && details != "";
  }

  function createRequiredPreparationInstrument() {
    fetch("http://localhost/api/v1/required_preparation_instruments/", {
      method: "POST",
      mode: "cors", // no-cors, cors, *same-origin
      credentials: "include",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({
        name,
        details
      })
    }).then(function(response) {
      if (response.status != 201) {
        console.error("something has gone awry");
      } else {
        name = "";
        details = "";
      }
    });
  }
</script>

<!-- RequiredPreparationInstruments.svelte -->
<form
  id="requiredPreparationInstrumentForm"
  on:submit|preventDefault={createRequiredPreparationInstrument}
  style="margin-top: 7.5%; text-align: center;">
  <p>
    name:
    <input
      bind:value={name}
      on:keyup={evaluateSubmission}
      type="text"
      name="name" />
  </p>
  <p>
    details:
    <input
      bind:value={details}
      on:keyup={evaluateSubmission}
      type="text"
      name="details" />
  </p>
  <input type="submit" value="create" disabled={!canSubmit} />
  <Link to="/required_preparation_instruments">required preparation instruments list</Link>
</form>
