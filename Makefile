SHELL=/bin/bash

.PHONY: clean
clean:
	@golangci-lint cache clean

.PHONY: setup
setup:
	@pre-commit install

.PHONY: pre-commit
pre-commit:
	@pre-commit run --all-files

.PHONY: lint
lint:
	@scripts/lint.sh

.PHONY: gen-proto
gen-proto:
	rm -f ./proto/helloworld.pb.go
	cp ./bazel-out/k8-fastbuild/bin/proto/proto_go_proto_/github.com/kujilabo/cocotola-1.21/proto/helloworld.pb.go ./proto/

.PHONY: work-init
work-init:
	@go work init
	@go work use -r .

.PHONY: gazelle
gazelle:
	@bazel run //:gazelle

.PHONY: gazelle-update-repos
gazelle-update-repos:
	@bazel run //:gazelle -- update-repos -from_file ./go.work

.PHONY: update-mod
update-mod:
	@pushd ./cocotola-api/ && \
		go get -u ./... && \
	popd

run-cocotola-api:
	@bazel run //cocotola-api/src

# https://github.com/bazel-contrib/rules_oci/blob/main/docs/go.md
# https://github.com/aspect-build/bazel-examples/blob/main/oci_go_image/BUILD.bazel
docker-build:
	bazel build //cocotola-api/src:tarball
	$(eval COCOTOLA_API_TARBALL := `bazel cquery --output=files //cocotola-api/src:tarball`)
	docker load --input $(COCOTOLA_API_TARBALL)

docker-run:
	docker run --rm gcr.io/cocotola/cocotola-api:latest

test:
	@bazel test //... --test_output=all --test_timeout=60

test-s:
    @go test -coverprofile="coverage.txt" -covermode=atomic ./... -count=1
