FROM centos:7

LABEL maintainer="rluckie@cisco.com"

COPY files/yum.repos.d/** /etc/yum.repos.d/
COPY files/python-requirements/** /tmp/python-requirements/
COPY files/provision-user /usr/local/bin/provision-user
COPY files/start-dockerd /usr/local/bin/start-dockerd
COPY files/awake /usr/local/bin/awake

RUN yum update --nogpgcheck -y && \
    yum install --nogpgcheck -y epel-release deltarpm yum-utils https://centos7.iuscommunity.org/ius-release.rpm && \
    yum group install --nogpgcheck -y "Development Tools" && \
    yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo && \
    yum install -y \
        bc \
        bind \
        bind-utils \
        colordiff \
        ctags \
        curl-devel \
        device-mapper-persistent-data \
        docker-ce \
        dos2unix \
        emacs \
        fontconfig \
        gettext-devel \
        glibc-static \
        google-cloud-sdk \
        graphviz \
        httpd-tools \
        htop \
        initscripts \
        kubectl \
        libevent \
        libevent-devel \
        lua \
        lua-devel \
        luajit \
        luajit-devel \
        lvm2 \
        most \
        ncurses \
        ncurses-devel \
        nmap \
        openssl-devel \
        perl \
        perl-CPAN \
        perl-devel \
        perl-ExtUtils-CBuilder \
        perl-ExtUtils-Embed \
        perl-ExtUtils-ParseXS \
        perl-ExtUtils-XSpp \
        pigz \
        python \
        python-devel \
        python-pip \
        python36u \
        python36u-devel \
        python36u-pip \
        ruby \
        ruby-devel \
        screen \
        socat \
        sudo \
        sysvinit-tools \
        tcl-devel \
        telnet \
        tmux2u \
        traceroute \
        tree \
        wget \
        which \
        xcb-util \
        xdotool \
        xorg-x11-font \
        xorg-x11-fonts* \
        xorg-x11-server-Xvfb \
        xorg-x11-twm \
        xorg-x11-xinit \
        xorg-x11-xinit-session \
        xterm \
        yadr \
        zlib-devel && \
    sudo ln -s /usr/bin/xsubpp /usr/share/perl5/ExtUtils/xsubpp && \
    yum clean all && \
    rm -rf /var/cache/yum

# zsh
RUN cd /tmp && curl -L https://sourceforge.net/projects/zsh/files/zsh/5.5.1/zsh-5.5.1.tar.gz/download | tar xz && \
    cd zsh-* && \
    ./configure --with-tcsetpgrp --prefix=/usr/local && make && sudo make install && \
    sudo echo "/usr/local/bin/zsh" >> /etc/shells && \
    cd .. && rm -fr zsh-*

# git
RUN cd /tmp && curl -L https://github.com/git/git/archive/v2.17.0.tar.gz | tar xz && \
    cd git-* && \
    make configure && \
    ./configure --prefix=/usr/local && \
    make && sudo make install && \
    cd .. && rm -fr git-*

# tmux
RUN cd /tmp && curl -L https://github.com/libevent/libevent/releases/download/release-2.0.22-stable/libevent-2.0.22-stable.tar.gz | tar xz && \
    cd libevent-* && \
    ./configure && make && sudo make install && \
    cd .. && rm -fr libevent-* && \
    cd /tmp && curl -L https://github.com/tmux/tmux/releases/download/2.7/tmux-2.7.tar.gz | tar xz && \
    cd tmux-* && \
    ./configure --prefix=/usr/local && \
    make && sudo make install && \
    cd .. && rm -fr tmux-*

# vim
RUN cd /tmp && curl -L https://github.com/vim/vim/archive/v8.0.1736.tar.gz | tar xz && \
    cd vim-* && \
    ./configure \
        --with-features=huge \
        --enable-cscope \
        --enable-rubyinterp \
        --enable-luainterp \
        --enable-perlinterp \
        --enable-pythoninterp \
        --enable-python3interp \
        --enable-tclinterp \
        --enable-multibyte \
        --prefix=/usr/local && \
    make VIMRUNTIMEDIR=/usr/local/share/vim/vim80 && \
    sudo make install && \
    cd .. && rm -fr vim-*

# jq
RUN cd /tmp && \
    curl -Lo jq https://github.com/stedolan/jq/releases/download/jq-1.5/jq-linux64 && \
    chmod +x jq && mv jq /usr/local/bin/

# kops
RUN curl -Lo kops-1.9.0 https://github.com/kubernetes/kops/releases/download/1.9.0/kops-linux-amd64 && \
    chmod +x kops-1.9.0 && mv kops-1.9.0 /usr/local/bin/ && \
    curl -Lo kops-1.8.1 https://github.com/kubernetes/kops/releases/download/1.8.1/kops-linux-amd64 && \
    chmod +x kops-1.8.1 && mv kops-1.8.1 /usr/local/bin/

# helm
RUN curl -L https://storage.googleapis.com/kubernetes-helm/helm-v2.8.2-linux-amd64.tar.gz | tar xz && \
    chmod +x linux-amd64/helm && mv linux-amd64/helm /usr/local/bin/ && \
    rm -fr linux-amd64

# minio - mc CLI
RUN cd /tmp && \
    curl -LO https://dl.minio.io/client/mc/release/linux-amd64/mc && \
    chmod +x mc && mv mc /usr/local/bin

# direnv
RUN cd /tmp && \
    curl -Lo direnv https://github.com/direnv/direnv/releases/download/v2.15.2/direnv.linux-amd64 && \
    chmod +x direnv && mv direnv /usr/local/bin

# terraform
RUN cd /tmp && \
    curl -Lo terraform.zip https://releases.hashicorp.com/terraform/0.11.7/terraform_0.11.7_linux_amd64.zip && \
    unzip terraform.zip && chmod +x terraform && mv terraform /usr/local/bin && rm -f terraform.zip

# golang
RUN curl -L https://dl.google.com/go/go1.10.1.linux-amd64.tar.gz | tar -C /usr/local -xz

# minikube
RUN curl -Lo minikube https://storage.googleapis.com/minikube/releases/v0.26.1/minikube-linux-amd64 && \
    chmod +x minikube && sudo mv minikube /usr/local/bin/

# testssl
RUN git clone --depth 1 https://github.com/drwetter/testssl.sh.git /usr/local/share/testssl.sh && \
    ln -s /usr/local/share/testssl.sh/testssl.sh /usr/local/bin/testssl && chmod +x /usr/local/bin/testssl

# easy-rsa
RUN curl -L https://github.com/OpenVPN/easy-rsa/releases/download/v3.0.4/EasyRSA-3.0.4.tgz | tar xz && \
    chmod +x EasyRSA-* && mv EasyRSA-* /usr/local/bin/easyrsa

# python requirements
RUN pip2.7 install -U pip && pip3.6 install -U pip && \
    pip2.7 install -r /tmp/python-requirements/pip2.7.txt --no-cache-dir --ignore-installed six && \
    pip3.6 install -r /tmp/python-requirements/pip3.6.txt --no-cache-dir

CMD ["/bin/bash"]
