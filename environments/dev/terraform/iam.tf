resource "google_project_iam_custom_role" "api_server_role" {
  role_id     = "api_server_role"
  title       = "API Server role"
  description = "An IAM role for the API server"
  permissions = [
    "secretmanager.versions.access",
    "cloudsql.instances.connect",
    "cloudsql.instances.get",
    "pubsub.topics.list",
    "pubsub.topics.publish",
  ]
}

resource "google_project_iam_custom_role" "data_changes_worker_role" {
  role_id     = "data_changes_worker_role"
  title       = "Data changes worker role"
  description = "An IAM role for the data changes worker"
  permissions = [
    "secretmanager.versions.access",
    "pubsub.topics.list",
    "pubsub.subscriptions.consume",
    "pubsub.subscriptions.create",
    "pubsub.subscriptions.delete",
  ]
}

resource "google_project_iam_custom_role" "meal_plan_finalizer_role" {
  role_id     = "meal_plan_finalizer_role"
  title       = "Meal Plan finalizer role"
  description = "An IAM role for the Meal Plan finalizer"
  permissions = [
    "secretmanager.versions.access",
    "cloudsql.instances.connect",
    "cloudsql.instances.get",
    "pubsub.topics.list",
    "pubsub.topics.publish",
    "pubsub.subscriptions.consume",
    "pubsub.subscriptions.create",
    "pubsub.subscriptions.delete",
  ]
}