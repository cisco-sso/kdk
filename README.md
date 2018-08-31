Kubernetes Development Kit (KDK)
===

## Quickstart (`TL;DR`)

```console
curl -sSL https://raw.githubusercontent.com/cisco-sso/kdk/master/files/install | bash && kdk init && kdk pull && kdk up && kdk ssh
```

## Background

The kdk repository may be used to create a docker container with all of the
tools that one would typically use in order to develop and operate kubernetes
clusters.

Getting setup to create and operate a Kubernetes cluster in AWS, Openstack, or
even locally may be painful because a user may be running Windows10 or OSX, and
one must configure 20+ tools for cluster automation to work effectively. We've
created a Docker Image to enable every one of us to work in the same
environment, with the same tools, at the same versions.

Tools include: docker, kubectl, helm, multihelm, kops, terraform, ansible,
minio-cli, aws-cli, direnv, golang, git, vi/vim, emacs, python 2/3, jq, zsh,
helm-s3, kafkacat, dig, ssh-keygen, gitslave, dep, gomplate, minikube, awscli,
docker-compose, neutronclient, openstackclient, supernova, virtualenv, yq,
colordiff, nmap, screen, tmux, yadm, and many others.

* Some example use cases include:
  * Operating Kubernetes clusters.
  * Deploying Kubernetes clusters to AWS using `kops`.
  * Developing and applying Helm Charts and mh Apps.
  * Developing docker containers.

## Setup

### OSX

1. Install [homebrew](https://brew.sh/)
```console
# Open a Terminal
<Spotlight_Search -> Terminal>
/usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
```
2. Install required utils

```console
brew install git
brew cask install keybase
```

3. Install Docker from [here](https://docs.docker.com/docker-for-mac/release-notes/)

### Windows

1. Install [chocolatey](https://chocolatey.org)

```console
# Open Powershell as Administrator
Set-ExecutionPolicy Bypass -Scope Process -Force; iex ((New-Object System.Net.WebClient).DownloadString('https://chocolatey.org/install.ps1'))
```

2. Install required utils

```console
# Open Powershell as Administrator
choco install -y openssh git curl sudo
```

3. Unintall a few utils which interfere with Docker for Windows
choco uninstall vagrant virtualbox

4. Install Docker from [here](https://docs.docker.com/docker-for-windows/release-notes/)

5. Install Keybase from [here](https://keybase.io/docs/the_app/install_windows)


## Dependency Configuration

### SSH

1. Open a terminal

```console
OSX: <Spotlight_Search -> Terminal>
Windows: <Windows_Search -> "Git Bash">
```
2. Start `ssh-agent`

```console
eval `ssh-agent`
```

3. Generate ssh key

```console
if [[ ! -f ~/.ssh/id_rsa ]]; then ssh-keygen -b 4096 -t rsa -f ~/.ssh/id_rsa -q -N ""; fi
```

4. Add generated ssh key to ssh-agent

```console
ssh-add ~/.ssh/id_rsa
```

5. Add ssh public key to github

* Add content of `~/.ssh/id_rsa.pub` to [here](https://github.com/settings/keys)

6. Add ssh public key to gitbucket

* Add content of `~/.ssh/id_rsa.pub` to https://<BITBUCKET-SERVER>/bitbucket/plugins/servlet/ssh/account/keys

### Keybase

```console
# Start Keybase
OSX: <Spotlight_Search -> Keybase>
Windows: <Windows_Search -> Keybase>

# Ensure you are registered on keybase with your full name and at least one verification.  
# Ask your team lead to add you to any relevant keybase teams.

# Ensure that keybaseFS is configured and mounted on your system
<Keybase -> Folders -> "Display in Explorer" or "Open Folder" -> "Repair">

# Verify that keybaseFS has been mounted on your system
OSX: ls /keybase
Windows: dir k:
```

## Download and Initialize the KDK

1. Open a terminal

```console
OSX: <Spotlight_Search -> Terminal>
Windows: <Windows_Search -> "Git Bash">
```

2. Install the KDK

```console
curl -sSL https://raw.githubusercontent.com/cisco-sso/kdk/master/files/install | bash
```

3. create KDK config [`~/.kdk/kdk/config.yaml`] and ssh keys 

```console
kdk init
```

4. Add KDK ssh key to ssh-agent

```console
ssh-add ~/.kdk/ssh/id_rsa
```

5. Edit KDK config [`~/.kdk/kdk/config.yaml`] to suit your needs.

## Use the KDK

```bash
# Pull the latest KDK image
kdk pull

# Start the KDK
kdk up

# Connect to the KDK
kdk ssh

# Destroy the KDK
kdk destroy
```

## Customization

* **NOTE:**  By default, the `KDK` uses a set of opinionated dotfiles. 
To customize, fork [this](https://github.com/cisco-sso/yadm-dotfiles) repo, make changes, and update `~/.kdk/kdk/config.yaml` to reference customized fork.

## Common Configurations
### AWS

* `~/.aws/config`: Ensure there is an entry for each AWS account that you must
  access.  Tools such as the aws-cli, kops, and helm depend on these settings.
  The name of each profile must match that listed in the http://go2/aws (Cisco
  only) index page.

```console
# EXAMPLE: ~/.aws/config

[profile account-foo]
output = json
region = us-west-1

[profile account-bar]
output = json
region = us-west-1
```

* `~/.aws/credentials`: Ensure there is an entry for each AWS account that you
  must access.  Tools such as the aws-cli, kops, and helm depend on these
  settings.  The name of each profile must match that listed in the
  http://go2/aws index page.  Be sure to replace your key_id and access_key for
  each entry.

```console
# EXAMPLE: ~/.aws/credentials
[account-foo]
aws_access_key_id = XXXXXXXXXXXXXXXXXXXX
aws_secret_access_key = YYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYY

[account-bar]
aws_access_key_id = XXXXXXXXXXXXXXXXXXXX
aws_secret_access_key = YYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYY
```

## Updating your KDK Machine

1. Update KDK bin
```console
curl -sSL https://raw.githubusercontent.com/cisco-sso/kdk/master/files/install | bash
```
2. Download latest KDK image

```console
kdk pull
```
3. Destroy previous KDK container

```console
kdk destroy
```

4. Recreate KDK config (may be optional).

```console
kdk init
```

5. Customize `~/.kdk/kdk/config.yaml` to suit your needs (if config was regenerated).

6. Start KDK container

```console
kdk up
```

## Running Multiple KDK Containers

You might have a need to run multiple KDK containers.  The KDK CLI can do that!

1. Create a new KDK config

  - **NOTE:** port and name arguments must be unique (no other container can have this name or port assignment) 
```console
kdk init --name kdk1 --port 2023
```

2. Start `kdk1` container

```console
kdk up --name kdk1
```

3. Connect to `kdk1` container

```console
kdk ssh --name kdk1
```


**NOTE:** There are many configuration options available in `kdk init`.See `kdk init --help` for details 

## Saving and Restoring snapshots

TODO: Finish this section

## Building the KDK from scratch

TODO: Finish this section


## KDK TODOS

* [x] (Dave) Fixed issue with new kdk where it bombs out on encountering old config file format.
* [x] (Ryan) Verbose output for KDK commands
* [x] (Dave) KDK init: Autodetect keybase dirs and ask if user wants to mount them in
* [x] (Ryan) Windows 10 instructions and testing
* [x] (Dave) Refactor kdk config.yaml file to directly use Docker lib structs
* [x] (Dave) KDK init: Prompt the user and ask if they want to mount additional directories with explanation
* [ ] (Ryan) Windows 10 keybase mounts
* [x] (Ryan) KDK init: Enable starting of more than one kdk (Needed for development)
* [x] (Ryan) Curl installation/upgrade script
* [ ] (???) KDK doctor (like brew doctor)  Verifies current dependencies
* [ ] (???) Image Upgrades: Check if a later version of an image exists, and ask user if they wish to download
* [ ] (???) KDK Upgrades: Check if a later version of a KDK binary exists, and ask user if they wish to download
* [ ] (???) KDK tool add: Install [goreleaser](https://github.com/goreleaser/goreleaser/releases).