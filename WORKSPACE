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
    commit = "f221b8435cfb71e54062f6c6e99e9ade30b124d5",
    importpath = "gopkg.in/yaml.v2",
)

go_repository(
    name = "org_uber_go_atomic",
    commit = "df976f2515e274675050de7b3f42545de80594fd",
    importpath = "go.uber.org/atomic",
)

go_repository(
    name = "org_uber_go_multierr",
    commit = "71e610a0e48dbda460d9c18adca8b5f2de04a7c1",
    importpath = "go.uber.org/multierr",
)
