#! /bin/sh

harbor_addr=$1
harbor_repo=$2
project=$3
version=$4
container_port=$5
host_port=$6

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

# shellcheck disable=SC2039
if [ "$tag" != "$version" ]; then
  docker rmi "$imageName"
fi

# 登录harbor
docker login -u admin -p Harbor12345 "$harbor_addr"

docker pull "$imageName"

docker run \
  --name="$project" \
  -d \
  --restart=always \
  --privileged=true \
  -p "${host_port}:${container_port}" \
  -v /root/www:/usr/src/app/www \
  -v /root/logs/"${project}":/usr/src/app/logs \
  "$imageName"

echo "SUCCESS: Container Created"
