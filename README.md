Dockerized Kubernetes Development Kit
===

A docker image with tools for Kubernetes, Helm and Docker DevOps.

## Dependencies

* [Docker for Mac](https://docs.docker.com/docker-for-mac/)
* [Python 3.x](http://docs.python-guide.org/en/latest/starting/install3/osx/) and [Virtualenv](https://virtualenv.pypa.io/en/stable/)
* [Python requirements](requirements.txt)

## Getting Started

1. Get KDK script

```bash
git clone git@github.com:cisco-sso/kdk.git
cd kdk

# In the future, the kdk will be a golang binary installed like this:
# curl -so kdk https://raw.githubusercontent.com/cisco-sso/dockerized-k8s-devkit/master/kdk && chmod +x kdk
```

2. Install Pre-Reqs

```bash
virtualenv -p python3 .venv
source .venv/bin/activate
pip install -r requirements.txt
```

3. Initilize KDK Configurations

The follow command generates a default working `~/.kdk/config.yaml` configuration file

```bash
./kdk init
```

4. Customize the Configuration

Customize `~/.kdk/config.yaml` to fit your needs.  Most people should add the
following to the `volumes` section.

```yaml
  volumes:
    /Users/<USERNAME>/<FROM_HOST_FOLDER>:
      bind: /home/<USERNAME>/<TO_GUEST_FOLDER>:
      mode: rw
    "/Volumes/Keybase (<YOUR_KEYBASE_ID>)/":
      bind: /keybase
      mode: rw
```

5. Start KDK container

```bash
./kdk up
```

6. Exec to KDK container

```bash
./kdk ssh
```

## Customization
* **NOTE:**  The `launch-kdk` script uses a set of opinionated dotfiles by default
* Fork [this](https://github.com/cisco-sso/yadm-dotfiles) repo, make changes, and update `launch-kdk` script accordingly to point to your customized fork.
