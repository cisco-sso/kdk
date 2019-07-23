# Packer builds of KDK

In addition to the Docker build of KDK, this directory enables Packer-based VM image (non-Docker) builds as well.

Supported targets/hypervisors:

* Vagrant
  * Mac/VirtualBox
  * Windows/HyperV

## Mac Virtualbox Build

```bash

# Enter the packer build directory
cd kdk/packer

# Validate the Packer template
make validate

# Build the box
make clean build_virtualbox
```

## Windows Hyper-V Build

```bash
# Start Powershell as Administrator
# Start a bash shell in git
C:\Program Files\git\bin\bash.exe

# Enter the packer build directory
cd kdk/packer

# Validate the Packer template
make validate

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
vagrant init kdk/ubuntu-18.04-test

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
