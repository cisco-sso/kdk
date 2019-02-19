#!/usr/bin/env bash
set -euo pipefail

check-go() {
    if ! which git &>/dev/null; then
        echo "You must install Git"
        return 1
    fi
    if ! which go &>/dev/null; then
        echo "You must install Go"
        return 1
    fi
}

check-docker() {
    if ! which docker &>/dev/null; then
        echo "You must install Docker"
        return 1
    fi
    return 0
}

check-publish() {
    # Figure out whether to release the docker image and executable binary
    #   The lack of setting PUBLISH to anything means its undefined

    tag_of_current_commit="$(git describe --exact-match --tags HEAD 2>/dev/null || true)"
    latest_tag_in_repo="$(git describe --tags | cut -d '-' -f1)"
    if [ "$tag_of_current_commit" != "$latest_tag_in_repo" ]; then
        echo "Check-Publish: Not publishing because current HEAD is not equal to the latest tag" >&2
        echo false
        return 1
    fi
    if [[ -z ${DOCKER_USERNAME-} ]]; then
        echo "Check-Publish: Not publishing because DOCKER_USERNAME is unset" >&2
        echo false
        return 1
    fi
    if [[ -z ${DOCKER_PASSWORD-} ]]; then
        echo "Check-Publish: Not publishing because DOCKER_USERNAME is unset" >&2
        echo false
        return 1
    fi
    if [[ ! -z ${TRAVIS_TAG-} ]]; then
        echo "Check-Publish: Publish because we are building a Tag on TravisCI" >&2
        echo true
        return 0
    fi
    if [[ -z ${CI-} ]]; then
        echo "Check-Publish: Publish because we are building on a local machine" >&2
        echo true
        return 0
    fi

    echo "Not publishing because of unmet conditions" >&2
    echo false
    return 1
}

needs-build?() {
    # Input: branch name, along with list of files or dirs
    # On local build or tag or code differences, returns true
    # Otherwise, returns false

    # Always need build on local non-ci machine
    if [[ -z ${CI-} ]]; then
        echo "Needs-Build: Should build because we are on a local machine" >&2
        echo true
        return 0
    fi

    # Always build for Travis Tags
    if [[ ! -z ${TRAVIS_TAG-} ]]; then
        echo "Needs-Build: Should build because TRAVIS_TAG is set" >&2
        echo true
        return 0
    fi

    # Otherwise, we are on CI, and should only build if there are differences
    if [[ $(git diff "$@") !=  "" ]]; then
        echo "Needs-Build: Should build because we are on CI with code differences" >&2
	echo true
	return 0
    fi

    echo "Needs-Build: Should Not build because no build conditions were met" >&2
    echo false
    return 0
}

deps() {
    if ! which dep &>/dev/null; then
        go get -u github.com/golang/dep/cmd/dep
    fi
    if ! which gox &>/dev/null; then
        go get -u github.com/mitchellh/gox
    fi
    if ! which goreleaser &>/dev/null; then
        curl -sSLO https://github.com/goreleaser/goreleaser/releases/download/v0.93.2/goreleaser_Linux_x86_64.tar.gz
        tar -C $GOPATH/bin -xzf goreleaser*.tar.gz goreleaser && rm goreleaser*.tar.gz
    fi

    dep ensure
}

version() {
    echo $(git describe --tags --long --dirty | sed 's/-0-........$//;')
}

case "$1" in
        checks)
            check-go
            check-docker
            ;;
        check-go)
            check-go
            ;;
        check-docker)
            check-docker
            ;;
        check-publish)
            check-publish
            ;;
        deps)
            deps
            ;;
        needs-build?)
            shift
            needs-build? $@
            ;;
        version)
            version
            ;;
        *)
            echo $"Usage: $0 {checks|check-go|check-docker|check-publish|deps|needs-build?|version}"
            exit 1
esac
