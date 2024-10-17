# Dependencies
load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")
load("@io_bazel_rules_docker//container:container.bzl", "container_image", "container_push")

package(default_visibility = ["//visibility:public"])

# Create Docker Container
container_image(
    name = "cloud_lambda_image",
    base = "@lambda_go//image",
    cmd = ["/lambda/file_upload"],
    directory = "/lambda",
    files = [
        "//packages/cloud/cmd/associate_file",
        "//packages/cloud/cmd/delete_file",
        "//packages/cloud/cmd/file_download",
        "//packages/cloud/cmd/file_upload",
        "//packages/cloud/cmd/get_file",
        "//packages/cloud/cmd/pre_signed_upload_url",
    ],
)

# Push to ECR
container_push(
    name = "cloud_lambda_image_push",
    format = "Docker",
    image = ":cloud_lambda_image",
    registry = "{BUILD_AWS_ACCOUNT_ID}.dkr.ecr.us-east-2.amazonaws.com",
    repository = "circulo/sonar-backend-cloud-lambda",
    tag = "{BUILD_GIT_COMMIT_SHA}",
)

go_library(
    name = "cloud_lib",
    srcs = ["main.go"],
    importpath = "github.com/circulohealth/sonar-backend/packages/cloud",
)

go_binary(
    name = "cloud",
    embed = [":cloud_lib"],
)
