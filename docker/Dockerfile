FROM ubuntu:16.04

MAINTAINER Ivy He <simple_hlw@163.com>

ENV CONFIG_BPF=y
ENV CONFIG_BPF_SYSCALL=y
ENV CONFIG_NET_CLS_BPF=m
ENV CONFIG_NET_ACT_BPF=m
ENV CONFIG_BPF_JIT=y
ENV CONFIG_HAVE_BPF_JIT=y
ENV CONFIG_HAVE_EBPF_JIT=y
ENV CONFIG_BPF_EVENTS=y
ENV CONFIG_NET_SCH_SFQ=m
ENV CONFIG_NET_ACT_POLICE=m
ENV CONFIG_NET_ACT_GACT=m
ENV CONFIG_DUMMY=m
ENV CONFIG_VXLAN=m

RUN apt-get update && apt-get -y install bison build-essential git cmake flex libedit-dev \
    libllvm3.7 llvm-3.7-dev libclang-3.7-dev python3.5 zlib1g-dev libelf-dev &&\
    apt-get -y install luajit luajit-5.1-dev &&\
    apt-get -y install vim wget curl  &&\
    ln -s /usr/bin/python3.5 /usr/bin/python &&\
    git clone https://github.com/iovisor/bcc.git &&\
    mkdir bcc/build && cd bcc/build &&\
    cmake .. &&\
    make && make install &&\
    cmake -DPYTHON_CMD=python .. &&\
    cd src/python/ &&\
    make && make install

RUN wget https://dl.google.com/go/go1.14.3.linux-amd64.tar.gz &&\
    mkdir -p /usr/local/go && tar -zxvf go1.14.3.linux-amd64.tar.gz -C /usr/local/go --strip-components 1 &&\
    rm go1.14.3.linux-amd64.tar.gz

ENV GO_HOME=/usr/local/go
ENV PATH=$GO_HOME/bin:$PATH

COPY ./main.go /
RUN git clone https://github.com/cloudflare/ebpf_exporter.git &&\
    mv /main.go  /ebpf_exporter/cmd/ebpf_exporter/main.go &&\
    cd ebpf_exporter && go get -u -v github.com/cloudflare/ebpf_exporter/...

