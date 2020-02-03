<script>
  import { Link } from "svelte-routing";

  let iterationMedia = {};

  const iterationMediaID = window.location.pathname.replace("/iteration_medias/", "");

  fetch(`/api/v1/iteration_medias/${ iterationMediaID }`)
    .then(response => response.json())
    .then(data => {
      iterationMedia = data;
    });
</script>

<!-- IterationMedias.svelte -->

<div style="text-align: center; margin: 0px auto;">
  <p>
    name:
    {@html iteration_media.name}
  </p>
  <p>
    details:
    {@html iteration_media.details}
  </p>

  <Link to="/iteration_medias">see all</Link>

</div>