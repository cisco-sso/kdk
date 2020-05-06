FROM ubuntu:focal AS build-cache-base

LABEL maintainer="rluckie@cisco.com" \
      kdk=""

COPY provision.sh /tmp/provision.sh

#######################################
# Install OS packages

RUN /tmp/provision.sh layer_install_os_packages

#######################################
# Install python-based utils and libs

RUN /tmp/provision.sh layer_install_python_based_utils_and_libs

#######################################
# Install apps (with pinned version) that are not provided by the OS packages.

RUN GITHUB_API_TOKEN=$(cat /run/secrets/github_api_token) /tmp/provision.sh layer_install_apps_not_provided_by_os_packages

#######################################
# go get installs

FROM build-cache-base AS build-cache-multistage-goinstall

COPY provision.sh /tmp/provision.sh

RUN /tmp/provision.sh layer_go_get_installs

#######################################
# Build apps that are not provided by the OS packages.

FROM build-cache-multistage-goinstall AS build-cache-multistage-compiler

COPY provision.sh /tmp/provision.sh

RUN echo "Clean out /usr/local" && rm -rf /usr/local/*
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
# Configure systemd and other Miscellaneous configuration bits
#  Mostly taken from:
#  https://github.com/dramaturg/docker-debian-systemd/blob/master/Dockerfile

RUN echo "Configuring systemd" && \
        cd /lib/systemd/system/sysinit.target.wants/ && \
        ls | grep -v systemd-tmpfiles-setup.service | xargs rm -f && \
        rm -f /lib/systemd/system/sockets.target.wants/*udev* && \
        systemctl mask -- \
            apt-daily-upgrade.timer \
            apt-daily.timer \
            e2scrub_all.timer \
            fstrim.timer \
            getty-static.service \
            getty.target \
            motd-news.timer \
            swap.swap \
            swap.target \
            systemd-ask-password-wall.path \
            systemd-logind.service \
            systemd-remount-fs.service \
            systemd-tmpfiles-setup.service \
            tmp.mount && \
        systemctl mask -- \
            cron.service \
            dbus.service \
            ntp.service && \
        systemctl set-default multi-user.target || true && \
        sed -ri /etc/systemd/journald.conf -e 's!^#?Storage=.*!Storage=volatile!' && \
        # Avoid port binding confict between dnsmasq and systemd-resolved && \
        sed -i 's/#DNSStubListener=yes/DNSStubListener=no/' /etc/systemd/resolved.conf && \
        # Set locale && \
        localedef -i en_US -f UTF-8 en_US.UTF-8 && \
        # Configure openssh-server && \
        sed -i 's/#Port 22/Port 2022/' /etc/ssh/sshd_config && \
        # Configure docker daemon to support docker in docker && \
        mkdir /etc/docker && echo '{"storage-driver": "vfs"}' > /etc/docker/daemon.json

#######################################
# Ensure systemd starts, which subsequently starts ssh and docker

EXPOSE 2022
CMD ["/lib/systemd/systemd"]
