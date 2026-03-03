#!/bin/bash
reso_addr='d2bz.cn-shanghai.personal.cr.aliyuncs.com/my-zero-im/social-rpc-dev'
tag='latest'

pod_idb="106.14.194.111"

container_name="chat-social-rpc-test"

docker stop ${container_name}

docker rm ${container_name}

docker rmi ${reso_addr}:${tag}

docker pull ${reso_addr}:${tag}


# 如果需要指定配置文件的
docker run -p 10001:10001 -e POD_IP=${pod_idb} --name=${container_name} -d ${reso_addr}:${tag}