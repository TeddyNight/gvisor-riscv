load("//tools:defs.bzl", "go_library")

package(
    default_applicable_licenses = ["//:license"],
    licenses = ["notice"],
)

go_library(
    name = "console",
    srcs = [
        "console.go",
        "pty_linux.go",
    ],
    visibility = [
        "//runsc:__subpackages__",
    ],
    deps = [
	# use creack pty to provide riscv support
        "@com_github_creack_pty//:go_default_library",
        "@org_golang_x_sys//unix:go_default_library",
    ],
)
