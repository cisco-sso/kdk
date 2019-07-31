Kubernetes Development Kit (KDK)
===

## Quickstart (`TL;DR`)

This Quickstart assumes that you have installed all of the
[dependencies](https://github.com/cisco-sso/kdk#installation-instructions).


### Mac and Linux

```console
curl -sSL https://raw.githubusercontent.com/cisco-sso/kdk/master/files/install | bash
kdk init && kdk ssh
```

### Windows

Please use Windows10 powershell for installation.

```console
Set-ExecutionPolicy Bypass -Scope Process -Force
iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/cisco-sso/kdk/master/files/install.ps1'))
kdk init ; kdk ssh
```

NOTE: After installation, Windows CMD prompt will work. The KDK has not been
tested with Cygwin, Mingw, or Windows Subsystem for Linux.

## Installation Instructions

Detailed installation instructions of the KDK along with all of its dependencies are found here:

* [Mac](https://kdf.csco.cloud/getting-started/mac/)
* [Windows](https://kdf.csco.cloud/getting-started/windows/)
* [Linux](https://kdf.csco.cloud/getting-started/linux/)

## Background

Getting your workstation setup to work with Kubernetes clusters may require the install and configuration of quite a few tools. Do it inconsistently among your team, and your automation and workflows may not work properly for everyone. Even if it works on your machine because the latest code you've written requires the latest version of `kubectl` and a new installation of `jq`, your teammates Billy on Windows and Jane on Mac are busy filing bugs against your latest PR because they haven't received the memo about updating their toolchains.

We've created the open-source Kubernetes Development Kit (KDK) in order to solve this problem. The KDK is a docker container or a vagrant virtual machine, which may be deployed on Mac, Windows, and Linux. It is a Linux-based environment which has over 30+ tools pre-installed and pre-configured. If your team uses the KDK, then you are guaranteed to have a similar development and operations environment, with the same tools, at the same versions.

A sampling of tools include: docker, kubectl, helm, helmfile, kops, kubetail, docker-compose, terraform, ansible, minio-cli, aws-cli, gcloud, drone-cli, direnv, golang, git, hub, jsonnet, vi/vim, emacs, python 2/3, pipenv, pyenv, jq, zsh, helm-s3, kafkacat, dig, ssh-keygen, dep, gomplate, minikube, neutronclient, openstackclient, supernova, virtualenv, yq, colordiff, nmap, screen, tmux, sshuttle, yadm, and many others.

The KDK may make your life easier if you often:

* Operate Kubernetes clusters.
* Deploy Kubernetes clusters to various clouds including AWS and GKE.
* Develop and deploy Helm Charts.
* Develop docker containers.


## Basic Usage

1. Create or re-create the config

```console
kdk init
```

2. Connect or reconnect to the KDK (will pull and start container if necessary)

```console
kdk ssh
```

3. Destroy the KDK

```console
kdk destroy
```

4. Update the KDK (binary, config, and container)

```console
kdk update
```

## Running Multiple KDK Containers

You might have a need to run multiple KDK containers.  The KDK CLI can do that!

1. Create a new KDK config

  - **NOTE:** name parameter must be unique (no other container can have this name)
```console
kdk init --name kdk1
```

2. Connect to `kdk1` container

```console
kdk ssh --name kdk1
```

**NOTE:** There are many configuration options available in `kdk init`.See `kdk init --help` for details 
