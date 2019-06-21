.DEFAULT_GOAL=help

# Required for globs to work correctly
SHELL=/bin/bash


build: validate  ## Build the vagrant images
	packer build ubuntu-18.04.json 2>&1 | tee build.log

build_debug:
	PACKER_LOG=1 packer build -on-error=ask ubuntu-18.04.json 2>&1 | tee build.log

add_box:
	vagrant box add output-vagrant/package.box --name kdk/ubuntu-18.04-test --force

validate:  ## Validate the packer json file
	packer validate ubuntu-18.04.json

clean:  ## Clean up the build dirs
	@rm -rf build.log output-vagrant

help:  ## Print list of Makefile targets
	@# Taken from https://github.com/spf13/hugo/blob/master/Makefile
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	  cut -d ":" -f1- | \
	  awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
