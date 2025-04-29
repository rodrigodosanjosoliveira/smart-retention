resource "aws_security_group" "backend" {
  name        = "${var.project_name}-backend-sg"
  description = "Allow ALB to access Backend service"
  vpc_id      = aws_vpc.this.id

  ingress {
    description = "Allow ALB"
    from_port       = 8080
    to_port         = 8080
    protocol        = "tcp"
    security_groups = [aws_security_group.alb.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
