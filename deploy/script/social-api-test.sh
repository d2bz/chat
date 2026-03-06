#!/bin/bash
reso_addr='d2bz.cn-shanghai.personal.cr.aliyuncs.com/my-zero-im/social-api-dev'
tag='latest'

container_name="chat-social-api-test"

docker stop ${container_name}

docker rm ${container_name}

docker rmi ${reso_addr}:${tag}

docker pull ${reso_addr}:${tag}


# 如果需要指定配置文件的
docker run -p 8881:8881 --name=${container_name} -d ${reso_addr}:${tag}