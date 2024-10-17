output "arn" {
  description = "Full ARN of the repository"
  value       = aws_ecr_repository.repo.arn
}

output "registry_id" {
  description = "The registry ID where the repository was created"
  value       = aws_ecr_repository.repo.registry_id
}

output "repository_url" {
  description = "The URL of the repository"
  value       = aws_ecr_repository.repo.repository_url
}

output "tags_all" {
  description = "A map of tags assigned to the resource"
  value       = aws_ecr_repository.repo.tags_all
}
