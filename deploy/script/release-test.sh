need_start_server_shell=(
   # rpc
    user-rpc-test.sh
    social-rpc-test.sh
    im-rpc-test.sh
    im-ws-test.sh

    # api
    user-api-test.sh
    social-api-test.sh
    im-api-test.sh

    #task
    task-mq-test.sh
)

for i in ${need_start_server_shell[*]} ; do
    chmod +x $i
    ./$i
done


docker ps

# 查找etcd中所有key，判断服务是否成功注册
docker exec -it etcd etcdctl get --prefix ""