resource "aws_ecr_repository" "backend" {
  name = "${var.project_name}-backend"
}
