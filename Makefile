DOCKER_REGISTRY   ?=
IMAGE_PREFIX      ?= ciscosso
SHORT_NAME        ?= kdk
TARGETS           ?= darwin/amd64 linux/amd64 linux/386 linux/arm linux/arm64 linux/ppc64le linux/s390x windows/amd64
DIST_DIRS         = find * -type d -exec

# go option
GO        ?= go
PKG       := $(shell dep ensure)
TAGS      :=
TESTS     := .
TESTFLAGS :=
LDFLAGS   := -w -s
GOFLAGS   :=
BINDIR    := $(CURDIR)/bin

# Required for globs to work correctly
SHELL=/bin/bash

.PHONY: all
all: build

.PHONY: build
build:
	GOBIN=$(BINDIR) $(GO) install $(GOFLAGS) -tags '$(TAGS)' -ldflags '$(LDFLAGS)' ./

# usage: make clean build-cross dist VERSION=1.0.0
.PHONY: build-cross
build-cross: LDFLAGS += -extldflags "-static"
build-cross:
	CGO_ENABLED=0 gox -parallel=3 -output="_dist/{{.OS}}-{{.Arch}}/{{.Dir}}" -osarch='$(TARGETS)' $(GOFLAGS) $(if $(TAGS),-tags '$(TAGS)',) -ldflags '$(LDFLAGS)' ./

.PHONY: dist
dist:
	( \
		cd _dist && \
		$(DIST_DIRS) cp ../LICENSE {} \; && \
		$(DIST_DIRS) cp ../README.md {} \; && \
		$(DIST_DIRS) tar -zcf $(SHORT_NAME)-${VERSION}-{}.tar.gz {} \; \
	)

.PHONY: check-docker
check-docker:
	@if [ -z $$(which docker) ]; then \
	  echo "Missing \`docker\` client which is required for development"; \
	  exit 2; \
	fi

.PHONY: docker-build
docker-build: check-docker
	docker build --rm -t ${IMAGE} -t ${MUTABLE_IMAGE} files

.PHONY: gofmt
gofmt:
	gofmt -w -s $$(find ./cmd ./pkg -type f -name '*.go')

.PHONY: clean
clean:
	@rm -rf $(BINDIR) ./_dist ./bin vendor

HAS_DEP := $(shell command -v dep;)
HAS_GIT := $(shell command -v git;)

.PHONY: bootstrap
bootstrap:
ifndef HAS_DEP
	go get -u github.com/golang/dep/cmd/dep
endif
ifndef HAS_GOX
	go get -u github.com/mitchellh/gox
endif
ifndef HAS_GIT
	$(error You must install Git)
endif
	dep ensure

include versioning.mk