FROM debian:9.5

LABEL maintainer="rluckie@cisco.com" \
      kdk=""


#######################################
# Copy local files
COPY tmp/python-requirements/** /tmp/python-requirements/
COPY usr/local/bin/** /usr/local/bin/

#######################################
# Install OS packages

RUN echo "Install OS packages" && \
    apt-get -y update && apt-get --no-install-recommends -y install \
        apache2-utils \
        apt-transport-https \
        autoconf \
        bash-completion \
        bc \
        bridge-utils \
        build-essential \
        ca-certificates \
        colordiff \
        ctags \
        curl \
        dc \
        dhcpdump \
        dialog \
        dnsutils \
        dos2unix \
        emacs25-nox \
        fonts-powerline \
        gettext \
        gnupg \
        gnupg2 \
        htop \
        less \
        libcurl4-openssl-dev \
        libevent-2.0-5 \
        libncurses5-dev \
        libz-dev \
        locales \
        man \
        moreutils \
        nmap \
        ntp \
        ntpdate \
        openssh-server \
        perl \
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
         apt-get update -y && apt-get --no-install-recommends -y install \
           docker-ce \
           google-cloud-sdk && \
   apt-get -y clean && apt-get -y autoremove && rm -rf /var/lib/apt/lists/*

#######################################
# Install/Build apps that are not provided by the OS packages.
# We sometimes install or build from source to get specific versions.

RUN echo "Install/build apps that are not provided by the OS packages." && \
    # git (needs to build first as a dependency) && \
        cd /tmp && curl -sSL https://github.com/git/git/archive/v2.18.0.tar.gz | tar xz && \
        cd git-* && make configure && ./configure --prefix=/usr/local && make && make install && \
        cd .. && rm -fr git-* && \
    # dep && \
        curl -sSLo dep https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 && \
        chmod a+x dep && mv dep /usr/local/bin && \
    # direnv && \
        curl -sSLo direnv https://github.com/direnv/direnv/releases/download/v2.17.0/direnv.linux-amd64 && \
        chmod +x direnv && mv direnv /usr/local/bin && \
    # drone cli && \
        curl -sSL https://github.com/drone/drone-cli/releases/download/v0.8.6/drone_linux_amd64.tar.gz | tar xz && \
        chmod +x drone && mv drone /usr/local/bin && \
    # easy-rsa && \
        curl -sSL https://github.com/OpenVPN/easy-rsa/releases/download/v3.0.4/EasyRSA-3.0.4.tgz | tar xz && \
        chmod +x EasyRSA-* && mv EasyRSA-* /usr/local/bin/easyrsa && \
    # gomplate && \
        curl -sSLo gomplate https://github.com/hairyhenderson/gomplate/releases/download/v2.7.0/gomplate_linux-amd64 && \
        chmod a+x gomplate && mv gomplate /usr/local/bin && \
    # gitslave && \
        curl -sSL https://github.com/joelpurra/gitslave/archive/v2.0.2-joelpurra-51-g203467d.tar.gz | tar xz && \
        chmod a+x gitslave-*/gits && cp gitslave-*/gits /usr/local/bin && rm -rf gitslave-* && \
    # golang && \
        curl -sSL https://dl.google.com/go/go1.10.3.linux-amd64.tar.gz | tar -C /usr/local -xz && \
        mkdir -p /go && chmod a+rw /go && \
    # helm && \
        curl -sSL https://storage.googleapis.com/kubernetes-helm/helm-v2.8.2-linux-amd64.tar.gz | tar xz && \
        chmod +x linux-amd64/helm && mv linux-amd64/helm /usr/local/bin/helm-2.8.2 && \
        rm -fr linux-amd64 && \
        curl -sSL https://storage.googleapis.com/kubernetes-helm/helm-v2.9.1-linux-amd64.tar.gz | tar xz && \
        chmod +x linux-amd64/helm && mv linux-amd64/helm /usr/local/bin/helm-2.9.1 && \
        rm -fr linux-amd64 && \
        curl -sSL https://storage.googleapis.com/kubernetes-helm/helm-v2.10.0-linux-amd64.tar.gz | tar xz && \
        chmod +x linux-amd64/helm && mv linux-amd64/helm /usr/local/bin/helm-2.10.0 && \
        rm -fr linux-amd64 && \
        ln -s /usr/local/bin/helm-2.10.0 /usr/local/bin/helm && \
    # hugo && \
        cd /tmp && curl -sSL https://github.com/gohugoio/hugo/releases/download/v0.48/hugo_0.48_Linux-ARM64.tar.gz | tar xz && \
        chmod +x hugo && mv hugo /usr/local/bin/hugo && \
    # jq && \
        curl -sSLo jq https://github.com/stedolan/jq/releases/download/jq-1.5/jq-linux64 && \
        chmod +x jq && mv jq /usr/local/bin/ && \
    # kops && \
        curl -sSLo kops-1.8.1 https://github.com/kubernetes/kops/releases/download/1.8.1/kops-linux-amd64 && \
        chmod +x kops-1.8.1 && mv kops-1.8.1 /usr/local/bin/ && \
        curl -sSLo kops-1.9.2 https://github.com/kubernetes/kops/releases/download/1.9.2/kops-linux-amd64 && \
        chmod +x kops-1.9.2 && mv kops-1.9.2 /usr/local/bin/ && \
        curl -sSLo kops-1.10.0 https://github.com/kubernetes/kops/releases/download/1.10.0/kops-linux-amd64 && \
        chmod +x kops-1.10.0 && mv kops-1.10.0 /usr/local/bin/ && \
        ln -s /usr/local/bin/kops-1.10.0 /usr/local/bin/kops && \
    # kubectl && \
        curl -sSLo /usr/local/bin/kubectl-1.10.6 https://storage.googleapis.com/kubernetes-release/release/v1.10.6/bin/linux/amd64/kubectl && \
        chmod a+x /usr/local/bin/kubectl-1.10.6 && \
        curl -sSLo /usr/local/bin/kubectl-1.11.2 https://storage.googleapis.com/kubernetes-release/release/v1.11.2/bin/linux/amd64/kubectl && \
        chmod a+x /usr/local/bin/kubectl-1.11.2 && \
        ln -s /usr/local/bin/kubectl-1.11.2 /usr/local/bin/kubectl && \
    # kubetail && \
        curl -sSLo kubetail.zip https://github.com/johanhaleby/kubetail/archive/1.6.1.zip && \
        unzip kubetail.zip && chmod +x kubetail-1.6.1/kubetail && mv kubetail-1.6.1/kubetail /usr/local/bin && \
        rm -f kubetail.zip && \
    # minikube && \
        curl -sSLo minikube https://storage.googleapis.com/minikube/releases/v0.28.2/minikube-linux-amd64 && \
        chmod +x minikube &&  mv minikube /usr/local/bin/ && \
    # minio - mc CLI && \
        curl -sSLO https://dl.minio.io/client/mc/release/linux-amd64/mc && \
        chmod +x mc && mv mc /usr/local/bin && \
    # redis-benchmark, redis-cli && \
        curl -sSLO http://download.redis.io/releases/redis-4.0.10.tar.gz && \
        mkdir -p redis && \
        tar --strip-components=1 -C redis -xzf redis-4.0.10.tar.gz && \
        make -C redis && \
        cp redis/src/redis-cli redis/src/redis-benchmark /usr/local/bin && \
        rm -rf redis && \
    # terraform && \
        curl -sSLo terraform.zip https://releases.hashicorp.com/terraform/0.11.8/terraform_0.11.8_linux_amd64.zip && \
        unzip terraform.zip && chmod +x terraform && mv terraform /usr/local/bin && rm -f terraform.zip && \
    # testssl && \
        git clone -q --depth 1 https://github.com/drwetter/testssl.sh.git /usr/local/share/testssl.sh && \
        ln -s /usr/local/share/testssl.sh/testssl.sh /usr/local/bin/testssl && chmod +x /usr/local/bin/testssl && \
    # tmux && \
        curl -sSL https://github.com/libevent/libevent/releases/download/release-2.0.22-stable/libevent-2.0.22-stable.tar.gz | tar xz && \
        cd libevent-* && \
        ./configure && make &&  make install && \
        cd .. && rm -fr libevent-* && \
        cd /tmp && curl -sSL https://github.com/tmux/tmux/releases/download/2.7/tmux-2.7.tar.gz | tar xz && \
        cd tmux-* && \
        ./configure --prefix=/usr/local && \
        make &&  make install && \
        cd .. && rm -fr tmux-* && \
    # vim && \
       curl -sSL https://github.com/vim/vim/archive/v8.1.0290.tar.gz | tar xz && cd vim-* && \
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
       make VIMRUNTIMEDIR=/usr/local/share/vim/vim81 && make install && cd .. && rm -fr vim-* && \
    # yadm && \
        curl -sSLo /usr/local/bin/yadm https://github.com/TheLocehiliosan/yadm/raw/master/yadm && \
        chmod a+x /usr/local/bin/yadm && \
    # zsh && \
        curl -sSL https://sourceforge.net/projects/zsh/files/zsh/5.5.1/zsh-5.5.1.tar.gz/download | tar xz && \
        cd zsh-* && \
        ./configure --with-tcsetpgrp --prefix=/usr/local && make && make install && \
         echo "/usr/local/bin/zsh" >> /etc/shells && \
        cd .. && rm -fr zsh-*

#######################################
# go get installs

ARG GOPATH=/go

RUN echo "go get installs" && \
    /usr/local/go/bin/go get github.com/spf13/cobra/cobra && \
    /usr/local/go/bin/go get github.com/kubernetes-incubator/cri-tools/cmd/crictl && \
    /usr/local/go/bin/go get github.com/golang/lint/golint && \
    /usr/local/go/bin/go get github.com/gpmgo/gopm && \
    /usr/local/go/bin/go get github.com/vmware/govmomi/govc && \
    /usr/local/go/bin/go get github.com/github/hub && \
    /usr/local/go/bin/go get github.com/cisco-sso/mh && ln -sf /go/bin/mh /go/bin/multihelm && \
    /usr/local/go/bin/go get github.com/mikefarah/yq

#######################################
# Install python-based utils and libs

RUN cd /tmp && \
    curl -sSL https://bootstrap.pypa.io/get-pip.py -o get-pip.py && \
    python2 get-pip.py && \
    python3 get-pip.py && \
    rm get-pip.py && \
    pip2.7 install -U setuptools && \
    pip3.5 install -U setuptools && \
    pip2.7 install -r /tmp/python-requirements/pip2.7.txt --no-cache-dir --ignore-installed six && \
    pip3.5 install -r /tmp/python-requirements/pip3.5.txt --no-cache-dir

#######################################
# Configure systemd and disable a bunch of systemd.units
#  Mostly taken from:
#  https://github.com/dramaturg/docker-debian-systemd/blob/master/Dockerfile

RUN cd /lib/systemd/system/sysinit.target.wants/ && \
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
        sed -ri /etc/systemd/journald.conf -e 's!^#?Storage=.*!Storage=volatile!'

#######################################
# Misc
RUN echo "Misc configuration" && \
    # Set locale && \
    localedef -i en_US -f UTF-8 en_US.UTF-8 && \
    # Configure openssh-server && \
    sed -i 's/#Port 22/Port 2022/' /etc/ssh/sshd_config && \
    # Configure docker daemon to support docker in docker && \
    mkdir /etc/docker && echo '{"storage-driver": "vfs"}' > /etc/docker/daemon.json && \
    # Generate bash-completion for kubectl and helm && \
    kubectl completion bash > /etc/bash_completion.d/kubectl_completion.sh && \
    helm completion bash > /etc/bash_completion.d/helm_completion.sh

#######################################
# Ensure systemd starts, which subsequently starts ssh and docker

EXPOSE 2022
CMD ["/lib/systemd/systemd"]