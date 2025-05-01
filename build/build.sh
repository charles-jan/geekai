#!/bin/bash

version=$1
arch=${2:-amd64}
registry="crpi-1rhmnyrson8h7iaj.ap-southeast-1.personal.cr.aliyuncs.com/charles_jan"

# build go api program
cd ../api
make clean $arch

# build web app
cd ../web
npm run build

cd ../build

# remove docker image if exists
docker rmi -f $registry/geekai-api:$version-$arch

# build docker image for geekai-go
docker build -t $registry/geekai-api:$version-$arch -f dockerfile-api-go ../

# build docker image for geekai-web
docker rmi -f $registry/geekai-web:$version-$arch
docker build -t $registry/geekai-web:$version-$arch -f dockerfile-web-$arch ../

if [ "$3" = "push" ]; then
  docker push $registry/geekai-api:$version-$arch
  docker push $registry/geekai-web:$version-$arch
fi
