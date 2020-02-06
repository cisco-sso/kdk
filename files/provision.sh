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
    apt-get -y -qq install dkms ukuu linux-headers-$(uname -r)

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
    apt-get -y -qq update && apt-get --no-install-recommends -y -qq install \
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
        iputils-ping \
        less \
        libcurl4-openssl-dev \
        locales \
        lsof \
        make \
        man \
        moreutils \
        nmap \
        ntp \
        ntpdate \
        openconnect \
        openjdk-11-jdk-headless \
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
        xauth \
        zip && \
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
    curl -sSfL https://deb.nodesource.com/setup_13.x | bash - && \
        curl -sSfL https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add - && \
        echo "deb https://dl.yarnpkg.com/debian/ stable main" | tee -a /etc/apt/sources.list.d/yarn.list && \
    apt-get update -y -qq && apt-get --no-install-recommends -y -qq install \
        docker-ce \
        google-cloud-sdk \
        nodejs \
        postgresql-client-10 \
        yarn && \
    # Add deps to enable pyenv-driven on-demand python installs on KDK \
    apt-get --no-install-recommends -y -qq install \
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
    apt-get -y -qq clean && apt-get -y -qq autoremove && rm -rf /var/lib/apt/lists/*
}

function layer_install_python_based_utils_and_libs() {
    echo "#### ${FUNCNAME[0]}"
    curl -sSfL https://bootstrap.pypa.io/get-pip.py | python3 && \
    pip3 install -qq --no-cache-dir -U setuptools && \
    pip3 install -qq \
         --no-cache-dir \
         'ansible==2.9.2' \
         'awscli==1.16.302' \
         'boto3==1.10.38' \
         'boto==2.49.0' \
         'click==7.0' \
         'docker-compose==1.25.0' \
         'ipython==7.10.1' \
         'ipdb==0.12.3' \
         'Jinja2==2.10.3' \
         'jinja2-cli[yaml]==0.7.0' \
         'jsonschema==3.2.0' \
         'openshift==0.10.0' \
         'peru==1.2.0' \
         'pipenv==2018.11.26' \
         'pytest==5.3.1' \
         'python-neutronclient==6.14.0' \
         'python-octaviaclient==1.11.0' \
         'python-openstackclient==3.19.0' \
         'pyvmomi==6.7.3' \
         'pyyaml==5.2' \
         'requests==2.22.0' \
         'sh==1.12.14' \
         'sshuttle==0.78.5' \
         'structlog==19.2.0' \
         'urllib3==1.25.7' \
         'virtualenv==16.7.8' \
         'yamllint==1.19.0' \
         'yapf==0.29.0' \
         'yq==2.9.2' && \
    rm -rf /root/.cache/pip
}

function layer_install_apps_not_provided_by_os_packages() {
    echo "#### ${FUNCNAME[0]}"
    echo "Install apps (with pinned version) that are not provided by the OS packages." && \
    echo "Install amtool." && \
        export ORG="prometheus" && export REPO="alertmanager" && export VERSION=$(get_latest_github_release_version "${ORG}" "${REPO}") && export ARTIFACT="amtool" && \
        curl -sSfL https://github.com/"${ORG}/${REPO}"/releases/download/v"${VERSION}"/"${REPO}"-"${VERSION}".linux-amd64.tar.gz | tar xz && \
        mv "${REPO}"*/"${ARTIFACT}" /usr/local/bin/"${ARTIFACT}" && rm -rf "${REPO}"* && \
    echo "Install consul cli." && \
        export ORG="hashicorp" && export REPO="consul" && export VERSION="1.6.2" && export ARTIFACT="${REPO}" && \
        curl -sSfLo "${ARTIFACT}".zip https://releases."${ORG}".com/"${ARTIFACT}"/"${VERSION}"/"${ARTIFACT}"_"${VERSION}"_linux_amd64.zip && \
        unzip -qq "${ARTIFACT}".zip && chmod a+x "${ARTIFACT}" && mv "${ARTIFACT}" /usr/local/bin && \
        rm -rf "${ARTIFACT}"* && \
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
    echo "Install fluxctl." && \
        export ORG="fluxcd" && export REPO="flux" && export VERSION=$(get_latest_github_release_version "${ORG}" "${REPO}") && export ARTIFACT="fluxctl" && \
        curl -sSfLo /usr/local/bin/"${ARTIFACT}" https://github.com/"${ORG}"/"${REPO}"/releases/download/"${VERSION}"/"${ARTIFACT}"_linux_amd64 && \
        chmod a+x /usr/local/bin/"${ARTIFACT}" && \
    echo "Install go-task." && \
        export ORG="go-task" && export REPO="task" && export VERSION=$(get_latest_github_release_version "${ORG}" "${REPO}") && export ARTIFACT="${REPO}" && \
        curl -sSfL https://github.com/"${ORG}"/"${REPO}"/releases/download/v"${VERSION}"/"${ARTIFACT}"_linux_amd64.tar.gz | tar -C /usr/local/bin -xz "${ARTIFACT}" && chmod a+x /usr/local/bin/"${ARTIFACT}" && \
    echo "Install gomplate." && \
        export ORG="hairyhenderson" && export REPO="gomplate" && export VERSION=$(get_latest_github_release_version "${ORG}" "${REPO}") && export ARTIFACT="${REPO}" && \
        curl -sSfLo "${ARTIFACT}" https://github.com/"${ORG}"/"${REPO}"/releases/download/v"${VERSION}"/"${ARTIFACT}"_linux-amd64 && \
        chmod a+x "${ARTIFACT}" && mv "${ARTIFACT}" /usr/local/bin && \
    echo "Install golang." && \
        export ORG="golang" && export REPO="go" && export VERSION="1.13.5" && export ARTIFACT="${REPO}" && \
        curl -sSfL https://dl.google.com/"${REPO}"/"${ARTIFACT}""${VERSION}".linux-amd64.tar.gz | tar -C /usr/local -xz && \
        mkdir -p /"${ARTIFACT}" && chmod a+rw /"${ARTIFACT}" && \
    echo "Install goreleaser." && \
        export ORG="goreleaser" && export REPO="goreleaser" && export VERSION=$(get_latest_github_release_version "${ORG}" "${REPO}") && export ARTIFACT="${REPO}" && \
        curl -sSfLO https://github.com/"${ORG}"/"${REPO}"/releases/download/v"${VERSION}"/"${ARTIFACT}"_Linux_x86_64.tar.gz && \
        tar -C /usr/local/bin -xzf "${ARTIFACT}"*.tar.gz "${ARTIFACT}" && rm "${ARTIFACT}"*.tar.gz && \
    echo "Install gradle." && \
        export VERSION="6.0.1" && ARTIFACT="gradle" && \
        curl -sSfLo "${ARTIFACT}".zip https://services."${ARTIFACT}".org/distributions/"${ARTIFACT}"-"${VERSION}"-bin.zip && \
        unzip -d /usr/local/share -qq "${ARTIFACT}".zip && \
        ln -sf /usr/local/share/"${ARTIFACT}"-"${VERSION}"/bin/gradle /usr/local/bin/"${ARTIFACT}" && \
        rm -rf "${ARTIFACT}"* && \
    echo "Install grpcurl." && \
        export ORG="fullstorydev" && export REPO="grpcurl" && export VERSION=$(get_latest_github_release_version "${ORG}" "${REPO}") && export ARTIFACT="${REPO}" && \
        curl -sSfL https://github.com/"${ORG}"/"${REPO}"/releases/download/v"${VERSION}"/"${ARTIFACT}"_"${VERSION}"_linux_x86_64.tar.gz | tar -C /usr/local/bin -xz "${ARTIFACT}" && chmod a+x /usr/local/bin/"${ARTIFACT}" && \
    echo "Install helm." && \
        export ORG="helm" && export REPO="helm" && export ARTIFACT="${REPO}" && \
        export VERSION="3.0.1" && curl -sSfL https://get.helm.sh/"${ARTIFACT}"-v"${VERSION}"-linux-amd64.tar.gz | tar xz && \
            chmod a+x linux-amd64/"${ARTIFACT}" && mv linux-amd64/"${ARTIFACT}" /usr/local/bin/"${ARTIFACT}"-"${VERSION}" && rm -fr linux-amd64 &&
            ln -sf /usr/local/bin/"${ARTIFACT}"-"${VERSION}" /usr/local/bin/"${ARTIFACT}"3 && \
            ln -sf /usr/local/bin/"${ARTIFACT}"-"${VERSION}" /usr/local/bin/"${ARTIFACT}" && \
        export VERSION="2.16.1" && curl -sSfL https://get.helm.sh/"${ARTIFACT}"-v"${VERSION}"-linux-amd64.tar.gz | tar xz && \
            chmod a+x linux-amd64/"${ARTIFACT}" && mv linux-amd64/"${ARTIFACT}" /usr/local/bin/"${ARTIFACT}"-"${VERSION}" && rm -fr linux-amd64 &&
            ln -sf /usr/local/bin/"${ARTIFACT}"-"${VERSION}" /usr/local/bin/"${ARTIFACT}"2 && \
    echo "Install helmfile" && \
        export ORG="roboll" && export REPO="helmfile" && export VERSION=$(get_latest_github_release_version "${ORG}" "${REPO}") && export ARTIFACT="${REPO}" && \
        curl -sSfLo "${ARTIFACT}" https://github.com/"${ORG}"/"${REPO}"/releases/download/v"${VERSION}"/"${ARTIFACT}"_linux_amd64 && \
        chmod a+x "${ARTIFACT}" && mv "${ARTIFACT}" /usr/local/bin && \
    echo "Install hugo." && \
        export ORG="gohugoio" && export REPO="hugo" && export VERSION=$(get_latest_github_release_version "${ORG}" "${REPO}") && export ARTIFACT="${REPO}" && \
        curl -sSfL https://github.com/"${ORG}"/"${REPO}"/releases/download/v"${VERSION}"/"${ARTIFACT}"_"${VERSION}"_Linux-64bit.tar.gz | tar xz && \
        chmod a+x "${ARTIFACT}" && mv "${ARTIFACT}" /usr/local/bin/"${ARTIFACT}" && \
    echo "Install jq." && \
        export ORG="stedolan" && export REPO="jq" && export VERSION=$(get_latest_github_release_version "${ORG}" "${REPO}") && export ARTIFACT="${REPO}" && \
        curl -sSfLo "${ARTIFACT}" https://github.com/"${ORG}"/"${REPO}"/releases/download/"${VERSION}"/"${ARTIFACT}"-linux64 && \
        chmod a+x "${ARTIFACT}" && mv "${ARTIFACT}" /usr/local/bin/ && \
    echo "Install jsonnet." && \
        export ORG="google" && export REPO="jsonnet" && export VERSION=$(get_latest_github_release_version "${ORG}" "${REPO}") && export ARTIFACT="jsonnet" && \
        curl -sSfL https://github.com/"${ORG}"/"${REPO}"/releases/download/v"${VERSION}"/"${ARTIFACT}"-bin-v"${VERSION}"-linux.tar.gz | tar -C /usr/local/bin -xz && \
    echo "Install kops." && \
        export ORG="kubernetes" && export REPO="kops" && export VERSION=$(get_latest_github_release_version "${ORG}" "${REPO}") && export ARTIFACT="${REPO}" && \
        curl -sSfLo /usr/local/bin/"${ARTIFACT}"-"${VERSION}" https://github.com/"${ORG}"/"${REPO}"/releases/download/v"${VERSION}"/"${ARTIFACT}"-linux-amd64 && \
        chmod a+x /usr/local/bin/"${ARTIFACT}"-"${VERSION}" && \
        ln -sf /usr/local/bin/"${ARTIFACT}"-"${VERSION}" /usr/local/bin/"${ARTIFACT}" && \
    echo "Install kubectl." && \
        export ORG="kubernetes" && export REPO="kubernetes" && export VERSION=$(get_latest_github_release_version "${ORG}" "${REPO}") && export ARTIFACT="kubectl" && \
        curl -sSfLo /usr/local/bin/"${ARTIFACT}"-"${VERSION}" https://storage.googleapis.com/kubernetes-release/release/v"${VERSION}"/bin/linux/amd64/"${ARTIFACT}" && \
        chmod a+x /usr/local/bin/"${ARTIFACT}"-"${VERSION}" && \
        ln -sf /usr/local/bin/"${ARTIFACT}"-"${VERSION}" /usr/local/bin/"${ARTIFACT}" && \
    echo "Install kubetail." && \
        export ORG="johanhaleby" && export REPO="kubetail" && export VERSION=$(get_latest_github_release_version "${ORG}" "${REPO}") && export ARTIFACT="${REPO}" && \
        curl -sSfLo "${ARTIFACT}".zip https://github.com/"${ORG}"/"${REPO}"/archive/"${VERSION}".zip && \
        unzip -qq "${ARTIFACT}".zip && chmod a+x "${ARTIFACT}"-"${VERSION}"/"${ARTIFACT}" && mv "${ARTIFACT}"-"${VERSION}"/"${ARTIFACT}" /usr/local/bin && \
        rm -rf "${ARTIFACT}"* && \
    echo "Install maven." && \
        export VERSION="3.6.3"
        curl -sSfL http://apache.mirrors.tds.net/maven/maven-3/"${VERSION}"/binaries/apache-maven-"${VERSION}"-bin.tar.gz | tar -C /usr/local/share -xz && \
        ln -sf /usr/local/share/apache-maven-"${VERSION}"/bin/mvn /usr/local/bin/mvn && \
    echo "Install mc." && \
        export ORG="minio" && export REPO="minio" && export ARTIFACT="mc" && \
        curl -sSfLo /usr/local/bin/"${ARTIFACT}" https://dl."${ORG}".io/client/"${ARTIFACT}"/release/linux-amd64/"${ARTIFACT}" && \
        chmod a+x /usr/local/bin/"${ARTIFACT}" && \
    echo "Install minikube." && \
        export ORG="kubernetes" && export REPO="minikube" && export VERSION=$(get_latest_github_release_version "${ORG}" "${REPO}") && export ARTIFACT="${REPO}" && \
        curl -sSfLo /usr/local/bin/"${ARTIFACT}"-"${VERSION}" https://github.com/"${ORG}"/"${REPO}"/releases/download/v"${VERSION}"/"${ARTIFACT}"-linux-amd64  && \
        chmod a+x /usr/local/bin/"${ARTIFACT}"-"${VERSION}" && \
        ln -sf /usr/local/bin/"${ARTIFACT}"-"${VERSION}" /usr/local/bin/"${ARTIFACT}" && \
    echo "Install neovim." && \
        export ORG="neovim" && export REPO="neovim" && export VERSION=$(get_latest_github_release_version "${ORG}" "${REPO}") && export ARTIFACT="nvim" && \
        mkdir -p /usr/local/share/"${ARTIFACT}" && curl -sSfL https://github.com/"${ORG}"/"${REPO}"/releases/download/v"${VERSION}"/"${ARTIFACT}"-linux64.tar.gz | tar -C /usr/local/share/"${ARTIFACT}"  --strip-components=1 -xz && \
        ln -sf /usr/local/share/"${ARTIFACT}"/bin/"${ARTIFACT}" /usr/local/bin/"${ARTIFACT}" && \
    echo "Install sops." && \
        export ORG="mozilla" && export REPO="sops" && export VERSION=$(get_latest_github_release_version "${ORG}" "${REPO}") && export ARTIFACT="${REPO}" && \
        curl -sSfLo /usr/local/bin/"${ARTIFACT}" https://github.com/"${ORG}"/"${REPO}"/releases/download/v"${VERSION}"/"${ARTIFACT}"-v"${VERSION}".linux && \
        chmod a+x /usr/local/bin/"${ARTIFACT}" && \
    echo "Install terraform." && \
        export ORG="hashicorp" && export REPO="terraform" && export VERSION=$(get_latest_github_release_version "${ORG}" "${REPO}") && export ARTIFACT="${REPO}" && \
        curl -sSfLo "${ARTIFACT}".zip https://releases."${ORG}".com/"${REPO}"/"${VERSION}"/"${ARTIFACT}"_"${VERSION}"_linux_amd64.zip && \
        unzip -qq "${ARTIFACT}".zip && chmod a+x "${ARTIFACT}" && mv "${ARTIFACT}" /usr/local/bin/"${ARTIFACT}"-"${VERSION}" && rm -f "${ARTIFACT}".zip && \
        ln -sf /usr/local/bin/"${ARTIFACT}"-"${VERSION}" /usr/local/bin/"${ARTIFACT}" && \
    echo "Install testssl." && \
        export ORG="drwetter" && export REPO="testssl.sh" && export VERSION="2.9.5-8" && export ARTIFACT="testssl" && \
        curl -sSfL https://github.com/"${ORG}"/"${REPO}"/archive/v"${VERSION}".tar.gz | tar xz && \
        mv "${ARTIFACT}"* /usr/local/share/"${ARTIFACT}" && ln -sf /usr/local/share/"${ARTIFACT}"/"${REPO}" /usr/local/bin/"${ARTIFACT}" && chmod a+x /usr/local/bin/"${ARTIFACT}" && \
    echo "Install vault cli." && \
        export ORG="hashicorp" && export REPO="vault" && export VERSION="1.3.0" && export ARTIFACT="${REPO}" && \
        curl -sSfLo "${ARTIFACT}".zip https://releases."${ORG}".com/"${ARTIFACT}"/"${VERSION}"/"${ARTIFACT}"_"${VERSION}"_linux_amd64.zip && \
        unzip -qq "${ARTIFACT}".zip && chmod a+x "${ARTIFACT}" && mv "${ARTIFACT}" /usr/local/bin && \
        rm -rf "${ARTIFACT}"* && \
    echo "Install yadm." && \
        export ORG="thelocehiliosan" && export REPO="yadm" && export VERSION="2.2.0" && export ARTIFACT="${REPO}" && \
        curl -sSfLo /usr/local/bin/"${ARTIFACT}" https://github.com/"${ORG}"/"${ARTIFACT}"/raw/"${VERSION}"/"${ARTIFACT}" && chmod a+x /usr/local/bin/"${ARTIFACT}"
}

function layer_go_get_installs() {
    echo "#### ${FUNCNAME[0]}"
    export GOPATH=/go && \
    echo "go get installs" && \
      apt-get -y -qq update && apt-get --no-install-recommends -y -qq install git && \
      apt-get -y -qq clean && apt-get -y -qq autoremove && rm -rf /var/lib/apt/lists/*
    /usr/local/go/bin/go get github.com/cloudflare/cfssl/cmd/cfssl
    /usr/local/go/bin/go get github.com/cloudflare/cfssl/cmd/cfssljson
    /usr/local/go/bin/go get github.com/spf13/cobra/cobra
    /usr/local/go/bin/go get github.com/kubernetes-incubator/cri-tools/cmd/crictl
    /usr/local/go/bin/go get golang.org/x/lint/golint
    /usr/local/go/bin/go get github.com/gpmgo/gopm
    /usr/local/go/bin/go get github.com/vmware/govmomi/govc
    /usr/local/go/bin/go get github.com/github/hub
    git clone https://github.com/cisco-sso/mh.git /tmp/mh && cd /tmp/mh && \
        /usr/local/go/bin/go mod init github.com/cisco-sso/mh && \
        /usr/local/go/bin/go build -o /go/bin/mh && \
        ln -sf /go/bin/mh /go/bin/multihelm
    GO111MODULE=on /usr/local/go/bin/go get github.com/mikefarah/yq/v2
    rm -rf /root/.cache/go-build
    rm -rf /go/src
}

function layer_build_apps_not_provided_by_os_packages() {
    echo "#### ${FUNCNAME[0]}"
    apt-get -y -qq update && apt-get --no-install-recommends -y -qq install \
        autoconf \
        build-essential \
        gcc \
        libgnutls28-dev \
        libncurses5-dev \
        libz-dev \
        texinfo \
        yodl && \
    apt-get -y -qq clean && apt-get -y -qq autoremove && rm -rf /var/lib/apt/lists/*

    echo "Install git (needs to build first as a dependency)." && \
        export ORG="git" && export REPO="git" && export VERSION="2.24.1" && export ARTIFACT="${REPO}" && \
        curl -sSfL https://github.com/"${ORG}"/"${REPO}"/archive/v"${VERSION}".tar.gz | tar xz && cd git-* && \
        make configure && ./configure --prefix=/usr/local && make && make install && cd .. && rm -fr git-*

    echo "Install bats" && \
    curl -sSfL https://github.com/sstephenson/bats/archive/v0.4.0.tar.gz | tar xz && cd bats-* && \
    ./install.sh /usr/local && cd .. && rm -fr bats-*

    echo "Install emacs." && \
        curl -sSfL http://mirrors.ibiblio.org/gnu/ftp/gnu/emacs/emacs-26.3.tar.gz | tar xz && cd emacs-* && \
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

    echo "Install pyenv with dependencies." && \
        curl -sSfLo pyenv-installer https://raw.githubusercontent.com/pyenv/pyenv-installer/master/bin/pyenv-installer && \
        chmod a+x pyenv-installer && mv pyenv-installer /usr/local/bin && \
        PYENV_ROOT=/usr/local/pyenv pyenv-installer && chmod -R a+rwx /usr/local/pyenv

    echo "Install vim." && \
        export ORG="vim" && export REPO="vim" && export VERSION="8.2.0" && export ARTIFACT="${REPO}" && \
        curl -sSfL https://github.com/"${ORG}"/"${REPO}"/archive/v"${VERSION}".tar.gz | tar xz && cd "${ARTIFACT}"-* && \
        ./configure \
            --with-features=huge \
            --enable-multibyte \
            --enable-rubyinterp=yes \
            --enable-pythoninterp=yes \
            --enable-python3interp=yes \
            --enable-perlinterp=yes \
            --enable-luainterp=yes \
            --enable-cscope \
            --prefix=/usr/local && \
        make VIMRUNTIMEDIR=/usr/local/share/vim/vim82 && make install && cd .. && rm -fr vim-* && \
        ln -sf /usr/local/bin/vim /usr/local/bin/vi

    echo "Install tmux." && \
        curl -sSfL https://github.com/libevent/libevent/releases/download/release-2.1.11-stable/libevent-2.1.11-stable.tar.gz | tar xz && cd libevent-* && \
        ./configure && make && make install && cd .. && rm -fr libevent-* && \
        curl -sSfL https://github.com/tmux/tmux/releases/download/3.0a/tmux-3.0a.tar.gz | tar xz && cd tmux-* && \
        ./configure --prefix=/usr/local && make && make install && cd .. && rm -fr tmux-*

    echo "Install zsh." && \
        export ORG="zsh-users" && export REPO="zsh" && export VERSION="5.7.1" && export ARTIFACT="${REPO}" && \
        curl -sSfL https://github.com/"${ORG}"/"${REPO}"/archive/"${ARTIFACT}"-"${VERSION}".tar.gz | tar xz && cd "${ARTIFACT}"-* && \
        ./Util/preconfig && \
        ./configure \
            --prefix=/usr/local \
            --mandir=/usr/local/share/man \
            --bindir=/usr/local/bin \
            --infodir=/usr/local/share/info \
            --enable-maildir-support \
            --enable-etcdir=/usr/local/etc/zsh \
            --enable-function-subdirs \
            --enable-site-fndir=/usr/local/share/zsh/site-functions \
            --enable-fndir=/usr/local/share/zsh/functions \
            --with-tcsetpgrp \
            --with-term-lib="ncursesw" \
            --enable-cap \
            --enable-pcre \
            --enable-readnullcmd=pager \
            --enable-custom-patchlevel=Debian \
            LDFLAGS="-Wl,--as-needed -g" && \
        make && make install && echo "/usr/local/bin/zsh" >> /etc/shells && cd .. && rm -fr zsh-*

    echo "Install redis-cli tools." && \
        curl -sSfL http://download.redis.io/releases/redis-5.0.7.tar.gz | tar xz && cd redis-* && \
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
