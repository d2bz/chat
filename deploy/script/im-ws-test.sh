#!/bin/bash
reso_addr='d2bz.cn-shanghai.personal.cr.aliyuncs.com/my-zero-im/im-ws-dev'
tag='latest'

pod_idb="106.14.194.111"

container_name="chat-im-ws-test"

docker stop ${container_name}

docker rm ${container_name}

docker rmi ${reso_addr}:${tag}

docker pull ${reso_addr}:${tag}


# 如果需要指定配置文件的
docker run -p 10090:10090 -e POD_IP=${pod_idb} --name=${container_name} -d ${reso_addr}:${tag}