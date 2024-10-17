load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

# Dependencies
load("@io_bazel_rules_docker//container:container.bzl", "container_image", "container_push")

package(default_visibility = ["//visibility:public"])

# Create Docker Container
container_image(
    name = "users_lambda_image",
    base = "@lambda_go//image",
    cmd = ["/lambda/send"],
    directory = "/lambda",
    files = [
        "//packages/users/cmd/create_organizations",
        "//packages/users/cmd/define_auth_challenge",
        "//packages/users/cmd/external_refresh",
        "//packages/users/cmd/external_sign_in",
        "//packages/users/cmd/external_sign_up",
        "//packages/users/cmd/get_organizations",
        "//packages/users/cmd/internal_post_authentication",
        "//packages/users/cmd/internal_post_confirmation",
        "//packages/users/cmd/post_confirmation",
        "//packages/users/cmd/pre_sign_up",
        "//packages/users/cmd/receive",
        "//packages/users/cmd/revoke_access",
        "//packages/users/cmd/send",
        "//packages/users/cmd/sync_cognito_users",
        "//packages/users/cmd/task_get_external_user",
        "//packages/users/cmd/task_get_internal_user",
        "//packages/users/cmd/update_user",
        "//packages/users/cmd/user_list",
    ],
)

# Push to ECR
container_push(
    name = "users_lambda_image_push",
    format = "Docker",
    image = ":users_lambda_image",
    registry = "{BUILD_AWS_ACCOUNT_ID}.dkr.ecr.us-east-2.amazonaws.com",
    repository = "circulo/sonar-backend-users-lambda",
    tag = "{BUILD_GIT_COMMIT_SHA}",
)

go_library(
    name = "users_lib",
    srcs = ["main.go"],
    importpath = "github.com/circulohealth/sonar-backend/packages/users",
)

go_binary(
    name = "users",
    embed = [":users_lib"],
)
