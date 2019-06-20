FROM ubuntu:bionic AS build-cache-base

LABEL maintainer="rluckie@cisco.com" \
      kdk=""

COPY provision.sh /tmp/provision.sh

#######################################
# Install OS packages

RUN /tmp/provision.sh layer_install_os_packages

#######################################
# Configure systemd and other Miscellaneous configuration bits
#  Mostly taken from:
#  https://github.com/dramaturg/docker-debian-systemd/blob/master/Dockerfile

RUN echo "Configuring systemd" && \
        cd /lib/systemd/system/sysinit.target.wants/ && \
        ls | grep -v systemd-tmpfiles-setup.service | xargs rm -f && \
        rm -f /lib/systemd/system/sockets.target.wants/*udev* && \
        systemctl mask -- \
            tmp.mount \
            etc-hostname.mount \
            etc-hosts.mount \
            etc-resolv.conf.mount \
            -.mount \
            swap.target \
            getty.target \
            getty-static.service \
            dev-mqueue.mount \
            cgproxy.service \
            systemd-tmpfiles-setup-dev.service \
            systemd-remount-fs.service \
            systemd-ask-password-wall.path \
            systemd-logind.service && \
                systemctl mask -- \
                        cron.service \
                        dbus.service \
                        exim4.service \
                        ntp.service && \
        systemctl set-default multi-user.target || true && \
        sed -ri /etc/systemd/journald.conf -e 's!^#?Storage=.*!Storage=volatile!' && \
        # Set locale && \
        localedef -i en_US -f UTF-8 en_US.UTF-8 && \
        # Configure openssh-server && \
        sed -i 's/#Port 22/Port 2022/' /etc/ssh/sshd_config && \
        # Configure docker daemon to support docker in docker && \
        mkdir /etc/docker && echo '{"storage-driver": "vfs"}' > /etc/docker/daemon.json


#######################################
# Install python-based utils and libs

RUN /tmp/provision.sh layer_install_python_based_utils_and_libs

#######################################
# Install apps (with pinned version) that are not provided by the OS packages.

RUN /tmp/provision.sh layer_install_apps_not_provided_by_os_packages

#######################################
# go get installs

FROM build-cache-base AS build-cache-multistage-goinstall

COPY provision.sh /tmp/provision.sh

RUN /tmp/provision.sh layer_go_get_installs

#######################################
# Build apps that are not provided by the OS packages.

FROM build-cache-multistage-goinstall AS build-cache-multistage-compiler

COPY provision.sh /tmp/provision.sh

RUN /tmp/provision.sh layer_build_apps_not_provided_by_os_packages

#######################################
# Piece together the final image

FROM build-cache-base AS the-final-image

COPY --from=build-cache-multistage-goinstall /go/bin/    /go/bin/
COPY --from=build-cache-multistage-compiler  /usr/local/ /usr/local/
COPY --from=build-cache-multistage-compiler  /etc/shells /etc/shells

#######################################
# Copy local files

COPY usr/local/bin/ /usr/local/bin/
COPY lib/systemd/system/ /lib/systemd/system/
RUN systemctl enable kdk-setup.service && \
    ldconfig

#######################################
# Ensure systemd starts, which subsequently starts ssh and docker

EXPOSE 2022
CMD ["/lib/systemd/systemd"]
