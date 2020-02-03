<script>
  import { Link, navigate } from "svelte-routing";

  import Table from "../../components/Table.svelte";

  const columns = [
    {
      title: "ID",
      key: "id"
    },
    {
      title: "ReportType",
      key: "report_type"
    },
    {
      title: "Concern",
      key: "concern"
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
  let reports = [];

  function deleteReport(row) {
    if (confirm("are you sure you want to delete this report?")) {
      fetch(`/api/v1/reports/${row.id}`, {
        method: "DELETE"
      }).then(response => {
        if (response.status != 204) {
          console.error("something has gone awry");
        }
        reports = reports.filter(report => {
          return report.id != row.id;
        });
      });
    }
  }

  function goToReport(row) {
    navigate(`/reports/${row.id}`, { replace: true });
  }

  fetch("/api/v1/reports/")
    .then(response => response.json())
    .then(data => {
      reports = data["reports"];
    });
</script>

<!-- Reports.svelte -->

<Table
  {columns}
  tableStyle={'margin: 0px auto;'}
  rows={ reports }
  rowClickFunc={goToReport}
  rowDeleteFunc={deleteReport} />