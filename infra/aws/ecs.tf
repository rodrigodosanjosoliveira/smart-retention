resource "aws_ecs_cluster" "this" {
  name = "${var.project_name}-ecs-cluster"
}

resource "aws_ecs_task_definition" "backend" {
  family                   = "${var.project_name}-task"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = "512"
  memory                   = "1024"
  execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn

  container_definitions = jsonencode([{
    name      = "backend"
    image     = aws_ecr_repository.backend.repository_url
    essential = true
    portMappings = [{
      containerPort = 8080
      hostPort      = 8080
    }],
    logConfiguration = {
      logDriver = "awslogs"
      options = {
        awslogs-group         = "/ecs/smart-retention"
        awslogs-region        = var.aws_region
        awslogs-stream-prefix = "ecs"
      }
    }

    environment = [
      {
        name      = "APP_ENV"
        valueFrom = "/smart-retention/APP_ENV"
      },
      {
        name      = "DB_HOST"
        valueFrom = "/smart-retention/DB_HOST"
      },
      {
        name      = "DB_USER"
        valueFrom = "/smart-retention/DB_USER"
      },
      {
        name      = "DB_PASSWORD"
        valueFrom = "/smart-retention/DB_PASSWORD"
      },
      {
        name      = "DB_NAME"
        valueFrom = "/smart-retention/DB_NAME"
      },
      {
        name      = "DB_PORT"
        value     = "/smart-retention/DB_PORT"
      }
    ]
  }])
}

resource "aws_ecs_service" "backend" {
  name            = "${var.project_name}-service"
  cluster         = aws_ecs_cluster.this.id
  task_definition = aws_ecs_task_definition.backend.arn
  launch_type     = "FARGATE"
  desired_count   = 1
  network_configuration {
    subnets          = [aws_subnet.public_a.id]
    assign_public_ip = true
    security_groups  = [aws_security_group.backend.id]
  }
  load_balancer {
    target_group_arn = aws_lb_target_group.backend.arn
    container_name   = "backend"
    container_port   = 8080
  }
}

resource "aws_cloudwatch_log_group" "ecs_logs" {
  name              = "/ecs/smart-retention"
  retention_in_days = 7
}

data "aws_caller_identity" "current" {}

resource "aws_iam_policy" "ssm_parameter_access" {
  name = "smart-retention-ssm-parameter-access"

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect   = "Allow",
        Action   = [
          "ssm:GetParameters",
          "ssm:GetParameter"
        ],
        Resource = "arn:aws:ssm:${var.aws_region}:${data.aws_caller_identity.current.account_id}:parameter/smart-retention/*"
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "ecs_task_execution_ssm" {
  role       = aws_iam_role.ecs_task_execution_role.name
  policy_arn = aws_iam_policy.ssm_parameter_access.arn
}