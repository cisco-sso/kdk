Dockerized Kubernetes Development Kit
===

A docker image with tools for Kubernetes, Helm and Docker DevOps.

## Dependencies

* Docker 

## Getting Started

1. Get KDK script

```bash
curl -so kdk https://raw.githubusercontent.com/cisco-sso/dockerized-k8s-devkit/master/scripts/kdk; chmod +x kdk
```
2. Init KDK

```bash
kdk init
```

Customize `~/.kdk/config.yaml` to fit your needs.

3. Start KDK container
```bash
kdk start
```

4. Connect to KDK container

```bash
docker exec -it kdk bash
```

5. Provision user

```bash
provision-user
```

### Usage
After the `provision-user` script has completed you should be ready to use all the KDK container has to offer.


## Customization
* **NOTE:**  The `launch-kdk` script uses a set of opinionated dotfiles by default
* Fork [this](https://github.com/rtluckie/work-dotfiles) repo, make changes, and update `launch-kdk` script accordingly to point to your customized fork.
