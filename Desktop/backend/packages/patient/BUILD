load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")
load("@io_bazel_rules_docker//container:container.bzl", "container_image", "container_push")

package(default_visibility = ["//visibility:public"])

# Create Docker Container
container_image(
    name = "patient_lambda_image",
    base = "@lambda_go//image",
    directory = "/lambda",
    files = [
        "//packages/patient/cmd/create_agency_providers",
        "//packages/patient/cmd/create_appointments",
        "//packages/patient/cmd/create_patients",
        "//packages/patient/cmd/delete_appointments",
        "//packages/patient/cmd/edit_agency_providers",
        "//packages/patient/cmd/edit_appointments",
        "//packages/patient/cmd/list_agency_providers",
        "//packages/patient/cmd/list_appointments",
        "//packages/patient/cmd/list_patients",
        "//packages/patient/cmd/patch_patient",
    ],
)

# Push to ECR
container_push(
    name = "patient_lambda_image_push",
    format = "Docker",
    image = ":patient_lambda_image",
    registry = "{BUILD_AWS_ACCOUNT_ID}.dkr.ecr.us-east-2.amazonaws.com",
    repository = "circulo/sonar-backend-patient-lambda",
    tag = "{BUILD_GIT_COMMIT_SHA}",
)

go_library(
    name = "patient_lib",
    srcs = ["main.go"],
    importpath = "github.com/circulohealth/sonar-backend/packages/patient",
    visibility = ["//visibility:private"],
)

go_binary(
    name = "patient",
    embed = [":patient_lib"],
    visibility = ["//visibility:public"],
)
