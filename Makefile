BEAT_NAME=pubsubbeat
BEAT_PATH=github.com/GoogleCloudPlatform/pubsubbeat
BEAT_GOPATH=$(firstword $(subst :, ,${GOPATH}))
BEAT_URL=https://${BEAT_PATH}
SYSTEM_TESTS=false
TEST_ENVIRONMENT=false
ES_BEATS?=./vendor/github.com/elastic/beats
GOPACKAGES=$(shell govendor list -no-status +local)
PREFIX?=.
NOTICE_FILE=NOTICE
GOBUILD_FLAGS=-i -ldflags "-X $(BEAT_PATH)/vendor/github.com/elastic/beats/libbeat/version.buildTime=$(NOW) -X $(BEAT_PATH)/vendor/github.com/elastic/beats/libbeat/version.commit=$(COMMIT_ID)"
GOX_OS=linux darwin windows ## @Building List of all OS to be supported by "make crosscompile".
GOX_FLAGS=-arch="arm64 amd64"
EXES=pubsubbeat-darwin-amd64 pubsubbeat-linux-amd64 pubsubbeat-linux-arm64 pubsubbeat-windows-amd64.exe
RELEASE_TEMPLATE_DIR=${BUILD_DIR}/releases/template

# Path to the libbeat Makefile
-include $(ES_BEATS)/libbeat/scripts/Makefile

# Initial beat setup
.PHONY: setup
setup: copy-vendor
	$(MAKE) update

# Copy beats into vendor directory
.PHONY: copy-vendor
copy-vendor:
	mkdir -p vendor/github.com/elastic/
	cp -R ${BEAT_GOPATH}/src/github.com/elastic/beats vendor/github.com/elastic/
	rm -rf vendor/github.com/elastic/beats/.git
	rm -R vendor/github.com/elastic/beats/auditbeat
	rm -R vendor/github.com/elastic/beats/filebeat
	rm -R vendor/github.com/elastic/beats/heartbeat
	rm -R vendor/github.com/elastic/beats/metricbeat
	rm -R vendor/github.com/elastic/beats/packetbeat
	rm -R vendor/github.com/elastic/beats/winlogbeat

.PHONY: git-init
git-init:
	git init
	git add README.md CONTRIBUTING.md
	git commit -m "Initial commit"
	git add LICENSE
	git commit -m "Add the LICENSE"
	git add .gitignore
	git commit -m "Add git settings"
	git add .
	git reset -- .travis.yml
	git commit -m "Add pubsubbeat"
	git add .travis.yml
	git commit -m "Add Travis CI"

# This is called by the beats packer before building starts
.PHONY: before-build
before-build:

# Collects all dependencies and then calls update
.PHONY: collect
collect:

.PHONY: pre-commit
pre-commit: clean fmt update unit

# Generates release archives without needing Docker
.PHONY: release
release: $(EXES)

$(EXES): crosscompile release-template
	@echo "Generating release: " $@

	mkdir -p ${BUILD_DIR}/releases/$@
	cp -r ${RELEASE_TEMPLATE_DIR}/. ${BUILD_DIR}/releases/$@
	cp ${BUILD_DIR}/bin/$@ ${BUILD_DIR}/releases/$@/${BEAT_NAME}$(suffix $@)

	tar -zcvf ${BUILD_DIR}/releases/$@.tar.gz -C ${BUILD_DIR}/releases $@

.PHONY: release-template
release-template: update
	mkdir -p ${RELEASE_TEMPLATE_DIR}

	cp {${BEAT_NAME}.yml,${BEAT_NAME}.reference.yml} ${RELEASE_TEMPLATE_DIR}
	cp {README.md,NOTICE,LICENSE,fields.yml} ${RELEASE_TEMPLATE_DIR}

	cp -r _meta/kibana ${RELEASE_TEMPLATE_DIR}/dashboards
