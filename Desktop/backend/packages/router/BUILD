load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

# Dependencies
load("@io_bazel_rules_docker//container:container.bzl", "container_image", "container_push")

package(default_visibility = ["//visibility:public"])

# Create Docker Container
container_image(
    name = "router_lambda_image",
    base = "@lambda_go//image",
    cmd = ["/lambda/connect"],
    directory = "/lambda",
    files = [
        "//packages/router/cmd/broadcast",
        "//packages/router/cmd/connect",
        "//packages/router/cmd/disconnect",
        "//packages/router/cmd/receive",
        "//packages/router/cmd/service_forward",
        "//packages/router/cmd/service_receive",
        "//packages/router/cmd/task_is_external_user_online",
        "//packages/router/cmd/unconfirmed_connect",
        "//packages/router/cmd/unconfirmed_disconnect",
    ],
)

# Push to ECR
container_push(
    name = "router_lambda_image_push",
    format = "Docker",
    image = ":router_lambda_image",
    registry = "{BUILD_AWS_ACCOUNT_ID}.dkr.ecr.us-east-2.amazonaws.com",
    repository = "circulo/sonar-backend-router-lambda",
    tag = "{BUILD_GIT_COMMIT_SHA}",
)

go_library(
    name = "router_lib",
    srcs = ["main.go"],
    importpath = "github.com/circulohealth/sonar-backend/packages/router",
)

go_binary(
    name = "router",
    embed = [":router_lib"],
)
