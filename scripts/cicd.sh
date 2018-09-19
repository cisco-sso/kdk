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

deps() {
    if ! which dep &>/dev/null; then
        go get -u github.com/golang/dep/cmd/dep
    fi
    if ! which gox &>/dev/null; then
        go get -u github.com/mitchellh/gox
    fi
    if ! which goreleaser &>/dev/null; then
        go get -u github.com/goreleaser/goreleaser
    fi

    dep ensure
}

publish?() {
    # Figure out whether to release the docker image and executable binary
    #   The lack of setting PUBLISH to anything means its undefined

    tag_of_current_commit="$(git describe --exact-match --tags HEAD &>/dev/null || true)"
    latest_tag_in_repo="$(git describe --tags | cut -d '-' -f1)"
    if [ "$tag_of_current_commit" != "$latest_tag_in_repo" ]; then
        echo "Not publishing because current HEAD is not equal to the latest tag" >&2
        echo false
        return 1
    fi
    if [[ -z ${DOCKER_USERNAME+x} ]]; then
        echo "Not publishing because DOCKER_USERNAME is unset" >&2
        echo false
        return 1
    fi
    if [[ -z ${DOCKER_PASSWORD+x} ]]; then
        echo "Not publishing because DOCKER_USERNAME is unset" >&2
        echo false
        return 1
    fi
    if [[ ! -z ${TRAVIS_TAG+x} ]]; then
        echo "Publish because we are building a Tag on TravisCI" >&2
        echo true
        return 0
    fi
    if [[ -z ${CI+x} ]]; then
        echo "Publish because we are building on a local machine" >&2
        echo true
        return 0
    fi

    echo "Not publishing because of unmet conditions" >&2
    echo false
    return 1
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
        deps)
            deps
            ;;
        publish?)
            publish?
            ;;
        version)
            version
            ;;
        *)
            echo $"Usage: $0 {checks|check-go|check-docker|deps|publish?|version}"
            exit 1
esac
