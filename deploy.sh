#! /bin/sh

harbor_addr=$1
harbor_repo=$2
project=$3
version=$4
container_port=$5
host_port=$6
environment=$7
deployment_server_ip=$8
static_folder_name=$9

# 确保没有container在运行
containerID=$(docker -H tcp://$deployment_server_ip:2375 ps -a | grep "${project}" | awk '{print $1}')

echo "History Container ID: $containerID"

if [ "$containerID" != "" ]; then
  docker -H tcp://$deployment_server_ip:2375 stop "$containerID"
  docker -H tcp://$deployment_server_ip:2375 rm "$containerID"
fi

# 清除同名image
imageIDRemote=$(docker -H tcp://$deployment_server_ip:2375 images | grep "${project}" | awk '{print $3}' | head -n 1)

echo "History Image ID Remote: $imageIDRemote"

if [ "$imageIDRemote" != "" ]; then
  docker -H tcp://$deployment_server_ip:2375 rmi "$imageIDRemote"
fi

# 即将部署的镜像
imageName="$harbor_addr/$harbor_repo/$project:$version"

echo "Image Name: $imageName"

# 远程登录harbor
docker -H tcp://$deployment_server_ip:2375 login -u admin -p Harbor12345 "$harbor_addr"

docker -H tcp://$deployment_server_ip:2375 pull "$imageName"

docker -H tcp://$deployment_server_ip:2375 run \
  --name="$project" \
  -d \
  --restart=always \
  --privileged=true \
  -p "${host_port}:${container_port}" \
  -e ENV="${environment}" \
  -v /root/static/"${static_folder_name}"/www:/usr/src/app/www \
  -v /root/logs/"${project}":/usr/src/app/log \
  "$imageName"

echo "SUCCESS: Container Created"

# 部署后，清除jenkins服务器产生的image
imageIDLocal=$(docker images | grep "${project}" | awk '{print $3}' | head -n 1)

echo "History Image ID Local: $imageIDLocal"

if [ "$imageIDLocal" != "" ]; then
  docker rmi -f "$imageIDLocal"
fi
