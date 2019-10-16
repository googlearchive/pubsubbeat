BEAT_NAME=pubsubbeat
BUILD_DIR=bin
BEAT_PATH=github.com/GoogleCloudPlatform/pubsubbeat
GOBUILD_FLAGS=-ldflags "-X $(BEAT_PATH)/vendor/github.com/elastic/beats/libbeat/version.buildTime=$(NOW) -X $(BEAT_PATH)/vendor/github.com/elastic/beats/libbeat/version.commit=$(COMMIT_ID)"
EXES=linux windows darwin
RELEASE_TEMPLATE_DIR=${BUILD_DIR}/releases/template

.PHONY: pre-commit
pre-commit: clean fmt vet update test

.PHONY: clean
clean:
	rm -rf ${BUILD_DIR}

.PHONY: test
test:
	go test ./...

.PHONY: fmt
fmt:
	gofmt -s -w $(shell find . -type f -name '*.go' | grep -v \./vendor/)

.PHONY: vet
vet:
	go vet ./...

.PHONY: build
build:
	go build -o ${BUILD_DIR}/${BEAT_NAME} ${GOBUILD_FLAGS}

.PHONY: update
update:
	@echo "TODO: update deps"

.PHONY: release
release: $(EXES)

$(EXES): release-template
	@echo "Generating release: " $@

	mkdir -p ${BUILD_DIR}/releases/$@
	cp -r ${RELEASE_TEMPLATE_DIR}/. ${BUILD_DIR}/releases/$@

	GOOS=$@ GOARCH=amd64 go build -o ${BUILD_DIR}/releases/$@/${BEAT_NAME} ${GOBUILD_FLAGS}

	tar -zcvf ${BUILD_DIR}/releases/$@.tar.gz -C ${BUILD_DIR}/releases $@

.PHONY: release-template
release-template:
	mkdir -p ${RELEASE_TEMPLATE_DIR}

	cp {${BEAT_NAME}.yml,${BEAT_NAME}.reference.yml} ${RELEASE_TEMPLATE_DIR}
	cp {README.md,NOTICE,LICENSE,fields.yml} ${RELEASE_TEMPLATE_DIR}

	cp -r _meta/kibana ${RELEASE_TEMPLATE_DIR}/dashboards
