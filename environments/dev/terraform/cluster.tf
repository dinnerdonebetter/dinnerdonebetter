# GKE cluster
data "google_container_engine_versions" "gke_version" {
  location       = local.gcp_region
  version_prefix = "1.27."
}

resource "google_container_cluster" "primary" {
  name     = local.environment
  location = local.gcp_region

  # We can't create a cluster with no node pool defined, but we want to only use
  # separately managed node pools. So we create the smallest possible default
  # node pool and immediately delete it.
  remove_default_node_pool = true
  initial_node_count       = 1
}

# Separately Managed Node Pool
resource "google_container_node_pool" "primary_nodes" {
  name     = google_container_cluster.primary.name
  location = local.gcp_region
  cluster  = google_container_cluster.primary.name

  version    = data.google_container_engine_versions.gke_version.release_channel_latest_version["STABLE"]
  node_count = 1

  node_locations = [
    "us-central1-a",
  ]

  node_config {
    oauth_scopes = [
      "googleapis.com/auth/cloud-platform",
    ]

    labels = {
      env = local.project_id
    }

    preemptible  = true
    machine_type = "e2-small"
    tags         = ["gke-node", local.environment]
    metadata = {
      disable-legacy-endpoints = "true"
    }
  }
}
