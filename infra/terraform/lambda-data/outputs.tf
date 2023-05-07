
output "secret_data" {
  value       = data.aws_secretsmanager_secret.this
  description = <<EOF
The secret data. It's required since it's going to be used to enable
lambda permissions to access the secret.
EOF
}


output "secret_arn" {
  value       = data.aws_secretsmanager_secret.this.arn
  description = <<EOF
The ARN of the secret. It's required since it's going to be used to enable
lambda permissions to access the secret.
EOF
}
