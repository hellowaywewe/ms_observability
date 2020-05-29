 #! /bin/sh
    
 cd $HOME/ebpf_mindspore

 DOCKER_NAME=${DOCKER_NAME:-ubuntu_bcc}
 TAG=${TAG:-master}

 # 将host主机的内核目录挂载到容器目录中
 docker run -d --name ${DOCKER_NAME} --privileged \
            -v $(pwd):/mnt/ms_observability \
            -v /lib/modules:/lib/modules:ro \
            -v /usr/src:/usr/src:ro \
            -v /boot/:/boot:ro \
            -v /sys/kernel/debug:/sys/kernel/debug \
            ${DOCKER_NAME}:${TAG} sleep 3600d
