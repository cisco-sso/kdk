SHORT_NAME        ?= kdk
TARGETS           ?= darwin/amd64 linux/amd64 linux/386 linux/arm linux/arm64 linux/ppc64le linux/s390x windows/amd64
DIST_DIRS         = find * -type d -exec
VERSION           ?= $(shell cat cmd/root.go| grep -e 'kdk.Version' | grep [0-9] | sed 's|.*"\(.*\)"|\1|g')

# go option
GO        ?= go
PKG       := $(shell dep ensure)
TAGS      :=
TESTS     := .
TESTFLAGS :=
LDFLAGS   := -w -s -X main.Version=$(shell git describe --tags --long --dirty | sed 's/-/+/2')
GOFLAGS   :=
BINDIR    := $(CURDIR)/bin

# Required for globs to work correctly
SHELL=/bin/bash

.PHONY: all
all: build

.PHONY: build
build:
	GOBIN=$(BINDIR) $(GO) install $(GOFLAGS) -tags '$(TAGS)' -ldflags '$(LDFLAGS)' ./

# usage: make clean build-cross dist VERSION=v2.0.0-alpha.3
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
		$(DIST_DIRS) tar -zcf $(SHORT_NAME)-${VERSION}-{}.tar.gz {} \; && \
		$(DIST_DIRS) zip -r $(SHORT_NAME)-${VERSION}-{}.zip {} \; \
	)

.PHONY: gofmt
gofmt:
	gofmt -w -s $$(find ./cmd ./internal -type f -name '*.go')

.PHONY: clean
clean:
	@rm -rf $(BINDIR) ./_dist ./bin

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
