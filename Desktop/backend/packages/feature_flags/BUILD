load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

# Dependencies
load("@io_bazel_rules_docker//container:container.bzl", "container_image", "container_push")

package(default_visibility = ["//visibility:public"])

# Create Docker Container
container_image(
    name = "feature_flags_lambda_image",
    base = "@lambda_go//image",
    cmd = ["/lambda/connect"],
    directory = "/lambda",
    files = [
        "//packages/feature_flags/cmd/create_flag",
        "//packages/feature_flags/cmd/delete_flag",
        "//packages/feature_flags/cmd/evaluate",
        "//packages/feature_flags/cmd/list_flags",
        "//packages/feature_flags/cmd/loop_evaluate",
        "//packages/feature_flags/cmd/patch_flag",
    ],
)

# Push to ECR
container_push(
    name = "feature_flags_lambda_image_push",
    format = "Docker",
    image = ":feature_flags_lambda_image",
    registry = "{BUILD_AWS_ACCOUNT_ID}.dkr.ecr.us-east-2.amazonaws.com",
    repository = "circulo/sonar-backend-feature_flags-lambda",
    tag = "{BUILD_GIT_COMMIT_SHA}",
)

go_library(
    name = "feature_flags_lib",
    srcs = ["main.go"],
    importpath = "github.com/circulohealth/sonar-backend/packages/feature_flags",
)

go_binary(
    name = "feature_flags",
    embed = [":feature_flags_lib"],
)
