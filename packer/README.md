# Packer builds of KDK

In addition to the Docker build of KDK, this directory enables Packer-based VM image (non-Docker) builds as well.

Supported targets/hypervisors:

* Vagrant
  * Mac/VirtualBox
  * Windows/HyperV

## Mac Virtualbox Build

```bash
# Set a Github API Token so that API call throttling limits are increased
export GITHUB_API_TOKEN=<token created from https://github.com/settings/tokens>

# Enter the packer build directory
cd kdk/packer

# Build the box
make clean build_virtualbox
```

## Windows Hyper-V Build

Before you start the Windows Hyper-V build:

* Disconnect from your corporate VPN.
* Ensure you are building from a wired connection.
  * Wireless connections almost always fail because the build is network-intensive.
* Configure Hyper-V network settings.
  * The packer build is configured to auto-select the network named 'Default Switch'.
  * Ensure the Hyper-V 'Default Switch' exists.

```bash
# Start Powershell as Administrator
# Start a bash shell in git (within the Administrator Powershell)
"C:\Program Files\git\bin\bash.exe"

# Set a Github API Token so that API call throttling limits are increased
export GITHUB_API_TOKEN=<token created from https://github.com/settings/tokens>

# Enter the packer build directory
cd kdk/packer

# Build the box
make clean build_hyperv
```

## Try the new image.

```bash
# Test the box by adding it locally as kdk/ubuntu-18.04-test
make add_box

# Create a vagrant file
mkdir -p vagrant-test
cd vagrant-test
vagrant init kdk/test

# Start the box
vagrant up

# Login to the box, and check it out
vagrant ssh

# Destroy the box
vagrant destroy -f
```

## Publish the new image.

```
# Upload the box to vagrant cloud (use your own account)
#   Upload this file: output-vagrant/package.box
https://app.vagrantup.com/dcwangmit01/boxes/kdk
```
