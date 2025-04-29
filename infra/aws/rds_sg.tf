resource "aws_security_group" "rds" {
  name        = "${var.project_name}-rds-sg"
  description = "Allow Backend to access RDS"
  vpc_id      = aws_vpc.this.id

  ingress {
    description     = "Allow Backend to PostgreSQL"
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = [aws_security_group.backend.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
