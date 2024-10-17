load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")
load("@bazel_gazelle//:def.bzl", "gazelle")

# gazelle:prefix github.com/circulohealth/sonar-backend
gazelle(name = "gazelle")

go_library(
    name = "sonar-backend_lib",
    srcs = ["main.go"],
    importpath = "github.com/circulohealth/sonar-backend",
    visibility = ["//visibility:private"],
)

go_binary(
    name = "sonar-backend",
    embed = [":sonar-backend_lib"],
    visibility = ["//visibility:public"],
)
