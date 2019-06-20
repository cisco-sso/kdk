#!/usr/bin/env bash
set -euo pipefail

function main() {
    export DEBIAN_FRONTEND="noninteractive"

    if [ "$#" -eq 0 ]; then
        # If this script was called with zero args, then we build for vagrant
        vagrant
    else
        # If this script was called with arguments, then we build for docker
        #   The following expects the first argument to be the function name
        $@
    fi
}

function vagrant() {
    layer_install_os_packages
    layer_install_python_based_utils_and_libs
    layer_install_apps_not_provided_by_os_packages
    layer_go_get_installs
    layer_build_apps_not_provided_by_os_packages
    mark_provisioned
}

function layer_install_os_packages() {
    # big items: gcc, python-dev

    echo "Install OS packages" && \
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
    curl -sSfL https://download.docker.com/linux/ubuntu/gpg | apt-key add - && \
       add-apt-repository \
         "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
         $(lsb_release -cs) \
         stable" && \
    export CLOUD_SDK_REPO="cloud-sdk-$(lsb_release -c -s)" && \
        echo "deb http://packages.cloud.google.com/apt $CLOUD_SDK_REPO main" | tee -a /etc/apt/sources.list.d/google-cloud-sdk.list && \
        curl -sSfL https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add - && \
    export POSTGRESQL_REPO="$(lsb_release -c -s)-pgdg" && \
        echo "deb https://apt.postgresql.org/pub/repos/apt $POSTGRESQL_REPO main" | tee -a /etc/apt/sources.list.d/pgdg.list && \
        curl -sSfL https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add - && \
    curl -sSfL https://deb.nodesource.com/setup_11.x | bash - && \
        curl -sSfL https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add - && \
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
}

function layer_install_python_based_utils_and_libs() {
    echo "Install python-based utils and libs" && \
    curl -sSfL https://bootstrap.pypa.io/get-pip.py | python2 && \
    pip install --no-cache-dir -U setuptools && \
    pip install \
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
    curl -sSfL https://bootstrap.pypa.io/get-pip.py | python3 && \
    pip3 install --no-cache-dir -U setuptools && \
    pip3 install \
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
}

function layer_install_apps_not_provided_by_os_packages() {
    echo "Install apps (with pinned version) that are not provided by the OS packages." && \
    echo "Install dep." && \
        curl -sSfLo dep https://github.com/golang/dep/releases/download/v0.5.1/dep-linux-amd64 && \
        chmod a+x dep && mv dep /usr/local/bin && \
    echo "Install direnv." && \
        curl -sSfLo direnv https://github.com/direnv/direnv/releases/download/v2.20.0/direnv.linux-amd64 && \
        chmod a+x direnv && mv direnv /usr/local/bin && \
    echo "Install drone-cli." && \
        curl -sSfL https://github.com/drone/drone-cli/releases/download/v1.1.0/drone_linux_amd64.tar.gz | tar xz && \
        chmod a+x drone && mv drone /usr/local/bin/drone-1.1.0 && \
        ln -s /usr/local/bin/drone-1.1.0 /usr/local/bin/drone && \
    echo "Install easy-rsa." && \
        curl -sSfL https://github.com/OpenVPN/easy-rsa/releases/download/v3.0.6/EasyRSA-unix-v3.0.6.tgz | tar xz && \
        chmod a+x EasyRSA-* && mv EasyRSA-* /usr/local/bin/easyrsa && \
    echo "Install go-task." && \
        curl -sSfL https://github.com/go-task/task/releases/download/v2.5.1/task_linux_amd64.tar.gz | tar -C /usr/local/bin -xz task && chmod a+x /usr/local/bin/task && \
    echo "Install gomplate." && \
        curl -sSfLo gomplate https://github.com/hairyhenderson/gomplate/releases/download/v3.4.0/gomplate_linux-amd64 && \
        chmod a+x gomplate && mv gomplate /usr/local/bin && \
    echo "Install golang." && \
        curl -sSfL https://dl.google.com/go/go1.12.4.linux-amd64.tar.gz | tar -C /usr/local -xz && \
        mkdir -p /go && chmod a+rw /go && \
    echo "Install goreleaser." && \
        curl -sSfLO https://github.com/goreleaser/goreleaser/releases/download/v0.106.0/goreleaser_Linux_x86_64.tar.gz && \
        tar -C /usr/local/bin -xzf goreleaser*.tar.gz goreleaser && rm goreleaser*.tar.gz && \
    echo "Install grpcurl." && \
        curl -sSfL https://github.com/fullstorydev/grpcurl/releases/download/v1.2.1/grpcurl_1.2.1_linux_x86_64.tar.gz | tar -C /usr/local/bin -xz grpcurl && chmod a+x /usr/local/bin/grpcurl && \
    echo "Install helm." && \
        curl -sSfL https://storage.googleapis.com/kubernetes-helm/helm-v2.11.0-linux-amd64.tar.gz | tar xz && \
          chmod a+x linux-amd64/helm && mv linux-amd64/helm /usr/local/bin/helm-2.11.0 && rm -fr linux-amd64 && \
        curl -sSfL https://storage.googleapis.com/kubernetes-helm/helm-v2.12.3-linux-amd64.tar.gz | tar xz && \
          chmod a+x linux-amd64/helm && mv linux-amd64/helm /usr/local/bin/helm-2.12.3 && rm -fr linux-amd64 && \
        curl -sSfL https://storage.googleapis.com/kubernetes-helm/helm-v2.13.1-linux-amd64.tar.gz | tar xz && \
          chmod a+x linux-amd64/helm && mv linux-amd64/helm /usr/local/bin/helm-2.13.1 && rm -fr linux-amd64 && \
        ln -s /usr/local/bin/helm-2.13.1 /usr/local/bin/helm && \
    echo "Install helmfile" && \
        curl -sSfLo helmfile https://github.com/roboll/helmfile/releases/download/v0.54.2/helmfile_linux_amd64 && \
        chmod a+x helmfile && mv helmfile /usr/local/bin && \
    echo "Install hugo." && \
        curl -sSfL https://github.com/gohugoio/hugo/releases/download/v0.55.4/hugo_0.55.4_Linux-64bit.tar.gz | tar xz && \
        chmod a+x hugo && mv hugo /usr/local/bin/hugo && \
    echo "Install jq." && \
        curl -sSfLo jq https://github.com/stedolan/jq/releases/download/jq-1.6/jq-linux64 && \
        chmod a+x jq && mv jq /usr/local/bin/ && \
    echo "Install kops." && \
        curl -sSfLo kops-1.10.1 https://github.com/kubernetes/kops/releases/download/1.10.1/kops-linux-amd64 && \
        chmod a+x kops-1.10.1 && mv kops-1.10.1 /usr/local/bin/ && \
        curl -sSfLo kops-1.11.1 https://github.com/kubernetes/kops/releases/download/1.11.1/kops-linux-amd64 && \
        chmod a+x kops-1.11.1 && mv kops-1.11.1 /usr/local/bin/ && \
        ln -s /usr/local/bin/kops-1.11.1 /usr/local/bin/kops && \
    echo "Install kubectl." && \
        curl -sSfLo /usr/local/bin/kubectl-1.14.1 https://storage.googleapis.com/kubernetes-release/release/v1.14.1/bin/linux/amd64/kubectl && \
        chmod a+x /usr/local/bin/kubectl-* && \
        ln -s /usr/local/bin/kubectl-1.14.1 /usr/local/bin/kubectl && \
    echo "Install kubetail." && \
        curl -sSfLo kubetail.zip https://github.com/johanhaleby/kubetail/archive/1.6.8.zip && \
        unzip -qq kubetail.zip && chmod a+x kubetail-1.6.8/kubetail && mv kubetail-1.6.8/kubetail /usr/local/bin && \
        rm -f kubetail.zip && \
    echo "Install mc." && \
        curl -sSfLo /usr/local/bin/mc https://dl.minio.io/client/mc/release/linux-amd64/archive/mc.RELEASE.2019-05-01T23-27-44Z && \
        chmod a+x /usr/local/bin/mc && \
    echo "Install minikube." && \
        curl -sSfLo minikube https://storage.googleapis.com/minikube/releases/v1.0.1/minikube-linux-amd64 && \
        chmod a+x minikube &&  mv minikube /usr/local/bin/ && \
    echo "Install terraform." && \
        curl -sSfLo terraform.zip https://releases.hashicorp.com/terraform/0.11.13/terraform_0.11.13_linux_amd64.zip && \
        unzip -qq terraform.zip && chmod a+x terraform && mv terraform /usr/local/bin && rm -f terraform.zip && \
    echo "Install testssl." && \
        curl -sSfL https://github.com/drwetter/testssl.sh/archive/v2.9.5-7.tar.gz | tar xz && \
        mv testssl* /usr/local/share/testssl && ln -s /usr/local/share/testssl/testssl.sh /usr/local/bin/testssl && chmod a+x /usr/local/bin/testssl && \
    echo "Install yadm." && \
        curl -sSfL https://github.com/TheLocehiliosan/yadm/archive/1.12.0.tar.gz | tar xz && \
        mv yadm* /usr/local/share/yadm && ln -s /usr/local/share/yadm/yadm /usr/local/bin/yadm && chmod a+x /usr/local/bin/yadm
}

function layer_go_get_installs() {
    export GOPATH=/go
    echo "go get installs" && \
      apt-get -y update && apt-get --no-install-recommends -y install git && \
      apt-get -y clean && apt-get -y autoremove && rm -rf /var/lib/apt/lists/*
    /usr/local/go/bin/go get github.com/spf13/cobra/cobra
    /usr/local/go/bin/go get github.com/kubernetes-incubator/cri-tools/cmd/crictl
    /usr/local/go/bin/go get golang.org/x/lint/golint
    /usr/local/go/bin/go get github.com/gpmgo/gopm
    /usr/local/go/bin/go get github.com/vmware/govmomi/govc
    /usr/local/go/bin/go get github.com/github/hub
    /usr/local/go/bin/go get github.com/cisco-sso/mh && ln -sf /go/bin/mh /go/bin/multihelm
    /usr/local/go/bin/go get github.com/mikefarah/yq
    rm -rf /root/.cache/go-build
}

function layer_build_apps_not_provided_by_os_packages() {
    echo "Install OS BUILD packages" && \
    apt-get -y update && apt-get --no-install-recommends -y install \
        autoconf \
        build-essential \
        libgnutls28-dev \
        libncurses5-dev \
        libz-dev && \
   apt-get -y clean && apt-get -y autoremove && rm -rf /var/lib/apt/lists/*

    echo "Clean out /usr/local" && \
    rm -rf /usr/local/*

    echo "Install git (needs to build first as a dependency)." && \
    curl -sSfL https://github.com/git/git/archive/v2.19.1.tar.gz | tar xz && cd git-* && \
    make configure && ./configure --prefix=/usr/local && make && make install && cd .. && rm -fr git-*

    echo "Install bats" && \
    curl -sSfL https://github.com/sstephenson/bats/archive/v0.4.0.tar.gz | tar xz && cd bats-* && \
    ./install.sh /usr/local && cd .. && rm -fr bats-*

    echo "Install emacs." && \
    curl -sSfL http://mirrors.ibiblio.org/gnu/ftp/gnu/emacs/emacs-26.1.tar.gz | tar xz && cd emacs-* && \
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

    echo "Install jsonnet" && \
    curl -sSfL https://github.com/google/jsonnet/archive/v0.12.1.tar.gz | tar xz && cd jsonnet-* && \
    make && chmod a+x jsonnet && mv jsonnet /usr/local/bin

    echo "Install pyenv with dependencies." && \
    curl -sSfLo pyenv-installer https://raw.githubusercontent.com/pyenv/pyenv-installer/master/bin/pyenv-installer && \
    chmod a+x pyenv-installer && mv pyenv-installer /usr/local/bin && \
    PYENV_ROOT=/usr/local/pyenv pyenv-installer && chmod -R a+rwx /usr/local/pyenv

    echo "Install vim." && \
    curl -sSfL https://github.com/vim/vim/archive/v8.1.0481.tar.gz | tar xz && cd vim-* && \
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

    echo "Install tmux." && \
    curl -sSfL https://github.com/libevent/libevent/releases/download/release-2.1.8-stable/libevent-2.1.8-stable.tar.gz | tar xz && cd libevent-* && \
    ./configure && make &&  make install && cd .. && rm -fr libevent-* && \
    curl -sSfL https://github.com/tmux/tmux/releases/download/2.8/tmux-2.8.tar.gz | tar xz && cd tmux-* && \
    ./configure --prefix=/usr/local && make && make install && cd .. && rm -fr tmux-*

    echo "Install zsh." && \
    curl -sSfL https://sourceforge.net/projects/zsh/files/zsh/5.6.2/zsh-5.6.2.tar.xz/download | tar Jx && cd zsh-* && \
    ./configure --with-tcsetpgrp --prefix=/usr/local && make && make install && echo "/usr/local/bin/zsh" >> /etc/shells && cd .. && rm -fr zsh-*

    echo "Install redis-cli tools." && \
    curl -sSfL http://download.redis.io/releases/redis-5.0.3.tar.gz | tar xz && cd redis-* && \
    make && cp src/redis-cli src/redis-benchmark /usr/local/bin && cd .. && rm -fr redis-*
}

function mark_provisioned() {
    sudo touch /var/lib/provisioned
}

#####################################################################
# Run the main program
main "$@"

## Report general success.
echo OK
