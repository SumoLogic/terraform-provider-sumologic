OS            := $(shell go env GOOS)
ARCH          := $(shell go env GOARCH)
PLUGIN_PATH   := ${HOME}/.terraform.d/plugins/${OS}_${ARCH}
PLUGIN_NAME   := terraform-provider-sumologic
DIST_PATH     := dist/${OS}_${ARCH}
GO_PACKAGES   := $(shell go list ./... | grep -v /vendor/)
GO_FILES      := $(shell find . -type f -name '*.go')
TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
PKG_NAME=sumologic
DIR=~/.terraform.d/plugins

default: build

all: test build

test-all:
	@TF_ACC=1 go test -v -race $(GO_PACKAGES)

${DIST_PATH}/${PLUGIN_NAME}: ${GO_FILES}
	mkdir -p $(DIST_PATH); \
	go build -o $(DIST_PATH)/${PLUGIN_NAME}

build: ${DIST_PATH}/${PLUGIN_NAME}

install: build
	mkdir -p $(PLUGIN_PATH); \
	rm -rf $(PLUGIN_PATH)/${PLUGIN_NAME}; \
	install -m 0755 $(DIST_PATH)/${PLUGIN_NAME} $(PLUGIN_PATH)/${PLUGIN_NAME}

clean:
	rm -rf ${DIST_PATH}/*

update:
	dep ensure -update -v

test: fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

.PHONY: all test-all build install clean update fmtcheck fmt testacc