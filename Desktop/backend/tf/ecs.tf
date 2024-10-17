resource "aws_ecs_cluster" "sonar_cluster" {
  name = "sonar"

  setting {
    name  = "containerInsights"
    value = "enabled"
  }
}
