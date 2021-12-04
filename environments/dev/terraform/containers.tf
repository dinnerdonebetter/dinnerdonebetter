data "aws_ecr_repository" "writer_worker" {
  name = "writer_worker"
}

data "aws_ecr_repository" "updater_worker" {
  name = "updater_worker"
}

data "aws_ecr_repository" "archiver_worker" {
  name = "archiver_worker"
}

data "aws_ecr_repository" "data_changes_worker" {
  name = "data_changes_worker"
}

data "aws_ecr_repository" "chore_worker" {
  name = "chore_worker"
}
