load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

# Dependencies
load("@io_bazel_rules_docker//container:container.bzl", "container_image", "container_push")

package(default_visibility = ["//visibility:public"])

# Create Docker Container
container_image(
    name = "global_lambda_image",
    base = "@lambda_go//image",
    cmd = ["/lambda/connect"],
    directory = "/lambda",
    files = [
        "//packages/global/cmd/connect",
        "//packages/global/cmd/disconnect",
        "//packages/global/cmd/external_oh_authorizer",
        "//packages/global/cmd/sonar_authorizer",
        "//packages/global/cmd/task_is_internal_user_online",
    ],
)

# Push to ECR
container_push(
    name = "global_lambda_image_push",
    format = "Docker",
    image = ":global_lambda_image",
    registry = "{BUILD_AWS_ACCOUNT_ID}.dkr.ecr.us-east-2.amazonaws.com",
    repository = "circulo/sonar-backend-global-lambda",
    tag = "{BUILD_GIT_COMMIT_SHA}",
)

go_library(
    name = "global_lib",
    srcs = ["main.go"],
    importpath = "github.com/circulohealth/sonar-backend/packages/global",
)

go_binary(
    name = "global",
    embed = [":global_lib"],
)
