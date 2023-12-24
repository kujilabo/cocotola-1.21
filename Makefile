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
	@bazelisk run //:gazelle -- -build_tags=small,medium,large

.PHONY: gazelle-update-repos
gazelle-update-repos:
	@bazelisk run //:gazelle -- update-repos -from_file ./go.work

.PHONY: update-mod
update-mod:
	@pushd ./cocotola-api/ && \
		go get -u ./... && \
	popd
	@pushd ./cocotola-auth/ && \
		go get -u ./... && \
	popd
	@pushd ./lib/ && \
		go get -u ./... && \
	popd

bazel-run-api:
	@bazelisk run //cocotola-api/src

bazel-run-auth:
	@bazelisk run //cocotola-auth/src

# https://github.com/bazel-contrib/rules_oci/blob/main/docs/go.md
# https://github.com/aspect-build/bazel-examples/blob/main/oci_go_image/BUILD.bazel
bazel-docker-load-api:
	$(eval COCOTOLA_API_TARBALL := `bazel cquery --output=files //cocotola-api/src:tarball`)
	docker load --input $(COCOTOLA_API_TARBALL)

bazel-docker-load-auth:
	$(eval COCOTOLA_AUTH_TARBALL := `bazel cquery --output=files //cocotola-auth/src:tarball`)
	docker load --input $(COCOTOLA_AUTH_TARBALL)

bazel-build-api:
	bazelisk build //cocotola-api/src:tarball

bazel-build-auth:
	bazelisk build //cocotola-auth/src:tarball

bazel-docker-push-api:
	bazelisk run //cocotola-api/src:push -- --tag $(REMOTE_TAG)

bazel-docker-push-auth:
	bazelisk run //cocotola-auth/src:push -- --tag $(REMOTE_TAG)

docker-run-api:
	docker run --rm asia.gcr.io/cocotola-001/cocotola-api:latest

docker-run-auth:
	docker run --rm asia.gcr.io/cocotola-001/cocotola-auth:latest

bazel-docker-run-api: bazel-build-auth bazel-docker-load-api docker-run-api

bazel-docker-run-auth: bazel-build-auth bazel-docker-load-auth docker-run-auth

# all
bazel-build: bazel-build-api bazel-build-auth

test:
	@bazelisk test //... --test_output=errors --test_timeout=60 --test_size_filters=small
	@bazelisk coverage //... --combined_report=lcov
	$(eval OUTPUT_PATH := `bazel info output_path`)
	cp "$(OUTPUT_PATH)/_coverage/_coverage_report.dat" ./coverage.lcov

bazel-test-s:
	@bazelisk test //... --test_output=errors --test_timeout=60 --test_size_filters=small

bazel-coverage-s:
	rm -f ./coverage.lcov
	@bazelisk coverage //... --combined_report=lcov
	$(eval OUTPUT_PATH := `bazel info output_path`)
	cp "$(OUTPUT_PATH)/_coverage/_coverage_report.dat" ./coverage.lcov

test-report:
	@bazelisk test //... --test_output=errors --test_timeout=60 --test_size_filters=small --@io_bazel_rules_go//go/config:race 
	@bazelisk coverage //... --combined_report=lcov
	$(eval OUTPUT_PATH := `bazel info output_path`)
	genhtml --branch-coverage --output genhtml "$(OUTPUT_PATH)/_coverage/_coverage_report.dat"

test-s:
	@pushd ./cocotola-api/ && \
		go test -coverprofile="coverage.txt" -covermode=atomic ./... -count=1 -race -tags=small && \
	popd
