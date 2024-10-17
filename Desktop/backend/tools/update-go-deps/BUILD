# update-deps multi-package configuration

# gazelle:prefix gitlab.com/circulohealth/sonar/backend/tools/update-deps
load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "update-go-deps_lib",
    srcs = [
        "helpers.go",
        "main.go",
    ],
    importpath = "gitlab.com/circulohealth/sonar/backend/tools/update-go-deps",
    visibility = ["//visibility:private"],
    deps = [
        "@org_golang_x_mod//modfile:go_default_library",
        "@org_golang_x_mod//semver:go_default_library",
    ],
)

go_binary(
    name = "update-go-deps",
    out = "update-go-deps",
    embed = [":update-deps_lib"],
    visibility = ["//visibility:public"],
)

go_library(
    name = "update-deps_lib",
    srcs = [
        "helpers.go",
        "main.go",
    ],
    importpath = "gitlab.com/circulohealth/sonar/backend/tools/update-deps",
    visibility = ["//visibility:private"],
    deps = [
        "@org_golang_x_mod//modfile",
        "@org_golang_x_mod//semver",
    ],
)
