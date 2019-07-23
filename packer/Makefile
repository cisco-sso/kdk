.DEFAULT_GOAL=help

# Required for globs to work correctly
SHELL=/bin/bash

deps_windows:
	choco install -y make packer git

build_virtualbox: validate  ## Build the virtualbox vagrant image
	PACKER_LOG=1 packer build packer_virtualbox.json 2>&1 | tee build-virtualbox.log

build_hyperv: validate  ## Build the hyperv vagrant images
	PACKER_LOG=1 packer build packer_hyperv.json     2>&1 | tee build-hyperv.log

add_box:  ## Add the box as kdk/ubuntu-18.04-test
	vagrant box remove kdk/ubuntu-18.04-test 2>&1 > /dev/null || true
	vagrant box add output-vagrant/package.box --name kdk/ubuntu-18.04-test --force

validate:  ## Validate the packer json files
	packer validate packer_virtualbox.json
	packer validate packer_hyperv.json

clean:  ## Clean up the build dirs
	@rm -rf *.log output-vagrant *~

help:  ## Print list of Makefile targets
	@# Taken from https://github.com/spf13/hugo/blob/master/Makefile
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	  cut -d ":" -f1- | \
	  awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
