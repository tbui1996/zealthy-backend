load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

# Dependencies
load("@io_bazel_rules_docker//container:container.bzl", "container_image", "container_push")

package(default_visibility = ["//visibility:public"])

# Create Docker Container
container_image(
    name = "support_lambda_image",
    base = "@lambda_go//image",
    cmd = ["/lambda/connect"],
    directory = "/lambda",
    files = [
        "//packages/support/cmd/assign_pending_chat_session",
        "//packages/support/cmd/chat_messages_get",
        "//packages/support/cmd/chat_session_create",
        "//packages/support/cmd/chat_session_star",
        "//packages/support/cmd/chat_session_update_status",
        "//packages/support/cmd/chat_sessions_get",
        "//packages/support/cmd/loop_chat_sessions_get",
        "//packages/support/cmd/loop_online_internal_users",
        "//packages/support/cmd/loop_send",
        "//packages/support/cmd/patients_get",
        "//packages/support/cmd/pending_chat_session_create",
        "//packages/support/cmd/receive",
        "//packages/support/cmd/send",
        "//packages/support/cmd/submit_feedback",
        "//packages/support/cmd/task_consume_send_message_event",
        "//packages/support/cmd/task_record_offline_message",
        "//packages/support/cmd/task_reset_offline_message",
        "//packages/support/cmd/task_trigger_offline_email",
        "//packages/support/cmd/update_chat_notes",
    ],
)

# Push to ECR
container_push(
    name = "support_lambda_image_push",
    format = "Docker",
    image = ":support_lambda_image",
    registry = "{BUILD_AWS_ACCOUNT_ID}.dkr.ecr.us-east-2.amazonaws.com",
    repository = "circulo/sonar-backend-support-lambda",
    tag = "{BUILD_GIT_COMMIT_SHA}",
)

go_library(
    name = "support_lib",
    srcs = ["main.go"],
    importpath = "github.com/circulohealth/sonar-backend/packages/support",
)

go_binary(
    name = "support",
    embed = [":support_lib"],
)
