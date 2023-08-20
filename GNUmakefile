TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
WEBSITE_REPO=github.com/hashicorp/terraform-website
PKG_NAME=sumologic
PLUGIN_DIR=~/.terraform.d/plugins
UNAME=$(shell uname -m)
DEBUG_FLAGS=-gcflags="all=-N -l"

default: build

tools:
	GO111MODULE=off go get -u github.com/client9/misspell/cmd/misspell
	GO111MODULE=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

build: fmtcheck
	go install

install: fmtcheck
	mkdir -vp $(PLUGIN_DIR)
	go build -o $(PLUGIN_DIR)/terraform-provider-sumologic

# separate command for installing in the right directory for testing development changes
install-dev: fmtcheck
	mkdir -vp $(PLUGIN_DIR)
ifeq ($(UNAME), arm64)
	go build ${DEBUG_FLAGS} -o $(PLUGIN_DIR)/sumologic.com/dev/sumologic/1.0.0/darwin_arm64/terraform-provider-sumologic 
else
	go build ${DEBUG_FLAGS} -o $(PLUGIN_DIR)/sumologic.com/dev/sumologic/1.0.0/darwin_amd64/terraform-provider-sumologic
endif

uninstall:
	@rm -vf $(PLUGIN_DIR)/terraform-provider-sumologic
ifeq ($(UNAME), arm64)
	@rm -vf $(PLUGIN_DIR)/sumologic.com/dev/sumologic/1.0.0/darwin_arm64/terraform-provider-sumologic
else
	@rm -vf $(PLUGIN_DIR)/sumologic.com/dev/sumologic/1.0.0/darwin_amd64/terraform-provider-sumologic
endif

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

lint:
	@echo "==> Checking source code against linters..."
	golangci-lint run ./...
	
test: fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

website:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

website-lint:
	@echo "==> Checking website against linters..."
	@misspell -error -source=text website/

website-test:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider-test PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

.PHONY: build install uninstall test testacc vet fmt fmtcheck errcheck lint tools test-compile website website-lint website-test