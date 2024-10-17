load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

# Dependencies
load("@io_bazel_rules_docker//container:container.bzl", "container_image", "container_push")

package(default_visibility = ["//visibility:public"])

# Create Docker Container
container_image(
    name = "forms_lambda_image",
    base = "@lambda_go//image",
    cmd = ["/lambda/list"],
    directory = "/lambda",
    files = [
        "//packages/forms/cmd/close",
        "//packages/forms/cmd/count",
        "//packages/forms/cmd/create",
        "//packages/forms/cmd/delete",
        "//packages/forms/cmd/edit",
        "//packages/forms/cmd/get",
        "//packages/forms/cmd/list",
        "//packages/forms/cmd/receive",
        "//packages/forms/cmd/response",
        "//packages/forms/cmd/send",
    ],
)

# Push to ECR
container_push(
    name = "forms_lambda_image_push",
    format = "Docker",
    image = ":forms_lambda_image",
    registry = "{BUILD_AWS_ACCOUNT_ID}.dkr.ecr.us-east-2.amazonaws.com",
    repository = "circulo/sonar-backend-forms-lambda",
    tag = "{BUILD_GIT_COMMIT_SHA}",
)

go_library(
    name = "forms_lib",
    srcs = ["main.go"],
    importpath = "github.com/circulohealth/sonar-backend/packages/forms",
)

go_binary(
    name = "forms",
    embed = [":forms_lib"],
)
