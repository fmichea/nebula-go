#!/usr/bin/env bash

PKG_NAME="$1"
shift

INTERFACES="$@"

function usage_and_quit {
    echo "usage: $0 PKG INT1 [INT2 [...]]" >&2
    exit 1
}

test -z "${PKG_NAME}" && usage_and_quit
test -z "${INTERFACES}" && usage_and_quit

package_basename=$(basename "${PKG_NAME}")

interfaces_code=""
for interface in ${INTERFACES}; do
    interfaces_code="${interfaces_code}        \"${interface}\",\n"
done

cat <<EOF
##############################################################################
# This code is generated, do not edit go_mocks_library or gomock manually, use
# ./tools/generate-bazel-mocks-code.sh to generate.
go_library(
    name = "go_mocks_library",
    srcs = ["${package_basename}_mocks.go"],
    importpath = "nebula-go/mocks/pkg/gbc/${package_basename}mocks",
    visibility = ["//visibility:public"],
    deps = [
        ":go_default_library",
        "@com_github_golang_mock//gomock:go_default_library",
    ],
)

# gazelle:exclude ${package_basename}_mocks.go
gomock(
    name = "mocks",
    out = "${package_basename}_mocks.go",
    interfaces = [
$(echo -e "${interfaces_code}")
    ],
    library = "//${PKG_NAME}:go_default_library",
    package = "${package_basename}mocks",
)
EOF