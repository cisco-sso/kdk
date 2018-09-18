DOCKER_REGISTRY           ?=
IMAGE_PREFIX              ?= ciscosso
SHORT_NAME                ?= kdk
TARGETS                   ?= darwin/amd64 linux/amd64 linux/386 linux/arm linux/arm64 linux/ppc64le linux/s390x windows/amd64
DIST_DIRS                 = find * -type d -exec
BASE_IMAGE_TAG            ?= $(IMAGE_PREFIX)/$(SHORT_NAME)
VERSION                   ?= $(shell git describe --tags --long --dirty)
LATEST_RELEASE            := $(shell curl -sSL "https://api.github.com/repos/cisco-sso/kdk/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
LATEST_RELEASE_IMAGE_TAG  ?= $(BASE_IMAGE_TAG):$(LATEST_RELEASE)
NEW_IMAGE_TAG             ?= $(BASE_IMAGE_TAG):$(VERSION)

#ifdef LATEST_RELEASE
#	echo $(LATEST_RELEASE)
#else
#	$(error latest release note defined)
#endif

# go option
GO        ?= go
PKG       := $(shell dep ensure)
TAGS      :=
TESTS     := .
TESTFLAGS :=
LDFLAGS   := -w -s
GOFLAGS   :=
BINDIR    := $(CURDIR)/bin

LDFLAGS += -X github.com/cisco-sso/kdk/pkg/kdk.Version=${VERSION}

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

.PHONY: docker-pull-latest-release-image
docker-pull-latest-release-image:
	docker pull ${LATEST_RELEASE_IMAGE_TAG}

.PHONY: docker-login
docker-login: check-docker
	echo "$DOCKER_PASSWORD" | docker login -u "DOCKER_USERNAME" --password-stdin

.PHONY: docker-build
docker-build: check-docker docker-pull-latest-release-image
	docker build -t ${NEW_IMAGE_TAG} --cache-from ${LATEST_RELEASE_IMAGE_TAG} -f files/Dockerfile files

.PHONY: docker-build-clean
docker-build-clean: check-docker
	docker build --rm -t ${NEW_IMAGE_TAG} -f files/Dockerfile files

.PHONY: docker-push
docker-push: docker-login
	docker push ${NEW_IMAGE_TAG}

.PHONY: ci
ci: docker-build

.PHONY: gofmt
gofmt:
	gofmt -w -s $$(find ./cmd ./pkg -type f -name '*.go')

.PHONY: clean
clean:
	@rm -rf $(BINDIR) ./_dist ./bin vendor

.PHONY: release
release:
	goreleaser --debug

HAS_GIT := $(shell command -v git;)
HAS_DEP := $(shell command -v dep;)
HAS_GORELEASER := $(shell command -v goreleaser;)

.PHONY: bootstrap
bootstrap:
ifndef HAS_GIT
	$(error You must install Git)
endif
ifndef HAS_DEP
	go get -u github.com/golang/dep/cmd/dep
endif
ifndef HAS_GOX
	go get -u github.com/mitchellh/gox
endif
ifndef HAS_GORELEASER
	go get -u github.com/goreleaser/goreleaser
endif
	dep ensure
