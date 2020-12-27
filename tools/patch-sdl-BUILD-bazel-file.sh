#! /usr/bin/env bash

set -e on

fail() {
    echo "error: $@"
    exit 1
}

BAZEL_TEMPDIR=$(bazel info output_base)
if [ ! -d "${BAZEL_TEMPDIR}" ]; then
    fail "bazel temporary dir does not exist: ${BAZEL_TEMPDIR}"
fi

SDL_BUILD_DIR="${BAZEL_TEMPDIR}/external/com_github_veandco_go_sdl2/sdl"
echo "SDL_BUILD_DIR: ${SDL_BUILD_DIR}"

if [ -d "${SDL_BUILD_DIR}" ]; then
  cat <<EOF | patch -N -i - -d "${SDL_BUILD_DIR}" >&2 || rm -f "${SDL_BUILD_DIR}/BUILD.bazel.rej"
--- BUILD.bazel	2020-01-26 19:04:41.000000000 +0100
+++ BUILD.bazel	2020-01-26 19:05:23.000000000 +0100
@@ -48,9 +48,14 @@
         "video.go",
     ],
     cgo = True,
+    copts = select({
+        "@io_bazel_rules_go//go/platform:darwin": [
+            "-I/usr/local/include/SDL2",
+        ],
+    }),
     clinkopts = select({
-        "@io_bazel_rules_go//go/platform:windows": [
-            "-lSDL2",
+        "@io_bazel_rules_go//go/platform:darwin": [
+            "-L/usr/local/lib -lSDL2",
         ],
         "//conditions:default": [],
     }),
EOF
else
  fail "SDL directory not found"
fi