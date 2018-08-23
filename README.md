Kubernetes Development Kit (KDK)
===

This Dockerized KDK has only been tested with OSX.  If you need Windows10
support, try [k8s-devkit](https://github.com/cisco-sso/k8s-devkit).

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

## Dependencies Setup

### OSX Specfic Setup Instructions

```bash
# Open a Terminal
<Spotlight_Search -> Terminal>

# Install Homebrew (https://brew.sh/)
/usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"

# Install Git, Keybase
#   If you are a Keybase user, the KeybaseFS may be mounted directly into the docker image.
brew install git
brew cask install keybase

# Install Docker for Mac
open https://docs.docker.com/docker-for-mac/install/
```

### Windows Specific Setup Instructions

```bash
# Open a Windows Powershell
<Windows_Search -> Powershell (Right click, Start as Administrator)>

# Install Chocolatey (https://chocolatey.org)
Set-ExecutionPolicy Bypass -Scope Process -Force; iex ((New-Object System.Net.WebClient).DownloadString('https://chocolatey.org/install.ps1'))

# Install Git, Keybase, and other utils
#   If you are a Keybase user, the KeybaseFS may be mounted directly into the docker image.
choco.exe install -y keybase openssh git curl
```

* Download and install docker from [here](https://docs.docker.com/docker-for-windows/release-notes/)

## Dependencies Configuration

## Configure Keybase

```bash
# Start Keybase
OSX: <Spotlight_Search -> Keybase>
Windows: <Windows_Search -> Keybase>

# Ensure you are registered on keybase with your full name and at least one
#   verification.  Keybase is the encrypted store used to share team secrets.
# Ask your team lead to add you to any relevant keybase teams.

# Ensure that keybaseFS is configured and mounted on your system
<Keybase -> Folders -> "Display in Explorer" or "Open Folder" -> "Repair">

# Verify that keybaseFS has been mounted on your system
OSX: ls /keybase
Windows: dir k:
```

## Configure SSH

```bash

# Open a bash shell and go to your home directory
OSX: <Spotlight_Search -> Terminal>
Windows: <Windows_Search -> "Git Bash">
cd ~/

# Ensure you have an ssh-key generated with default settings
#   Paste the following into bash
if [ ! -e ~/.ssh/id_rsa ]; then ssh-keygen -b 4096; done

# Provision the ssh key in your github.com account
#   Your new public key is here: ~/.ssh/id_rsa.pub
#   https://github.com/settings/keys

# Provision the ssh key in your bitbucket account
#   Your new public key is here: ~/.ssh/id_rsa.pub
#   https://<BITBUCKET-SERVER>/bitbucket/plugins/servlet/ssh/account/keys
```

## Download and Configure the Dockerized KDK

```bash

# Open a bash shell and go to your home directory
OSX: <Spotlight_Search -> Terminal>
Windows: <Windows_Search -> "Git Bash">
cd ~/

# If you want to save your files in between VM creation and destroy, create a
# ~/Dev directory which will be auto-mounted into the virtualmachine from the
# host.  This currently is NOT RECOMMENDED FOR WINDOWS, because git cloned
# symlinks and file line-endings do not work well on a windows host-mounted fs.
# Using a host mounted ~/Dev directory is favorable so that you are able to
# edit source code on the host machine using your host editor.
OSX: mkdir ~/Dev; cd ~/Dev
Windows: <Ignore This>



# Download the binary for your kdk (OSX)
curl -sSL https://github.com/cisco-sso/kdk/releases/download/0.5.3/kdk-0.5.3-darwin-amd64.tar.gz \
  | tar xz && chmod +x darwin-amd64/kdk && sudo mv darwin-amd64/kdk /usr/local/bin/kdk && rm -rf darwin-amd64

# Create your ~/.kdk/config.yaml
kdk init

# Start ssh-agent and load your key
eval `ssh-agent`
ssh-add ~/.kdk/ssh/id_rsa
ssh-add -l  # verify that the key has been loaded


# Edit your ~/.kdk/config.yaml with additional volume mounts if you would like
#   to mount host directories into the docker machine.  Additional volume bind
#   mount examples include:
  binds:
    - source: /Users/<YOUR_USERID>/<YOUR_PROJECT_DIR>
      target: /home/<YOUR_USERID>/<YOUR_PROJECT_DIR>
    - source:  "/Volumes/Keybase (<YOUR_KEYBASE_ID)/"
      target: /keybase

```

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

* **NOTE:**  The `kdk up` binary uses a set of opinionated dotfiles by default
* Fork [this](https://github.com/cisco-sso/yadm-dotfiles) repo, make changes,
  and update `launch-kdk` script accordingly to point to your customized fork.


## Configuring your KDK Machine

* `~/.aws/config`: Ensure there is an entry for each AWS account that you must
  access.  Tools such as the aws-cli, kops, and helm depend on these settings.
  The name of each profile must match that listed in the http://go2/aws (Cisco
  only) index page.

```bash
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

```bash
# EXAMPLE: ~/.aws/credentials
[account-foo]
aws_access_key_id = XXXXXXXXXXXXXXXXXXXX
aws_secret_access_key = YYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYY

[account-bar]
aws_access_key_id = XXXXXXXXXXXXXXXXXXXX
aws_secret_access_key = YYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYY
```

## Using your KDK Machine

```
cd ~/
git clone git@github.com:cisco-sso/k8s-deploy.git  # Or platform-deploy within Cisco
cd k8s-deploy
direnv allow

# All of your work must be done from a cluster directory.  Upon entering a
#   cluster directory, `direnv` will automatically set your enviromental
#   configurations.  Upon entering a cluster directory for the first time, you
#   must run `direnv allow` to permanently record that direnv is allowed to
#   execute the .envrc script.

# Activate cluster1 settings by entering the directory
cd clusters/cluster1.domain.com
direnv allow

# Ensure that aws cli works
#   Upon failure, check your ~/.aws config files
aws ec2 describe-instances

# Check that kops works
kops validate cluster

# Check that kubectl works
kubectl cluster-info

# Check that helm works
helm ls

# Activate cluster3 settings by entering the directory
cd ../cluster3.domain.com
direnv allow
... <do the same thing above to verify that you can access cluster3>
```

## Updating your KDK Machine

```bash
# Re-install by downloading the latest binary (OSX)
curl -sSL https://github.com/cisco-sso/kdk/releases/download/0.5.3/kdk-0.5.3-darwin-amd64.tar.gz \
  | tar xz && chmod +x darwin-amd64/kdk && sudo mv darwin-amd64/kdk /usr/local/bin/kdk && rm -rf darwin-amd64

# Download the latest image
kdk pull
kdk destroy
kdk up
```

## Saving and Restoring snapshots

It is often useful to save a snapshot of the vagrant machine.

TODO: Finish this section

## Building the KDK from scratch

TODO: Finish this section


## KDK TODOS

* [x] (Dave) Fixed issue with new kdk where it bombs out on encountering old config file format.
* [x] (Ryan) Verbose output for KDK commands
* [x] (Dave) KDK init: Autodetect keybase dirs and ask if user wants to mount them in
* [ ] (Ryan) Windows 10 instructions and testing
* [ ] (Dave) Refactor kdk config.yaml file to directly use Docker lib structs
* [ ] (Dave) Windows 10 keybase mounts
* [ ] (???) KDK init: Enable starting of more than one kdk (Needed for development)
* [ ] (???) KDK init: Prompt the user and ask if they want to mount additional directories with explanation
* [ ] (???) KDK init: If a dir doesn't exist upon startup it should warn or error with an explanation to check that keybaseFS is mounted.
* [ ] (???) Curl installation/upgrade script
* [ ] (???) KDK doctor (like brew doctor)  Verifies current dependencies
* [ ] (???) Image Upgrades: Check if a later version of an image exists, and ask user if they wish to download
* [ ] (???) KDK Upgrades: Check if a later version of a KDK binary exists, and ask user if they wish to download
