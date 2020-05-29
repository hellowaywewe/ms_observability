# MindSpore Observability

#### Experimental notice: This project is still experimental, only shows how to run eBPF code in container to trace kernel, and provides some simple examples to teach you how to use [BCC](https://github.com/iovisor/bcc.git) to develop eBPF tools and use [ebpf_exporter](https://github.com/cloudflare/ebpf_exporter.git) to visualize the system tracing metrics. Right now here's an simple attempt to combine [MindSpore](https://github.com/mindspore-ai/mindspore.git) with eBPF, real practical examples are expected in the near future.

- [MindSpore Observability](#mindspore-observability)
  - [Introduction of ms_observability](#introduction-of-ms-observability)
  - [Getting Started](#getting-started)
    - [Run eBPF code in container to probe kernel metrics](#run-ebpf-code-in-container-to-probe-kernel-metrics)
    - [Visualize kernel metrics in the unified format of Prometheus](#visualize-kernel-metrics-in-the-unified-format-of-prometheus)
    - [A simple attempt to combine MindSpore with eBPF](#a-simple-attempt-to-combine-mindspore-with-ebpf)
  - [Future Work](#future-work)

## Introduction of ms_observability
 
MindSpore is a new open source deep learning training/inference framework that
could be used for mobile, edge and cloud scenarios. MindSpore is designed to
provide development experience with friendly design and efficient execution for
the data scientists and algorithmic engineers, native support for Ascend AI
processor, and software hardware co-optimization.

Currently, the problem with all deep learning job is that the AI training process
is invisible. While running a AI job by using the MindSpore, we don't know how
it is layered, don’t know which CPU core it runs on , even don’t know what kernel
functions it calls and how to jump. Once the task has bottlenecks, developers tend
to choose to use some common monitoring tools to analyze, but these usually have
blind spots and they are inflexible, such as: they can get long-lived processes
information, but for some short-lived processes, often can't capture which leads
to loss of information, a lot of these processes are actually on the consumption
of resources.

To solve the gap, the project ms_observability combines the MindSpore with the 
new technology eBPF to improve the observability of the AI kernel throughout the
training and reasoning process. eBPF can make the kernel fully programmable and
dynamically run a mini programs on a wide variety of kernel events, which can
empower non-kernel developers to customize their own tracing codes to solve real
problems they met, which means that it can keep watch over the whole kernel states
of the AI job to provide more detailed context to further analyze your system and 
application.

## Getting Started

### Prerequisites
- [Ubuntu](http://releases.ubuntu.com/16.04/): `16.04.6 LTS`
- [docker](https://github.com/docker/docker-ce/tags): `v19.03.8`
- [MindSpore](https://github.com/mindspore-ai/mindspore/releases/tag/v0.2.0-alpha): `v0.2.0-alpha`


### Run eBPF code in container to probe kernel metrics

#### Download ms_observability code
```shell
cd $HOME
git clone https://github.com/hellowaywewe/ms_observability.git
```

#### Build and run ebpf_bcc_exporter container
```shell
cd $HOME/ms_observability/docker
docker build -f Dockerfile -t ebpf_bcc_exporter:latest .
cd $HOME/ms_observability
DOCKER_NAME=ebpf_bcc_exporter TAG=latest ./run_docker.sh   // mount the host kernel to the container
```

#### Show the kernel queue IO latency metrics (simple example, showing how to use bcc to develop eBPF code and probe kernel)
```shell
docker exec -it ebpf_bcc_exporter /bin/bash     // Container interactive operation
cd /mnt/ms_observability/ebpf_example && ./io-latency.py 1 2
```

### Visualize kernel metrics in the unified format of Prometheus

#### Show the kernel queue IO latency metrics (simple example, show how to use ebpf_exporter to configure and visualize metrics)
```shell
~/go/bin/ebpf_exporter --config.file=/mnt/ms_observability/exporter_example/io-latency.yaml
```

#### Use the `curl` command to verify that the visual metrics are properly captured
```shell
docker inspect ebpf_bcc_exporter | grep IPAddress  // Query the IP of the container
curl http://<yourContainerIP>:9435/metrics
```

### A simple attempt to combine MindSpore with eBPF

When executing MindSpore LENET job in the host, if the kernel function
“blk_account_io_done” is called, the words “Hello World” will be printed,
if not, print nothing.

#### Run the lenet-io.py code in container


```shell
docker exec -it ebpf_bcc_exporter /bin/bash
cd /mnt/ms_observability/ebpf_example
./lenet-io.py
``` 

#### Run the MindSpore lenet training job in the host (Required MindSpore v0.2.0-alpha Env)
```shell
cd $HOME && git clone https://github.com/mindspore-ai/docs.git
conda activate mindspore && cd $HOME/docs/tutorials/tutorial_code/
python lenet.py --device_target="CPU"
```

## Future Work

Currently the ms_observability is in the early stages of experiment, in the
future, most importantly, we should analyze what to do in AI scenarios and
which can be used and traced from the thousands of available kernel events. 
And then collaborate with other open source communities:
1. Work with the iovisor/bcc project to develop AI observability tools based on eBPF.
2. Enable MindSpore to support eBPF AI observability tools.
3. Work with the Prometheus and ebpf_exporter project to visualize the AI kernel metrics.

