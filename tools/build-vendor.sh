#!/usr/bin/env bash

# This is a stop-gap solution for generating the vendor directory based on bazel output. Unfortunately this isn't
# available today. Bazel itself can compile everything, but IDEs cannot index generated files/external dependencies
# properly. To handle that, we have this stop-gap script. This is a very fragile solution however, so we should migrate
# as fast as possible to a better solution.
#
# References:
# https://stackoverflow.com/questions/44375665/bazel-build-protobuf-and-code-completion
# https://github.com/bazelbuild/rules_go/issues/512

PROJECT_ROOT=$(git rev-parse --show-toplevel)

cd "${PROJECT_ROOT}"

################################################################################
# Install dependencies managed by dep (Bazel uses a different vendor directory).
dep ensure -update

###############################################################################
# Find all of the mocks, and copy them under their correct paths in vendor too.
mocks_library_id="mocks"

build_files=$(find pkg -name BUILD.bazel)
build_files=$(grep -l "name = \"${mocks_library_id}\"," ${build_files})

bazel_pkgs=""
for build_file in ${build_files}; do
    bazel_pkgs="${bazel_pkgs} $(echo "//${build_file}:${mocks_library_id}" | sed -e 's/\/BUILD.bazel//')"
done

# If there are no packages with mocks, we are done.
test -z "${bazel_pkgs}" && exit 0

# Build all of the generated files.
bazel build ${bazel_pkgs}

# Now find all of the mocks inside of the bazel build directory, and copy them to vendor directory under the usual
# expected path.
mocks_gofiles=$(find bazel-bin -follow -name '*_mocks.go')

for mocks_gofile in ${mocks_gofiles}; do
    target=$(echo "${mocks_gofile}" | sed -e 's/^bazel-bin\///')

    target_filename=$(basename "${target}")
    target_dir="vendor/nebula-go/mocks/$(dirname "${target}")mocks"

    mkdir -p "${target_dir}"
    cp "${mocks_gofile}" "${target_dir}/${target_filename}"
done