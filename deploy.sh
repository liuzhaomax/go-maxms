#! /bin/sh

harbor_addr=$1
harbor_repo=$2
project=$3
version=$4
container_port=$5
host_port=$6

imageName=$harbor_addr/$harbor_repo/$project:$version

echo $imageName

containerID=`docker ps -a | grep ${project} | awk "{print $1}"`

echo $containerID

if [ "$containerID" != "" ]; then
  docker stop $containerID
  docker rm $containerID
fi

tag=`docker images | grep ${project} | awk "{print $2}"`

echo $tag

# shellcheck disable=SC2039
if [[ $tag =~ $version ]]; then
  docker rmi $imageName
fi

docker login -u admin -p Harbor12345 $harbor_addr

docker pull $imageName


# shellcheck disable=SC2039
# shellcheck disable=SC2053
if [[ $tag == $version ]]; then
  docker stop project
  docker rm project
fi

docker run --name=${project} -d --restart=always --privileged=true -p ${host_port}:${container_port} \
  -v /home/opc/static:/usr/src/app/static \
  /home/opc/logs/${project}/golog.txt:/usr/src/app/golog.txt \
  ${imageName}

echo "SUCCESS: Container Created"