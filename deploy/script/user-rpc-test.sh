#!/bin/bash
reso_addr='d2bz.cn-shanghai.personal.cr.aliyuncs.com/my-zero-im/user-rpc-dev'
tag='latest'

pod_idb="106.14.194.111"

container_name="chat-user-rpc-test"

docker stop ${container_name}

docker rm ${container_name}

docker rmi ${reso_addr}:${tag}

docker pull ${reso_addr}:${tag}


# 如果需要指定配置文件的
# -e POD_IP=${pod_idb} 是用来在容器内设置一个名为 POD_IP 的环境变量，用于etcd配置
docker run -p 10000:10000 -e POD_IP=${pod_idb} --name=${container_name} -d ${reso_addr}:${tag}