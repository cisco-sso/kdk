# kdk devel

## Publish a release using goreleaser

1. Install [goreleaser](https://github.com/goreleaser/goreleaser/releases).
1. Create a semver tag on the main repo.
1. Clone the main repo, `cd` into it and run `dep ensure`.
1. Export a `repo`-scoped personal GitHub token (e.g. `export GITHUB_TOKEN=foo`).
1. Run `goreleaser` to build and publish to GitHub Releases.

ref: https://goreleaser.com/quick-start/
