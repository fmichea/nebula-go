all: build

gazelle:
	@bazel run //:gazelle

setup-vendor: gazelle
	@./tools/build-vendor.sh

add-dependency:
	@dep ensure -add $(PACKAGE)
	@bazel run //:gazelle -- update-repo -from_file=Gopkg.lock

build:
	@bazel build //...

test:
	@bazel test --nocache_test_results //...

clean:
	@bazel clean

# How to use: $(make run-gbc-cmd) ./data/rom.gb
run-gbc-cmd:
	@echo "bazel run //bins/nebula-gbc-go:nebula-gbc-go --"