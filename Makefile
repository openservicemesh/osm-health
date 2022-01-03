#!make

TARGETS := darwin/amd64 linux/amd64 windows/amd64
BINNAME ?= osm-health
DIST_DIRS := find * -type d -exec
VERSION ?= dev
BUILD_DATE=$$(date +%F)
GIT_SHA=$$(git rev-parse HEAD)
BUILD_DATE_VAR := github.com/openservicemesh/osm-health/pkg/version.BuildDate
BUILD_VERSION_VAR := github.com/openservicemesh/osm-health/pkg/version.Version
BUILD_GITCOMMIT_VAR := github.com/openservicemesh/osm-health/pkg/version.GitCommit

GOX    = go run github.com/mitchellh/gox

LDFLAGS ?= "-X $(BUILD_DATE_VAR)=$(BUILD_DATE) -X $(BUILD_VERSION_VAR)=$(VERSION) -X $(BUILD_GITCOMMIT_VAR)=$(GIT_SHA) -s -w"

.PHONY: build-ci
build-ci: build-osm-health

.PHONY: build
build: build-osm-health

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
go-test-coverage:
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

.PHONY: kind-up
kind-up:
	./scripts/kind-with-registry.sh

# -------------------------------------------
#  release targets below
# -------------------------------------------

.PHONY: build-cross
build-cross: cmd
	GO111MODULE=on CGO_ENABLED=0 $(GOX) -ldflags $(LDFLAGS) -parallel=3 -output="_dist/{{.OS}}-{{.Arch}}/$(BINNAME)" -osarch='$(TARGETS)' ./cmd

.PHONY: dist
dist:
	( \
		cd _dist && \
		$(DIST_DIRS) cp ../LICENSE {} \; && \
		$(DIST_DIRS) cp ../README.md {} \; && \
		$(DIST_DIRS) tar -zcf osm-${VERSION}-{}.tar.gz {} \; && \
		$(DIST_DIRS) zip -r osm-${VERSION}-{}.zip {} \; && \
		sha256sum osm-* > sha256sums.txt \
	)

.PHONY: release-artifacts
release-artifacts: build-cross dist
