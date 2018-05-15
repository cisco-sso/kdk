KDK (Kubernetes Development Kit) [dockerized]
===

A docker image with tools for Kubernetes, Helm and Docker DevOps.

## Dependencies

* Docker 

## Getting Started

### Quick Start

```bash
bash <(curl -s https://raw.githubusercontent.com/cisco-sso/k8s-devkit/master/docker/scripts/launch-kdk)
```
**OR**

### Customized
```bash
curl -fL https://raw.githubusercontent.com/cisco-sso/k8s-devkit/master/docker/scripts/launch-kdk && chmod +x launch-kdk
```
Customize to fit your needs, then ...

```bash
./launch-kdk
```

### Provision User
After the container has started...

```bash
provision-user
```

### Usage
After the `provision-user` script has completed you should be ready to use all the KDK container has to offer.

## Saving State and Reusing KDK Image

You can save the state of your running KDK container by using `docker commit`

Example:

```bash
docker commit kdk kdk-snapshot
```

You can use the snapshot image with the `launch-kdk` script by setting `KDK_IMAGE` environment variable

Example

```bash
export KDK_IMAGE=kdk-snapshot
./launch-kdk

```

## Customization
* **NOTE:**  The `launch-kdk` script uses a set of opinionated dotfiles by default
* Fork [this](https://github.com/rtluckie/work-dotfiles) repo, make changes, and update `launch-kdk` script accordingly to point to your customized fork.
