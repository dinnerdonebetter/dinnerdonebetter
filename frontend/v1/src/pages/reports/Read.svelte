<script>
  import { Link } from "svelte-routing";

  let report = {};

  const reportID = window.location.pathname.replace("/reports/", "");

  fetch(`/api/v1/reports/${ reportID }`)
    .then(response => response.json())
    .then(data => {
      report = data;
    });
</script>

<!-- Reports.svelte -->

<div style="text-align: center; margin: 0px auto;">
  <p>
    name:
    {@html report.name}
  </p>
  <p>
    details:
    {@html report.details}
  </p>

  <Link to="/reports">see all</Link>

</div>