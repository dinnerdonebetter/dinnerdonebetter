resource "aws_ecr_repository" "writer" {
  name = "writer"
  # do not set image_tag_mutability to "IMMUTABLE", or else we cannot use :latest tags.

  image_scanning_configuration {
    scan_on_push = true
  }
}

resource "aws_ecr_repository" "updater" {
  name = "updater"
  # do not set image_tag_mutability to "IMMUTABLE", or else we cannot use :latest tags.

  image_scanning_configuration {
    scan_on_push = true
  }
}

resource "aws_ecr_repository" "archiver" {
  name = "archiver"
  # do not set image_tag_mutability to "IMMUTABLE", or else we cannot use :latest tags.

  image_scanning_configuration {
    scan_on_push = true
  }
}

resource "aws_ecr_repository" "data_change_observer" {
  name = "data_change_observer"
  # do not set image_tag_mutability to "IMMUTABLE", or else we cannot use :latest tags.

  image_scanning_configuration {
    scan_on_push = true
  }
}

resource "aws_ecr_repository" "chore_worker" {
  name = "chore_worker"
  # do not set image_tag_mutability to "IMMUTABLE", or else we cannot use :latest tags.

  image_scanning_configuration {
    scan_on_push = true
  }
}
