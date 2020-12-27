all: build

gazelle:
	@bazel run //:gazelle

setup-vendor: gazelle
	@./tools/build-vendor.sh

add-dependency:
	@dep ensure -add $(PACKAGE)
	@bazel run //:gazelle -- update-repos -from_file=Gopkg.lock

update-dependencies:
	@dep ensure -update

fetch-dependencies:
	@bazel fetch //... >&2
	@./tools/patch-sdl-BUILD-bazel-file.sh >&2

build: fetch-dependencies
	@bazel build //...

test: fetch-dependencies
	@bazel test --test_output=errors //...

clean:
	@bazel clean

# How to use: $(make run-gbc-cmd) ./data/rom.gb
run-gbc-cmd: fetch-dependencies
	@echo "bazel run //bins/nebula-gbc-go:nebula-gbc-go --sandbox_debug --"
