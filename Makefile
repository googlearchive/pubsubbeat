BEAT_NAME=pubsubbeat
BUILD_DIR=bin
BEAT_PATH=github.com/GoogleCloudPlatform/pubsubbeat

NOW             ?= $(shell date --iso-8601=seconds)
COMMIT_ID       ?= $(shell git rev-parse HEAD)
VERSION_LDFLAGS := \
  -X github.com/elastic/beats/libbeat/version.buildTime=$(NOW) \
  -X github.com/elastic/beats/libbeat/version.commit=$(COMMIT_ID)
GOBUILD_FLAGS=-ldflags "$(VERSION_LDFLAGS)"

.PHONY: build
build:
	go build -o ${BEAT_NAME} ${GOBUILD_FLAGS}

.PHONY: pre-commit
pre-commit: clean fmt vet update test

.PHONY: clean
clean:
	rm -rf ${BUILD_DIR}
	rm -f ${BEAT_NAME}

.PHONY: test
test:
	go test ./...

.PHONY: fmt
fmt:
	gofmt -s -w $(shell find . -type f -name '*.go' | grep -v \./vendor/)

.PHONY: vet
vet:
	go vet ./...

.PHONY: update
update:
	@echo "TODO: update deps"

.PHONY: release
release: dashboards
	goreleaser --snapshot --skip-publish --rm-dist

dashboards:
	cp -r _meta/kibana dashboards
