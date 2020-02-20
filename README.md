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

#### Docker KDK

Please use Windows10 powershell for installation.

```console
Set-ExecutionPolicy Bypass -Scope Process -Force
iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/cisco-sso/kdk/master/files/install.ps1'))
kdk init ; kdk ssh
```

NOTE: After installation, Windows CMD prompt will work. The KDK has not been
tested with Cygwin, Mingw, or Windows Subsystem for Linux.

#### Vagrant Hyper-V or Virtualbox KDK

```console
git clone git@github.com:cisco-sso/kdk.git # or https://github.com/cisco-sso/kdk.git
cd kdk
# Edit Vagrantfile: You may want to tune memory, network settings, or host-mounted directories.
vagrant up  # Starts the KDK
vagrant ssh -- -A -D 8000  # Connect to the KDK (-A ssh-agent forwarding, -D socks proxy forwarding)
# Use the KDK
vagrant destroy
```

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

## Saving State between Resetting your KDK Enviroment

The KDK is meant to be ephemeral.  You should be able to `kdk destroy && kdk ssh` whenever you need to reset your enviroment.  Resetting should be done often, because over time your environment will diverge from original state as you use it.

Here are a few approaches for saving state in between resets.

### Customing your `.bash_profile`

If you have installed [Keybase](https://keybase.io/), executing `kdk init` will ask if you would like to mount the keybase directory into the KDK.  If you do, you may customize your .bash_profile within the KDK by creating the file `/keybase/private/<user-keybase-id>/.bash_profile_private`.  The default KDK [`.bash_profile`](https://github.com/cisco-sso/yadm-dotfiles/blob/master/.bash_profile#L128) will automatically source the `.bash_profile_private` file from keybase, if it exists.

This method may be used to set enviroment variables as well as create entire dotfiles, such as `~/.aws/credentials` and `~/.aws/config`.  See [here for an example](https://github.com/cisco-sso/yadm-dotfiles#customizing-your-setup).

### Mounting Directories Directly into the KDK

Upon KDK init, you will be prompted to mount additional directories from your host system into the KDK system.  Typically this is used to mount code directories from the host machine to the KDK, but it can also be used to mount configuration directories.

Here's an example of mounting the `~/.aws` directory from an OSX machine to a location within the KDK.

```
Would you like to mount additional docker host directories into the KDK? [y/n] y
Please enter the docker host source directory (e.g. /Users/<username>/Projects) /Users/mcboats/.aws
INFO[0022] Entered host source directory mount /Users/mcboats/.aws
Please enter the docker container target directory (e.g. /home/<username>/Projects) /home/mcboats/.aws
INFO[0026] Entered container target directory mount /home/mcboats/.aws
```

### SSH-Agent

If you are using OSX, then you can use ssh-agent to automatically forward your SSH keys into the KDK.  This will allow you to access SSH resources without physically copying your keys into the KDK machine, which lowers security.  To set this up, you may manually start ssh-agent and load the SSH keys into the agent.  Or, you can copy the following [lines](https://github.com/cisco-sso/yadm-dotfiles/blob/master/.bash_profile#L19-L44) into your OSX ~/.bash_profile for a more automatic method.

### Customizing your dotfiles

If you have your own yadm dotfiles repository, you may `kdk init` with the option:
```
--dotfiles-repo string      KDK Dotfiles Repo (default "https://github.com/cisco-sso/yadm-dotfiles.git")
```

**NOTE:** There are many configuration options available in `kdk init`.See `kdk init --help` for details 
