# Packer builds of KDK

In addition to the Docker build of KDK, this directory enables Packer-based VM image (non-Docker) builds as well.

Supported targets/hypervisors:

* Vagrant
  * VirtualBox

## Validate the Packer template and build a new KDK image.

```bash
cd kdk/packer

packer validate ubuntu-18.04.json

packer build ubuntu-18.04.json
```

## Add the new image to Vagrant.

```bash
vagrant box add output-vagrant/package.box --name kdk/ubuntu-18.04
```

## Try the new image.

```bash
mkdir ~/vagrant/kdk

cd ~/vagrant/kdk

vagrant init kdk/ubuntu-18.04

vagrant up

vagrant ssh
```

## TODO: Publish the new image.
