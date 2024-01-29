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
containerID=$(docker ps -a | grep "${project}" | awk '{print $1}')

echo "Container ID: $containerID"

if [ "$containerID" != "" ]; then
  docker stop "$containerID"
  docker rm "$containerID"
fi

# 确保没有同名image
imageName="$harbor_addr/$harbor_repo/$project:$version"

echo "Image Name: $imageName"

tag=$(docker images | grep "${project}" | awk '{print $2}')

echo "Image Tag: $tag"

if [ "$tag" != "$version" ]; then
  docker rmi "$imageName"
fi

# 创建远程连接
docker-machine create \
  --driver generic \
 --generic-ip-address="${deployment_server_ip}" \
 --generic-ssh-key /var/ssh/max.pem \
 tmp_deployment_mechine

# 远程登录harbor
docker-machine ssh tmp_deployment_mechine \
  docker login -u admin -p Harbor12345 "$harbor_addr"

docker-machine ssh tmp_deployment_mechine \
  docker pull "$imageName"

docker-machine ssh tmp_deployment_mechine \
  docker run \
  --name="$project" \
  -d \
  --restart=always \
  --privileged=true \
  -p "${host_port}:${container_port}" \
  -e ENV="${environment}" \
  -v /root/www:/usr/src/app/www \
  -v /root/logs/"${project}":/usr/src/app/logs \
  "$imageName"

docker-machine stop tmp_deployment_mechine
docker-machine rm tmp_deployment_mechine

echo "SUCCESS: Container Created"
