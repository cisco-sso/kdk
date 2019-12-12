#!/usr/bin/env bash
set -euo pipefail

function main() {
    export DEBIAN_FRONTEND="noninteractive"
    if [ "$#" -eq 0 ]; then
        # If this script was called with zero args, then we build for ubuntu bionic
        ubuntu_bionic
    else
        # If this script was called with arguments, then we build for docker
        #   The following expects the first argument to be the function name
        $@
    fi
}

function get_latest_github_release_version() {
    if [[ ! "$#" -eq 2 ]]; then
        # If this function was not called with 2 args. then complain and exitcalled with zero args, then we build for ubuntu bionic
        echo "Function must be called with 2 args (ORG and REPO)"
        exit 1
    else
        export ORG=$1
        export REPO=$2
        export VERSION=$(curl -sSfL https://api.github.com/repos/"${ORG}/${REPO}"/releases/latest 2>/dev/null | grep tag_name | cut -d '"' -f 4 | sed 's|v||g')
        if [[ -z "${VERSION}" ]]; then
            echo "Failed to get version"
            exit 1
        else
            echo "${VERSION}"
        fi
    fi
}

function vagrant() {
    exit_if_provisioned

    pushd /tmp
    # Remove this workaround after bento releases new hyperv box
    vagrant_disable_ssh_password_logins
    vagrant_upgrade_kernel_workaround_sshuttle_kernel_bug
    vagrant_bento_workaround_openssl_bug
    layer_install_os_packages
    layer_install_python_based_utils_and_libs
    layer_install_apps_not_provided_by_os_packages
    layer_go_get_installs
    layer_build_apps_not_provided_by_os_packages
    vagrant_fix_permissions
    mark_provisioned
    rm -rf /tmp/* && popd
}

function ubuntu_bionic() {
    exit_if_provisioned

    pushd /tmp
    layer_install_os_packages
    layer_install_python_based_utils_and_libs
    layer_install_apps_not_provided_by_os_packages
    layer_go_get_installs
    layer_build_apps_not_provided_by_os_packages
    mark_provisioned
    rm -rf /tmp/* && popd
}

function vagrant_disable_ssh_password_logins() {
    # Vagrant machines are ssh-able via user/pass vagrant/vagrant.
    #   Disable this, since some boxes may run with bridged networking by default
    sed -i 's@#PasswordAuthentication yes@PasswordAuthentication no@g' \
        /etc/ssh/sshd_config
}

function vagrant_upgrade_kernel_workaround_sshuttle_kernel_bug() {
    # https://github.com/sshuttle/sshuttle/issues/208
    echo "#### ${FUNCNAME[0]}"

    # Install a kernel upgrade helper
    #   Dkms will automatically recompile kmods upon kernel update
    apt-add-repository -y ppa:teejee2008/ppa
    apt-get update
    apt-get -y install dkms ukuu linux-headers-$(uname -r)

    # Update virtualbox guest additions.  This will rebuild kernel modules
    wget -q -O /tmp/additions.iso \
      http://download.virtualbox.org/virtualbox/6.0.4/VBoxGuestAdditions_6.0.4.iso
    mkdir -p /cdrom
    mount -o loop /tmp/additions.iso /cdrom
    /cdrom/VBoxLinuxAdditions.run || true # always errors
    umount /cdrom
    rm -rf /cdrom /tmp/additions.iso

    # Install a newer kernel
    #   dkms will kick in and rebuild modules for this new kernel
    ukuu --list
    ukuu --install v4.20.17
    ukuu --list-installed
}

function vagrant_bento_workaround_openssl_bug() {
    # https://github.com/chef/bento/issues/1201#issuecomment-503060115
    echo "#### ${FUNCNAME[0]}"
    DEBIAN_FRONTEND=noninteractive dpkg-reconfigure libc6
    DEBIAN_FRONTEND=noninteractive dpkg-reconfigure libssl1.1
    apt-get update
    DEBIAN_FRONTEND=noninteractive apt-get install -y libssl1.1
}

function layer_install_os_packages() {
    # big items: gcc, python-dev

    echo "#### ${FUNCNAME[0]}"
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
        python3-openssl && \
    apt-get -y clean && apt-get -y autoremove && rm -rf /var/lib/apt/lists/*
}

function layer_install_python_based_utils_and_libs() {
    echo "#### ${FUNCNAME[0]}"
    curl -sSfL https://bootstrap.pypa.io/get-pip.py | python3 && \
    pip3 install --no-cache-dir -U setuptools && \
    pip3 install \
         --no-cache-dir \
         'ansible==2.7.12' \
         'awscli==1.16.218' \
         'boto3==1.9.208' \
         'boto==2.49.0' \
         'docker-compose==1.24.1' \
         'Jinja2==2.10.1' \
         'jinja2-cli[yaml]==0.7.0' \
         'jsonschema' \
         'openshift==0.9.0' \
         'peru==1.2.0' \
         'pipenv==2018.11.26' \
         'python-neutronclient==6.12.0' \
         'python-octaviaclient==1.9.0' \
         'python-openstackclient==3.19.0' \
         'pyvmomi==6.7.1.2018.12' \
         'sh==1.12.14' \
         'sshuttle==0.78.5' \
         'structlog==19.1.0' \
         'urllib3==1.22' \
         'virtualenv==16.7.2' \
         'yamllint' \
         'yapf' \
         'yq==2.7.2' && \
    rm -rf /root/.cache/pip
}

function layer_install_apps_not_provided_by_os_packages() {
    echo "#### ${FUNCNAME[0]}"
    echo "Install apps (with pinned version) that are not provided by the OS packages." && \
    echo "Install amtool." && \
        export ORG="prometheus" && export REPO="alertmanager" && export VERSION=$(get_latest_github_release_version "${ORG}" "${REPO}") && export ARTIFACT="amtool" && \
        curl -sSfL https://github.com/"${ORG}/${REPO}"/releases/download/v"${VERSION}"/"${REPO}"-"${VERSION}".linux-amd64.tar.gz | tar xz && \
        mv "${REPO}"*/"${ARTIFACT}" /usr/local/bin/"${ARTIFACT}" && rm -rf "${REPO}"* && \
    echo "Install dep." && \
        export ORG="golang" && export REPO="dep" && export VERSION=$(get_latest_github_release_version "${ORG}" "${REPO}") && export ARTIFACT="${REPO}"  && \
        curl -sSfLo "${ARTIFACT}" https://github.com/"${ORG}"/"${REPO}"/releases/download/v"${VERSION}"/"${ARTIFACT}"-linux-amd64 && \
        chmod a+x "${ARTIFACT}" && mv "${ARTIFACT}" /usr/local/bin && \
    echo "Install direnv." && \
        export ORG="direnv" && export REPO="direnv" && export VERSION=$(get_latest_github_release_version "${ORG}" "${REPO}") && export ARTIFACT="${REPO}" && \
        curl -sSfLo "${ARTIFACT}" https://github.com/"${ORG}"/"${REPO}"/releases/download/v"${VERSION}"/"${REPO}".linux-amd64 && \
        chmod a+x "${ARTIFACT}" && mv "${ARTIFACT}" /usr/local/bin && \
    echo "Install drone-cli." && \
        export ORG="drone" && export REPO="drone-cli" && export VERSION=$(get_latest_github_release_version "${ORG}" "${REPO}") && export ARTIFACT="${ORG}" && \
        curl -sSfL https://github.com/"${ORG}"/"${REPO}"/releases/download/v"${VERSION}"/"${ARTIFACT}"_linux_amd64.tar.gz | tar xz && \
        chmod a+x "${ARTIFACT}" && mv "${ARTIFACT}" /usr/local/bin/"${ARTIFACT}"-"${VERSION}" && \
        ln -sf /usr/local/bin/"${ARTIFACT}"-"${VERSION}" /usr/local/bin/"${ARTIFACT}" && \
    echo "Install easy-rsa." && \
        export ORG="OpenVPN" && export REPO="easy-rsa" && export VERSION=$(get_latest_github_release_version "${ORG}" "${REPO}") && export ARTIFACT="EasyRSA" && \
        curl -sSfL https://github.com/"${ORG}"/"${REPO}"/releases/download/v"${VERSION}"/"${ARTIFACT}"-unix-v"${VERSION}".tgz | tar xz && \
        chmod a+x "${ARTIFACT}"-v"${VERSION}"/$(echo "${ARTIFACT}" | tr '[:upper:]' '[:lower:]') && mv "${ARTIFACT}"-v"${VERSION}" /usr/local/share/"${ARTIFACT}"-"${VERSION}" && \
        ln -s /usr/local/share/"${ARTIFACT}"-"${VERSION}"/$(echo "${ARTIFACT}" | tr '[:upper:]' '[:lower:]') /usr/local/bin/ && \
    echo "Install etcdctl." && \
        export ORG="etcd-io" && export REPO="etcd" && export VERSION=$(get_latest_github_release_version "${ORG}" "${REPO}") && export ARTIFACT="etcdctl" && \
        curl -sSfL https://github.com/"${ORG}"/"${REPO}"/releases/download/v"${VERSION}"/"${REPO}"-v"${VERSION}"-linux-amd64.tar.gz | tar xz && \
        mv "${REPO}"*/${ARTIFACT} /usr/local/bin/${ARTIFACT} && rm -rf "${REPO}"* && \
    echo "Install go-task." && \
        export ORG="go-task" && export REPO="task" && export VERSION=$(get_latest_github_release_version "${ORG}" "${REPO}") && export ARTIFACT="${REPO}" && \
        curl -sSfL https://github.com/"${ORG}"/"${REPO}"/releases/download/v"${VERSION}"/"${ARTIFACT}"_linux_amd64.tar.gz | tar -C /usr/local/bin -xz "${ARTIFACT}" && chmod a+x /usr/local/bin/"${ARTIFACT}" && \
    echo "Install gomplate." && \
        export ORG="hairyhenderson" && export REPO="gomplate" && export VERSION=$(get_latest_github_release_version "${ORG}" "${REPO}") && export ARTIFACT="${REPO}" && \
        curl -sSfLo "${ARTIFACT}" https://github.com/"${ORG}"/"${REPO}"/releases/download/v"${VERSION}"/"${ARTIFACT}"_linux-amd64 && \
        chmod a+x "${ARTIFACT}" && mv "${ARTIFACT}" /usr/local/bin && \
    echo "Install golang." && \
        curl -sSfL https://dl.google.com/go/go1.13.1.linux-amd64.tar.gz | tar -C /usr/local -xz && \
        mkdir -p /go && chmod a+rw /go && \
    echo "Install goreleaser." && \
        curl -sSfLO https://github.com/goreleaser/goreleaser/releases/download/v0.113.1/goreleaser_Linux_x86_64.tar.gz && \
        tar -C /usr/local/bin -xzf goreleaser*.tar.gz goreleaser && rm goreleaser*.tar.gz && \
    echo "Install grpcurl." && \
        curl -sSfL https://github.com/fullstorydev/grpcurl/releases/download/v1.3.1/grpcurl_1.3.1_linux_x86_64.tar.gz | tar -C /usr/local/bin -xz grpcurl && chmod a+x /usr/local/bin/grpcurl && \
    echo "Install helm." && \
        curl -sSfL https://storage.googleapis.com/kubernetes-helm/helm-v2.14.2-linux-amd64.tar.gz | tar xz && \
          chmod a+x linux-amd64/helm && mv linux-amd64/helm /usr/local/bin/helm-2.14.2 && rm -fr linux-amd64 && \
        ln -sf /usr/local/bin/helm-2.14.2 /usr/local/bin/helm && \
    echo "Install helmfile" && \
        curl -sSfLo helmfile https://github.com/roboll/helmfile/releases/download/v0.80.2/helmfile_linux_amd64 && \
        chmod a+x helmfile && mv helmfile /usr/local/bin && \
    echo "Install hugo." && \
        curl -sSfL https://github.com/gohugoio/hugo/releases/download/v0.56.3/hugo_0.56.3_Linux-64bit.tar.gz | tar xz && \
        chmod a+x hugo && mv hugo /usr/local/bin/hugo && \
    echo "Install jq." && \
        curl -sSfLo jq https://github.com/stedolan/jq/releases/download/jq-1.6/jq-linux64 && \
        chmod a+x jq && mv jq /usr/local/bin/ && \
    echo "Install kops." && \
        curl -sSfLo kops-1.11.1 https://github.com/kubernetes/kops/releases/download/1.11.1/kops-linux-amd64 && \
        chmod a+x kops-1.11.1 && mv kops-1.11.1 /usr/local/bin/ && \
        curl -sSfLo kops-1.12.2 https://github.com/kubernetes/kops/releases/download/1.12.2/kops-linux-amd64 && \
        chmod a+x kops-1.12.2 && mv kops-1.12.2 /usr/local/bin/ && \
        ln -sf /usr/local/bin/kops-1.12.2 /usr/local/bin/kops && \
    echo "Install kubectl." && \
        curl -sSfLo /usr/local/bin/kubectl-1.15.4 https://storage.googleapis.com/kubernetes-release/release/v1.15.4/bin/linux/amd64/kubectl && \
        chmod a+x /usr/local/bin/kubectl-* && \
        curl -sSfLo /usr/local/bin/kubectl-1.16.0 https://storage.googleapis.com/kubernetes-release/release/v1.16.0/bin/linux/amd64/kubectl && \
        chmod a+x /usr/local/bin/kubectl-* && \
        ln -sf /usr/local/bin/kubectl-1.15.4 /usr/local/bin/kubectl && \
    echo "Install kubetail." && \
        curl -sSfLo kubetail.zip https://github.com/johanhaleby/kubetail/archive/1.6.8.zip && \
        unzip -qq kubetail.zip && chmod a+x kubetail-1.6.8/kubetail && mv kubetail-1.6.8/kubetail /usr/local/bin && \
        rm -rf kubetail* && \
    echo "Install mc." && \
        curl -sSfLo /usr/local/bin/mc https://dl.minio.io/client/mc/release/linux-amd64/archive/mc.RELEASE.2019-05-01T23-27-44Z && \
        chmod a+x /usr/local/bin/mc && \
    echo "Install minikube." && \
        curl -sSfLo minikube https://storage.googleapis.com/minikube/releases/v1.0.1/minikube-linux-amd64 && \
        chmod a+x minikube &&  mv minikube /usr/local/bin/ && \
    echo "Install terraform." && \
        curl -sSfLo terraform.zip https://releases.hashicorp.com/terraform/0.11.14/terraform_0.11.14_linux_amd64.zip && \
        unzip -qq terraform.zip && chmod a+x terraform && mv terraform /usr/local/bin/terraform-0.11.14 && rm -f terraform.zip && \
        curl -sSfLo terraform.zip https://releases.hashicorp.com/terraform/0.12.6/terraform_0.12.6_linux_amd64.zip && \
        unzip -qq terraform.zip && chmod a+x terraform && mv terraform /usr/local/bin/terraform-0.12.6 && rm -f terraform.zip && \
        ln -sf /usr/local/bin/terraform-0.12.6 /usr/local/bin/terraform && \
    echo "Install testssl." && \
        curl -sSfL https://github.com/drwetter/testssl.sh/archive/v2.9.5-7.tar.gz | tar xz && \
        mv testssl* /usr/local/share/testssl && ln -sf /usr/local/share/testssl/testssl.sh /usr/local/bin/testssl && chmod a+x /usr/local/bin/testssl && \
    echo "Install yadm." && \
        curl -sSfL https://github.com/TheLocehiliosan/yadm/archive/1.12.0.tar.gz | tar xz && \
        mv yadm* /usr/local/share/yadm && ln -sf /usr/local/share/yadm/yadm /usr/local/bin/yadm && chmod a+x /usr/local/bin/yadm
}

function layer_go_get_installs() {
    echo "#### ${FUNCNAME[0]}"
    export GOPATH=/go && \
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
    rm -rf /go/src
}

function layer_build_apps_not_provided_by_os_packages() {
    echo "#### ${FUNCNAME[0]}"
    apt-get -y update && apt-get --no-install-recommends -y install \
        autoconf \
        build-essential \
        libgnutls28-dev \
        libncurses5-dev \
        libz-dev && \
    apt-get -y clean && apt-get -y autoremove && rm -rf /var/lib/apt/lists/*

    echo "Install git (needs to build first as a dependency)." && \
    curl -sSfL https://github.com/git/git/archive/v2.22.0.tar.gz | tar xz && cd git-* && \
    make configure && ./configure --prefix=/usr/local && make && make install && cd .. && rm -fr git-*

    echo "Install bats" && \
    curl -sSfL https://github.com/sstephenson/bats/archive/v0.4.0.tar.gz | tar xz && cd bats-* && \
    ./install.sh /usr/local && cd .. && rm -fr bats-*

    echo "Install emacs." && \
    curl -sSfL http://mirrors.ibiblio.org/gnu/ftp/gnu/emacs/emacs-26.2.tar.gz | tar xz && cd emacs-* && \
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
    curl -sSfL https://github.com/google/jsonnet/archive/v0.13.0.tar.gz | tar xz && cd jsonnet-* && \
    make && chmod a+x jsonnet jsonnetfmt && mv jsonnet /usr/local/bin && mv jsonnetfmt /usr/local/bin \
    && cd .. && rm -fr jsonnet-*

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
    curl -sSfL https://github.com/tmux/tmux/releases/download/2.9a/tmux-2.9a.tar.gz | tar xz && cd tmux-* && \
    ./configure --prefix=/usr/local && make && make install && cd .. && rm -fr tmux-*

    echo "Install zsh." && \
    curl -sSfL https://sourceforge.net/projects/zsh/files/zsh/5.6.2/zsh-5.6.2.tar.xz/download | tar Jx && cd zsh-* && \
    ./configure --with-tcsetpgrp --prefix=/usr/local && make && make install && echo "/usr/local/bin/zsh" >> /etc/shells && cd .. && rm -fr zsh-*

    echo "Install redis-cli tools." && \
    curl -sSfL http://download.redis.io/releases/redis-5.0.3.tar.gz | tar xz && cd redis-* && \
    make && cp src/redis-cli src/redis-benchmark /usr/local/bin && cd .. && rm -fr redis-*
}

function exit_if_provisioned() {
    echo "#### ${FUNCNAME[0]}"
    if [ -f /var/lib/provisioned ]; then
        echo "Already provisioned since exists: /var/lib/provisioned"
        exit 0
    fi
}

function vagrant_fix_permissions() {
    chown -R vagrant:vagrant /home/vagrant/.ssh
    chown -R vagrant:vagrant /home/vagrant/.cache
    chown -R vagrant:vagrant /home/vagrant/.wget-hsts
}

function mark_provisioned() {
    echo "#### ${FUNCNAME[0]}"
    sudo touch /var/lib/provisioned
}

#####################################################################
# Run the main program
main "$@"

## Report general success.
echo OK
