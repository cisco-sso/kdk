FROM debian:9.8 AS build-cache-base

LABEL maintainer="rluckie@cisco.com" \
      kdk=""

#######################################
# Install OS packages

# big items: gcc, python-dev

RUN echo "Install OS packages" && \
    apt-get -y update && apt-get --no-install-recommends -y install \
        apache2-utils \
        apt-transport-https \
        bash-completion \
        bc \
        bridge-utils \
        ca-certificates \
        colordiff \
        ctags \
        curl \
        dc \
        dhcpdump \
        dialog \
        dnsmasq \
        dnsutils \
        dos2unix \
        fonts-powerline \
        fio \
        gcc \
        gettext \
        gnupg \
        gnupg2 \
        htop \
        less \
        libcurl4-openssl-dev \
        libevent-2.0-5 \
        locales \
        make \
        man \
        moreutils \
        nmap \
        ntp \
        ntpdate \
        openconnect \
        openssh-server \
        perl \
        proxychains \
        python \
        python-dev \
        python3 \
        python3-dev \
        qemu-user-static \
        ruby \
        screen \
        socat \
        software-properties-common \
        sudo \
        systemd \
        systemd-sysv \
        tcl \
        telnet \
        traceroute \
        tree \
        unzip \
        wget \
        whois \
        xauth && \
    curl -sSL https://download.docker.com/linux/debian/gpg | apt-key add - && \
       add-apt-repository \
         "deb [arch=amd64] https://download.docker.com/linux/debian \
         $(lsb_release -cs) \
         stable" && \
    export CLOUD_SDK_REPO="cloud-sdk-$(lsb_release -c -s)" && \
        echo "deb http://packages.cloud.google.com/apt $CLOUD_SDK_REPO main" | tee -a /etc/apt/sources.list.d/google-cloud-sdk.list && \
        curl -sSL https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add - && \
    export POSTGRESQL_REPO="$(lsb_release -c -s)-pgdg" && \
        echo "deb https://apt.postgresql.org/pub/repos/apt $POSTGRESQL_REPO main" | tee -a /etc/apt/sources.list.d/pgdg.list && \
        curl -sSL https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add - && \
    curl -sL https://deb.nodesource.com/setup_11.x | bash - && \
        curl -sL https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add - && \
        echo "deb https://dl.yarnpkg.com/debian/ stable main" | tee -a /etc/apt/sources.list.d/yarn.list && \
    apt-get update -y && apt-get --no-install-recommends -y install \
        docker-ce \
        google-cloud-sdk \
        nodejs \
        postgresql-client-10 \
        yarn && \
    # Add deps to enable pyenv-driven on-demand python installs on KDK \
    apt-get --no-install-recommends -y install \
        make build-essential \
        libssl-dev \
        zlib1g-dev \
        libbz2-dev \
        libreadline-dev \
        libsqlite3-dev \
        wget \
        curl \
        llvm \
        libncurses5-dev \
        libncursesw5-dev \
        xz-utils \
        tk-dev \
        libffi-dev \
        liblzma-dev \
        python-openssl && \
    apt-get -y clean && apt-get -y autoremove && rm -rf /var/lib/apt/lists/*


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

RUN echo "Install python-based utils and libs" && \
    curl -sSL https://bootstrap.pypa.io/get-pip.py | python2 && \
    pip2.7 install --no-cache-dir -U setuptools && \
    pip2.7 install \
        --no-cache-dir \
        --ignore-installed six \
        'ansible==2.6.4' \
        'awscli==1.16.14' \
        'boto==2.49.0' \
        'boto3==1.9.4' \
        'docker-compose==1.22.0' \
        'idna==2.6' \
        'Jinja2==2.10' \
        'jinja2-cli[yaml]==0.6.0' \
        'openshift==0.7.1' \
        'passlib==1.7.1' \
        'python-neutronclient==6.10.0' \
        'python-octaviaclient==1.7.0' \
        'python-openstackclient==3.16.1' \
        'pyvmomi==6.7.0.2018.9' \
        'urllib3==1.22' \
        'virtualenv==16.0.0' \
        'yq==2.7.0' && \
    curl -sSL https://bootstrap.pypa.io/get-pip.py | python3 && \
    pip3.5 install --no-cache-dir -U setuptools && \
    pip3.5 install \
         --no-cache-dir \
         'ansible==2.6.4' \
         'awscli==1.16.14' \
         'boto==2.49.0' \
         'boto3==1.9.4' \
         'docker-compose==1.22.0' \
         'idna==2.6' \
         'Jinja2==2.10' \
         'jinja2-cli[yaml]==0.6.0' \
         'openshift==0.7.1' \
         'passlib==1.7.1' \
         'peru==1.2.0' \
         'pipenv==2018.11.26' \
         'python-neutronclient==6.10.0' \
         'python-octaviaclient==1.7.0' \
         'python-openstackclient==3.16.1' \
         'pyvmomi==6.7.0.2018.9' \
         'urllib3==1.22' \
         'virtualenv==16.0.0' \
         'yq==2.7.0' && \
    rm -rf /root/.cache/pip


#######################################
# Install apps (with pinned version) that are not provided by the OS packages.

RUN echo "Install apps (with pinned version) that are not provided by the OS packages." && \
    echo "Install dep." && \
        curl -sSLo dep https://github.com/golang/dep/releases/download/v0.5.1/dep-linux-amd64 && \
        chmod a+x dep && mv dep /usr/local/bin && \
    echo "Install direnv." && \
        curl -sSLo direnv https://github.com/direnv/direnv/releases/download/v2.22.0/direnv.linux-amd64 && \
        chmod a+x direnv && mv direnv /usr/local/bin && \
    echo "Install drone-cli." && \
        curl -sSL https://github.com/drone/drone-cli/releases/download/v1.1.0/drone_linux_amd64.tar.gz | tar xz && \
        chmod a+x drone && mv drone /usr/local/bin/drone-1.1.0 && \
        ln -s /usr/local/bin/drone-1.1.0 /usr/local/bin/drone && \
    echo "Install easy-rsa." && \
        curl -sSL https://github.com/OpenVPN/easy-rsa/releases/download/v3.0.6/EasyRSA-unix-v3.0.6.tgz | tar xz && \
        chmod a+x EasyRSA-* && mv EasyRSA-* /usr/local/bin/easyrsa && \
    echo "Install go-task." && \
        curl -sSL https://github.com/go-task/task/releases/download/v2.5.1/task_linux_amd64.tar.gz | tar -C /usr/local/bin -xz task && chmod a+x /usr/local/bin/task && \
    echo "Install gomplate." && \
        curl -sSLo gomplate https://github.com/hairyhenderson/gomplate/releases/download/v3.4.0/gomplate_linux-amd64 && \
        chmod a+x gomplate && mv gomplate /usr/local/bin && \
    echo "Install golang." && \
        curl -sSL https://dl.google.com/go/go1.12.4.linux-amd64.tar.gz | tar -C /usr/local -xz && \
        mkdir -p /go && chmod a+rw /go && \
    echo "Install goreleaser." && \
        curl -sSLO https://github.com/goreleaser/goreleaser/releases/download/v0.101.0/goreleaser_Linux_x86_64.tar.gz && \
        tar -C /usr/local/bin -xzf goreleaser*.tar.gz goreleaser && rm goreleaser*.tar.gz && \
    echo "Install grpcurl." && \
        curl -sSL https://github.com/fullstorydev/grpcurl/releases/download/v1.1.0/grpcurl_1.1.0_linux_x86_64.tar.gz | tar -C /usr/local/bin -xz grpcurl && chmod a+x /usr/local/bin/grpcurl && \
    echo "Install helm." && \
        curl -sSL https://storage.googleapis.com/kubernetes-helm/helm-v2.11.0-linux-amd64.tar.gz | tar xz && \
          chmod a+x linux-amd64/helm && mv linux-amd64/helm /usr/local/bin/helm-2.11.0 && rm -fr linux-amd64 && \
        curl -sSL https://storage.googleapis.com/kubernetes-helm/helm-v2.12.3-linux-amd64.tar.gz | tar xz && \
          chmod a+x linux-amd64/helm && mv linux-amd64/helm /usr/local/bin/helm-2.12.3 && rm -fr linux-amd64 && \
        ln -s /usr/local/bin/helm-2.12.3 /usr/local/bin/helm && \
    echo "Install helmfile" && \
        curl -sSLo helmfile https://github.com/roboll/helmfile/releases/download/v0.54.2/helmfile_linux_amd64 && \
        chmod a+x helmfile && mv helmfile /usr/local/bin && \
    echo "Install hugo." && \
        curl -sSL https://github.com/gohugoio/hugo/releases/download/v0.54.0/hugo_0.54.0_Linux-64bit.tar.gz | tar xz && \
        chmod a+x hugo && mv hugo /usr/local/bin/hugo && \
    echo "Install jq." && \
        curl -sSLo jq https://github.com/stedolan/jq/releases/download/jq-1.6/jq-linux64 && \
        chmod a+x jq && mv jq /usr/local/bin/ && \
    echo "Install kops." && \
        curl -sSLo kops-1.10.1 https://github.com/kubernetes/kops/releases/download/1.10.1/kops-linux-amd64 && \
        chmod a+x kops-1.10.1 && mv kops-1.10.1 /usr/local/bin/ && \
        curl -sSLo kops-1.11.0 https://github.com/kubernetes/kops/releases/download/1.11.0/kops-linux-amd64 && \
        chmod a+x kops-1.11.0 && mv kops-1.11.0 /usr/local/bin/ && \
        ln -s /usr/local/bin/kops-1.11.0 /usr/local/bin/kops && \
    echo "Install kubectl." && \
        curl -sSLo /usr/local/bin/kubectl-1.10.13 https://storage.googleapis.com/kubernetes-release/release/v1.10.13/bin/linux/amd64/kubectl && \
        curl -sSLo /usr/local/bin/kubectl-1.11.7 https://storage.googleapis.com/kubernetes-release/release/v1.11.7/bin/linux/amd64/kubectl && \
        curl -sSLo /usr/local/bin/kubectl-1.12.6 https://storage.googleapis.com/kubernetes-release/release/v1.12.6/bin/linux/amd64/kubectl && \
        curl -sSLo /usr/local/bin/kubectl-1.13.3 https://storage.googleapis.com/kubernetes-release/release/v1.13.3/bin/linux/amd64/kubectl && \
        chmod a+x /usr/local/bin/kubectl-* && \
        ln -s /usr/local/bin/kubectl-1.13.3 /usr/local/bin/kubectl && \
    echo "Install kubetail." && \
        curl -sSLo kubetail.zip https://github.com/johanhaleby/kubetail/archive/1.6.7.zip && \
        unzip -qq kubetail.zip && chmod a+x kubetail-1.6.7/kubetail && mv kubetail-1.6.7/kubetail /usr/local/bin && \
        rm -f kubetail.zip && \
    echo "Install mc." && \
        curl -sSLo /usr/local/bin/mc https://dl.minio.io/client/mc/release/linux-amd64/mc.RELEASE.2019-02-20T22-21-50Z && \
        chmod a+x /usr/local/bin/mc && \
    echo "Install minikube." && \
        curl -sSLo minikube https://storage.googleapis.com/minikube/releases/v0.28.2/minikube-linux-amd64 && \
        chmod a+x minikube &&  mv minikube /usr/local/bin/ && \
    echo "Install terraform." && \
        curl -sSLo terraform.zip https://releases.hashicorp.com/terraform/0.11.11/terraform_0.11.11_linux_amd64.zip && \
        unzip -qq terraform.zip && chmod a+x terraform && mv terraform /usr/local/bin && rm -f terraform.zip && \
    echo "Install testssl." && \
        curl -sSL https://github.com/drwetter/testssl.sh/archive/v2.9.5-7.tar.gz | tar xz && \
        mv testssl* /usr/local/share/testssl && ln -s /usr/local/share/testssl/testssl.sh /usr/local/bin/testssl && chmod a+x /usr/local/bin/testssl && \
    echo "Install yadm." && \
        curl -sSL https://github.com/TheLocehiliosan/yadm/archive/1.12.0.tar.gz | tar xz && \
        mv yadm* /usr/local/share/yadm && ln -s /usr/local/share/yadm/yadm /usr/local/bin/yadm && chmod a+x /usr/local/bin/yadm

#######################################
# go get installs

FROM build-cache-base AS build-cache-multistage-goinstall

ARG GOPATH=/go
RUN echo "go get installs" && \
   apt-get -y update && apt-get --no-install-recommends -y install git && \
   apt-get -y clean && apt-get -y autoremove && rm -rf /var/lib/apt/lists/*
RUN /usr/local/go/bin/go get github.com/spf13/cobra/cobra
RUN /usr/local/go/bin/go get github.com/kubernetes-incubator/cri-tools/cmd/crictl
RUN /usr/local/go/bin/go get golang.org/x/lint/golint
RUN /usr/local/go/bin/go get github.com/gpmgo/gopm
RUN /usr/local/go/bin/go get github.com/vmware/govmomi/govc
RUN /usr/local/go/bin/go get github.com/github/hub
RUN /usr/local/go/bin/go get github.com/cisco-sso/mh && ln -sf /go/bin/mh /go/bin/multihelm
RUN /usr/local/go/bin/go get github.com/mikefarah/yq
RUN rm -rf /root/.cache/go-build

#######################################
# Build apps that are not provided by the OS packages.

FROM build-cache-multistage-goinstall AS build-cache-multistage-compiler

RUN echo "Install OS BUILD packages" && \
    apt-get -y update && apt-get --no-install-recommends -y install \
        autoconf \
        build-essential \
        libgnutls28-dev \
        libncurses5-dev \
        libz-dev && \
   apt-get -y clean && apt-get -y autoremove && rm -rf /var/lib/apt/lists/*

RUN echo "Clean out /usr/local" && \
    rm -rf /usr/local/*

RUN echo "Install git (needs to build first as a dependency)." && \
    curl -sSL https://github.com/git/git/archive/v2.19.1.tar.gz | tar xz && cd git-* && \
    make configure && ./configure --prefix=/usr/local && make && make install && cd .. && rm -fr git-*

RUN echo "Install bats" && \
    curl -sSL https://github.com/sstephenson/bats/archive/v0.4.0.tar.gz | tar xz && cd bats-* && \
    ./install.sh /usr/local && cd .. && rm -fr bats-*

RUN echo "Install emacs." && \
    curl -sSL http://mirrors.ibiblio.org/gnu/ftp/gnu/emacs/emacs-26.1.tar.gz | tar xz && cd emacs-* && \
    CANNOT_DUMP=yes ./configure \
        --prefix=/usr/local \
        --disable-build-details \
        --without-all \
        --without-x \
        --without-x-toolkit \
        --without-sound \
        --with-xml2 \
        --with-zlib \
        --with-modules \
        --with-file-notification \
        --with-gnutls \
        --with-compress-install && \
    make && make install && cd .. && rm -fr emacs-*

RUN echo "Install jsonnet" && \
    curl -sSL https://github.com/google/jsonnet/archive/v0.12.1.tar.gz | tar xz && cd jsonnet-* && \
    make && chmod a+x jsonnet && mv jsonnet /usr/local/bin

RUN echo "Install pyenv with dependencies." && \
    curl -sSLo pyenv-installer https://raw.githubusercontent.com/pyenv/pyenv-installer/master/bin/pyenv-installer && \
    chmod a+x pyenv-installer && mv pyenv-installer /usr/local/bin && \
    PYENV_ROOT=/usr/local/pyenv pyenv-installer && chmod -R a+rwx /usr/local/pyenv

RUN echo "Install vim." && \
    curl -sSL https://github.com/vim/vim/archive/v8.1.0481.tar.gz | tar xz && cd vim-* && \
    ./configure \
       --with-features=huge \
       --enable-multibyte \
       --enable-rubyinterp=yes \
       --enable-pythoninterp=yes \
       --with-python-config-dir=/usr/lib/python2.7/config-x86_64-linux-gnu \
       --enable-python3interp=yes \
       --with-python3-config-dir=/usr/lib/python3.5/config-3.5m-x86_64-linux-gnu \
       --enable-perlinterp=yes \
       --enable-luainterp=yes \
       --enable-cscope \
      --prefix=/usr/local && \
    make VIMRUNTIMEDIR=/usr/local/share/vim/vim81 && make install && cd .. && rm -fr vim-*

RUN echo "Install tmux." && \
    curl -sSL https://github.com/libevent/libevent/releases/download/release-2.1.8-stable/libevent-2.1.8-stable.tar.gz | tar xz && cd libevent-* && \
    ./configure && make &&  make install && cd .. && rm -fr libevent-* && \
    curl -sSL https://github.com/tmux/tmux/releases/download/2.8/tmux-2.8.tar.gz | tar xz && cd tmux-* && \
    ./configure --prefix=/usr/local && make && make install && cd .. && rm -fr tmux-*

RUN echo "Install zsh." && \
    curl -sSL https://sourceforge.net/projects/zsh/files/zsh/5.6.2/zsh-5.6.2.tar.xz/download | tar Jx && cd zsh-* && \
    ./configure --with-tcsetpgrp --prefix=/usr/local && make && make install && echo "/usr/local/bin/zsh" >> /etc/shells && cd .. && rm -fr zsh-*

RUN echo "Install redis-cli tools." && \
    curl -sSL http://download.redis.io/releases/redis-5.0.3.tar.gz | tar xz && cd redis-* && \
    make && cp src/redis-cli src/redis-benchmark /usr/local/bin && cd .. && rm -fr redis-*


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
