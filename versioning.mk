MUTABLE_VERSION := canary

GIT_COMMIT = $(shell git rev-parse HEAD)
GIT_SHA    = $(shell git rev-parse --short HEAD)
GIT_TAG    = $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
GIT_DIRTY  = $(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")

ifdef VERSION
	DOCKER_VERSION = $(VERSION)
	BINARY_VERSION = $(VERSION)
endif

DOCKER_VERSION ?= git-${GIT_SHA}
BINARY_VERSION ?= ${GIT_TAG}

# Only set Version if building a tag or VERSION is set
ifneq ($(BINARY_VERSION),)
	LDFLAGS += -X github.com/cisco-sso/kdk/pkg/kdk/version.Version=${BINARY_VERSION}
endif

# Clear the "unreleased" string in BuildMetadata
ifneq ($(GIT_TAG),)
	LDFLAGS += -X github.com/cisco-sso/kdk/pkg/kdk/version.BuildMetadata=
endif
LDFLAGS += -X github.com/cisco-sso/kdk/pkg/kdk/version.GitCommit=${GIT_COMMIT}
LDFLAGS += -X github.com/cisco-sso/kdk/pkg/kdk/version.GitTreeState=${GIT_DIRTY}

IMAGE                := ${IMAGE_PREFIX}/${SHORT_NAME}:${DOCKER_VERSION}
MUTABLE_IMAGE        := ${IMAGE_PREFIX}/${SHORT_NAME}:${MUTABLE_VERSION}

info:
	@echo "Version:           ${VERSION}"
	@echo "Git Tag:           ${GIT_TAG}"
	@echo "Git Commit:        ${GIT_COMMIT}"
	@echo "Git Tree State:    ${GIT_DIRTY}"
	@echo "Docker Version:    ${DOCKER_VERSION}"
