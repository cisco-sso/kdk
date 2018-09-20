.DEFAULT_GOAL=help

DOCKER_REGISTRY   ?=
IMAGE_PREFIX      ?= ciscosso
SHORT_NAME        ?= kdk
TARGETS           ?= darwin/amd64 linux/amd64 windows/amd64
VERSION           := $(shell ./scripts/cicd.sh version)
BASE_IMAGE        ?= $(IMAGE_PREFIX)/$(SHORT_NAME)
NEW_IMAGE_TAG     ?= $(BASE_IMAGE):$(VERSION)

# go option
GO        ?= go
PKG       :=
TAGS      :=
TESTS     := .
TESTFLAGS :=
LDFLAGS   := -w -s
GOFLAGS   :=
BINDIR    := $(CURDIR)/bin

LDFLAGS += -X github.com/cisco-sso/kdk/pkg/kdk.Version=${VERSION}
LDFLAGS += -extldflags "-static"

# Required for globs to work correctly
SHELL=/bin/bash

#####################################################################

ifeq ($(shell ./scripts/cicd.sh needs-build? master ./files),true)
NEEDS_BUILD_DOCKER=true
else
$(info Will skip *Docker* build since no files changed)
endif

ifeq ($(shell ./scripts/cicd.sh needs-build? master $(shell find main.go ./cmd ./pkg -type f -name '*.go')),true)
NEEDS_BUILD_BIN=true
else
$(info Will skip *bin* build since no files changed)
endif

#####################################################################

.PHONY: checks check-go check-docker check-publish deps gofmt ci build \
	build-cross docker-build docker-push bin-build bin-push clean help

checks: check-go check-docker  ## Check the entire system before building

check-go:  ## Check the system for go builds
	./scripts/cicd.sh check-go

check-docker:  ## Check the system for docker builds
	./scripts/cicd.sh check-docker

check-publish:  ## Check ENV vars set for push to docker and github
	./scripts/cicd.sh check-publish

deps:    ## Ensure dependencies are installed
	./scripts/cicd.sh deps

gofmt:   ## Format all golang code
	gofmt -w -s $$(find ./cmd ./pkg -type f -name '*.go')

ci: checks bin-build docker-build docker-push bin-push  ## Run the CICD build, and publish depending on circumstances

build: check-go deps  ## Build locally for local os/arch creating bin in ./
	GOBIN=$(BINDIR) $(GO) install $(GOFLAGS) -tags '$(TAGS)' -ldflags '$(LDFLAGS)' ./

build-cross: check-go deps  ## Build locally for all os/arch combinations in ./_dist
	@# # usage: make clean build-cross dist VERSION=1.0.0
	CGO_ENABLED=0 gox -parallel=3 \
	  -output="_dist/{{.OS}}-{{.Arch}}/{{.Dir}}" \
	  -osarch='$(TARGETS)' $(GOFLAGS) $(if $(TAGS),-tags '$(TAGS)',) \
	  -ldflags '$(LDFLAGS)' ./

docker-build: check-docker  ## Build the docker image
ifdef NEEDS_BUILD_DOCKER
	@# Work around the fact that multistage builds do not implicitly cache
	@#   https://github.com/moby/moby/issues/34715
	@#   Once the above issue is resolved, then the below condenses to a single docker build command line on the Dockerfile
	@#   docker build -t $(NEW_IMAGE_TAG) --cache-from $(BASE_IMAGE):latest -f files/Dockerfile files

	@# Populate the build cache
	docker pull $(BASE_IMAGE):build-cache-base || true
	docker pull $(BASE_IMAGE):build-cache-multistage-goinstall || true
	docker pull $(BASE_IMAGE):build-cache-multistage-compiler || true
	docker pull $(BASE_IMAGE):latest || true

	@# The option '--cache-from' order is significant
	docker build \
	  --target build-cache-base \
	  --tag $(BASE_IMAGE):build-cache-base \
	  --cache-from $(BASE_IMAGE):build-cache-base \
	  files/
	docker build \
	  --target build-cache-multistage-goinstall \
	  --tag $(BASE_IMAGE):build-cache-multistage-goinstall \
	  --cache-from $(BASE_IMAGE):build-cache-multistage-goinstall \
	  --cache-from $(BASE_IMAGE):build-cache-base \
	  files/
	docker build \
	  --target build-cache-multistage-compiler \
	  --tag $(BASE_IMAGE):build-cache-multistage-compiler \
	  --cache-from $(BASE_IMAGE):build-cache-multistage-compiler \
	  --cache-from $(BASE_IMAGE):build-cache-multistage-goinstall \
	  --cache-from $(BASE_IMAGE):build-cache-base \
	  files/
	docker build \
	  --tag $(BASE_IMAGE):latest \
	  --cache-from $(BASE_IMAGE):latest \
	  --cache-from $(BASE_IMAGE):build-cache-multistage-compiler \
	  --cache-from $(BASE_IMAGE):build-cache-multistage-goinstall \
	  --cache-from $(BASE_IMAGE):build-cache-base \
	  files/

	@# Then retag as the new version
	docker tag $(BASE_IMAGE):latest $(NEW_IMAGE_TAG)
endif

docker-push: check-publish check-docker  ## Publish the docker image
ifdef NEEDS_BUILD_DOCKER
	@echo "Executing docker push for build"
	echo "$${DOCKER_PASSWORD}" | docker login -u "$${DOCKER_USERNAME}" --password-stdin

	@# Push cached build layers first
	docker push $(BASE_IMAGE):build-cache-base
	docker push $(BASE_IMAGE):build-cache-multistage-compiler
	docker push $(BASE_IMAGE):build-cache-multistage-goinstall
	docker push $(BASE_IMAGE):latest
	docker push $(NEW_IMAGE_TAG)
endif

bin-build: check-go  ## Build the binary executable
ifdef NEEDS_BUILD_BIN
	$(MAKE) build-cross
endif

bin-push: check-publish check-go deps  # Publish the binary executable
ifdef NEEDS_BUILD_BIN
	@echo "Executing bin push for build"
	git status
	git reset --hard HEAD
	goreleaser --rm-dist --debug
endif

clean:  ## Clean up the build dirs
	@rm -rf $(BINDIR) ./_dist ./bin vendor .vendor-new .venv

help:  ## Print list of Makefile targets
	@# Taken from https://github.com/spf13/hugo/blob/master/Makefile
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	  cut -d ":" -f1- | \
	  awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
