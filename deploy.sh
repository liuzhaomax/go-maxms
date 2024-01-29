#! /bin/sh

harbor_addr=$1
harbor_repo=$2
project=$3
version=$4
container_port=$5
host_port=$6
environment=$7
deployment_server_ip=$8

# 确保没有container在运行
containerID=$(docker -H tcp://$deployment_server_ip:2375 ps -a | grep "${project}" | awk '{print $1}')

echo "Container ID: $containerID"

if [ "$containerID" != "" ]; then
  docker -H tcp://$deployment_server_ip:2375 stop "$containerID"
  docker -H tcp://$deployment_server_ip:2375 rm "$containerID"
fi

# 确保没有同名image
imageName="$harbor_addr/$harbor_repo/$project:$version"

echo "Image Name: $imageName"

tagRemote=$(docker -H tcp://$deployment_server_ip:2375 images | grep "${project}" | awk '{print $2}')
tagLocal=$(docker images | grep "${project}" | awk '{print $2}')

echo "Image Remote Tag: $tagRemote"
echo "Image Local Tag: $tagLocal"

if [ "$tagRemote" != "$version" ]; then
  docker rmi -f "$tagLocal"
  docker -H tcp://$deployment_server_ip:2375 rmi "$tagRemote"
fi

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
  -v /root/www:/usr/src/app/www \
  -v /root/logs/"${project}":/usr/src/app/logs \
  "$imageName"

echo "SUCCESS: Container Created"
