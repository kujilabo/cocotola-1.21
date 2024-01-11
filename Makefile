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

.PHONY: gen-code
gen-code:
	@pushd ./cocotola-auth/ && \
		mockery	&& \
	popd

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
	GOPROXy=direct bazelisk run //:gazelle -- update-repos github.com/kujilabo/redstart

.PHONY: go-mod-tidy
go-mod-tidy:
	@pushd ./cocotola-app/ && \
		go mod tidy && \
	popd
	@pushd ./cocotola-core/ && \
		go mod tidy && \
	popd
	@pushd ./cocotola-auth/ && \
		go mod tidy && \
	popd
	@pushd ./lib/ && \
		go mod tidy && \
	popd

.PHONY: update-mod
update-mod:
	@pushd ./cocotola-app/ && \
		GOPROXY=direct go get -u github.com/kujilabo/redstart && \
		go get -u ./... && \
	popd
	@pushd ./cocotola-core/ && \
		GOPROXY=direct go get -u github.com/kujilabo/redstart && \
		go get -u ./... && \
	popd
	@pushd ./cocotola-auth/ && \
		GOPROXY=direct go get -u github.com/kujilabo/redstart && \
		go get -u ./... && \
	popd
	@pushd ./lib/ && \
		go get -u ./... && \
	popd

bazel-run-app:
	@bazelisk run //cocotola-app/src

bazel-run-core:
	@bazelisk run //cocotola-core/src

bazel-run-auth:
	@bazelisk run //cocotola-auth/src

# https://github.com/bazel-contrib/rules_oci/blob/main/docs/go.md
# https://github.com/aspect-build/bazel-examples/blob/main/oci_go_image/BUILD.bazel
bazel-docker-load-app:
	$(eval COCOTOLA_APP_TARBALL := `bazel cquery --output=files //cocotola-app/src:tarball`)
	docker load --input $(COCOTOLA_APP_TARBALL)

bazel-docker-load-core:
	$(eval COCOTOLA_CORE_TARBALL := `bazel cquery --output=files //cocotola-core/src:tarball`)
	docker load --input $(COCOTOLA_CORE_TARBALL)

bazel-docker-load-auth:
	$(eval COCOTOLA_AUTH_TARBALL := `bazel cquery --output=files //cocotola-auth/src:tarball`)
	docker load --input $(COCOTOLA_AUTH_TARBALL)

bazel-build-app:
	bazelisk build //cocotola-app/src:tarball

bazel-build-core:
	bazelisk build //cocotola-core/src:tarball

bazel-build-auth:
	bazelisk build //cocotola-auth/src:tarball

bazel-docker-push-app:
	bazelisk run //cocotola-app/src:push -- --tag $(REMOTE_TAG)

bazel-docker-push-core:
	bazelisk run //cocotola-core/src:push -- --tag $(REMOTE_TAG)

bazel-docker-push-auth:
	bazelisk run //cocotola-auth/src:push -- --tag $(REMOTE_TAG)

docker-run-app:
	docker run --rm -p 8080:8080 asia.gcr.io/cocotola-001/cocotola-app:latest

docker-run-core:
	docker run --rm asia.gcr.io/cocotola-001/cocotola-core:latest

docker-run-auth:
	docker run --rm asia.gcr.io/cocotola-001/cocotola-auth:latest

bazel-docker-run-app: bazel-build-app bazel-docker-load-app docker-run-app

bazel-docker-run-core: bazel-build-core bazel-docker-load-core docker-run-core

bazel-docker-run-auth: bazel-build-auth bazel-docker-load-auth docker-run-auth

# all
bazel-build: bazel-build-app bazel-build-core bazel-build-auth

build-web:
	mkdir -p ./cocotola-app/src/web_dist
	rm -rf ./cocotola-app/src/web_dist/*
	@pushd ./cocotola-web/ && \
		npm install && \
		npm run build && \
		cp -rf ./dist/* ../cocotola-app/src/web_dist/ && \
	popd

test:
	rm -f ./coverage.lcov
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
	@pushd ./cocotola-core/ && \
		go test -coverprofile="coverage.txt" -covermode=atomic ./... -count=1 -race -tags=small && \
	popd

.PHONY: dev-docker-up
dev-docker-up:
	@docker compose -f docker/development/docker-compose.yml up -d

.PHONY: dev-docker-down
dev-docker-down:
	@docker compose -f docker/development/docker-compose.yml down

.PHONY: test-docker-up
test-docker-up:
	@docker compose -f docker-compose-test.yml up -d

.PHONY: test-docker-down
test-docker-down:
	@docker compose -f docker-compose-test.yml down
