#!make

TARGETS         := darwin/amd64 linux/amd64 windows/amd64
BINNAME         ?= osm-health
DIST_DIRS       := find * -type d -exec

GOPATH = $(shell go env GOPATH)
GOBIN  = $(GOPATH)/bin
GOX    = go run github.com/mitchellh/gox

VERSION ?= dev
BUILD_DATE ?=
GIT_SHA=$$(git rev-parse HEAD)
BUILD_DATE_VAR := github.com/openservicemesh/osm/pkg/version.BuildDate
BUILD_VERSION_VAR := github.com/openservicemesh/osm/pkg/version.Version
BUILD_GITCOMMIT_VAR := github.com/openservicemesh/osm/pkg/version.GitCommit

LDFLAGS ?= "-X $(BUILD_DATE_VAR)=$(BUILD_DATE) -X $(BUILD_VERSION_VAR)=$(VERSION) -X $(BUILD_GITCOMMIT_VAR)=$(GIT_SHA) -s -w"

# These two values are combined and passed to go test
E2E_FLAGS ?= -installType=KindCluster
E2E_FLAGS_DEFAULT := -test.v -ginkgo.v -ginkgo.progress -ctrRegistry $(CTR_REGISTRY) -osmImageTag $(CTR_TAG)

# Installed Go version
# This is the version of Go going to be used to compile this project.
# It will be compared with the minimum requirements for OSM.
GO_VERSION_MAJOR = $(shell go version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f1)
GO_VERSION_MINOR = $(shell go version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f2)
GO_VERSION_PATCH = $(shell go version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f3)
ifeq ($(GO_VERSION_PATCH),)
GO_VERSION_PATCH := 0
endif

.PHONY: build-osm-health
build-osm-health:
	CGO_ENABLED=0  go build -v -o ./bin/osm-health -ldflags ${LDFLAGS} ./cmd

.PHONY: go-checks
go-checks: go-lint go-fmt go-mod-tidy check-mocks

.PHONY: go-vet
go-vet:
	go vet ./...

.PHONY: go-lint
go-lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run --config .golangci.yml

.PHONY: go-fmt
go-fmt:
	go fmt ./...

.PHONY: go-mod-tidy
go-mod-tidy:
	./scripts/go-mod-tidy.sh

.PHONY: go-test
go-test:
	./scripts/go-test.sh

.PHONY: go-test-coverage
go-test-coverage: embed-files
	./scripts/test-w-coverage.sh

.PHONY: shellcheck
shellcheck:
	shellcheck -x $(shell find . -name '*.sh')

.PHONY: install-git-pre-push-hook
install-git-pre-push-hook:
	./scripts/install-git-pre-push-hook.sh

.PHONY: run-collection
run-collection: build-osm-health
	./bin/osm-health collect
