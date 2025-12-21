need_start_server_shell=(
  # rpc
  user-rpc-test.sh
)

for i in ${need_start_server_shell[*]} ; do
    chmod +x $i
    ./$i
done


docker ps

# 查找etcd中所有key，判断服务是否成功注册
docker exec -it etcd etcdctl get --prefix ""