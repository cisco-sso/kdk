.DEFAULT_GOAL=help

DOCKER_REGISTRY   ?=
IMAGE_PREFIX      ?= ciscosso
SHORT_NAME        ?= kdk
TARGETS           ?= darwin/amd64 linux/amd64 linux/386 linux/arm linux/arm64 linux/ppc64le linux/s390x windows/amd64
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

PUBLISH := $(shell ./scripts/cicd.sh publish?)

#####################################################################

.PHONY: checks check-go check-docker deps gofmt \
	ci build build-cross \
	docker-build docker-push \
	bin-build bin-push \
	clean help

checks: check-go check-docker  ## Check the entire system before building

check-go:  ## Check the system for go builds
	./scripts/cicd.sh check-go

check-docker:  ## Check the system for docker builds
	./scripts/cicd.sh check-docker

deps:    ## Ensure dependencies are installed
	./scripts/cicd.sh deps

gofmt:   ## Format all golang code
	gofmt -w -s $$(find ./cmd ./pkg -type f -name '*.go')

ci: checks docker-build bin-build docker-push bin-push  ## Run the CICD build, and publish depending on circumstances

build: check-go deps  ## Build locally for local os/arch creating bin in ./
	GOBIN=$(BINDIR) $(GO) install $(GOFLAGS) -tags '$(TAGS)' -ldflags '$(LDFLAGS)' ./

build-cross: check-go deps  ## Build locally for all os/arch combinations in ./_dist
	@# # usage: make clean build-cross dist VERSION=1.0.0
	CGO_ENABLED=0 gox -parallel=3 \
	  -output="_dist/{{.OS}}-{{.Arch}}/{{.Dir}}" \
	  -osarch='$(TARGETS)' $(GOFLAGS) $(if $(TAGS),-tags '$(TAGS)',) \
	  -ldflags '$(LDFLAGS)' ./

docker-build: check-docker  ## Build the docker image
	docker pull $(BASE_IMAGE):latest
	docker build -t $(NEW_IMAGE_TAG) --cache-from $(BASE_IMAGE):latest -f files/Dockerfile files

docker-push: check-docker  ## Publish the docker image
ifeq ($(PUBLISH),true)
	@echo "Executing docker push for build"
	echo "$${DOCKER_PASSWORD}" | docker login -u "$${DOCKER_USERNAME}" --password-stdin
	docker push $(NEW_IMAGE_TAG)
	docker tag $(NEW_IMAGE_TAG) $(BASE_IMAGE):latest
	docker push $(BASE_IMAGE):latest
else
	@echo "Skipping docker push"
endif

bin-build: build-cross  ## Build the binary executable

bin-push: check-go deps  # Publish the binary executable
ifeq ($(PUBLISH),true)
	@echo "Executing bin push for build"
	git status
	git reset --hard HEAD
	goreleaser release --rm-dist --debug
else
	@echo "Skipping bin push"
endif

clean:  ## Clean up the build dirs
	@rm -rf $(BINDIR) ./_dist ./bin vendor .vendor-new .venv

help:  ## Print list of Makefile targets
	@# Taken from https://github.com/spf13/hugo/blob/master/Makefile
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	  cut -d ":" -f1- | \
	  awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
