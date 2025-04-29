# IAM Policy for ECS Task to access SSM Parameters
resource "aws_iam_policy" "ssm_parameter_access" {
  name = "smart-retention-ssm-parameter-access"

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Action = [
          "ssm:GetParameters",
          "ssm:GetParameter"
        ],
        Resource = "arn:aws:ssm:${var.aws_region}:${data.aws_caller_identity.current.account_id}:parameter/smart-retention/*"
      }
    ]
  })
}

# Attach policy to ECS execution role
resource "aws_iam_role_policy_attachment" "attach_ssm" {
  role       = aws_iam_role.ecs_task_execution_role.name
  policy_arn = aws_iam_policy.ssm_parameter_access.arn
}