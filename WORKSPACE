load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "io_bazel_rules_go",
    urls = [
        "https://storage.googleapis.com/bazel-mirror/github.com/bazelbuild/rules_go/releases/download/v0.19.5/rules_go-v0.19.5.tar.gz",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.19.5/rules_go-v0.19.5.tar.gz",
    ],
    sha256 = "513c12397db1bc9aa46dd62f02dd94b49a9b5d17444d49b5a04c5a89f3053c1c",
)

load("@io_bazel_rules_go//go:deps.bzl", "go_rules_dependencies", "go_register_toolchains")
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "io_bazel_rules_go",
    urls = [
        "https://storage.googleapis.com/bazel-mirror/github.com/bazelbuild/rules_go/releases/download/v0.19.5/rules_go-v0.19.5.tar.gz",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.19.5/rules_go-v0.19.5.tar.gz",
    ],
    sha256 = "513c12397db1bc9aa46dd62f02dd94b49a9b5d17444d49b5a04c5a89f3053c1c",
)

load("@io_bazel_rules_go//go:deps.bzl", "go_rules_dependencies", "go_register_toolchains")

go_rules_dependencies()

go_register_toolchains()

http_archive(
    name = "bazel_gazelle",
    urls = [
        "https://storage.googleapis.com/bazel-mirror/github.com/bazelbuild/bazel-gazelle/releases/download/0.18.2/bazel-gazelle-0.18.2.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/0.18.2/bazel-gazelle-0.18.2.tar.gz",
    ],
    sha256 = "7fc87f4170011201b1690326e8c16c5d802836e3a0d617d8f75c3af2b23180c4",
)

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

gazelle_dependencies()

go_rules_dependencies()

go_register_toolchains()

go_repository(
    name = "com_github_jmhodges_bazel_gomock",
    importpath = "github.com/jmhodges/bazel_gomock",
    sum = "h1:8jQePzYVW8vRoTcadW25gGzpFSe0kTHIXyF4tDa/EU8=",
    version = "v0.0.0-20191008215159-92c80571fc80",
)

go_repository(
    name = "com_github_golang_mock",
    commit = "9fa652df1129bef0e734c9cf9bf6dbae9ef3b9fa",
    importpath = "github.com/golang/mock",
)

go_repository(
    name = "com_github_davecgh_go_spew",
    commit = "8991bc29aa16c548c550c7ff78260e27b9ab7c73",
    importpath = "github.com/davecgh/go-spew",
)

go_repository(
    name = "com_github_pmezard_go_difflib",
    commit = "792786c7400a136282c1664665ae0a8db921c6c2",
    importpath = "github.com/pmezard/go-difflib",
)

go_repository(
    name = "com_github_stretchr_testify",
    commit = "221dbe5ed46703ee255b1da0dec05086f5035f62",
    importpath = "github.com/stretchr/testify",
)

go_repository(
    name = "in_gopkg_yaml_v2",
    commit = "1f64d6156d11335c3f22d9330b0ad14fc1e789ce",
    importpath = "gopkg.in/yaml.v2",
)

go_repository(
    name = "org_uber_go_atomic",
    commit = "40ae6a40a970ef4cdbffa7b24b280e316db8accc",
    importpath = "go.uber.org/atomic",
)

go_repository(
    name = "org_uber_go_multierr",
    commit = "824d08f79702fe5f54aca8400aa0d754318786e7",
    importpath = "go.uber.org/multierr",
)

go_repository(
    name = "com_github_veandco_go_sdl2",
    commit = "d59fa3b143886176f7c484340f3c1952aed89699",
    importpath = "github.com/veandco/go-sdl2",
)

go_repository(
    name = "co_honnef_go_tools",
    commit = "afd67930eec2a9ed3e9b19f684d17a062285f16a",
    importpath = "honnef.co/go/tools",
)

go_repository(
    name = "com_github_benbjohnson_clock",
    commit = "e7ca2eb904dce84a821b09d9285dc600d6174b8f",
    importpath = "github.com/benbjohnson/clock",
)

go_repository(
    name = "com_github_burntsushi_toml",
    commit = "3012a1dbe2e4bd1391d42b32f0577cb7bbc7f005",
    importpath = "github.com/BurntSushi/toml",
)

go_repository(
    name = "org_golang_x_lint",
    commit = "fdd1cda4f05fd1fd86124f0ef9ce31a0b72c8448",
    importpath = "golang.org/x/lint",
)

go_repository(
    name = "org_golang_x_tools",
    commit = "774c71fcf11405d0a5ce0aba75dc113822d62178",
    importpath = "golang.org/x/tools",
)

go_repository(
    name = "org_uber_go_tools",
    commit = "2cfd321de3ee5d5f8a5fda2521d1703478334d98",
    importpath = "go.uber.org/tools",
)

go_repository(
    name = "com_github_pkg_profile",
    commit = "acd64d450fd45fb2afa41f833f3788c8a7797219",
    importpath = "github.com/pkg/profile",
)
