resource "google_project_iam_custom_role" "dev_cluster_role" {
  role_id     = "dev_cluster_role"
  title       = "Dev cluster role"
  description = "An IAM role for the dev cluster"
  permissions = [
    "secretmanager.versions.access",
    "cloudsql.instances.connect",
    "cloudsql.instances.get",
    "pubsub.topics.list",
    "pubsub.topics.publish",
    "cloudtrace.traces.patch",
    "logging.buckets.create",
    "logging.buckets.write",
    "logging.buckets.list",
    "logging.buckets.get",
    "storage.objects.list",
    "storage.objects.get",
    "storage.objects.update",
    "storage.objects.create",
    "storage.objects.delete",
    "storage.objects.get",
  ]
}

resource "google_service_account" "dev_cluster_service_account" {
  account_id   = "dev-cluster"
  display_name = "dev Cluster Service Account"
}

resource "google_project_iam_member" "dev_cluster" {
  project = local.project_id
  role    = google_project_iam_custom_role.dev_cluster_role.id
  member  = format("serviceAccount:%s", google_service_account.dev_cluster_service_account.email)
}

resource "google_container_cluster" "primary" {
  name     = local.environment
  location = local.gcp_region

  # We can't create a cluster with no node pool defined, but we want to only use
  # separately managed node pools. So we create the smallest possible default
  # node pool and immediately delete it.
  remove_default_node_pool = true
  initial_node_count       = 1
  deletion_protection      = false

  network = google_compute_network.private_network.name
  # subnetwork = google_compute_subnetwork.subnet.name
}

# Separately Managed Node Pool
resource "google_container_node_pool" "primary_nodes" {
  location   = local.gcp_region
  name       = google_container_cluster.primary.name
  cluster    = google_container_cluster.primary.name
  node_count = 1

  node_locations = [
    local.gcp_main_zone,
  ]

  autoscaling {
    total_max_node_count = 1
    location_policy = "BALANCED"
  }

  management {
    auto_repair = true
    auto_upgrade = true
  }

  network_config {
    
  }

  node_config {
    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform",
    ]

    labels = {
      env = local.project_id
    }

    preemptible  = true
    machine_type = "e2-medium"
    tags         = ["gke-node", local.environment]
    metadata = {
      disable-legacy-endpoints = "true"
    }
  }
}

# resource "google_gke_backup_backup_plan" "basic" {
#   name = "basic-plan"
#   cluster = google_container_cluster.primary.id
#   location = local.gcp_region
#   backup_config {
#     include_volume_data = true
#     include_secrets = true
#     all_namespaces = true
#   }
# }
