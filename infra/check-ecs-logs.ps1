$clusterName = "smart-retention-ecs-cluster"
$serviceName = "smart-retention-service"
$logGroupName = "/ecs/smart-retention"

# 1. Obtém a task mais recente em execução
$taskArn = aws ecs list-tasks `
  --cluster $clusterName `
  --service-name $serviceName `
  --desired-status STOPPED `
  --query "taskArns[-1]" `
  --output text

if (-not $taskArn -or $taskArn -eq "None") {
    Write-Host "❌ Nenhuma task em execução encontrada."
    exit
}

$taskId = ($taskArn -split "/")[-1]

# 2. Encontra o log stream correto
$logStreams = aws logs describe-log-streams `
  --log-group-name $logGroupName `
  --query "logStreams[?contains(logStreamName, '$taskId')].logStreamName" `
  --output text

if (-not $logStreams) {
    Write-Host "❌ Log stream não encontrado para a task ID: $taskId"
    exit
}

# 3. Mostra os logs
Write-Host "`n📄 Logs da task: $taskId`n"
aws logs get-log-events `
  --log-group-name $logGroupName `
  --log-stream-name $logStreams `
  --limit 50 `
  --query "events[*].message" `
  --output text
